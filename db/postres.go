package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres Connection
)

type postgresRepo struct{}

// NewPostgresRepo func
func NewPostgresRepo() DatabaseRepo {
	return &postgresRepo{}
}

// Connect func
func (*postgresRepo) Connect(host string, port int, user string, password string, dbName string, sslMode string) (*sqlx.DB, error) {
	connectionStr := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		host,
		port,
		user,
		dbName,
		password,
		sslMode,
	)

	log.Println(connectionStr)

	dbConn, err := sqlx.Open("postgres", connectionStr)
	if err != nil {
		log.Println("Error Postgres Database Connection..")
		log.Println(err)
		os.Exit(3)
	}

	dbConn.SetConnMaxLifetime(time.Minute * 3)
	dbConn.SetMaxOpenConns(10)
	dbConn.SetMaxIdleConns(2)

	err = dbConn.Ping()
	if err != nil {
		log.Println("Could't Connect to Postgres Database..")
		os.Exit(3)
	}

	log.Println("Postgres Database Connected..")

	return dbConn, nil
}
