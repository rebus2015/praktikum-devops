package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rebus2015/praktikum-devops/internal/storage"
	"github.com/stretchr/testify/assert"
)

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
	MemStats = new(storage.MemStorage)
	MemStats.Init()
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter()
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

	MemStats = new(storage.MemStorage)
	MemStats.Init()

	for _, tt := range tests {

		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter()
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
		counters []storage.MetricStr
		gauges   []storage.MetricStr
		method   string
		wantcode int
		path     string
	}{
		{
			name:     "Positive test #1",
			counters: []storage.MetricStr{{Name: "cnt1", Val: "123"}, {Name: "cnt2", Val: "64"}},
			gauges:   []storage.MetricStr{{Name: "gauge1", Val: "12.003"}, {Name: "gauge2", Val: "-164"}},
			method:   http.MethodGet,
			wantcode: http.StatusOK,
			path: "/",
		},
	}

	for _, tt := range tests {

		MemStats = new(storage.MemStorage)
		MemStats.Init()
		for _, c := range tt.counters {
			MemStats.AddCounter(c.Name, c.Val)
		}

		for _, g := range tt.gauges {
			MemStats.AddGauge(g.Name, g.Val)
		}

		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter()
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
