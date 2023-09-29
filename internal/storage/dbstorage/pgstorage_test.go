package dbstorage

import (
	"context"
	"database/sql"
	"reflect"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

func TestPostgreSQLStorage_Ping(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	mock.ExpectPing().WillDelayFor(time.Second * 3)
	pgs := &PostgreSQLStorage{
		connection: mock,
		Sync:       false,
	}
	if err := pgs.Ping(ctx); err != nil {
		t.Errorf("PostgreSQLStorage.Ping() error = %v", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func Test_restoreDB(t *testing.T) {

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	mock.ExpectPing().WillDelayFor(time.Second * 3)
	mock.ExpectExec(regexp.QuoteMeta(restoreDBscript)).WillReturnResult(pgxmock.NewResult("", 0))

	pgs := &PostgreSQLStorage{
		connection: mock,
		Sync:       false,
	}
	if err := pgs.restoreDB(ctx); err != nil {
		t.Errorf("PostgreSQLStorage.Ping() error = %v", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestPostgreSQLStorage_Restore(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	mockDB := mock.NewRows([]string{"name", "type", "value", "delta"}).
		AddRow("metric1", "gauge", "231.12", nil).
		AddRow("metric2", "counter", nil, "101")

	mock.ExpectQuery(regexp.QuoteMeta(GetMetricsQuery)).WillReturnRows(mockDB)

	pgs := &PostgreSQLStorage{
		connection: mock,
		Sync:       false,
	}

	gauges := map[string]float64{
		"metric1": 231.12,
	}
	counters := map[string]int64{
		"metric2": 101,
	}

	ms, err := pgs.Restore(ctx)
	if err != nil {
		t.Errorf("PostgreSQLStorage.Ping() error = %v", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.True(t, reflect.DeepEqual(gauges, ms.Gauges))
	assert.True(t, reflect.DeepEqual(counters, ms.Counters))
}

func TestPostgreSQLStorage_Save(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ctx := context.Background()

	defer func() {
		mock.Close()
	}()

	mock.ExpectBegin()
	args := pgx.NamedArgs{
		"name":  "metric1",
		"type":  "gauge",
		"value": 231.12,
		"delta": sql.NullInt64{Valid: true},
	}
	mock.ExpectExec(regexp.QuoteMeta(SetMetricQuery)).WithArgs(args).WillReturnResult(pgxmock.NewResult("", 0))
	args2 := pgx.NamedArgs{
		"name":  "metric2",
		"type":  "counter",
		"value": sql.NullFloat64{Valid: true},
		"delta": int64(101),
	}
	mock.ExpectExec(regexp.QuoteMeta(SetMetricQuery)).WithArgs(args2).WillReturnResult(pgxmock.NewResult("", 0))
	mock.ExpectCommit()

	pgs := &PostgreSQLStorage{
		connection: mock,
		Sync:       false,
	}

	ms := memstorage.MemStorage{
		Gauges: map[string]float64{
			"metric1": 231.12},
		Counters: map[string]int64{
			"metric2": 101,
		},
		Mux: &sync.RWMutex{},
	}
	ms.Mux.Lock()
	if errSave := pgs.Save(ctx, &ms); errSave != nil {
		t.Errorf("Save memstorage error = %v", err)
	}
	ms.Mux.Unlock()
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgreSQLStorage_SyncMode(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()
	pgs := &PostgreSQLStorage{
		connection: mock,
		Sync:       false,
	}

	assert.Equal(t, pgs.SyncMode(), pgs.Sync)

}

func TestPostgreSQLStorage_Close(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	pgs := &PostgreSQLStorage{
		connection: mock,
		Sync:       false,
	}
	pgs.Close()

	assert.Error(t, pgs.connection.Ping(context.Background()))
}
