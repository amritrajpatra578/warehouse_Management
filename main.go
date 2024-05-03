package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var products []Product

// CRUD

// create
// read by id
// read all
// delete
// update

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("!! welcome to the store !!\n"))
}

func createProduct(w http.ResponseWriter, r *http.Request) { // create()
	var product Product

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil { //decodeing from json to go lang
		log.Println("failed to decode:", err)
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}

	for _, value := range products {
		if value.Id == product.Id {
			w.WriteHeader(http.StatusConflict) //409
			return
		}
	}

	products = append(products, product)

	w.WriteHeader(http.StatusCreated) //201
}

func viewAllproduct(w http.ResponseWriter, r *http.Request) { // getALL()
	j, err := json.Marshal(products) //encoding from go lang to json
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		return
	}

	w.WriteHeader(http.StatusOK) //200
	if _, err := w.Write(j); err != nil {
		log.Println("failed to write response:", err)
		return
	}

}
func viewproduct(w http.ResponseWriter, r *http.Request) { // getByID()
	vars := mux.Vars(r)

	idStr := vars["id"]
	id, err := strconv.Atoi(idStr) //converting from string to int
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400
		log.Printf("failed to decode ID('%s'): %v\n", idStr, err)
		return
	}

	var product Product
	for _, currentProduct := range products {
		if currentProduct.Id == id {
			w.WriteHeader(http.StatusOK)                               //200
			if err := json.NewEncoder(w).Encode(product); err != nil { //encoding from go lang to json
				log.Println("failed to decode:", err)
			}
			return
		}
	}

	w.WriteHeader(http.StatusNotFound) //404
}

func deleteProduct(w http.ResponseWriter, r *http.Request) { // delete()
	vars := mux.Vars(r)

	idstr := vars["id"]
	id, err := strconv.Atoi(idstr) //converting from string to int
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400
		log.Println("falied to decode Id: ", err)
		return
	}

	var product Product
	for index, currentProduct := range products {
		if currentProduct.Id == id {
			products = append(products[:index], products[index+1:]...)
			product = currentProduct
			break
		}
	}

	if product.Id == 0 {
		w.WriteHeader(http.StatusNotFound) //404
		return
	}
	w.WriteHeader(http.StatusNoContent) //204

}

func updateProduct(w http.ResponseWriter, r *http.Request) { // update()
	vars := mux.Vars(r)

	idstr := vars["id"]
	id, err := strconv.Atoi(idstr) //got the id value from path
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400			//then we are checking for error while converting from string to int
		log.Println("falied to decode Id: ", err)
		return
	}

	var updatedProduct Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil { //decoding from json to go lang
		w.WriteHeader(http.StatusBadRequest) //400
		log.Println("failed to decode:", err)
		return
	}

	for index, currentProduct := range products {
		if currentProduct.Id == id {
			updatedProduct.Id = id // assign the request body id value with the id that
			products[index] = updatedProduct

			if err := json.NewEncoder(w).Encode(products[index]); err != nil { //encoding from go lang to json
				log.Println("failed to encode", err)
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound) //404
}

func handleRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/home", Home).Methods("GET")
	r.HandleFunc("/products", createProduct).Methods("POST")
	r.HandleFunc("/products", viewAllproduct).Methods("GET")
	r.HandleFunc("/products/{id}", viewproduct).Methods("GET")
	r.HandleFunc("/products/{id}", deleteProduct).Methods("DELETE")
	r.HandleFunc("/products/{id}", updateProduct).Methods("PUT")

	return r
}

func oldmain() {
	r := handleRoutes()

	err := http.ListenAndServe(":4000", r)
	log.Println("http server exiting", err)
}

func connectPostgres(user, password, address, dbName string) *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, address, dbName)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	return db
}

func main() {
	db := connectPostgres("postgres", "postgres", "127.0.0.1:5432", "productsdb")

	if _, err := db.Exec("select 1"); err != nil {
		log.Fatalln("failed to connect to db:", err)
	}

	repo := NewPostgresRepo(db)
	svc := NewProductServiceImpl(repo)
	transport := NewhttpTransport(svc)

	httpHandler := buildHttpHandler(transport)

	err := http.ListenAndServe(":5000", httpHandler)
	log.Println("http server exiting:", err)
}
