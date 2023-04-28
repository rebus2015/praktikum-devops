package dbstorage

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib" // init db driver for postgeSQl\
	log "github.com/sirupsen/logrus"

	"github.com/rebus2015/praktikum-devops/internal/storage"
	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

var _ storage.SecondaryStorage = new(PostgreSQLStorage)

type PostgreSQLStorage struct {
	connection *sql.DB
	Sync       bool
	context    context.Context
}

type SQLStorage interface {
	Ping(ctx context.Context) error
	Close()
}

func (pgs *PostgreSQLStorage) SyncMode() bool {
	return pgs.Sync
}

func NewStorage(ctx context.Context, connectionString string, sync bool) (*PostgreSQLStorage, error) {
	db, err := restoreDB(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	return &PostgreSQLStorage{connection: db, context: ctx, Sync: sync}, nil
}

func (pgs *PostgreSQLStorage) Close() {
	pgs.connection.Close()
}

func (pgs *PostgreSQLStorage) Ping(ctx context.Context) error {
	if pgs == nil {
		return fmt.Errorf("cannot ping database because connection is nil")
	}
	if err := pgs.connection.PingContext(ctx); err != nil {
		log.Printf("Cannot ping database because %s", err)
		return fmt.Errorf("cannot ping database because %w", err)
	}
	return nil
}

func (pgs *PostgreSQLStorage) Save(ms *memstorage.MemStorage) error {
	ctx, cancel := context.WithCancel(pgs.context)
	defer cancel()

	tx, err := pgs.connection.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		return err
	}
	defer func() {
		rberr := tx.Rollback()
		if rberr != nil {
			log.Printf("failed to rollback transaction err: %v", rberr)
		}
	}()

	for metric, val := range ms.Gauges {
		args := pgx.NamedArgs{
			"name":  metric,
			"type":  "gauge",
			"value": val,
			"delta": sql.NullInt64{Valid: true},
		}
		if _, errg := tx.ExecContext(ctx, SetMetricQuery, args); errg != nil {
			log.Printf("Error update gauge:[%v:%v] query '%s' error: %v", metric, val, SetMetricQuery, err)
			return fmt.Errorf("error update gauge:[%v:%v] query '%s' error: %w", metric, val, SetMetricQuery, err)
		}
	}

	for metric, val := range ms.Counters {
		args := pgx.NamedArgs{
			"name":  metric,
			"type":  "counter",
			"value": sql.NullFloat64{Valid: true},
			"delta": val,
		}
		if _, errc := tx.ExecContext(ctx, SetMetricQuery, args); errc != nil {
			log.Printf("Error update counter:[%v:%v] query '%s' error: %v", metric, val, SetMetricQuery, err)
			return fmt.Errorf("error update counter:[%v:%v] query '%s' error: %w", metric, val, SetMetricQuery, err)
		}
	}
	// шаг 4 — сохраняем изменения
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to execute transaction %w", err)
	}

	return nil
}

func (pgs *PostgreSQLStorage) Restore() (*memstorage.MemStorage, error) {
	ctx, cancel := context.WithTimeout(pgs.context, time.Second*5)
	defer cancel()
	counters := make(map[string]int64)
	gauges := make(map[string]float64)
	rows, err := pgs.connection.QueryContext(ctx, GetMetricsQuery)
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
		errs := pgs.Save(ms)
		if errs != nil {
			log.Printf("FileStorage Save error: %v", errs)
		}
	}
}

func restoreDB(ctx context.Context, connectionString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Printf("Unable to open connection to database connection:'%v'  error %s", connectionString, err)
		return nil, fmt.Errorf("unable to connect to database because %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		log.Printf("Cannot ping database because %s", err)
		return nil, fmt.Errorf("cannot ping database because %w", err)
	}

	_, err = db.ExecContext(ctx, restoreDBscript)
	if err != nil {
		log.Printf("Fail to invoke %s: %v", restoreDBscript, err)
		return nil, fmt.Errorf("fail to invoke %s: %w", restoreDBscript, err)
	}
	return db, nil
}
