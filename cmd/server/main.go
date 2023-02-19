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
	handlers.MemStats = &s
	r := handlers.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}
