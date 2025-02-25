package service

import (
	"IntershipExercise/internal/models"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type ApiError struct {
	Message string `json:"message"`
}
type Handler struct {
	db *sql.DB
}

type CountrySwiftCodesResponse struct {
	CountryISO2 string                       `json:"countryISO2"`
	CountryName string                       `json:"countryName"`
	SwiftCodes  []SwiftCodeByCountryResponse `json:"swiftCodes"`
}

type SwiftCodeByCountryResponse struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type SwiftCodeResponse struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.Handle("/swift-codes/{swift_code}", h.getSwiftCode()).Methods("GET")
	router.Handle("/swift-codes/country/{country_iso2}", h.getSwiftCodesByCountry()).Methods("GET")
	router.Handle("/swift-codes", h.CreateSwiftCode()).Methods("POST")
	router.Handle("/swift-codes/{swift_code}", h.DeleteSwiftCode()).Methods("DELETE")
}

func (h *Handler) getSwiftCode() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		swiftCode := mux.Vars(r)["swift_code"]

		var basicInfo SwiftCodeResponse
		err := h.db.QueryRow(`
            SELECT s.swift_code, s.address, s.bank_name, s.country_iso2,
                   c.name as country_name, s.is_headquarter
            FROM swift_codes s
            JOIN countries c ON s.country_iso2 = c.iso2_code
            WHERE s.swift_code = $1
        `, swiftCode).Scan(
			&basicInfo.SwiftCode,
			&basicInfo.Address,
			&basicInfo.BankName,
			&basicInfo.CountryISO2,
			&basicInfo.CountryName,
			&basicInfo.IsHeadquarter,
		)
		if errors.Is(err, sql.ErrNoRows) {
			WriteJSON(w, http.StatusNotFound, ApiError{Message: "Swift code not found"})
			return
		}
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ApiError{Message: "Database error"})
			return
		}

		if !basicInfo.IsHeadquarter {
			WriteJSON(w, http.StatusOK, basicInfo)
			return
		}

		response := models.Headquarter{
			SwiftCode:     basicInfo.SwiftCode,
			Address:       basicInfo.Address,
			BankName:      basicInfo.BankName,
			CountryISO2:   basicInfo.CountryISO2,
			CountryName:   basicInfo.CountryName,
			IsHeadquarter: basicInfo.IsHeadquarter,
		}

		prefix := swiftCode[:7]
		rows, err := h.db.Query(`
            SELECT s.swift_code, s.address, s.bank_name, s.country_iso2,
                   c.name as country_name, s.is_headquarter
            FROM swift_codes s
            JOIN countries c ON s.country_iso2 = c.iso2_code
            WHERE s.swift_code LIKE $1 || '%'
            AND s.swift_code != $2
            AND s.is_headquarter = false
        `, prefix, swiftCode)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ApiError{Message: "Database error"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var branch models.Branch
			if err := rows.Scan(
				&branch.SwiftCode,
				&branch.Address,
				&branch.BankName,
				&branch.CountryISO2,
				&branch.CountryName,
				&branch.IsHeadquarter,
			); err != nil {
				WriteJSON(w, http.StatusInternalServerError, ApiError{Message: "Database error"})
				return
			}
			response.Branches = append(response.Branches, branch)
		}

		WriteJSON(w, http.StatusOK, response)
	})
}

func (h *Handler) getSwiftCodesByCountry() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		countryISO2 := mux.Vars(r)["country_iso2"]

		var countryName string
		err := h.db.QueryRow(`SELECT name FROM countries WHERE iso2_code = $1`,
			countryISO2).Scan(&countryName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				WriteJSON(w, http.StatusNotFound, ApiError{Message: "Country not found"})
				return
			}
			WriteJSON(w, http.StatusInternalServerError, ApiError{Message: "Database error"})
			return
		}

		response := CountrySwiftCodesResponse{
			CountryISO2: countryISO2,
			CountryName: countryName,
			SwiftCodes:  []SwiftCodeByCountryResponse{},
		}

		rows, err := h.db.Query(`
            SELECT swift_code, address, bank_name, country_iso2, is_headquarter
            FROM swift_codes
            WHERE country_iso2 = $1
        `, countryISO2)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ApiError{Message: "Database error"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var branch SwiftCodeByCountryResponse
			if err := rows.Scan(
				&branch.SwiftCode,
				&branch.Address,
				&branch.BankName,
				&branch.CountryISO2,
				&branch.IsHeadquarter,
			); err != nil {
				WriteJSON(w, http.StatusInternalServerError, ApiError{Message: "Database error"})
				return
			}
			response.SwiftCodes = append(response.SwiftCodes, branch)
		}

		WriteJSON(w, http.StatusOK, response)
	})
}

func (h *Handler) CreateSwiftCode() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var swiftCode models.Branch
		if err := json.NewDecoder(r.Body).Decode(&swiftCode); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Message: "Invalid request payload"})
			return
		}

		_, err := h.db.Exec(`
			INSERT INTO swift_codes (swift_code, address, bank_name, country_iso2, is_headquarter)
			VALUES ($1, $2, $3, $4, $5)
		`, swiftCode.SwiftCode, swiftCode.Address, swiftCode.BankName, swiftCode.CountryISO2, swiftCode.IsHeadquarter)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ApiError{Message: "Database error"})
			return
		}

		WriteJSON(w, http.StatusCreated, ApiError{Message: "Swift code created successfully"})
	})
}

func (h *Handler) DeleteSwiftCode() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		swiftCode := mux.Vars(r)["swift_code"]

		tx, err := h.db.Begin()
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ApiError{Message: "Database error"})
			return
		}
		defer tx.Rollback()

		var isHeadquarter bool
		err = tx.QueryRow("SELECT is_headquarter FROM swift_codes WHERE swift_code = $1", swiftCode).Scan(&isHeadquarter)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				WriteJSON(w, http.StatusNotFound, ApiError{Message: "Swift code not found"})
				return
			}
			WriteJSON(w, http.StatusInternalServerError, ApiError{Message: "Database error"})
			return
		}

		if isHeadquarter {
			prefix := swiftCode[:7]
			_, err = tx.Exec("DELETE FROM swift_codes WHERE swift_code LIKE $1 || '%' AND swift_code != $2", prefix, swiftCode)
			if err != nil {
				WriteJSON(w, http.StatusInternalServerError, ApiError{Message: "Database error"})
				return
			}
		}

		_, err = tx.Exec("DELETE FROM swift_codes WHERE swift_code = $1", swiftCode)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ApiError{Message: "Database error"})
			return
		}

		if err = tx.Commit(); err != nil {
			WriteJSON(w, http.StatusInternalServerError, ApiError{Message: "Database error"})
			return
		}

		WriteJSON(w, http.StatusOK, ApiError{Message: "Swift code deleted successfully"})
	})
}
