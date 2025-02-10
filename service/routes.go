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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get swift code by swift_code
	})
}

func (h *Handler) getSwiftCodesByCountry() http.Handler {

}

func (h *Handler) CreateSwiftCode() http.Handler {

}

func (h *Handler) DeleteSwiftCode() http.Handler {

}
