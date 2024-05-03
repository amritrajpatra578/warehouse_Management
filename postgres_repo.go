package main

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type PostgresRepo struct {
	db *bun.DB
}

func NewPostgresRepo(db *bun.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (p *PostgresRepo) Create(product Product) error {
	_, err := p.db.NewInsert().Model(&product).Exec(context.Background())

	if err != nil {
		var pgdriverErr pgdriver.Error
		if errors.As(err, &pgdriverErr) {
			sqlErrCode := pgdriverErr.Field('C')
			if sqlErrCode == "23505" {
				return errDuplicateId
			}
		}
		return err
	}
	return nil
}

func (p *PostgresRepo) Update(product Product) error {
	result, err := p.db.NewUpdate().
		Model(&product).
		Column("id", "brand", "category", "quantity", "price", "updated_at").
		Where("id = ?", product.Id).
		Exec(context.Background())

	if err != nil {
		log.Println("error while update product in postgres:", err)
		return err
	}

	rowsAffect, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		log.Println("error while getting rows affected:", rowsErr)
		return rowsErr
	}
	if rowsAffect == 0 {
		return errProductNotFound
	}
	return nil
}

func (p *PostgresRepo) GetById(id int) (Product, error) {
	var product Product
	if err := p.db.NewSelect().Model(&product).Where("id = ?", id).Scan(context.Background()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Product{}, errProductNotFound
		}
		return Product{}, err
	}

	return product, nil
}

func (p *PostgresRepo) GetAll() ([]Product, error) {
	products := []Product{}

	err := p.db.NewSelect().Model(&products).Scan(context.Background())
	if err != nil {
		return []Product{}, err
	}

	return products, nil
}

func (p *PostgresRepo) Delete(id int) error {
	var product Product
	result, err := p.db.NewDelete().Model(&product).Where("id = ?", id).Exec(context.Background())
	if err != nil {
		log.Println("err while deleting Product: ", err)
		return err
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		log.Println("error while getting rows affected:", rowsErr)
		return rowsErr
	}
	if rowsAffected == 0 {
		return errProductNotFound
	}
	return nil
}
