package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type httpTransport struct {
	service ProductService
}

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func writeError(w http.ResponseWriter, statusCode int, errors ...string) {
	w.WriteHeader(statusCode)

	if len(errors) == 0 {
		return
	}

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Println("failed to encode:", err)
	}
}

func handleError(w http.ResponseWriter, err error) {
	var se *json.SyntaxError
	if errors.As(err, &se) {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if errors.Is(err, errDuplicateId) {
		writeError(w, http.StatusConflict, "product exists")
		return
	}
	if errors.Is(err, errProductNotFound) {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	var ve *validationError
	if errors.As(err, &ve) {
		writeError(w, http.StatusBadRequest, ve.failures...)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

func NewhttpTransport(svc ProductService) *httpTransport {
	return &httpTransport{
		service: svc,
	}
}

func buildHttpHandler(t *httpTransport) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/products", t.Create).Methods("POST")
	r.HandleFunc("/products/{id}", t.Update).Methods("PUT")
	r.HandleFunc("/products", t.GetAll).Methods("GET")
	r.HandleFunc("/products/{id}", t.GetById).Methods("GET")
	r.HandleFunc("/products/{id}", t.Delete).Methods("DELETE")
	r.HandleFunc("/ws", t.wsEndpoint)
	return r
}

func (t *httpTransport) Create(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Println("failed to decode:", err)
		handleError(w, err)
		return
	}

	if err := t.service.Create(product); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	gotProduct, gotProductErr := t.service.GetById(product.Id)

	if gotProductErr != nil {
		log.Println("unable to get product:", gotProductErr)
		handleError(w, gotProductErr)
		return
	}

	if err := json.NewEncoder(w).Encode(gotProduct); err != nil {
		log.Println("failed to encode:", err)
		return
	}
}

func (t *httpTransport) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idstr := vars["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid ID")
		return
	}

	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		handleError(w, err)
		return
	}

	product.Id = id
	if err := t.service.Update(product); err != nil {
		handleError(w, err)
		return
	}

	gotProduct, err := t.service.GetById(product.Id)

	if err != nil {
		log.Println("unable to get product:", err)
		handleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(gotProduct); err != nil {
		log.Println("failed to encode:", err)
		return
	}
}

func (t *httpTransport) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idstr := vars["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	product, err := t.service.GetById(id)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Println("failed to encode:", err)
		return
	}
}

func (t *httpTransport) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := t.service.GetAll()
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Println("failed to encode:", err)
		return
	}
}

func (t *httpTransport) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idstr := vars["id"]
	id, err := strconv.Atoi(idstr)

	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := t.service.Delete(id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

var upgrader = websocket.Upgrader{ //struct
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (t *httpTransport) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error while upgrading HTTP to Web Socket:", err)
		return
	}

	defer func() {
		if err := ws.Close(); err != nil {
			log.Println("err while closing websocket conn:", err)
		}
	}()

	log.Println("client connected", ws.RemoteAddr())

	sub := &subscribeWebSocket{w: ws}
	if err := t.service.subscribe(sub); err != nil {
		log.Println("error while subscribe with web socket:", err)
		return
	}

	defer func() {
		log.Println("unsubscribing", ws.RemoteAddr())
		if err := t.service.unsubscribe(sub); err != nil {
			log.Println("failed to subscribe:", err)
		}
	}()

	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			log.Println("failed to read", err)
			return
		}
	}
}

type subscribeWebSocket struct {
	w *websocket.Conn
}

func (s *subscribeWebSocket) Update(products []Product) {
	if err := s.w.WriteJSON(products); err != nil {
		log.Println("error while encoding", err)
		return
	}
}

func (s *subscribeWebSocket) Id() string {
	return s.w.RemoteAddr().String()
}
