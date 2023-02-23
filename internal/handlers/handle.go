package handlers

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/rebus2015/praktikum-devops/internal/storage"
)

//var MemStats storage.Repository

func CreateRepository() storage.Repository {
	var MemStats = new(storage.MemStorage)
	MemStats.Init()
	return MemStats
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

func NewRouter(MetricStorage storage.Repository) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	//r.Use(middleware.AllowContentType("plain/text"))
	r.Use(middleware.Recoverer)

	r.Get("/", getAllHandler(MetricStorage))

	r.Route("/update", func(r chi.Router) {
		r.Route("/{mtype}/{name}/{val}", func(r chi.Router) {
			r.Post("/", UpdateMetricHandlerFunc(MetricStorage))
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.Route("/{mtype}/{name}", func(r chi.Router) {
			r.Get("/", getMetricHandlerFunc(MetricStorage))
		})
	})

	return r
}

func UpdateMetricHandlerFunc(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mtype := chi.URLParam(r, "mtype")
		name := chi.URLParam(r, "name")
		val := chi.URLParam(r, "val")
		var err error
		switch mtype {
		case "gauge":
			err = metricStorage.AddGauge(name, val)
		case "counter":
			err = metricStorage.AddCounter(name, val)
		default:
			{
				http.Error(w, "Unknown metric Type", http.StatusNotImplemented)
				return
			}
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) //400
		}
		// устанавливаем статус-код 200
		w.WriteHeader(http.StatusOK)

	}
}

func getMetricHandlerFunc(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		mtype := chi.URLParam(r, "mtype")

		var val string
		var err error

		switch mtype {
		case "gauge":
			{
				val, err = metricStorage.GetGauge(name)
			}
		case "counter":
			{
				val, err = metricStorage.GetCounter(name)
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
}

func getAllHandler(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {

		metrics, err := metricStorage.GetView()
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
}
