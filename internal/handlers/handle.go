package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rebus2015/praktikum-devops/internal/storage"
)

var MemStats storage.MemStorage

// HelloWorld — обработчик запроса.
func UpdateCounterHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("content-type") != "text/plain" {
		http.Error(w, "Invalid Context-type!Please use text/plain", http.StatusUnsupportedMediaType)
		return
	}

	data := strings.Split(r.URL.Path, "/")
	if len(data) < 5 {
		http.Error(w, fmt.Sprintf("counter handler panic: url path wrong format %v", data), http.StatusNotFound)
		return
	}

	c := storage.CMetric{}
	err := c.TryParse(data[3], data[4])
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	MemStats.AddCounter(c)
	// устанавливаем статус-код 200
	w.WriteHeader(http.StatusOK)
}

func UpdateGaugeHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("content-type") != "text/plain" {
		http.Error(w, "Invalid Context-type!Please use text/plain", http.StatusUnsupportedMediaType)
		return
	}

	var data = strings.Split(r.URL.Path, "/")
	if len(data) < 5 {
		http.Error(w, fmt.Sprintf("gauge handler panic: url path wrong format %v", data), http.StatusNotFound)
		return
	}

	g := storage.GMetric{}
	err := g.TryParse(data[3], data[4])
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	MemStats.AddGauge(g)
	// устанавливаем статус-код 200
	w.WriteHeader(http.StatusOK)
}

func ErrorHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	w.Write([]byte("Wrong url path, metric type not found!"))
	http.Error(w, "Path not found", http.StatusNotImplemented)
}
