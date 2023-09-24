package dbstorage

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // init db driver for postgeSQl\

	log "github.com/sirupsen/logrus"

	"github.com/rebus2015/praktikum-devops/internal/storage"
	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

var _ storage.SecondaryStorage = new(PostgreSQLStorage)

type PgxPoolIface interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Ping(ctx context.Context) error
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Close()
}

type PostgreSQLStorage struct {
	connection PgxPoolIface //*pgxpool.Pool
	Sync       bool
}

type SQLStorage interface {
	Ping(ctx context.Context) error
	Close()
}

func (pgs *PostgreSQLStorage) SyncMode() bool {
	return pgs.Sync
}

func NewStorage(ctx context.Context, connectionString string, sync bool) (*PostgreSQLStorage, error) {
	db, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		log.Printf("Unable to open connection to database connection:'%v'  error %s", connectionString, err)
		return nil, fmt.Errorf("unable to connect to database because %w", err)
	}
	pgs := &PostgreSQLStorage{connection: db, Sync: sync}
	err1 := pgs.restoreDB(ctx)
	if err1 != nil {
		return nil, err1
	}
	return pgs, nil
}

func (pgs *PostgreSQLStorage) Close() {
	pgs.connection.Close()
}

func (pgs *PostgreSQLStorage) Ping(ctx context.Context) error {
	if pgs == nil {
		return fmt.Errorf("cannot ping database because connection is nil")
	}
	if err := pgs.connection.Ping(ctx); err != nil {
		log.Printf("Cannot ping database because %s", err)
		return fmt.Errorf("cannot ping database because %w", err)
	}
	return nil
}

func (pgs *PostgreSQLStorage) Save(ctx context.Context, ms *memstorage.MemStorage) (err error) {
	tx, err := pgs.connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("Error: [PostgreSQLStorage] failed connection transaction err: %v", err)
		return err
	}

	defer func() {
		if err != nil {
			errtx := tx.Rollback(ctx)
			if errtx != nil {
				log.Printf("Error: [PostgreSQLStorage] failed to Rollback transaction err: %v", err)
			}

		}
	}()

	for metric, val := range ms.Gauges {
		args := pgx.NamedArgs{
			"name":  metric,
			"type":  "gauge",
			"value": val,
			"delta": sql.NullInt64{Valid: true},
		}
		if _, errg := tx.Exec(ctx, SetMetricQuery, args); errg != nil {
			log.Printf("Error update gauge:[%v:%v] query '%s' error: %v", metric, val, SetMetricQuery, errg)
			return fmt.Errorf("error update gauge:[%v:%v] query '%s' error: %w", metric, val, SetMetricQuery, errg)
		}
	}

	for metric, val := range ms.Counters {
		args := pgx.NamedArgs{
			"name":  metric,
			"type":  "counter",
			"value": sql.NullFloat64{Valid: true},
			"delta": val,
		}
		if _, errc := tx.Exec(ctx, SetMetricQuery, args); errc != nil {
			log.Printf("Error update counter:[%v:%v] query '%s' error: %v", metric, val, SetMetricQuery, errc)
			return fmt.Errorf("error update counter:[%v:%v] query '%s' error: %w", metric, val, SetMetricQuery, errc)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Error failed to Commit transaction %v", err)
		return fmt.Errorf("failed to Commit transaction %w", err)
	}
	return nil
}

func (pgs *PostgreSQLStorage) Restore(ctx context.Context) (*memstorage.MemStorage, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*25)
	defer cancel()
	counters := make(map[string]int64)
	gauges := make(map[string]float64)
	rows, err := pgs.connection.Query(ctx, GetMetricsQuery)
	if err != nil {
		log.Printf("Error trying to get all metircs, query: '%s' error: %v", SetMetricQuery, err)
		return nil, fmt.Errorf("error trying to get all metircs, query: '%s' error: %w", SetMetricQuery, err)
	}
	for rows.Next() {
		var m dbMetric
		err = rows.Scan(&m.Name, &m.MType, &m.Value, &m.Delta)
		if err != nil {
			log.Printf("Error trying to Scan Rows error: %v", err)
			return nil, fmt.Errorf("error trying to Scan Rows error: %w", err)
		}
		switch m.MType.String {
		case "gauge":
			gauges[m.Name.String] = m.Value.Float64
		case "counter":
			counters[m.Name.String] = m.Delta.Int64
		default:
			return nil, fmt.Errorf("error parsing metric type '%v'", m)
		}
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		log.Printf("Error trying to get all metircs, query: '%s' error: %v", SetMetricQuery, err)
		return nil, err
	}
	return &memstorage.MemStorage{
			Counters: counters,
			Gauges:   gauges,
			Mux:      &sync.RWMutex{},
		},
		nil
}

func (pgs *PostgreSQLStorage) SaveTicker(storeint time.Duration, ms *memstorage.MemStorage) {
	ticker := time.NewTicker(storeint)
	for range ticker.C {
		ctx := context.Background()
		errs := pgs.Save(ctx, ms)
		if errs != nil {
			log.Printf("PostgreSQLStorage SaveTicker error: %v", errs)
		}
	}
}

func (pgs *PostgreSQLStorage) restoreDB(ctx context.Context) error {
	if err := pgs.Ping(ctx); err != nil {
		log.Printf("Cannot ping database because %s", err)
		return fmt.Errorf("cannot ping database because %w", err)
	}

	if _, err := pgs.connection.Exec(ctx, restoreDBscript); err != nil {
		log.Printf("Fail to invoke %s: %v", restoreDBscript, err)
		return fmt.Errorf("fail to invoke %s: %w", restoreDBscript, err)
	}
	return nil
}
