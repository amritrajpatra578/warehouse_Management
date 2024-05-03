package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHttpTransport_create(t *testing.T) {
	existing := []Product{
		{
			Id:        1,
			Brand:     "A",
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
			Price:     20,
			CreatedAt: time.Date(2023, 04, 26, 17, 00, 00, 00, time.UTC),
			UpdatedAt: time.Date(2023, 04, 26, 18, 00, 00, 00, time.UTC),
		},
	}

	tests := []struct {
		name           string
		productJSON    string
		wantStatusCode int
		wantResponse   string
	}{
		{
			name: "sucess",
			productJSON: `{
				"id": 3,
				"brand": "C",
				"category": "C",
				"quantity": 3,
				"price": 30
			}`,
			wantStatusCode: http.StatusCreated,
			wantResponse: `{
				"id": 3,
				"brand": "C",
				"category": "C",
				"quantity": 3,
				"price": 30 
			}`,
		},
		{
			name: "failed to decode",
			productJSON: `{
				"id": ada,
				"brand": "C",
				"category": "C",
				"quantity": 3,
				"price": 30 
			}`,
			wantStatusCode: http.StatusBadRequest,
			wantResponse: `
			{
				"errors": [
					"invalid json"
				]
			}`,
		},
		{
			name: "duplicate ID",
			productJSON: `{
				"id": 1,
				"brand": "A",
				"category": "A",
				"quantity": 1,
				"price": 10 
			}`,
			wantStatusCode: http.StatusConflict,
			wantResponse: `{
				"errors":[ 
					"product exists"
				]
			}`,
		},
		{
			name: "validation error",
			productJSON: `{
				"id": -3,
				"brand": "C",
				"category": "C",
				"quantity": -3,
				"price": 30 
			}`,
			wantStatusCode: http.StatusBadRequest,
			wantResponse: `
			{
				"errors": [
					"Id should not be less than 0",
					"Quantity should not be less than 0"
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupInMemoryRepo(existing)
			svc := NewProductServiceImpl(repo)
			httpTransport := NewhttpTransport(svc)

			handler := buildHttpHandler(httpTransport)
			body := strings.NewReader(tt.productJSON)
			r := httptest.NewRequest("POST", "/products", body)
			w := httptest.NewRecorder()

			start := time.Now()

			handler.ServeHTTP(w, r)

			end := time.Now()

			response := w.Result()
			assert.Equal(t, tt.wantStatusCode, response.StatusCode, "expect same status code")

			responseBytes, err := io.ReadAll(response.Body)
			assert.NoError(t, err, "read response body should succeed")

			var product Product
			err = json.Unmarshal(responseBytes, &product)
			assert.NoError(t, err, "expect no err while conv json to obj type")

			if tt.wantStatusCode == http.StatusCreated {
				var wantProduct Product
				err := json.Unmarshal([]byte(tt.wantResponse), &wantProduct)
				assert.NoError(t, err, "expect no err while conv json to obj type")

				assert.Equal(t, wantProduct.Id, product.Id, "expect same product id")
				assert.Equal(t, wantProduct.Brand, product.Brand, "expect same product brand")
				assert.Equal(t, wantProduct.Category, product.Category, "expect same product category")
				assert.Equal(t, wantProduct.Quantity, product.Quantity, "expect same product quantity")
				assert.Equal(t, wantProduct.Price, product.Price, "expect same product price")

				assertTimestampBetween(t, start, end, product.UpdatedAt)
				assertTimestampBetween(t, start, end, product.CreatedAt)
			} else {
				assert.JSONEq(t, tt.wantResponse, string(responseBytes), "expect same product response")
			}
		})
	}
}

func TestHttpTransport_Update(t *testing.T) {
	existing := []Product{
		{
			Id:        1,
			Brand:     "A",
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
			Price:     20,
			CreatedAt: time.Date(2023, 04, 26, 17, 00, 00, 00, time.UTC),
			UpdatedAt: time.Date(2023, 04, 26, 18, 00, 00, 00, time.UTC),
		},
	}
	tests := []struct {
		name           string
		productID      int
		productJSON    string
		wantStatusCode int
		wantResponse   string
	}{
		{
			name:      "success",
			productID: 2,
			productJSON: `{
				"id": 2,
				"brand": "C",
				"category": "C",
				"quantity": 3,
				"price": 30
			}`,
			wantStatusCode: http.StatusOK,
			wantResponse: `
			{
				"id": 2,
				"brand": "C",
				"category": "C",
				"quantity": 3,
				"price": 30,
				"createdAt": "2023-04-26T17:00:00Z"
			}`,
		},
		{
			name:      "invalid request json",
			productID: 2,
			productJSON: `{
				"id": adad,
				"brand": "C",
				"category": "C",
				"quantity": 3,
				"price": 30 
			}`,
			wantStatusCode: http.StatusBadRequest,
			wantResponse: `
			{
				"errors": [
					"invalid json"
				]
			}`,
		},
		{
			name:      "validation error",
			productID: 1,
			productJSON: `
			{
				"id": -1,
				"brand": "",
				"category": "C",
				"quantity": 3,
				"price": -30 
			}`,
			wantStatusCode: http.StatusBadRequest,
			wantResponse: `
			{
				"errors": [
					"Brand should not be empty",
					"Price should not be less than 0"
				]
			}`,
		},
		{
			name:      "product not found",
			productID: 12,
			productJSON: `
			{
				"id": 12,
				"brand": "C",
				"category": "C",
				"quantity": 3,
				"price": 30
			}
			`,
			wantStatusCode: http.StatusNotFound,
			wantResponse: `{
				"errors": [
					"product not found"
					]
				}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupInMemoryRepo(existing)
			svc := NewProductServiceImpl(repo)
			httpTransport := NewhttpTransport(svc)

			handler := buildHttpHandler(httpTransport) //setsup route for http server
			body := strings.NewReader(tt.productJSON)  //body contains the product in json format which will be used as http request body

			url := fmt.Sprintf("/products/%d", tt.productID) //this url is used to hit/call update api and by the product id given, it gives the refernce on what product to update

			r := httptest.NewRequest("PUT", url, body) //this creates an http req with url of update api ,its http mthod(PUT),and json of the product to be updated as its req body
			w := httptest.NewRecorder()                // initailzie a new recorder to record the res of http server

			start := time.Now()

			handler.ServeHTTP(w, r) //it will call http request and record the response in the w

			end := time.Now()

			response := w.Result()
			assert.Equal(t, tt.wantStatusCode, response.StatusCode, "expect same status code")

			responseBytes, err := io.ReadAll(response.Body) //converts readcloser type to []byte type
			assert.NoError(t, err, "read response body should succeed")

			var product Product
			err = json.Unmarshal(responseBytes, &product)
			assert.NoError(t, err, "expect no err while conv json to obj type")

			if tt.wantStatusCode == http.StatusOK {
				var wantProduct Product
				err := json.Unmarshal([]byte(tt.wantResponse), &wantProduct)
				assert.NoError(t, err, "expect no err while conv json to obj type")

				assert.Equal(t, wantProduct.Id, product.Id, "expect same product id")
				assert.Equal(t, wantProduct.Brand, product.Brand, "expect same product brand")
				assert.Equal(t, wantProduct.Category, product.Category, "expect same product category")
				assert.Equal(t, wantProduct.Quantity, product.Quantity, "expect same product quantity")
				assert.Equal(t, wantProduct.Price, product.Price, "expect same product price")
				assert.Equal(t, wantProduct.CreatedAt, product.CreatedAt, "expect same product createdAt")

				assertTimestampBetween(t, start, end, product.UpdatedAt)
			} else {
				assert.JSONEq(t, tt.wantResponse, string(responseBytes), "expect same product response")
			}
		})
	}
}

func TestHttpTransport_GetById(t *testing.T) {
	existing := []Product{
		{
			Id:        1,
			Brand:     "A",
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
			Price:     20,
			CreatedAt: time.Date(2023, 04, 26, 17, 00, 00, 00, time.UTC),
			UpdatedAt: time.Date(2023, 04, 26, 18, 00, 00, 00, time.UTC),
		},
	}
	tests := []struct {
		name           string
		productId      int
		wantStatusCode int
		wantResponse   string
	}{
		{
			name:           "product found",
			productId:      2,
			wantStatusCode: http.StatusOK,
			wantResponse: `
			{
				"id": 2,
				"brand": "B",
				"category": "B",
				"quantity":2,
				"price": 20,
				"createdAt": "2023-04-26T17:00:00Z",
				"updatedAt": "2023-04-26T18:00:00Z"
			}`,
		},
		{
			name:           "product not found",
			productId:      3,
			wantStatusCode: http.StatusNotFound,
			wantResponse: `
			{
				"errors": [
					"product not found"
					]
				}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupInMemoryRepo(existing)
			svc := NewProductServiceImpl(repo)
			httpTransport := NewhttpTransport(svc)

			handler := buildHttpHandler(httpTransport)

			url := fmt.Sprintf("/products/%d", tt.productId)
			r := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			response := w.Result()
			assert.Equal(t, tt.wantStatusCode, response.StatusCode, "expect same status code")

			responseBytes, err := io.ReadAll(response.Body)
			assert.NoError(t, err, "read response body should succeed")

			var product Product
			productErr := json.Unmarshal(responseBytes, &product)
			assert.NoError(t, productErr, "expect no while conv []byte to product obj")

			if tt.wantStatusCode == http.StatusOK {
				var wantProduct Product
				wantProductErr := json.Unmarshal([]byte(tt.wantResponse), &wantProduct)
				assert.NoError(t, wantProductErr, "expect no while conv []byte to product obj")

				assert.Equal(t, wantProduct, product, "expect products in the response to be same as expected")
			} else {
				assert.JSONEq(t, tt.wantResponse, string(responseBytes), "expect same product response")
			}
		})
	}
}

func TestHttpTransport_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		wantStatusCode int
		existing       []Product
		wantResponse   string
	}{
		{
			name:           "products found",
			wantStatusCode: http.StatusOK,
			existing: []Product{
				{
					Id:        1,
					Brand:     "A",
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
					Price:     20,
					CreatedAt: time.Date(2023, 04, 26, 17, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(2023, 04, 26, 18, 00, 00, 00, time.UTC),
				},
			},
			wantResponse: `[
			{
				"id": 1,
				"brand": "A",
				"category": "A",
				"quantity": 1,
				"price": 10,
				"createdAt": "2023-04-26T15:00:00Z",
				"updatedAt": "2023-04-26T16:00:00Z"
			},
			{
				"id": 2,
				"brand": "B",
				"category": "B",
				"quantity": 2,
				"price": 20,
				"createdAt": "2023-04-26T17:00:00Z",
				"updatedAt": "2023-04-26T18:00:00Z"
			}
		]`,
		},
		{
			name:           "no products",
			wantStatusCode: http.StatusOK,
			existing:       []Product{},
			wantResponse:   `[]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupInMemoryRepo(tt.existing)
			svc := NewProductServiceImpl(repo)
			httpTransport := NewhttpTransport(svc)
			handler := buildHttpHandler(httpTransport)

			body := strings.NewReader(tt.wantResponse)
			r := httptest.NewRequest("GET", "/products", body)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			response := w.Result()
			assert.Equal(t, tt.wantStatusCode, response.StatusCode, "expect same status code")

			responseBytes, err := io.ReadAll(response.Body)
			assert.NoError(t, err, "read response body should succeed")
			assert.JSONEq(t, tt.wantResponse, string(responseBytes), "expect same response")
		})
	}
}

func TestHttpTransport_Delete(t *testing.T) {
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
	tests := []struct {
		name           string
		productId      int
		wantStatusCode int
		wantResponse   string
	}{
		{
			name:           "success",
			productId:      2,
			wantStatusCode: http.StatusOK,
			wantResponse:   "",
		},
		{
			name:           "product not found",
			productId:      3,
			wantStatusCode: http.StatusNotFound,
			wantResponse: `
			{
				"errors": [
					"product not found"
				]
			}
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupInMemoryRepo(existing)
			svc := NewProductServiceImpl(repo)
			httpTransport := NewhttpTransport(svc)
			handler := buildHttpHandler(httpTransport)

			url := fmt.Sprintf("/products/%d", tt.productId)
			r := httptest.NewRequest("DELETE", url, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			response := w.Result()

			assert.Equal(t, tt.wantStatusCode, response.StatusCode, "expect same status code")

			responseBytes, err := io.ReadAll(response.Body)
			assert.NoError(t, err, "read response body should succeed")
			if tt.wantResponse != "" {
				assert.JSONEq(t, tt.wantResponse, string(responseBytes))
			} else {
				assert.Equal(t, []byte{}, responseBytes, "expect empty response")
			}
		})
	}
}

func startHttpServer(t *testing.T, handler http.Handler) string {
	listner, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal("failed to listen:", err)
	}

	url := listner.Addr().String()

	server := &http.Server{
		Handler: handler,
	}

	errC := make(chan error)
	go func() {
		errC <- server.Serve(listner)
	}()

	t.Cleanup(func() {
		if err := server.Close(); err != nil {
			t.Log("unable to close the server after the completion of test:", err)
		}
		t.Log("http server exited:", <-errC)
	})
	return url
}

func TestWebSocket(t *testing.T) {

	t.Run("create", func(t *testing.T) {
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

		repo := setupInMemoryRepo(existing)
		svc := NewProductServiceImpl(repo)
		httpTransport := NewhttpTransport(svc)
		handler := buildHttpHandler(httpTransport)

		url := startHttpServer(t, handler)
		t.Log("setup for http server completed")
		websocketURL := fmt.Sprintf("ws://%s/ws", url)

		conn, httpresponse, err := websocket.DefaultDialer.Dial(websocketURL, nil)
		if err != nil {
			t.Fatal("unable to establish websocket connection :", err)
		}
		assert.Equal(t, http.StatusSwitchingProtocols, httpresponse.StatusCode, "expect 101 statuscode while switchig from http to websocket")
		t.Log("websocket connection established")

		t.Cleanup(func() {
			if err := conn.Close(); err != nil {
				t.Log("err while closing websocket conn:", err)
			}
		})

		productJSON := `
		{
		"id": 3,
		"brand": "C",
		"category": "C",
		"quantity": 3,
		"price": 30
		}
		`
		body := strings.NewReader(productJSON)
		createURL := fmt.Sprintf("http://%s/products", url)
		req, reqErr := http.NewRequest("POST", createURL, body)
		if reqErr != nil {
			t.Fatal("failed to build request:", err)
		}

		res, resErr := http.DefaultClient.Do(req)
		if resErr != nil {
			t.Fatal("err in calling create product api:", resErr)
		}
		t.Log("Send request successfully")

		assert.Equal(t, http.StatusCreated, res.StatusCode, "expect same status code while calling create api")
		t.Log("got success response")

		wantProductsNotification := []Product{
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

		var gotProductsNotification []Product
		if err := conn.ReadJSON(&gotProductsNotification); err != nil {
			t.Fatal("unable to read the msg from websocket:", err)
		}
		assert.Equal(t, wantProductsNotification, gotProductsNotification, "expect products notification when create product is successful")
	})

	t.Run("update", func(t *testing.T) {
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

		repo := setupInMemoryRepo(existing)
		svc := NewProductServiceImpl(repo)
		httpTransport := NewhttpTransport(svc)
		handler := buildHttpHandler(httpTransport)
		t.Log("setup for http server completed")

		url := startHttpServer(t, handler)

		websocketURL := fmt.Sprintf("ws://%s/ws", url)

		conn, httpresponse, err := websocket.DefaultDialer.Dial(websocketURL, nil)
		if err != nil {
			t.Fatal("failed to establish websocket conn:", err)
		}

		assert.Equal(t, http.StatusSwitchingProtocols, httpresponse.StatusCode, "expect 101 status code while switching http to websocket")
		t.Log("websocket connection established")

		t.Cleanup(func() {
			if err := conn.Close(); err != nil {
				t.Log("err while closing websocket conn:", err)
			}
		})

		productJSON := `{
			"id": 1,
			"brand": "C",
			"category": "C",
			"quantity": 3,
			"price": 30
			}
			`
		body := strings.NewReader(productJSON)
		updateURL := fmt.Sprintf("http://%s/products/1", url)
		req, reqErr := http.NewRequest("PUT", updateURL, body)
		if reqErr != nil {
			t.Fatal("unable to build request:", err)
		}

		res, resErr := http.DefaultClient.Do(req)
		if resErr != nil {
			t.Fatal("error while calling update http request:", resErr)
		}
		t.Log("Send request successfully")

		assert.Equal(t, http.StatusOK, res.StatusCode, "expect status code to be same while calling create api")
		t.Log("update product succeeded")

		wantProductsNotification := []Product{
			{
				Id:       1,
				Brand:    "C",
				Category: "C",
				Quantity: 3,
				Price:    30,
			},
			{
				Id:       2,
				Brand:    "B",
				Category: "B",
				Quantity: 2,
				Price:    20,
			},
		}
		var gotProductsNotification []Product
		if err := conn.ReadJSON(&gotProductsNotification); err != nil {
			t.Fatal("unable to read the msg from websocket:", err)
		}
		assert.Equal(t, wantProductsNotification, gotProductsNotification, "expect same mssg while updating")
	})

	t.Run("delete", func(t *testing.T) {
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

		repo := setupInMemoryRepo(existing)
		svc := NewProductServiceImpl(repo)
		httpTransport := NewhttpTransport(svc)
		handler := buildHttpHandler(httpTransport)

		url := startHttpServer(t, handler)
		t.Log("setup for http server completed")

		websocketURL := fmt.Sprintf("ws://%s/ws", url)

		conn, httpresponse, err := websocket.DefaultDialer.Dial(websocketURL, nil)
		if err != nil {
			t.Fatal("unable to upgrade from http to websocket:", err)
		}
		assert.Equal(t, http.StatusSwitchingProtocols, httpresponse.StatusCode, "expect 101 status code while switch to websocket")

		deleteURL := fmt.Sprintf("http://%s/products/1", url)
		req, reqErr := http.NewRequest("DELETE", deleteURL, nil)
		if reqErr != nil {
			t.Fatal("failed to build request:", reqErr)
		}

		res, resErr := http.DefaultClient.Do(req)
		if resErr != nil {
			t.Fatal("unable to get response:", res)
		}
		t.Log("Send request successfully")

		assert.Equal(t, http.StatusOK, res.StatusCode, "expect same status code while calling delete api")
		t.Log("delete operation successfull")

		t.Cleanup(func() {
			if err := conn.Close(); err != nil {
				t.Log("err while closing websocket conn:", err)
			}
		})

		wantProductsNotification := []Product{
			{
				Id:       2,
				Brand:    "B",
				Category: "B",
				Quantity: 2,
				Price:    20,
			},
		}

		var gotProductsNotification []Product
		if err := conn.ReadJSON(&gotProductsNotification); err != nil {
			t.Fatal("unable toread mssg from websocket: ", err)
		}
		assert.Equal(t, wantProductsNotification, gotProductsNotification, "expect same mssg after delete")
	})

	t.Run("delete op failed", func(t *testing.T) {
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

		repo := setupInMemoryRepo(existing)
		svc := NewProductServiceImpl(repo)
		httpTransport := NewhttpTransport(svc)
		handler := buildHttpHandler(httpTransport)

		url := startHttpServer(t, handler)
		t.Log("setup for http server completed")

		websocketURL := fmt.Sprintf("ws://%s/ws", url)

		conn, httpresponse, err := websocket.DefaultDialer.Dial(websocketURL, nil)
		if err != nil {
			t.Fatal("unable to upgrade from http to websocket:", err)
		}
		assert.Equal(t, http.StatusSwitchingProtocols, httpresponse.StatusCode, "expect 101 status code while switch to websocket")

		deleteURL := fmt.Sprintf("http://%s/products/3", url)
		req, reqErr := http.NewRequest("DELETE", deleteURL, nil)
		if reqErr != nil {
			t.Fatal("failed to builed request:", reqErr)
		}

		res, resErr := http.DefaultClient.Do(req)
		if resErr != nil {
			t.Fatal("unable to get response:", res)
		}
		t.Log("Send request successfully")

		assert.Equal(t, http.StatusNotFound, res.StatusCode, "expect same status code while calling delete api")
		t.Log("got response successfully")

		t.Cleanup(func() {
			if err := conn.Close(); err != nil {
				t.Log("err while closing websocket conn:", err)
			}
		})

		if readDeadlineErr := conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond)); readDeadlineErr != nil {
			t.Fatal("err while reading websocket mssg:", readDeadlineErr)
		}

		var gotProductsNotification []Product
		readMssgErr := conn.ReadJSON(&gotProductsNotification)
		assert.Error(t, readMssgErr, "expect err while read mssg")
	})
}
