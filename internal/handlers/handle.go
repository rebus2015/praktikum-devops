package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
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

type metricContextKey struct {
	key string
}

func NewRouter(metricStorage storage.Repository) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", getAllHandler(metricStorage))

	r.Route("/update", func(r chi.Router) {
		r.With(metricContextBody).
			Post("/", updateJSONMetricHandlerFunc(metricStorage))
		r.Route("/{mtype}/{name}/{val}", func(r chi.Router) {
			r.Post("/", updateMetricHandlerFunc(metricStorage))
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.Use(metricContextBody)
		r.Post("/", getJSONMetricHandlerFunc(metricStorage))
		r.Route("/{mtype}/{name}", func(r chi.Router) {
			r.Get("/", getMetricHandlerFunc(metricStorage))
		})
	})

	return r
}

func updateJSONMetricHandlerFunc(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		data := r.Context().Value(metricContextKey{"metric"}).(*model.Metrics)
		mtype := data.MType
		name := data.ID
		//log.Printf("updateJSONMetricHandlerFunc: %v read from context", data)
		var err error
		var metric model.Metrics

		switch mtype {
		case "gauge":
			{
				if data.Value == nil {
					http.Error(w, "Not valid metric Value", http.StatusBadRequest)
					return
				}
				metric, err = metricStorage.AddGauge(name, data.Value)
			}
		case "counter":
			{
				if data.Delta == nil {
					http.Error(w, "Not valid metric Value", http.StatusBadRequest)
					return
				}
				metric, err = metricStorage.AddCounter(name, data.Delta)
			}
		default:
			{
				//log.Printf("updateJSONMetricHandlerFunc: exited with status %v", http.StatusNotImplemented)
				http.Error(w, "Unknown metric Type", http.StatusNotImplemented)
				return
			}
		}

		//log.Printf("updateJSONMetricHandlerFunc: update metric error: %v", err)
		switch {
		case err == io.EOF:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case err != nil:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&metric); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		//fmt.Println("New JSON Post message came!")
	}
}

func updateMetricHandlerFunc(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
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

func getJSONMetricHandlerFunc(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metric := r.Context().Value(metricContextKey{key: "metric"}).(*model.Metrics)
		//log.Printf("getJSONMetricHandlerFunc: %v read from context", metric)
		err := metricStorage.FillMetric(metric)
		if err != nil {
			//log.Printf(" error %v metric: %v getJSONMetricHandlerFunc", err, metric)
			log.Printf("GetJSONMetric find metric error: %v", err.Error())
			JSONError(w, err, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&metric); err != nil {
			//log.Printf("Encoder exited with error: %v metric: %v", err, metric)
			//http.Error(w, err.Error(), http.StatusBadRequest)
			JSONError(w, err, http.StatusBadRequest)
			log.Printf("GetJSONMetric encoder error: %v", err.Error())
			return
		}

		//fmt.Println("New JSON Get message came!")
	}
}

func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

func metricContextBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			http.Error(w, "not valid content-type", http.StatusBadRequest)
		}
		metric := &model.Metrics{}
		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			switch {
			case err == io.EOF:
				http.Error(w, err.Error(), http.StatusNotFound)
			case err != nil:
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}
		if metric.ID == "" {
			http.Error(w, "metric.ID is empty", http.StatusBadRequest)
		}
		if metric.MType == "" {
			http.Error(w, "metric.MType is empty", http.StatusBadRequest)
		}
		ctx := context.WithValue(r.Context(), metricContextKey{key: "metric"}, metric)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
