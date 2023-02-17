package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			request.Header.Add("content-type", tt.contentType)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(UpdateCounterHandlerFunc)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, tt.args.code, res.StatusCode)
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
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			request.Header.Add("content-type", tt.contentType)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(UpdateGaugeHandlerFunc)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, tt.args.code, res.StatusCode)
		})
	}
}
