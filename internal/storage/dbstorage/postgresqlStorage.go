package dbstorage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgreSQLStorage struct {
	connection *sql.DB
}

func NewPostgreSQLStorage(connectionString string) (*PostgreSQLStorage, error) {
	pgStore := PostgreSQLStorage{}
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Printf("Unable to connect to database because %s", err)
		return &pgStore, fmt.Errorf("unable to connect to database because %w", err)
	}
	pgStore.connection = db
	return &pgStore, nil
}

func (pgs *PostgreSQLStorage) Close() {
	pgs.connection.Close()
}

func (pgs *PostgreSQLStorage) Ping() error {
	if err := pgs.connection.Ping(); err != nil {
		log.Printf("Cannot ping database because %s", err)
		return fmt.Errorf("cannot ping database because %w", err)
	}
	return nil
}
