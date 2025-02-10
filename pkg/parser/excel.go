package parser

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
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
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	data := &ParsedData{}
	for i, row := range rows {
		if i == 0 {
			continue
		}

		adress := row[4]
		if adress == "" {
			adress = row[5]
		}
		swiftRecord := SwiftRecord{
			CountryISO2:   row[0],
			SwiftCode:     row[1],
			BankName:      row[3],
			Address:       strings.TrimSpace(adress),
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
