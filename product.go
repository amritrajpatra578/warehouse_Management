package main

import (
	"fmt"
	"time"
)

type Product struct {
	Id        int       `json:"id"`
	Brand     string    `json:"brand"`
	Category  string    `json:"category"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type validationError struct {
	failures []string
}

func (e *validationError) Error() string {
	return fmt.Sprintf("validation errors :%v", e.failures)
}

func validateProduct(product Product) error {
	failures := make([]string, 0)

	if product.Id < 0 {
		failures = append(failures, "Id should not be less than 0")
	}
	if product.Brand == "" {
		failures = append(failures, "Brand should not be empty")
	}
	if product.Category == "" {
		failures = append(failures, "Category should not be empty")
	}
	if product.Quantity < 0 {
		failures = append(failures, "Quantity should not be less than 0")
	}
	if product.Price < 0 {
		failures = append(failures, "Price should not be less than 0")
	}

	if len(failures) == 0 {
		return nil
	}
	return &validationError{failures: failures}
}
