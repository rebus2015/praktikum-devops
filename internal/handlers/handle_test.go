package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/rebus2015/praktikum-devops/internal/config"
	"github.com/rebus2015/praktikum-devops/internal/model"
	"github.com/rebus2015/praktikum-devops/internal/storage"
	"github.com/rebus2015/praktikum-devops/internal/storage/dbstorage"
	"github.com/rebus2015/praktikum-devops/internal/storage/filestorage"
	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_UpdateCounterHandlerFunc(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name        string
		request     string
		method      string
		contentType string
		args        want
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
	var metricStorage storage.Repository = storage.NewRepositoryWrapper(
		memstorage.NewStorage(), nil)
	dbStorage := &sqlStorageMock{}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(metricStorage, dbStorage, cfg)
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
		request     string
		method      string
		contentType string
		args        want
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

	var metricStorage storage.Repository = storage.NewRepositoryWrapper(
		memstorage.NewStorage(), nil)
	dbStorage := &sqlStorageMock{}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(metricStorage, dbStorage, config.Config{})
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, _ := testRequest(t, ts, tt.method, tt.request)
			// проверяем код ответа
			assert.Equal(t, tt.args.code, statusCode)
		})
	}
}

func TestGetMetricHandlerFunc(t *testing.T) {
	type args struct {
		mtype string
		name  string
		val   string
	}
	tests := []struct {
		want     args
		name     string
		method   string
		counters []memstorage.MetricStr
		gauges   []memstorage.MetricStr
		wantcode int
	}{
		{
			name:     "Positive test #1",
			counters: []memstorage.MetricStr{{Name: "cnt1", Val: "123"}, {Name: "cnt2", Val: "64"}},
			gauges:   []memstorage.MetricStr{{Name: "gauge1", Val: "12.003"}, {Name: "gauge2", Val: "-164"}},
			method:   http.MethodGet,
			wantcode: http.StatusOK,
			want: args{
				"counter",
				"cnt1",
				"123",
			},
		},
		{
			name:     "Positive test #1",
			counters: []memstorage.MetricStr{{Name: "cnt1", Val: "123"}, {Name: "cnt2", Val: "64"}},
			gauges:   []memstorage.MetricStr{{Name: "gauge1", Val: "12.003"}, {Name: "gauge2", Val: "-164"}},
			method:   http.MethodGet,
			wantcode: http.StatusNotFound,
			want: args{
				"gauge",
				"cnt1",
				"GetGauge error: cauge with name 'cnt1' is not found",
			},
		},
	}

	for _, tt := range tests {
		var metricStorage storage.Repository = storage.NewRepositoryWrapper(
			memstorage.NewStorage(), nil)
		dbStorage := &sqlStorageMock{}
		for _, c := range tt.counters {
			_, err := metricStorage.AddCounter(c.Name, c.Val)
			if err != nil {
				log.Printf("Test_GetAllHandler error:%v", err)
			}
		}

		for _, g := range tt.gauges {
			_, err := metricStorage.AddGauge(g.Name, g.Val)
			if err != nil {
				log.Printf("TestGetMetricHandlerFunc error:%v", err)
			}
		}

		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(metricStorage, dbStorage, config.Config{})
			ts := httptest.NewServer(r)
			defer ts.Close()
			p, _ := url.JoinPath("/value", tt.want.mtype, tt.want.name)
			statusCode, val := testRequest(t, ts, tt.method, p)
			// проверяем код ответа
			assert.Equal(t, tt.wantcode, statusCode)
			assert.EqualValues(t, val, tt.want.val)
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

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("error: %v", err)
		}
	}()

	return resp.StatusCode, string(respBody)
}

func ptr[T any](v T) *T {
	return &v
}

func testRequestJSONstring(t *testing.T, ts *httptest.Server, method, path string, metric string) (int, []byte) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer([]byte(metric)))
	assert.NoError(t, err)
	req.Header.Add("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("error reponce body close: %v", err)
		}
	}()
	respBody = respBody[:len(respBody)-1]
	return resp.StatusCode, respBody
}

func Test_UpdateJSONMetricHandlerFunc(t *testing.T) {
	type wantArgs struct {
		data    string
		code    int
		wantErr bool
	}
	type requestArgs struct {
		data        string
		path        string
		method      string
		contentType string
	}
	tests := []struct {
		name    string
		request requestArgs
		want    wantArgs
	}{
		{
			name: "positive add gauge test #1",
			want: wantArgs{
				code:    200,
				wantErr: false,
				data:    "{\"id\":\"G1\",\"type\":\"gauge\",\"value\":100.47}",
			},
			request: requestArgs{
				data:        "{\"id\":\"G1\",\"type\":\"gauge\",\"value\":100.47}",
				path:        "/update",
				method:      http.MethodPost,
				contentType: "application/json",
			},
		},
		{
			name: "positive add counter test #2",
			want: wantArgs{
				code:    200,
				wantErr: false,
				data:    "{\"id\":\"C1\",\"type\":\"counter\",\"delta\":47}",
			},
			request: requestArgs{
				data:        "{\"id\":\"C1\",\"type\":\"counter\",\"delta\":47}",
				path:        "/update",
				method:      http.MethodPost,
				contentType: "application/json",
			},
		},
		{
			name: "negative type mismatch test #2",
			want: wantArgs{
				code:    400,
				wantErr: true,
			},
			request: requestArgs{
				data:        "{\"id\":\"C1\",\"type\":\"unk\",\"delta\":1.47}",
				path:        "/update",
				method:      http.MethodPost,
				contentType: "application/json",
			},
		},
		{
			name: "negative add nil counter test #1",
			want: wantArgs{
				code:    400,
				wantErr: true,
			},
			request: requestArgs{
				data:        "{\"id\":\"C1\",\"type\":\"counter\"}",
				path:        "/update",
				method:      http.MethodPost,
				contentType: "application/json",
			},
		},
	}

	var metricStorage storage.Repository = storage.NewRepositoryWrapper(
		memstorage.NewStorage(), nil)
	dbStorage := &sqlStorageMock{}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(metricStorage, dbStorage, config.Config{})
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, body := testRequestJSONstring(t, ts, tt.request.method, tt.request.path, tt.request.data)
			// проверяем код ответа
			assert.Equal(t, tt.want.code, statusCode)
			if tt.want.wantErr {
				return
			}
			assert.EqualValues(t, body, []byte(tt.want.data))
		})
	}
}

func Test_UpdateJSONMultipleMetricHandlerFunc(t *testing.T) {
	type wantArgs struct {
		data    string
		code    int
		wantErr bool
	}
	type requestArgs struct {
		data        string
		path        string
		method      string
		contentType string
	}
	tests := []struct {
		name    string
		request requestArgs
		want    wantArgs
	}{
		{
			name: "positive add gauge test #1",
			want: wantArgs{
				code:    200,
				wantErr: false,
				data:    "[{\"id\":\"G1\",\"type\":\"gauge\",\"value\":100.47},{\"id\":\"C1\",\"type\":\"counter\",\"delta\":47}]",
			},
			request: requestArgs{
				data:        "[{\"id\":\"G1\",\"type\":\"gauge\",\"value\":100.47},{\"id\":\"C1\",\"type\":\"counter\",\"delta\":47}]",
				path:        "/updates",
				method:      http.MethodPost,
				contentType: "application/json",
			},
		},
		{
			name: "negative type mismatch test #2",
			want: wantArgs{
				code:    400,
				wantErr: true,
			},
			request: requestArgs{
				data:        "[{\"id\":\"G1\",\"type\":\"gauge\",\"value\":100.47},{\"id\":\"C1\",\"type\":\"unk\",\"delta\":1.47}]",
				path:        "/updates",
				method:      http.MethodPost,
				contentType: "application/json",
			},
		},
	}

	var metricStorage storage.Repository = storage.NewRepositoryWrapper(
		memstorage.NewStorage(), filestorage.NewStorage(context.Background(), &config.Config{}))
	dbStorage := &sqlStorageMock{}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(metricStorage, dbStorage, config.Config{})
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, body := testRequestJSONstring(t, ts, tt.request.method, tt.request.path, tt.request.data)
			// проверяем код ответа
			assert.Equal(t, tt.want.code, statusCode)
			if tt.want.wantErr {
				return
			}
			assert.EqualValues(t, body, []byte(tt.want.data))
		})
	}
}

func Test_checkMetric(t *testing.T) {
	type args struct {
		metric *model.Metrics
		key    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"test 1 negative",
			args{
				&model.Metrics{
					ID:    "TotalMemory",
					MType: "gauge",
					Delta: nil,
					Value: ptr(float64(7268679680)),
					Hash:  "a3afa5537e12b6a83f982b8a286031b03d26ee12eb43f24c83beacff3ed81f87",
				},
				"SuperSecretKey",
			},
			false,
			true,
		},
		{
			"test 2 positive",
			args{
				&model.Metrics{
					ID:    "TotalMemory",
					MType: "gauge",
					Delta: nil,
					Value: ptr(float64(7268679680)),
					Hash:  "aed2c239d3824aa6a0013258d4c748a786b5de3dfa739919c5b741eb1bd30af9",
				},
				"SuperSecretKey",
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkMetric(tt.args.metric, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}
}

type sqlStorageMock struct {
	isOpened bool
}

func (s *sqlStorageMock) Ping(ctx context.Context) error {
	if s.isOpened {
		return nil
	}
	return fmt.Errorf("Erroe: sqlStorageMock connection state is closed.")
}

func (s *sqlStorageMock) Close() {
	s.isOpened = false
}

func TestGetDBConnState(t *testing.T) {
	type args struct {
		sqlStorage dbstorage.SQLStorage
		path       string
		method     string
	}

	type result struct {
		wantErr bool
		code    int
	}

	tests := []struct {
		name string
		args args
		want result
	}{
		{
			"ping test positive",
			args{
				sqlStorage: &sqlStorageMock{
					isOpened: true,
				},
				path:   "/ping",
				method: http.MethodGet,
			},
			result{
				wantErr: false,
				code:    200,
			},
		},
		{
			"ping test negative",
			args{
				sqlStorage: &sqlStorageMock{
					isOpened: false,
				},
				path:   "/ping",
				method: http.MethodGet,
			},
			result{
				wantErr: false,
				code:    500,
			},
		},
	}
	for _, tt := range tests {
		var metricStorage storage.Repository = storage.NewRepositoryWrapper(
			memstorage.NewStorage(), filestorage.NewStorage(context.Background(), &config.Config{}))

		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(metricStorage, tt.args.sqlStorage, config.Config{})
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, _ := testRequest(t, ts, tt.args.method, tt.args.path)
			// проверяем код ответа
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func Test_getAllHandler(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		counters []memstorage.MetricStr
		gauges   []memstorage.MetricStr
		wantcode int
	}{
		{
			name:     "Positive test #1",
			counters: []memstorage.MetricStr{{Name: "cnt1", Val: "123"}, {Name: "cnt2", Val: "64"}},
			gauges:   []memstorage.MetricStr{{Name: "gauge1", Val: "12.003"}, {Name: "gauge2", Val: "-164"}},
			method:   http.MethodGet,
			wantcode: http.StatusOK,
			path:     "/",
		},
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
		var metricStorage storage.Repository = storage.NewRepositoryWrapper(
			memstorage.NewStorage(), filestorage.NewStorage(context.Background(), &config.Config{}))
		dbStorage := &sqlStorageMock{}
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
			r := NewRouter(metricStorage, dbStorage, config.Config{})
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, _ := testRequest(t, ts, tt.method, tt.path)
			// проверяем код ответа
			assert.Equal(t, tt.wantcode, statusCode)
		})
	}
}

func TestUpdateMetricHandlerFunc(t *testing.T) {
	type args struct {
		mtype string
		name  string
		val   string
	}
	tests := []struct {
		name     string
		method   string
		path     string
		want     args
		wantcode int
	}{
		{
			name:     "Positive test #1",
			method:   http.MethodPost,
			wantcode: http.StatusOK,
			path:     "/update",
			want: args{
				"counter",
				"cnt1",
				"5055",
			},
		},
		{
			name:     "Positive test #2",
			method:   http.MethodPost,
			wantcode: http.StatusOK,
			want: args{
				"gauge",
				"gauge1",
				"0.0012",
			},
		},
		{
			name:     "Negative test #1 value type mismatch",
			method:   http.MethodPost,
			wantcode: http.StatusBadRequest,
			want: args{
				"counter",
				"gauge1",
				"0.0012",
			},
		},
		{
			name:     "Negative test #1 unk metric type",
			method:   http.MethodPost,
			wantcode: http.StatusNotImplemented,
			want: args{
				"unknown",
				"metricXX",
				"0.0012",
			},
		},
	}

	for _, tt := range tests {
		var metricStorage storage.Repository = storage.NewRepositoryWrapper(
			memstorage.NewStorage(), nil)
		dbStorage := &sqlStorageMock{}

		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(metricStorage, dbStorage, config.Config{})
			ts := httptest.NewServer(r)
			defer ts.Close()
			p, _ := url.JoinPath("/update", tt.want.mtype, tt.want.name, tt.want.val)
			statusCode, _ := testRequest(t, ts, tt.method, p)
			// проверяем код ответа
			assert.Equal(t, tt.wantcode, statusCode)
		})
	}
}
func TestGzipMiddleware(t *testing.T) {
	// Create a mock handler
	expectedBody := []byte("test body with respect to unzip process")
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the request context has the expected value
		buf, _ := r.Context().Value(bodyContextKey{}).([]byte)
		assert.Equal(t, buf, expectedBody)
	})

	var gzipbody bytes.Buffer
	gz := gzip.NewWriter(&gzipbody)
	if _, err := gz.Write(expectedBody); err != nil {
		log.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		log.Fatal(err)
	}

	compressedRequest, _ := http.NewRequest(http.MethodPost, "/", &gzipbody)
	compressedRequest.Header.Set("Content-Encoding", "gzip")

	// Use the gzipMiddleware with the mock handler
	handler := gzipMiddleware(mockHandler)

	// Perform the test request
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, compressedRequest)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
	}
}
