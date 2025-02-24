package db

import (
	"IntershipExercise/pkg/parser"
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func OpenDatabase() error {
	var err error

	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	dbInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	DB, err = sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}
	return nil
}
func CloseDatabase() error {
	return DB.Close()
}

func SaveParsedData(data *parser.ParsedData) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			fmt.Printf("failed to rollback transaction: %v", err)
		}
	}(tx)

	// Insert countries
	for _, hq := range data.Headquarters {
		_, err := tx.Exec(`
            INSERT INTO countries (iso2_code, name)
            VALUES ($1, $2)
            ON CONFLICT (iso2_code) DO UPDATE SET name = EXCLUDED.name
        `, hq.CountryISO2, hq.CountryName)
		if err != nil {
			return fmt.Errorf("failed to insert country %s: %w", hq.CountryISO2, err)
		}
	}

	// Insert headquarters
	for _, hq := range data.Headquarters {
		_, err := tx.Exec(`
            INSERT INTO swift_codes (
                swift_code, address, bank_name, country_iso2,
                is_headquarter
            )
            VALUES ($1, $2, $3, $4, true)
            ON CONFLICT (swift_code) DO UPDATE SET
                address = EXCLUDED.address,
                bank_name = EXCLUDED.bank_name,
                country_iso2 = EXCLUDED.country_iso2,
                is_headquarter = true
        `, hq.SwiftCode, hq.Address, hq.BankName, hq.CountryISO2)
		if err != nil {
			return fmt.Errorf("failed to insert headquarter %s: %w", hq.SwiftCode, err)
		}
	}

	// Insert branches
	for _, branch := range data.Branches {
		_, err := tx.Exec(`
            INSERT INTO swift_codes (
                swift_code, address, bank_name, country_iso2,
                is_headquarter
            )
            VALUES ($1, $2, $3, $4, false)
            ON CONFLICT (swift_code) DO UPDATE SET
                address = EXCLUDED.address,
                bank_name = EXCLUDED.bank_name,
                country_iso2 = EXCLUDED.country_iso2,
                is_headquarter = false
        `, branch.SwiftCode, branch.Address, branch.BankName, branch.CountryISO2)
		if err != nil {
			return fmt.Errorf("failed to insert branch %s: %w", branch.SwiftCode, err)
		}
	}

	return tx.Commit()
}
