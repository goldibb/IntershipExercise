package main

import (
	"IntershipExercise/cmd/api"
	"IntershipExercise/internal/db"
	"IntershipExercise/pkg/parser"
	"log"
)

func main() {
	err := db.OpenDatabase()
	if err != nil {
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer func() {
		err := db.CloseDatabase()
		if err != nil {
			log.Fatalf("Error: Unable to close database connection: %v", err)
		}
	}()

	data, err := parser.ParsedExcelFile("c:/data/Interns_2025_SWIFT_CODES.xlsx")
	if err != nil {
		log.Fatalf("Error parsing Excel file: %v", err)
	}

	// Save the parsed data to database
	if err := db.SaveParsedData(data); err != nil {
		log.Fatalf("Error saving data to database: %v", err)
	}

	server := api.NewAPIServer(":8080", db.DB)
	log.Println("Starting server on port 8080")
	if err := server.Run(); err != nil {
		log.Fatalf("Error: Unable to start server: %v", err)
	}

}
