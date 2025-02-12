package api

import (
	"IntershipExercise/service"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	address string
	db      *sql.DB
}

func NewAPIServer(address string, db *sql.DB) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/v1").Subrouter()
	apiHandler := service.NewHandler()
	apiHandler.RegisterRoutes(subrouter)
	return http.ListenAndServe(s.address, router)
}
