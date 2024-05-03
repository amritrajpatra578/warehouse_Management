package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProduct_Marshal(t *testing.T) {
	type args struct {
		product Product
	}
	tests := []struct {
		name     string
		args     args
		wantJSON string
		wantErr  error
	}{
		{
			name: "product marshal",
			args: args{
				product: Product{
					Id:        1,
					Brand:     "A",
					Category:  "A",
					Quantity:  1,
					Price:     10,
					CreatedAt: time.Date(2023, 04, 25, 1, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 25, 1, 00, 00, 00, time.UTC),
				},
			},
			wantJSON: `
			{
				"id": 1,
				"brand": "A",
				"category": "A",
				"quantity": 1,
				"price": 10,
				"createdAt":"2023-04-25T01:00:00Z",
				"updatedAt":"2023-04-25T01:00:00Z"
			}
			`,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, err := json.Marshal(tt.args.product)
			assert.NoError(t, err)

			assert.JSONEq(t, tt.wantJSON, string(bs))
		})
	}
}
