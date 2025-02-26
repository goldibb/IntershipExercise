# SWIFT Codes API

A RESTful API service for managing SWIFT/BIC codes for banks and their branches.

## Features

- Fetch SWIFT code details with branch information
- List all SWIFT codes by country
- Create new SWIFT codes
- Delete SWIFT codes (including branch cascade deletion)
- Automatic data import from Excel file
- Containerized deployment with Docker

## Tech Stack

- Go 1.21+
- PostgreSQL 16
- Docker & Docker Compose
- Gorilla Mux (HTTP router)

## API Endpoints

### GET `/v1/swift-codes/{swift_code}`
Get details of a specific SWIFT code. Returns branch information if the code is a headquarter.

### GET `/v1/swift-codes/country/{country_iso2}`
List all SWIFT codes for a specific country.

### POST `/v1/swift-codes`
Create a new SWIFT code entry.
```bash
curl -X POST http://localhost:8080/v1/swift-codes -H "Content-Type: application/json" -d '{"swiftCode":"BNPAFRPP123","address":"16 Boulevard des Italiens","bankName":"BNP Paribas","countryISO2":"PL","isHeadquarter":false}'
```

### DELETE `/v1/swift-codes/{swift_code}`
Delete a SWIFT code. If the code is a headquarter, all associated branches will be deleted.
```Bash
curl -X DELETE http://localhost:8080/v1/swift-codes/BNPAFRPP123
```
## Setup

1. Clone the repository:
```bash
git clone <repository_url>
```
2.Place your SWIFT codes Excel file in the data directory:
```bash
mkdir data
cp Interns_2025_SWIFT_CODES.xlsx data/
```
3.Start the application using Docker Compose:
```bash
docker-compose up --build
```

The API will be available at http://localhost:8080.

### Database Schema
The application uses two main tables:  
- swift_codes: Stores SWIFT/BIC codes and bank information
- countries: Stores country information

### Environment Variables
- DB_HOST: PostgreSQL host (default: db)
- DB_USER: Database user (default: postgres)
- DB_PASSWORD: Database password (default: postgres)
- DB_NAME: Database name (default: swift_codes)
- DB_PORT: Database port (default: 5432)
### Local Development
To run the application locally:
1. Set up environment variables:
```bash
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=swift_codes
export DB_PORT=5432
```
2. Create and initialize the database:
```bash
createdb swift_codes
psql -d swift_codes -f migrations/init.sql
```
3. Run the application:
```bash
go run cmd/main.go 
```

### Error Handling
The API returns standard HTTP status codes:  
- 200: Success
- 201: Created
- 400: Bad Request
- 404: Not Found
- 500: Internal Server Error

### Container Management
Stop the application:
```bash
docker-compose down
```
Remove volumes:
```bash
docker-compose down -v
```
View logs:
```bash
docker-compose logs -f
```