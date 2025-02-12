package service

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.Handle("/swift-codes/{swift_code}", h.getSwiftCode()).Methods("GET")
	router.Handle("/swift-codes/country/{country_iso2}", h.getSwiftCodesByCountry()).Methods("GET")
	router.Handle("/swift-codes", h.CreateSwiftCode()).Methods("POST")
	router.Handle("/swift-codes/{swift_code}", h.DeleteSwiftCode()).Methods("DELETE")
}

func (h *Handler) getSwiftCode() http.Handler {
	return nil
}

func (h *Handler) getSwiftCodesByCountry() http.Handler {
	return nil
}

func (h *Handler) CreateSwiftCode() http.Handler {
	return nil
}

func (h *Handler) DeleteSwiftCode() http.Handler {
	return nil
}
