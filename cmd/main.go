package main

import (
	"IntershipExercise/cmd/api"
	"IntershipExercise/internal/db"
	"IntershipExercise/pkg/parser"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func findExcelFile() (string, error) {
	fileName := "Interns_2025_SWIFT_CODES.xlsx"
	searchPaths := []string{
		"/data",   // Docker mounted volume
		"C:/data", // Windows local path
		"./data",
		"c:/data", // Local relative path
	}

	for _, basePath := range searchPaths {
		fullPath := filepath.Join(basePath, fileName)
		if _, err := os.Stat(fullPath); err == nil {
			log.Printf("Found file at: %s", fullPath)
			return fullPath, nil
		}
	}

	return "", fmt.Errorf("excel file %s not found in search paths", fileName)
}
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

	filePath, err := findExcelFile()
	if err != nil {
		log.Fatalf("Error finding Excel file: %v", err)
	}

	log.Printf("Found Excel file at: %s", filePath)
	data, err := parser.ParsedExcelFile(filePath)
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
