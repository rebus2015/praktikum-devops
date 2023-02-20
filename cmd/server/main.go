package main

//import "../../internal/storage"
import (
	"log"
	"net/http"

	"github.com/rebus2015/praktikum-devops/internal/handlers"
)

func main() {
	handlers.CreateRepository()
	r := handlers.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}
