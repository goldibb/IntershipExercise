package parser

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

type SwiftRecord struct {
	SwiftCode     string
	Address       string
	BankName      string
	CountryISO2   string
	CountryName   string
	IsHeadquarter bool
}

type ParsedData struct {
	Headquarters []SwiftRecord
	Branches     []SwiftRecord
}

func ParsedExcelFile(filePath string) (*ParsedData, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(f *excelize.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("failed to close file: %v", err)
		}
	}(f)

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	data := &ParsedData{}
	for i, row := range rows {
		if i == 0 {
			continue
		}

		// Validate SWIFT code format
		if len(row[1]) != 11 {
			return nil, fmt.Errorf("invalid SWIFT code length in row %d: %s", i+1, row[1])
		}

		adress := strings.TrimSpace(row[4])
		if adress == "" {
			adress = strings.TrimSpace(row[5]) + ", " + strings.TrimSpace(row[6])
		}
		swiftRecord := SwiftRecord{
			CountryISO2:   row[0],
			SwiftCode:     row[1],
			BankName:      row[3],
			Address:       adress,
			CountryName:   row[6],
			IsHeadquarter: isHeadquarterCode(row[1]),
		}
		if swiftRecord.IsHeadquarter {
			data.Headquarters = append(data.Headquarters, swiftRecord)
		} else {
			data.Branches = append(data.Branches, swiftRecord)
		}
	}
	return data, nil
}
func isHeadquarterCode(swiftCode string) bool {
	return swiftCode[8:] == "XXX"
}
