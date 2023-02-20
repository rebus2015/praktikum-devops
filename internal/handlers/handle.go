package handlers

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rebus2015/praktikum-devops/internal/storage"
)

var MemStats storage.Repository

func CreateRepository() {
	MemStats = new(storage.MemStorage)
	MemStats.Init()
}

const templ = `{{define "metrics"}}
<!doctype html>
<html lang="ru" class="h-100">
<head>
    <title>Title</title>
</head>
<body class="d-flex flex-column h-100">
    <table class="table">
        <th>name</th>
        <th>value</th>
        {{range .}}
        <tr>
            <td>{{.Name}}</td>
            <td>{{.Val}}</td>
        </tr>
        {{end}}
    </table>
</body>
</html>
{{end}}`

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	//r.Use(middleware.AllowContentType("plain/text"))
	r.Use(middleware.Recoverer)

	r.Get("/", getAllHandler)

	r.Route("/update", func(r chi.Router) {
		r.Route("/counter/{name}/{val}", func(r chi.Router) {
			r.Post("/", UpdateCounterHandlerFunc)
		})
		r.Route("/gauge/{name}/{val}", func(r chi.Router) {
			r.Post("/", UpdateGaugeHandlerFunc)

		})
	})

	r.Route("/value", func(r chi.Router) {
		r.Route("/{mtype}/{name}", func(r chi.Router) {
			r.Get("/", getMetricHandlerFunc)
		})
	})

	return r
}

// HelloWorld — обработчик запроса.
func UpdateCounterHandlerFunc(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	val := chi.URLParam(r, "val")
	err := MemStats.AddCounter(name, val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //400
	}
	// устанавливаем статус-код 200
	w.WriteHeader(http.StatusOK)
}

func UpdateGaugeHandlerFunc(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	val := chi.URLParam(r, "val")
	err := MemStats.AddGauge(name, val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //400
	}
	// устанавливаем статус-код 200
	w.WriteHeader(http.StatusOK)
}

func getMetricHandlerFunc(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	mtype := chi.URLParam(r, "mtype")

	var val string
	var err error

	switch mtype {
	case "gauge":
		{
			val, err = MemStats.GetGauge(name)
		}
	case "counter":
		{
			val, err = MemStats.GetCounter(name)
		}
	default:
		{
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("Unknown metric type"))
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(val))
}

func getAllHandler(w http.ResponseWriter, r *http.Request) {
	metrics, err := MemStats.GetView()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //500
		return
	}
	template, err := template.New("metrics").Parse(templ)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	err = template.ExecuteTemplate(w, "metrics", metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //500
		return
	}
}
