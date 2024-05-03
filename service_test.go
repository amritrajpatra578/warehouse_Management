package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testSubscriber struct {
	id       string
	products []Product
	count    int
}

func (t *testSubscriber) Update(products []Product) {
	t.count = t.count + 1
	t.products = products
}

func (t *testSubscriber) Id() string {
	return t.id
}

func assertProductsEqual(t *testing.T, expectedProducts, actualProducts []Product, id int, start, end time.Time) {
	for index, expected := range expectedProducts {
		actual := actualProducts[index]
		assert.Equal(t, expected.Id, actual.Id, "expect same product ID")
		assert.Equal(t, expected.Brand, actual.Brand, "expect same product Brand")
		assert.Equal(t, expected.Category, actual.Category, "expect same product Category")
		assert.Equal(t, expected.Quantity, actual.Quantity, "expect same product Quantity")
		assert.Equal(t, expected.Price, actual.Price, "expect same product Price")

		if expected.Id == id {
			assertTimestampBetween(t, start, end, actual.CreatedAt)
			assertTimestampBetween(t, start, end, actual.UpdatedAt)

		} else {
			assert.Equal(t, expected.CreatedAt, actual.CreatedAt, "expect same created at for existing products")
			assert.Equal(t, expected.UpdatedAt, actual.UpdatedAt, "expect same updated at for existing products")
		}
	}
}

func assertProductsEqualUpdate(t *testing.T, expectedProducts, actualProducts []Product, id int, start, end time.Time) {
	for index, expected := range expectedProducts {
		actual := actualProducts[index]
		assert.Equal(t, expected.Id, actual.Id, "expect same product ID")
		assert.Equal(t, expected.Brand, actual.Brand, "expect same product Brand")
		assert.Equal(t, expected.Category, actual.Category, "expect same product Category")
		assert.Equal(t, expected.Quantity, actual.Quantity, "expect same product Quantity")
		assert.Equal(t, expected.Price, actual.Price, "expect same product Price")
		assert.Equal(t, expected.CreatedAt, actual.CreatedAt, "expect same created at for existing products")

		if expected.Id == id {
			assertTimestampBetween(t, start, end, actual.UpdatedAt)
		} else {
			assert.Equal(t, expected.UpdatedAt, actual.UpdatedAt, "expect same updated at for existing products")
		}
	}
}

func assertTimestampBetween(t *testing.T, before, after, ts time.Time) {
	assert.True(t, before.Before(ts), "before should be before ts")
	assert.True(t, after.After(ts), "after should be after ts")
}

func TestProductServiceImpl_Create(t *testing.T) {
	existing := []Product{
		{
			Id:        1,
			Brand:     "A",
			Category:  "A",
			Quantity:  1,
			Price:     10,
			CreatedAt: time.Date(2023, 04, 26, 15, 00, 00, 00, time.UTC),
			UpdatedAt: time.Date(2023, 04, 26, 15, 00, 00, 00, time.UTC),
		},
		{
			Id:        2,
			Brand:     "B",
			Category:  "B",
			Quantity:  2,
			Price:     20,
			CreatedAt: time.Date(2023, 04, 26, 16, 00, 00, 00, time.UTC),
			UpdatedAt: time.Date(2023, 04, 26, 16, 00, 00, 00, time.UTC),
		},
	}
	type args struct {
		product Product
	}
	tests := []struct {
		name             string
		args             args
		wantProducts     []Product
		wantFailures     []string
		wantNotification bool
		wantError        error
	}{
		{
			name: "success",
			args: args{
				product: Product{
					Id:       3,
					Brand:    "C",
					Category: "C",
					Quantity: 2,
					Price:    20,
				},
			},
			wantProducts: []Product{
				{
					Id:        1,
					Brand:     "A",
					Category:  "A",
					Quantity:  1,
					Price:     10,
					CreatedAt: time.Date(2023, 04, 26, 15, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 15, 00, 00, 00, time.UTC),
				},
				{
					Id:        2,
					Brand:     "B",
					Category:  "B",
					Quantity:  2,
					Price:     20,
					CreatedAt: time.Date(2023, 04, 26, 16, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 16, 00, 00, 00, time.UTC),
				},
				{
					Id:       3,
					Brand:    "C",
					Category: "C",
					Quantity: 2,
					Price:    20,
				},
			},
			wantNotification: true,
			wantError:        nil,
		},
		{
			name: "validation errors",
			args: args{
				product: Product{
					Id:       -3,
					Brand:    "C",
					Category: "",
					Quantity: -1,
					Price:    20,
				},
			},
			wantProducts: []Product{
				{
					Id:        1,
					Brand:     "A",
					Category:  "A",
					Quantity:  1,
					Price:     10,
					CreatedAt: time.Date(2023, 04, 26, 15, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 15, 00, 00, 00, time.UTC),
				},
				{
					Id:        2,
					Brand:     "B",
					Category:  "B",
					Quantity:  2,
					Price:     20,
					CreatedAt: time.Date(2023, 04, 26, 16, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 16, 00, 00, 00, time.UTC),
				},
			},
			wantFailures: []string{
				"Id should not be less than 0",
				"Category should not be empty",
				"Quantity should not be less than 0",
			},
			wantNotification: false,
			wantError:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupInMemoryRepo(existing)
			svc := NewProductServiceImpl(repo)

			start := time.Now()

			subscriber := &testSubscriber{id: "B"}
			subscriberErr := svc.subscribe(subscriber)
			assert.NoError(t, subscriberErr, "subscribe should succeed")

			err := svc.Create(tt.args.product)

			end := time.Now()

			if len(tt.wantFailures) > 0 {
				var ve *validationError
				assert.ErrorAs(t, err, &ve, "error should be of ValidationError type")
				assert.Equal(t, tt.wantFailures, ve.failures, "expect failures to be same")
				assert.Equal(t, tt.wantProducts, repo.products, "expect same products")
			} else {
				assert.ErrorIs(t, err, tt.wantError, "error should match")
			}

			assertProductsEqual(t, tt.wantProducts, repo.products, tt.args.product.Id, start, end)

			if tt.wantNotification {
				assertProductsEqual(t, tt.wantProducts, subscriber.products, tt.args.product.Id, start, end)
				assert.Equal(t, 1, subscriber.count, "expect 1 notification")
			} else {
				assert.Nil(t, subscriber.products, "expect no products for subscriber")
				assert.Equal(t, 0, subscriber.count, "expect no notification")
			}
		})
	}
}
func TestProductServiceImpl_Update(t *testing.T) {
	existing := []Product{
		{
			Id:        1,
			Brand:     "Aa",
			Category:  "A",
			Quantity:  1,
			Price:     10,
			CreatedAt: time.Date(2023, 04, 26, 15, 00, 00, 00, time.UTC),
			UpdatedAt: time.Date(2023, 04, 26, 16, 00, 00, 00, time.UTC),
		},
		{
			Id:        2,
			Brand:     "Bb",
			Category:  "B",
			Quantity:  2,
			Price:     20,
			CreatedAt: time.Date(2023, 04, 26, 17, 00, 00, 00, time.UTC),
			UpdatedAt: time.Date(2023, 04, 26, 18, 00, 00, 00, time.UTC),
		},
	}

	type args struct {
		product Product
	}

	tests := []struct {
		name             string
		args             args
		wantProducts     []Product
		wantFailures     []string
		wantNotification bool
		wantError        error
	}{
		{
			name: "product updated",
			args: args{
				product: Product{
					Id:       2,
					Brand:    "B",
					Category: "B",
					Quantity: 2,
					Price:    22,
				},
			},
			wantProducts: []Product{
				{
					Id:        1,
					Brand:     "Aa",
					Category:  "A",
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
					Price:     22,
					CreatedAt: time.Date(2023, 04, 26, 17, 00, 00, 00, time.UTC),
				},
			},
			wantFailures:     []string{},
			wantNotification: true,
			wantError:        nil,
		},
		{
			name: "invalid id, price, brand",
			args: args{
				product: Product{
					Id:       -2,
					Brand:    "",
					Category: "B",
					Quantity: 2,
					Price:    -13,
				},
			},
			wantProducts: []Product{
				{
					Id:        1,
					Brand:     "Aa",
					Category:  "A",
					Quantity:  1,
					Price:     10,
					CreatedAt: time.Date(2023, 04, 26, 15, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 16, 00, 00, 00, time.UTC),
				},
				{
					Id:        2,
					Brand:     "Bb",
					Category:  "B",
					Quantity:  2,
					Price:     20,
					CreatedAt: time.Date(2023, 04, 26, 17, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 18, 00, 00, 00, time.UTC),
				},
			},
			wantFailures: []string{
				"Id should not be less than 0",
				"Brand should not be empty",
				"Price should not be less than 0",
			},
			wantNotification: false,
			wantError:        nil,
		},
		{
			name: "product not found",
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
					Brand:     "Aa",
					Category:  "A",
					Quantity:  1,
					Price:     10,
					CreatedAt: time.Date(2023, 04, 26, 15, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 16, 00, 00, 00, time.UTC),
				},
				{
					Id:        2,
					Brand:     "Bb",
					Category:  "B",
					Quantity:  2,
					Price:     20,
					CreatedAt: time.Date(2023, 04, 26, 17, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 18, 00, 00, 00, time.UTC),
				},
			},
			wantFailures:     []string{},
			wantNotification: false,
			wantError:        errProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupInMemoryRepo(existing)
			svc := NewProductServiceImpl(repo)

			start := time.Now()
			subscriber := &testSubscriber{id: "C"}
			subscriberErr := svc.subscribe(subscriber)
			assert.NoError(t, subscriberErr, "expect no error while subscribing")

			err := svc.Update(tt.args.product)

			end := time.Now()

			if len(tt.wantFailures) > 0 {
				var ve *validationError
				assert.ErrorAs(t, err, &ve, "error should be of ValidationError type")
				assert.Equal(t, tt.wantFailures, ve.failures, "expect failures to be same")
				assert.Equal(t, tt.wantProducts, repo.products, "expect products should be equal")
			} else {
				assert.ErrorIs(t, err, tt.wantError, "error should match")
			}

			assertProductsEqualUpdate(t, tt.wantProducts, repo.products, tt.args.product.Id, start, end)

			if tt.wantNotification {
				assertProductsEqualUpdate(t, tt.wantProducts, subscriber.products, tt.args.product.Id, start, end)
				assert.Equal(t, 1, subscriber.count, "expect 1 notification")
			} else {
				assert.Nil(t, subscriber.products)
				assert.Equal(t, 0, subscriber.count, "expect no notification")
			}
		})
	}
}
func TestProductServiceImpl_GetById(t *testing.T) {
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
				id: 2,
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
				id: 12,
			},
			wantProduct: Product{},
			wantErr:     errProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupInMemoryRepo(existing)
			svc := NewProductServiceImpl(repo)

			gotProduct, err := svc.GetById(tt.args.id)

			assert.Equal(t, tt.wantProduct, gotProduct, "expect same product")
			assert.ErrorIs(t, err, tt.wantErr, "expect same error")
		})
	}
}

func TestProductServiceImpl_GetAll(t *testing.T) {

	tests := []struct {
		name         string
		existing     []Product
		wantProducts []Product
		wantError    error
	}{
		{
			name: "products found",
			existing: []Product{
				{
					Id:       1,
					Brand:    "a",
					Category: "a",
					Quantity: 1,
					Price:    10,
				},
				{
					Id:       2,
					Brand:    "b",
					Category: "b",
					Quantity: 2,
					Price:    20,
				},
			},
			wantProducts: []Product{
				{
					Id:       1,
					Brand:    "a",
					Category: "a",
					Quantity: 1,
					Price:    10,
				},
				{
					Id:       2,
					Brand:    "b",
					Category: "b",
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
			svc := NewProductServiceImpl(repo)

			gotProducts, err := svc.GetAll()

			assert.Equal(t, tt.wantProducts, gotProducts, "expect same products")
			assert.ErrorIs(t, err, tt.wantError, "expect same error")
		})
	}
}

func TestProductServiceImpl_Delete(t *testing.T) {
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
			Quantity: 2,
			Price:    20,
		},
	}

	type args struct {
		id int
	}

	tests := []struct {
		name             string
		args             args
		wantProducts     []Product
		wantNotification bool
		wantError        error
	}{
		{
			name: "product deleted",
			args: args{
				id: 3,
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
			wantNotification: true,
			wantError:        nil,
		},
		{
			name: "product not found",
			args: args{
				id: 23434,
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
					Quantity: 2,
					Price:    20,
				},
			},
			wantNotification: false,
			wantError:        errProductNotFound,
		},
		{
			name: "notify subscribers product deleted",
			args: args{
				id: 1,
			},
			wantProducts: []Product{
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
					Quantity: 2,
					Price:    20,
				},
			},
			wantNotification: true,
			wantError:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupInMemoryRepo(existing)
			svc := NewProductServiceImpl(repo)

			subscriber := &testSubscriber{id: "A"}
			subscriberErr := svc.subscribe(subscriber)
			assert.NoError(t, subscriberErr, "subscribe should succeed")

			err := svc.Delete(tt.args.id)

			assert.ErrorIs(t, err, tt.wantError, "expect same error")
			assert.Equal(t, tt.wantProducts, repo.products, "expect products after delete")

			if tt.wantNotification {
				assert.Equal(t, tt.wantProducts, subscriber.products, "expect same notification")
				assert.Equal(t, 1, subscriber.count, "expect notification")
			} else {
				assert.Nil(t, subscriber.products, "expect no products for subscriber")
				assert.Equal(t, 0, subscriber.count, "expect no notification")
			}
		})
	}
}

func TestProductServiceImpl_subscribe(t *testing.T) {
	type args struct {
		subscriber Subscriber
	}

	tests := []struct {
		name           string
		args           args
		wantSubsribers []Subscriber
		wantError      bool
	}{
		{
			name: "subscriber added",
			args: args{
				subscriber: &testSubscriber{id: "A"},
			},
			wantSubsribers: []Subscriber{
				&testSubscriber{id: "A"},
			},
			wantError: false,
		},
		{
			name: "subscriber with empty id",
			args: args{
				subscriber: &testSubscriber{id: ""},
			},
			wantSubsribers: nil,
			wantError:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProductServiceImpl(nil)

			err := svc.subscribe(tt.args.subscriber)
			if !tt.wantError {
				assert.NoError(t, err, "expect no error")
			} else {
				assert.Error(t, err, "expect error")
			}
			assert.Equal(t, tt.wantSubsribers, svc.subscribers, "expect same subscriber")
		})
	}
}

func TestProductImpl_unsubscribe(t *testing.T) {
	existingSubscribers := []Subscriber{
		&testSubscriber{
			id: "A",
		},
		&testSubscriber{
			id: "B",
		},
	}
	type args struct {
		subscriber Subscriber
	}
	tests := []struct {
		name            string
		args            args
		wantSubscribers []Subscriber
		wantError       bool
	}{
		{
			name: "subscriber removed",
			args: args{
				subscriber: &testSubscriber{id: "A"},
			},
			wantSubscribers: []Subscriber{
				&testSubscriber{id: "B"},
			},
			wantError: false,
		},
		{
			name: "subscriber with empty id",
			args: args{
				subscriber: &testSubscriber{id: ""},
			},
			wantSubscribers: []Subscriber{
				&testSubscriber{
					id: "A",
				},
				&testSubscriber{
					id: "B",
				},
			},
			wantError: true,
		},
		{
			name: "subscriber not found",
			args: args{
				subscriber: &testSubscriber{
					id: "D",
				},
			},
			wantSubscribers: []Subscriber{
				&testSubscriber{
					id: "A",
				},
				&testSubscriber{
					id: "B",
				},
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProductServiceImpl(nil)
			svc.subscribers = append(svc.subscribers, existingSubscribers...)

			err := svc.unsubscribe(tt.args.subscriber)

			if !tt.wantError {
				assert.NoError(t, err, "expect no error")
			} else {
				assert.Error(t, err, "expect error")
			}
			assert.Equal(t, tt.wantSubscribers, svc.subscribers, "expect no subscriber")
		})
	}
}

func TestProductServiceImpl_notify(t *testing.T) {
	existing := []Product{
		{
			Id:       1,
			Brand:    "A",
			Category: "A",
			Quantity: 10,
			Price:    10,
		},
		{
			Id:       2,
			Brand:    "B",
			Category: "B",
			Quantity: 20,
			Price:    20,
		},
		{
			Id:       3,
			Brand:    "C",
			Category: "C",
			Quantity: 30,
			Price:    30,
		},
	}

	tests := []struct {
		name        string
		subscribers []Subscriber
		wantCount   int
	}{
		{
			name: "product removed",
			subscribers: []Subscriber{
				&testSubscriber{id: "A"},
				&testSubscriber{id: "B"},
			},
			wantCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupInMemoryRepo(existing)
			svc := NewProductServiceImpl(repo)
			svc.subscribers = append(svc.subscribers, tt.subscribers...)

			svc.notify()

			for _, subscriber := range svc.subscribers {
				s := subscriber.(*testSubscriber)
				assert.Equal(t, repo.products, s.products, "expect same products")
				assert.Equal(t, tt.wantCount, s.count, "expect same no. of count")
			}
		})
	}
}
