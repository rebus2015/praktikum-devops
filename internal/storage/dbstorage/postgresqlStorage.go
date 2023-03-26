package dbstorage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" // init db driver for postgeSQl
)

type PostgreSQLStorage struct {
	connection *sql.DB
}

type SQLStorage interface {
	Ping(ctx context.Context) error
	Close()
}

func NewPostgreSQLStorage(connectionString string) (*PostgreSQLStorage, error) {
	pgStore := PostgreSQLStorage{}
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Printf("Unable to open connection to database connection:'%v'  error %s", connectionString, err)
		return &pgStore, fmt.Errorf("unable to connect to database because %w", err)
	}
	pgStore.connection = db
	return &pgStore, nil
}

func (pgs *PostgreSQLStorage) Close() {
	pgs.connection.Close()
}

func (pgs *PostgreSQLStorage) Ping(ctx context.Context) error {
	if err := pgs.connection.PingContext(ctx); err != nil {
		log.Printf("Cannot ping database because %s", err)
		return fmt.Errorf("cannot ping database because %w", err)
	}
	return nil
}
