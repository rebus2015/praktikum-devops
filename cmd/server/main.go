package main

//import "../../internal/storage"
import (
	"log"
	"net/http"

	"github.com/rebus2015/praktikum-devops/internal/handlers"
	"github.com/rebus2015/praktikum-devops/internal/storage"
)

func main() {
	s := storage.MemStorage{}
	s.Init()
	handlers.MemStats = s
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handlers.ErrorHandleFunc)
	mux.HandleFunc("/update/counter", handlers.UpdateCounterHandlerFunc)
	mux.HandleFunc("/update/gauge", handlers.UpdateGaugeHandlerFunc)
	// конструируем свой сервер
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
