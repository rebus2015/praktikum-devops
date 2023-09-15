package dbstorage

import (
	"context"
	"reflect"
	"regexp"

	// "sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"

	// "github.com/rebus2015/praktikum-devops/internal/storage/memstorage"

	"github.com/stretchr/testify/assert"
)

func TestPostgreSQLStorage_Ping(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	mock.ExpectPing().WillDelayFor(time.Second * 3)
	pgs := &PostgreSQLStorage{
		connection: db,
		Sync:       false,
		context:    ctx,
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

	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	mock.ExpectPing().WillDelayFor(time.Second * 3)
	mock.ExpectExec(regexp.QuoteMeta(restoreDBscript)).WillReturnResult(sqlmock.NewResult(0, 0))

	pgs := &PostgreSQLStorage{
		connection: db,
		Sync:       false,
		context:    ctx,
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
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	mockDB := mock.NewRows([]string{"name", "type", "value", "delta"}).
		AddRow("metric1", "gauge", "231.12", nil).
		AddRow("metric2", "counter", nil, "101")

	mock.ExpectQuery(regexp.QuoteMeta(GetMetricsQuery)).WillReturnRows(mockDB)

	pgs := &PostgreSQLStorage{
		connection: db,
		Sync:       false,
		context:    ctx,
	}

	gauges := map[string]float64{
		"metric1": 231.12,
	}
	counters := map[string]int64{
		"metric2": 101,
	}

	ms, err := pgs.Restore()
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

// func TestPostgreSQLStorage_Save(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	ctx := context.Background()

// 	defer func() {
// 		db.Close()
// 	}()

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(SetMetricQuery)).WithArgs("metric1", "gauge", 231.12, nil) //.WillReturnResult(sqlmock.NewResult(1, 1))
// 	mock.ExpectExec(regexp.QuoteMeta(SetMetricQuery)).WithArgs("metric2", "counter", nil, 101)  //.WillReturnResult(sqlmock.NewResult(1, 1))
// 	mock.ExpectCommit()
// 	mock.ExpectClose()

// 	pgs := &PostgreSQLStorage{
// 		connection: db,
// 		Sync:       false,
// 		context:    ctx,
// 	}

// 	ms := memstorage.MemStorage{
// 		Gauges: map[string]float64{
// 			"metric1": 231.12},
// 		Counters: map[string]int64{
// 			"metric2": 101,
// 		},
// 		Mux: &sync.RWMutex{},
// 	}
// 	ms.Mux.Lock()
// 	if errSave := pgs.Save(&ms); errSave != nil {
// 		t.Errorf("Save memstorage error = %v", err)
// 	}
// 	ms.Mux.Unlock()
// 	// we make sure that all expectations were met
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }

func TestPostgreSQLStorage_SyncMode(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	pgs := &PostgreSQLStorage{
		connection: db,
		Sync:       false,
		context:    ctx,
	}

	assert.Equal(t, pgs.SyncMode(), pgs.Sync)

}

func TestPostgreSQLStorage_Close(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	pgs := &PostgreSQLStorage{
		connection: db,
		Sync:       false,
		context:    ctx,
	}
	pgs.Close()
	assert.Error(t, pgs.connection.Ping())
}
