package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rebus2015/praktikum-devops/internal/config"
	"github.com/rebus2015/praktikum-devops/internal/model"
	"github.com/rebus2015/praktikum-devops/internal/storage"
	"github.com/rebus2015/praktikum-devops/internal/storage/filestorage"
	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

type testSQLdbStorage struct{}

func (db *testSQLdbStorage) Ping(ctx context.Context) error {
	return nil
}
func (db *testSQLdbStorage) Close() {}

func Test_UpdateCounterHandlerFunc(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name        string
		args        want
		request     string
		method      string
		contentType string
	}{
		{
			name: "positive test #1",
			args: want{
				code: 200,
			},
			request:     "/update/counter/cnt/3",
			method:      http.MethodPost,
			contentType: "text/plain",
		},
		{
			name: "negative test #2",
			args: want{
				code: 400,
			},
			request:     "/update/counter/cnt/test",
			method:      http.MethodPost,
			contentType: "text/plain",
		},
		{
			name: "negative test #2",
			args: want{
				code: 404,
			},
			request:     "/update/counter/",
			method:      http.MethodPost,
			contentType: "text/plain",
		},
	}
	cfg := config.Config{StoreFile: "", ConnectionString: ""}
	var metricStorage storage.Repository = storage.NewRepositoryWrapper(*memstorage.NewStorage(), filestorage.NewStorage(&cfg))
	dbStorage := &testSQLdbStorage{}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(&metricStorage, dbStorage, cfg)
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, _ := testRequest(t, ts, tt.method, tt.request)
			// проверяем код ответа
			assert.Equal(t, tt.args.code, statusCode)
		})
	}
}

func Test_UpdateGaugeHandlerFunc(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name        string
		args        want
		request     string
		method      string
		contentType string
	}{
		{
			name: "positive test #1",
			args: want{
				code: 200,
			},
			request:     "/update/gauge/gg/3",
			method:      http.MethodPost,
			contentType: "text/plain",
		},
		{
			name: "negative test #2",
			args: want{
				code: 404,
			},
			request:     "/update/gauge/",
			method:      http.MethodPost,
			contentType: "text/plain",
		},
		{
			name: "negative test #3",
			args: want{
				code: 400,
			},
			request:     "/update/gauge/gg/xx",
			method:      http.MethodPost,
			contentType: "text/plain",
		},
	}

	var metricStorage storage.Repository = storage.NewRepositoryWrapper(*memstorage.NewStorage(), filestorage.NewStorage(&config.Config{}))
	dbStorage := &testSQLdbStorage{}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(&metricStorage, dbStorage, config.Config{})
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, _ := testRequest(t, ts, tt.method, tt.request)
			// проверяем код ответа
			assert.Equal(t, tt.args.code, statusCode)
		})
	}
}

func Test_getAllHandler(t *testing.T) {
	tests := []struct {
		name     string
		counters []memstorage.MetricStr
		gauges   []memstorage.MetricStr
		method   string
		wantcode int
		path     string
	}{
		{
			name:     "Positive test #1",
			counters: []memstorage.MetricStr{{Name: "cnt1", Val: "123"}, {Name: "cnt2", Val: "64"}},
			gauges:   []memstorage.MetricStr{{Name: "gauge1", Val: "12.003"}, {Name: "gauge2", Val: "-164"}},
			method:   http.MethodGet,
			wantcode: http.StatusOK,
			path:     "/",
		},
	}

	for _, tt := range tests {
		var metricStorage storage.Repository = storage.NewRepositoryWrapper(*memstorage.NewStorage(), filestorage.NewStorage(&config.Config{}))
		dbStorage := &testSQLdbStorage{}
		for _, c := range tt.counters {
			_, err := metricStorage.AddCounter(c.Name, c.Val)
			if err != nil {
				log.Printf("Test_GetAllHandler error:%v", err)
			}
		}

		for _, g := range tt.gauges {
			_, err := metricStorage.AddGauge(g.Name, g.Val)
			if err != nil {
				log.Printf("Test_GetAllHandler error:%v", err)
			}
		}

		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(&metricStorage, dbStorage, config.Config{})
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, _ := testRequest(t, ts, tt.method, tt.path)
			// проверяем код ответа
			assert.Equal(t, tt.wantcode, statusCode)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func testRequestJSON(t *testing.T, ts *httptest.Server, method, path string, metric model.Metrics) (int, []byte) {
	data, err := json.Marshal(metric)
	if err != nil {
		log.Panic(err)
	}
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(data))
	assert.NoError(t, err)
	req.Header.Add("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	defer resp.Body.Close()
	return resp.StatusCode, respBody
}

func ptr[T any](v T) *T {
	return &v
}

func Test_UpdateJSONMetricHandlerFunc(t *testing.T) {
	type wantArgs struct {
		code int
		data *model.Metrics
	}
	type requestArgs struct {
		data        *model.Metrics
		path        string
		method      string
		contentType string
	}
	tests := []struct {
		name    string
		want    wantArgs
		request requestArgs
	}{
		{
			name: "positive add gauge test #1",
			want: wantArgs{
				code: 200,
				data: &model.Metrics{ID: "G1", MType: "gauge", Value: ptr(100.47)},
			},
			request: requestArgs{
				data:        &model.Metrics{ID: "G1", MType: "gauge", Value: ptr(100.47)},
				path:        "/update",
				method:      http.MethodPost,
				contentType: "application/json",
			},
		},
		{
			name: "positive add counter test #2",
			want: wantArgs{
				code: 200,
				data: &model.Metrics{ID: "C1", MType: "counter", Delta: ptr(int64(147))},
			},
			request: requestArgs{
				data:        &model.Metrics{ID: "C1", MType: "counter", Delta: ptr(int64(147))},
				path:        "/update",
				method:      http.MethodPost,
				contentType: "application/json",
			},
		},
	}

	var metricStorage storage.Repository = storage.NewRepositoryWrapper(*memstorage.NewStorage(), filestorage.NewStorage(&config.Config{}))
	dbStorage := &testSQLdbStorage{}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(&metricStorage, dbStorage, config.Config{})
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, body := testRequestJSON(t, ts, tt.request.method, tt.request.path, *tt.request.data)
			// проверяем код ответа
			assert.Equal(t, tt.want.code, statusCode)
			var resp model.Metrics
			err := json.Unmarshal(body, &resp)
			assert.NoError(t, err)
			assert.EqualValues(t, *tt.want.data, resp)
		})
	}
}

func Ptr[T any](v T) *T {
	return &v
}
