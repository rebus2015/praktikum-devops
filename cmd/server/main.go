package main

//import "../../internal/storage"
import (
	"log"
	"net/http"

	"github.com/rebus2015/praktikum-devops/internal/handlers"
)

func main() {
	storage:= handlers.CreateRepository()
	r := handlers.NewRouter(storage)
	log.Fatal(http.ListenAndServe(":8080", r))
}
