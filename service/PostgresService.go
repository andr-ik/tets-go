package service

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresService struct {
	DB *sql.DB
}

func NewPostgresService(host string, port int, user, password, dbname string) PostgresService {
	connection := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	if password != "" {
		connection += fmt.Sprintf(" password=%s", password)
	}
	db, err := sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	service := PostgresService{
		DB: db,
	}

	return service
}

func (postgresService *PostgresService) Close() {
	postgresService.DB.Close()
}
