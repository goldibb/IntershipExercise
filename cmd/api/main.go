package main

import (
	"IntershipExercise/pkg/parser"
	"fmt"
)

func main() {
	parsedData, err := parser.ParsedExcelFile("data/swift_codes.xlsx")
	if err != nil {
		panic(fmt.Errorf("failed to parse Excel: %w", err))
	}

}
