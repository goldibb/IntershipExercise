package api

import ()

type Branch struct {
	address       string
	bankName      string
	countryISO2   string
	countryName   string
	isHeadquarter bool
	swiftCode     string
}

type Headquarter struct {
	address       string
	bankName      string
	countryISO2   string
	countryName   string
	isHeadquarter bool
	swiftCode     string
	branches      []Branch
}
