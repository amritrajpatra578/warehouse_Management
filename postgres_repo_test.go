package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
)

func setupPostgres(t *testing.T, fixtureFileName string) *bun.DB {
	db := connectPostgres("postgres", "postgres", "127.0.0.1:5432", "productsdb")
	db.RegisterModel((*Product)(nil))
	t.Cleanup(func() {
		t.Log("closing db", db.Close())
	})
	fixture := dbfixture.New(db, dbfixture.WithTruncateTables())

	err := fixture.Load(context.Background(), os.DirFS("testdata"), fixtureFileName)
	if err != nil {
		t.Fatal("error while loading fixture from testdata:", err)
	}
	return db
}

func TestPostgresRepo_Create(t *testing.T) {
	type args struct {
		product Product
	}
	tests := []struct {
		name         string
		args         args
		wantProducts []Product
		wantErr      error
	}{
		{
			name: "product created",
			args: args{
				product: Product{
					Id:       5,
					Brand:    "B",
					Category: "B",
					Quantity: 2,
					Price:    20,
				},
			},
			wantProducts: []Product{
				{
					Id:        10,
					Brand:     "J",
					Category:  "J",
					Quantity:  10,
					Price:     100,
					CreatedAt: time.Date(2023, 04, 28, 10, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 11, 00, 00, 00, time.UTC),
				},
				{
					Id:        20,
					Brand:     "K",
					Category:  "K",
					Quantity:  20,
					Price:     200,
					CreatedAt: time.Date(2023, 04, 28, 12, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 13, 00, 00, 00, time.UTC),
				},
				{
					Id:       5,
					Brand:    "B",
					Category: "B",
					Quantity: 2,
					Price:    20,
				},
			},
			wantErr: nil,
		},
		{
			name: "product with duplicate id",
			args: args{
				product: Product{
					Id:       10,
					Brand:    "hJ",
					Category: "hJ",
					Quantity: 1024,
					Price:    100,
				},
			},
			wantProducts: []Product{
				{
					Id:        10,
					Brand:     "J",
					Category:  "J",
					Quantity:  10,
					Price:     100,
					CreatedAt: time.Date(2023, 04, 28, 10, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 11, 00, 00, 00, time.UTC),
				},
				{
					Id:        20,
					Brand:     "K",
					Category:  "K",
					Quantity:  20,
					Price:     200,
					CreatedAt: time.Date(2023, 04, 28, 12, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 13, 00, 00, 00, time.UTC),
				},
			},
			wantErr: errDuplicateId,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, "existingData.yaml")
			repo := NewPostgresRepo(db)

			err := repo.Create(tt.args.product)

			assert.ErrorIs(t, err, tt.wantErr, "error while creating product should match the expected error")

			var gotProducts []Product
			gotErr := db.NewSelect().Model(&gotProducts).Scan(context.Background())
			assert.NoError(t, gotErr, "expect no error while getting products")

			assert.Equal(t, tt.wantProducts, gotProducts, "expect same response product")
		})
	}
}

func TestPostgresRepo_Update(t *testing.T) {
	type args struct {
		product Product
	}
	tests := []struct {
		name         string
		args         args
		wantProducts []Product
		wantErr      error
	}{
		{
			name: "update ok",
			args: args{
				product: Product{
					Id:        10,
					Brand:     "A",
					Category:  "A",
					Quantity:  1,
					Price:     10,
					CreatedAt: time.Date(2023, 04, 28, 10, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 11, 00, 00, 00, time.UTC),
				},
			},
			wantProducts: []Product{
				{
					Id:        20,
					Brand:     "K",
					Category:  "K",
					Quantity:  20,
					Price:     200,
					CreatedAt: time.Date(2023, 04, 28, 12, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 13, 00, 00, 00, time.UTC),
				},
				{
					Id:        10,
					Brand:     "A",
					Category:  "A",
					Quantity:  1,
					Price:     10,
					CreatedAt: time.Date(2023, 04, 28, 10, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 11, 00, 00, 00, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name: "product not found",
			args: args{
				product: Product{
					Id:       119,
					Brand:    "C",
					Category: "C",
					Quantity: 9,
					Price:    90,
				},
			},
			wantProducts: []Product{
				{
					Id:        10,
					Brand:     "J",
					Category:  "J",
					Quantity:  10,
					Price:     100,
					CreatedAt: time.Date(2023, 04, 28, 10, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 11, 00, 00, 00, time.UTC),
				},
				{
					Id:        20,
					Brand:     "K",
					Category:  "K",
					Quantity:  20,
					Price:     200,
					CreatedAt: time.Date(2023, 04, 28, 12, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 13, 00, 00, 00, time.UTC),
				},
			},
			wantErr: errProductNotFound,
		},
		{
			name: "created at should not change on update",
			args: args{
				product: Product{
					Id:        10,
					Brand:     "J",
					Category:  "J",
					Quantity:  10,
					Price:     100,
					CreatedAt: time.Date(2023, 04, 28, 15, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 20, 00, 00, 00, time.UTC),
				},
			},
			wantProducts: []Product{
				{
					Id:        20,
					Brand:     "K",
					Category:  "K",
					Quantity:  20,
					Price:     200,
					CreatedAt: time.Date(2023, 04, 28, 12, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 13, 00, 00, 00, time.UTC),
				},
				{
					Id:        10,
					Brand:     "J",
					Category:  "J",
					Quantity:  10,
					Price:     100,
					CreatedAt: time.Date(2023, 04, 28, 10, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 20, 00, 00, 00, time.UTC),
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, "existingData.yaml")
			repo := NewPostgresRepo(db)

			err := repo.Update(tt.args.product)

			assert.ErrorIs(t, err, tt.wantErr, "error while updating product should match the expected error")

			var gotProducts []Product
			gotErr := db.NewSelect().Model(&gotProducts).Scan(context.Background())
			assert.NoError(t, gotErr, "expect no error while getting products from repo")

			assert.Equal(t, tt.wantProducts, gotProducts, "expect same response product")
		})
	}
}

func TestPostgresRepo_GetById(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name        string
		args        args
		wantProduct Product
		wantErr     error
	}{
		{
			name: "product found",
			args: args{
				id: 10,
			},
			wantProduct: Product{
				Id:        10,
				Brand:     "J",
				Category:  "J",
				Quantity:  10,
				Price:     100,
				CreatedAt: time.Date(2023, 04, 28, 10, 00, 00, 00, time.UTC),
				UpdatedAt: time.Date(2023, 04, 28, 11, 00, 00, 00, time.UTC),
			},
			wantErr: nil,
		},
		{
			name: "product not found",
			args: args{
				id: 99,
			},
			wantProduct: Product{},
			wantErr:     errProductNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, "existingData.yaml")
			repo := NewPostgresRepo(db)

			product, err := repo.GetById(tt.args.id)

			assert.ErrorIs(t, err, tt.wantErr, "error while getting product should match the expected error")
			assert.Equal(t, tt.wantProduct, product, "expect same product ")
		})
	}
}
func TestPostgresRepo_GetAll(t *testing.T) {
	tests := []struct {
		name            string
		fixtureFileName string
		wantProducts    []Product
		wantErr         error
	}{
		{
			name:            "got products successfully",
			fixtureFileName: "existingData.yaml",
			wantProducts: []Product{
				{
					Id:        10,
					Brand:     "J",
					Category:  "J",
					Quantity:  10,
					Price:     100,
					CreatedAt: time.Date(2023, 04, 28, 10, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 11, 00, 00, 00, time.UTC),
				},
				{
					Id:        20,
					Brand:     "K",
					Category:  "K",
					Quantity:  20,
					Price:     200,
					CreatedAt: time.Date(2023, 04, 28, 12, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 13, 00, 00, 00, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name:            "empty products",
			fixtureFileName: "emptyData.yaml",
			wantProducts:    []Product{},
			wantErr:         nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, tt.fixtureFileName)
			repo := NewPostgresRepo(db)

			products, err := repo.GetAll()

			assert.ErrorIs(t, err, tt.wantErr, "error while geting products should match the expected error")
			assert.Equal(t, tt.wantProducts, products, "expect products in the postgres repo to be same as expected products")
		})
	}
}

func TestPostgresRepo_Delete(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name         string
		args         args
		wantProducts []Product
		wantErr      error
	}{
		{
			name: "product deleted successfully",
			args: args{
				id: 10,
			},
			wantProducts: []Product{
				{
					Id:        20,
					Brand:     "K",
					Category:  "K",
					Quantity:  20,
					Price:     200,
					CreatedAt: time.Date(2023, 04, 28, 12, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 13, 00, 00, 00, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name: "product not found",
			args: args{
				id: 11,
			},
			wantProducts: []Product{
				{
					Id:        10,
					Brand:     "J",
					Category:  "J",
					Quantity:  10,
					Price:     100,
					CreatedAt: time.Date(2023, 04, 28, 10, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 11, 00, 00, 00, time.UTC),
				},
				{
					Id:        20,
					Brand:     "K",
					Category:  "K",
					Quantity:  20,
					Price:     200,
					CreatedAt: time.Date(2023, 04, 28, 12, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 28, 13, 00, 00, 00, time.UTC),
				},
			},
			wantErr: errProductNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupPostgres(t, "existingData.yaml")
			repo := NewPostgresRepo(db)

			err := repo.Delete(tt.args.id)

			assert.ErrorIs(t, err, tt.wantErr, "error while geting products should match the expected error")

			var gotProducts []Product
			gotErr := db.NewSelect().Model(&gotProducts).Scan(context.Background())

			assert.NoError(t, gotErr, "expect no error while getting products from repo")
			assert.Equal(t, tt.wantProducts, gotProducts, "expect products in the postgres repo to be same as expected products")
		})
	}
}
