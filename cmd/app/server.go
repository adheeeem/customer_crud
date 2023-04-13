package app

import (
	"customer_crud/pkg/customers"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	mux          *http.ServeMux
	customersSvc *customers.Service
}

func NewServer(mux *http.ServeMux, customersSvc *customers.Service) *Server {
	return &Server{mux: mux, customersSvc: customersSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}
func (s *Server) Init() {
	s.mux.HandleFunc("/customers.getById", s.handleGetCustomersByID)
	s.mux.HandleFunc("/customers.getAll", s.handleGetAllCustomers)
	s.mux.HandleFunc("/customers.getAllActive", s.handleGetAllActiveCustomers)
	s.mux.HandleFunc("/customers.save", s.handleSaveCustomer)
	s.mux.HandleFunc("/customers.removeById", s.handleRemoveCustomerById)
	s.mux.HandleFunc("/customers.blockById", s.handleBlockCustomerById)
	s.mux.HandleFunc("/customers.unblockById", s.handleUnblockCustomerById)
}

func (s *Server) handleGetCustomersByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

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
	idParam := request.FormValue("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
	}
	customer := &customers.Customer{}

	customer.ID = id
	customer.Name = request.FormValue("name")
	customer.Phone = request.FormValue("phone")
	if request.FormValue("active") == "true" || request.FormValue("active") == "" {
		customer.Active = true
	} else if request.FormValue("active") == "false" {
		customer.Active = false
	}
	layout := "2006-01-02T15:04:05.000Z"
	if request.FormValue("created") == "" {
		customer.Created = time.Now()
	} else {
		customer.Created, err = time.Parse(layout, request.FormValue("created"))
	}
	if err != nil {
		log.Print(err)
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
	idParam := request.FormValue("id")
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
	idParam := request.FormValue("id")
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
	idParam := request.FormValue("id")
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
