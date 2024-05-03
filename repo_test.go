package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupInMemoryRepo(existing []Product) *InMemoryRepo {
	repo := NewInMemoryRepo()
	repo.products = append(repo.products, existing...)
	return repo
}

func TestInMemoryRepo_Create(t *testing.T) {
	existing := []Product{
		{
			Id:       1,
			Brand:    "A",
			Category: "A",
			Quantity: 1,
			Price:    10,
		},
		{
			Id:       2,
			Brand:    "B",
			Category: "B",
			Quantity: 2,
			Price:    20,
		},
	}

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
			name: "product with same id",
			args: args{
				product: Product{
					Id:       2,
					Brand:    "B",
					Category: "B",
					Quantity: 2,
					Price:    20,
				},
			},
			wantProducts: []Product{
				{
					Id:       1,
					Brand:    "A",
					Category: "A",
					Quantity: 1,
					Price:    10,
				},
				{
					Id:       2,
					Brand:    "B",
					Category: "B",
					Quantity: 2,
					Price:    20,
				},
			},
			wantErr: errDuplicateId,
		},
		{
			name: "product created",
			args: args{
				product: Product{
					Id:       3,
					Brand:    "A",
					Category: "A",
					Quantity: 1,
					Price:    10,
				},
			},
			wantProducts: []Product{
				{
					Id:       1,
					Brand:    "A",
					Category: "A",
					Quantity: 1,
					Price:    10,
				}, {
					Id:       2,
					Brand:    "B",
					Category: "B",
					Quantity: 2,
					Price:    20,
				}, {
					Id:       3,
					Brand:    "A",
					Category: "A",
					Quantity: 1,
					Price:    10,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupInMemoryRepo(existing)
			err := repo.Create(tt.args.product)

			assert.ErrorIs(t, err, tt.wantErr, "error should match")
			assert.Equal(t, tt.wantProducts, repo.products, "expect same products in slice")

		})
	}
}

func TestInMemoryRepo_Update(t *testing.T) {
	existing := []Product{
		{
			Id:        1,
			Brand:     "A",
			Category:  "B",
			Quantity:  1,
			Price:     10,
			CreatedAt: time.Date(2023, 04, 26, 15, 00, 00, 00, time.UTC),
			UpdatedAt: time.Date(2023, 04, 26, 16, 00, 00, 00, time.UTC),
		},
		{
			Id:        2,
			Brand:     "A",
			Category:  "B",
			Quantity:  1,
			Price:     10,
			CreatedAt: time.Date(2023, 04, 26, 17, 00, 00, 00, time.UTC),
			UpdatedAt: time.Date(2023, 04, 26, 18, 00, 00, 00, time.UTC),
		},
	}
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
			name: "product updated",
			args: args{
				product: Product{
					Id:        2,
					Brand:     "B",
					Category:  "B",
					Quantity:  2,
					Price:     20,
					CreatedAt: time.Date(2023, 04, 26, 17, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 19, 00, 00, 00, time.UTC),
				},
			},
			wantProducts: []Product{
				{
					Id:        1,
					Brand:     "A",
					Category:  "B",
					Quantity:  1,
					Price:     10,
					CreatedAt: time.Date(2023, 04, 26, 15, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 16, 00, 00, 00, time.UTC),
				},
				{
					Id:        2,
					Brand:     "B",
					Category:  "B",
					Quantity:  2,
					Price:     20,
					CreatedAt: time.Date(2023, 04, 26, 17, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 19, 00, 00, 00, time.UTC),
				},
			},
			wantErr: nil,
		},
		{
			name: "id not found",
			args: args{
				product: Product{
					Id:       4,
					Brand:    "C",
					Category: "C",
					Quantity: 3,
					Price:    30,
				},
			},
			wantProducts: []Product{
				{
					Id:        1,
					Brand:     "A",
					Category:  "B",
					Quantity:  1,
					Price:     10,
					CreatedAt: time.Date(2023, 04, 26, 15, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 16, 00, 00, 00, time.UTC),
				},
				{
					Id:        2,
					Brand:     "A",
					Category:  "B",
					Quantity:  1,
					Price:     10,
					CreatedAt: time.Date(2023, 04, 26, 17, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 18, 00, 00, 00, time.UTC),
				},
			},
			wantErr: errProductNotFound,
		},
		{
			name: "created at should not change on update",
			args: args{
				product: Product{
					Id:        1,
					Brand:     "A",
					Category:  "B",
					Quantity:  1,
					Price:     10,
					CreatedAt: time.Date(2023, 04, 26, 10, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 16, 00, 00, 00, time.UTC),
				},
			},
			wantProducts: []Product{
				{
					Id:        1,
					Brand:     "A",
					Category:  "B",
					Quantity:  1,
					Price:     10,
					CreatedAt: time.Date(2023, 04, 26, 15, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 16, 00, 00, 00, time.UTC),
				},
				{
					Id:        2,
					Brand:     "A",
					Category:  "B",
					Quantity:  1,
					Price:     10,
					CreatedAt: time.Date(2023, 04, 26, 17, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 18, 00, 00, 00, time.UTC),
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupInMemoryRepo(existing)

			err := repo.Update(tt.args.product)

			assert.ErrorIs(t, err, tt.wantErr, "error should match")
			assert.Equal(t, tt.wantProducts, repo.products, "expect same products in slice")
		})
	}
}

func TestInMemoryRepo_Delete(t *testing.T) {
	existing := []Product{
		{
			Id:       1,
			Brand:    "A",
			Category: "A",
			Quantity: 1,
			Price:    10,
		},
		{
			Id:       2,
			Brand:    "B",
			Category: "B",
			Quantity: 2,
			Price:    20,
		},
		{
			Id:       3,
			Brand:    "C",
			Category: "C",
			Quantity: 3,
			Price:    30,
		},
	}

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
			name: "product deleted",
			args: args{
				id: 2,
			},
			wantProducts: []Product{
				{
					Id:       1,
					Brand:    "A",
					Category: "A",
					Quantity: 1,
					Price:    10,
				},
				{
					Id:       3,
					Brand:    "C",
					Category: "C",
					Quantity: 3,
					Price:    30,
				},
			},
			wantErr: nil,
		},
		{
			name: "product not found",
			args: args{
				id: 213,
			},
			wantProducts: []Product{
				{
					Id:       1,
					Brand:    "A",
					Category: "A",
					Quantity: 1,
					Price:    10,
				},
				{
					Id:       2,
					Brand:    "B",
					Category: "B",
					Quantity: 2,
					Price:    20,
				},
				{
					Id:       3,
					Brand:    "C",
					Category: "C",
					Quantity: 3,
					Price:    30,
				},
			},
			wantErr: errProductNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupInMemoryRepo(existing)

			err := repo.Delete(tt.args.id)

			assert.ErrorIs(t, err, tt.wantErr, "error should match")
			assert.Equal(t, tt.wantProducts, repo.products, "products should be equal")

		})
	}

}

func TestInMemoryRepo_GetById(t *testing.T) {
	existing := []Product{
		{
			Id:       1,
			Brand:    "A",
			Category: "A",
			Quantity: 1,
			Price:    10,
		},
		{
			Id:       2,
			Brand:    "B",
			Category: "B",
			Quantity: 2,
			Price:    20,
		},
	}
	type args struct {
		Id int
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
				Id: 2,
			},
			wantProduct: Product{
				Id:       2,
				Brand:    "B",
				Category: "B",
				Quantity: 2,
				Price:    20,
			},
			wantErr: nil,
		},
		{
			name: "product not found",
			args: args{
				Id: 3,
			},
			wantProduct: Product{},
			wantErr:     errProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			repo := setupInMemoryRepo(existing)
			gotProduct, err := repo.GetById(tt.args.Id)

			assert.ErrorIs(t, err, tt.wantErr, "error should match")
			assert.Equal(t, tt.wantProduct, gotProduct, "expect same product in slice")

		})
	}

}
func TestInMemoryRepo_GetAll(t *testing.T) {
	tests := []struct {
		name         string
		existing     []Product
		wantProducts []Product
		wantError    error
	}{
		{
			name: "all products found",
			existing: []Product{

				{
					Id:       1,
					Brand:    "A",
					Category: "A",
					Quantity: 1,
					Price:    10,
				},
				{
					Id:       2,
					Brand:    "B",
					Category: "B",
					Quantity: 2,
					Price:    20,
				},
			},
			wantProducts: []Product{
				{
					Id:       1,
					Brand:    "A",
					Category: "A",
					Quantity: 1,
					Price:    10,
				},
				{
					Id:       2,
					Brand:    "B",
					Category: "B",
					Quantity: 2,
					Price:    20,
				},
			},

			wantError: nil,
		},
		{
			name:         "no products",
			existing:     []Product{},
			wantProducts: []Product{},
			wantError:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupInMemoryRepo(tt.existing)
			gotProducts, err := repo.GetAll()

			assert.ErrorIs(t, err, tt.wantError, "error should match")
			assert.Equal(t, tt.wantProducts, gotProducts, "same products should be present")
		})
	}

}
