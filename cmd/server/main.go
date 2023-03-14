package main

//import "../../internal/storage"
import (
	"log"
	"net/http"

	"github.com/rebus2015/praktikum-devops/internal/handlers"
	"github.com/rebus2015/praktikum-devops/internal/storage"
)

func main() {
	//cfg := getConfig()
	storage := storage.CreateRepository()
	r := handlers.NewRouter(storage)
	log.Fatal(http.ListenAndServe(":8080", r))
}
