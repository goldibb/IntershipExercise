# Dockerfile
FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN apk add --no-cache postgresql-client

RUN go mod download
RUN go build -o main ./cmd/main.go

EXPOSE 8080

COPY scripts/wait-for-postgres.sh /wait-for-postgres.sh
RUN chmod +x /wait-for-postgres.sh

CMD ["/bin/sh", "/wait-for-postgres.sh", "./main"]