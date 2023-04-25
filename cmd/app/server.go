package app

import (
	"customer_crud/cmd/app/middleware"
	"customer_crud/pkg/customers"
	"customer_crud/pkg/security"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	mux          *mux.Router
	customersSvc *customers.Service
	securitySvc  *security.Service
}

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

func NewServer(mux *mux.Router, customersSvc *customers.Service, securitySvc *security.Service) *Server {
	return &Server{mux: mux, customersSvc: customersSvc, securitySvc: securitySvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	chMd := middleware.Basic(s.securitySvc.Auth)
	//s.mux.HandleFunc("/customers.getById", s.handleGetCustomersByID)
	//s.mux.HandleFunc("/customers.getAll", s.handleGetAllCustomers)
	//s.mux.HandleFunc("/customers.getAllActive", s.handleGetAllActiveCustomers)
	//s.mux.HandleFunc("/customers.save", s.handleSaveCustomer)
	//s.mux.HandleFunc("/customers.removeById", s.handleRemoveCustomerById)
	//s.mux.HandleFunc("/customers.blockById", s.handleBlockCustomerById)
	//s.mux.HandleFunc("/customers.unblockById", s.handleUnblockCustomerById)
	s.mux.Handle("/customers", chMd(http.HandlerFunc(s.handleGetAllCustomers))).Methods(GET)
	s.mux.Handle("/customers/block/{id}", chMd(http.HandlerFunc(s.handleBlockCustomerById))).Methods(POST)
	s.mux.Handle("/customers/unblock/{id}", chMd(http.HandlerFunc(s.handleUnblockCustomerById))).Methods(POST)
	s.mux.Handle("/customers/active", chMd(http.HandlerFunc(s.handleGetAllActiveCustomers))).Methods(GET)
	//s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods(GET)
	s.mux.Handle("/customers/{id}", chMd(http.HandlerFunc(s.handleGetCustomersByID))).Methods(GET)
	s.mux.Handle("/customers", chMd(http.HandlerFunc(s.handleSaveCustomer))).Methods(POST)
	s.mux.Handle("/customers/{id}", chMd(http.HandlerFunc(s.handleRemoveCustomerById))).Methods(DELETE)
}

func (s *Server) handleGetCustomersByID(writer http.ResponseWriter, request *http.Request) {
	//idParam := request.URL.Query().Get("id")
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
	}

	item, err := s.customersSvc.ByID(request.Context(), id)

	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)

	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) handleGetAllCustomers(writer http.ResponseWriter, request *http.Request) {
	items, err := s.customersSvc.All(request.Context())

	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)

	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) handleGetAllActiveCustomers(writer http.ResponseWriter, request *http.Request) {
	items, err := s.customersSvc.GetAllActive(request.Context())

	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)

	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) handleSaveCustomer(writer http.ResponseWriter, request *http.Request) {
	var customer *customers.Customer

	err := json.NewDecoder(request.Body).Decode(&customer)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	items, err := s.customersSvc.Save(request.Context(), customer)

	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)

	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) handleRemoveCustomerById(writer http.ResponseWriter, request *http.Request) {
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
	}
	customer := &customers.Customer{}

	customer, err = s.customersSvc.RemoveCustomerById(request.Context(), id)

	if err != nil {
		log.Print(err)
		return
	}

	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(customer)
	if err != nil {
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)

	if err != nil {
		log.Print(err)
		return
	}

}

func (s *Server) handleBlockCustomerById(writer http.ResponseWriter, request *http.Request) {
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
	}
	customer := &customers.Customer{}

	customer, err = s.customersSvc.BlockCustomerById(request.Context(), id)

	if err != nil {
		log.Print(err)
		return
	}

	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(customer)
	if err != nil {
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)

	if err != nil {
		log.Print(err)
		return
	}

}

func (s *Server) handleUnblockCustomerById(writer http.ResponseWriter, request *http.Request) {
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
	}
	customer := &customers.Customer{}

	customer, err = s.customersSvc.UnblockCustomerById(request.Context(), id)

	if err != nil {
		log.Print(err)
		return
	}

	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(customer)
	if err != nil {
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)

	if err != nil {
		log.Print(err)
		return
	}
}
