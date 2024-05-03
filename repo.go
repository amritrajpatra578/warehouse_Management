package main

import (
	"errors"
)

var (
	errProductNotFound = errors.New("product not found")
	errDuplicateId     = errors.New("found duplicate id")
	errEmptyId         = errors.New("id should not be empty")
)

type Repo interface {
	Create(Product) error
	Update(Product) error
	GetById(id int) (Product, error)
	GetAll() ([]Product, error)
	Delete(id int) error
}

type InMemoryRepo struct {
	products []Product
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		products: make([]Product, 0),
	}
}

func (r *InMemoryRepo) Create(product Product) error {
	for _, currentProduct := range r.products {
		if currentProduct.Id == product.Id {
			return errDuplicateId
		}
	}
	r.products = append(r.products, product)
	return nil
}

func (r *InMemoryRepo) Update(product Product) error {
	for idx, currentProduct := range r.products {
		if currentProduct.Id == product.Id {
			product.CreatedAt = currentProduct.CreatedAt
			r.products[idx] = product
			return nil
		}
	}
	return errProductNotFound
}

func (r *InMemoryRepo) GetById(id int) (Product, error) {
	for _, currentProduct := range r.products {
		if currentProduct.Id == id {
			return currentProduct, nil
		}
	}
	return Product{}, errProductNotFound
}

func (r *InMemoryRepo) GetAll() ([]Product, error) {
	products := make([]Product, len(r.products))
	copy(products, r.products)
	return products, nil
}

func (r *InMemoryRepo) Delete(id int) error {
	for idx, currentProduct := range r.products {
		if currentProduct.Id == id {
			r.products = append(r.products[:idx], r.products[idx+1:]...)
			return nil
		}
	}
	return errProductNotFound
}
