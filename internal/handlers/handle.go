package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/rebus2015/praktikum-devops/internal/model"
	"github.com/rebus2015/praktikum-devops/internal/storage"
)

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

func NewRouter(metricStorage storage.Repository) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	//r.Use(middleware.AllowContentType("plain/text"))
	r.Use(middleware.Recoverer)

	r.Get("/", getAllHandler(metricStorage))

	r.Route("/update", func(r chi.Router) {
		r.Post("/", UpdateJsonMetricHandlerFunc(metricStorage))
		r.Route("/{mtype}/{name}/{val}", func(r chi.Router) {
			r.Post("/", UpdateMetricHandlerFunc(metricStorage))
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/", getJsonMetricHandlerFunc(metricStorage))
		r.Route("/{mtype}/{name}", func(r chi.Router) {
			r.Get("/", getMetricHandlerFunc(metricStorage))
		})
	})

	return r
}

func UpdateJsonMetricHandlerFunc(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			http.Error(w, "not valid content-type", http.StatusUnsupportedMediaType)
		}
		var data model.Metrics
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var err error
		var metric model.Metrics
		mtype := data.MType
		name := data.ID

		switch mtype {
		case "gauge":
			metric, err = metricStorage.AddGauge(name, data.Value)
		case "counter":
			metric, err = metricStorage.AddCounter(name, data.Delta)
		default:
			{
				http.Error(w, "Unknown metric Type", http.StatusNotImplemented)
				return
			}
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) //400
		}
		w.Header().Add("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&metric); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Println("New JSON Post message came!")
	}
}
func UpdateMetricHandlerFunc(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mtype := chi.URLParam(r, "mtype")
		name := chi.URLParam(r, "name")
		val := chi.URLParam(r, "val")
		var err error
		switch mtype {
		case "gauge":
			_, err = metricStorage.AddGauge(name, val)
		case "counter":
			_, err = metricStorage.AddCounter(name, val)
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

func getJsonMetricHandlerFunc(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			http.Error(w, "not valid content-type", http.StatusUnsupportedMediaType)
		}

		var metric model.Metrics

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := metricStorage.FillMetric(&metric)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) //400
		}

		w.Header().Add("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(&metric); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Println("New JSON Get message came!")
	}
}
func getMetricHandlerFunc(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		mtype := chi.URLParam(r, "mtype")

		var val string

		switch mtype {
		case "gauge":
			{

				g, err := metricStorage.GetGauge(name)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(err.Error()))
					return
				}
				val = fmt.Sprintf("%v", g)
			}
		case "counter":
			{
				c, err := metricStorage.GetCounter(name)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(err.Error()))
					return
				}
				val = fmt.Sprintf("%v", c)
			}
		default:
			{
				w.WriteHeader(http.StatusNotImplemented)
				w.Write([]byte("Unknown metric type"))
			}
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
