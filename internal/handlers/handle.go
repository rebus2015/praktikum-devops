package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

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

func NewRouter(metricStorage *storage.Repository) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", GetAllHandler(*metricStorage))

	r.Route("/update", func(r chi.Router) {
		r.With(MetricContextBody).
			Post("/", UpdateJSONMetricHandlerFunc(*metricStorage))
		r.Route("/{mtype}/{name}/{val}", func(r chi.Router) {
			r.Post("/", UpdateMetricHandlerFunc(*metricStorage))
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.With(MetricContextBody).
			Post("/", GetJSONMetricHandlerFunc(*metricStorage))
		r.Route("/{mtype}/{name}", func(r chi.Router) {
			r.Get("/", GetMetricHandlerFunc(*metricStorage))
		})
	})

	return r
}


func UpdateJSONMetricHandlerFunc(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metric, ok := r.Context().Value(metricContextKey{key: "metric"}).(*model.Metrics)
		if !ok {
			log.Printf("Error: [updateJSONMetricHandlerFunc] Metric info not found in context status-'500'")
			http.Error(w, "Metric info not found in context", http.StatusInternalServerError)
			return
		}

		retval := &model.Metrics{
			ID:    metric.ID,
			MType: metric.MType,
		}

		switch metric.MType {
		case "counter":
			{
				if metric.Delta == nil {
					log.Printf("Error: [updateJSONMetricHandlerFunc] Counter not found status- 400")
					http.Error(w, "Counter not found", http.StatusBadRequest)
					return
				}

				delta, err := metricStorage.AddCounter(metric.ID, metric.Delta)
				if err != nil {
					log.Printf("Error: [updateJSONMetricHandlerFunc] Update counter error: %v", err)
					http.Error(w, fmt.Sprintf("Update counter error: %v", err), http.StatusInternalServerError)
					return
				}
				retval.Delta = &delta
			}
		case "gauge":
			{
				if metric.Value == nil {
					log.Printf("Error: [updateJSONMetricHandlerFunc] gauge not found status- 400")
					http.Error(w, "gauge not found", http.StatusBadRequest)
					return
				}

				value, err := metricStorage.AddGauge(metric.ID, metric.Value)
				if err != nil {
					log.Printf("Error: [updateJSONMetricHandlerFunc] Update gauge error: %v", err)
					http.Error(w, fmt.Sprintf("Update counter error: %v", err), http.StatusInternalServerError)
					return
				}
				retval.Value = &value
			}
		default:
			{
				log.Printf("Error: [updateJSONMetricHandlerFunc] Unknown metric type status - 500")
				http.Error(w, "Unknown metric type", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		err := encoder.Encode(retval)
		if err != nil {
			log.Printf("Error: [updateJSONMetricHandlerFunc] Result Json encode error :%v", err)
			http.Error(w, "Result Json encode error", http.StatusInternalServerError)
		}
		log.Printf("Возвращаем UpdateJSON result :%v", retval)
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

func GetJSONMetricHandlerFunc(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		metric, ok := r.Context().Value(metricContextKey{key: "metric"}).(*model.Metrics)
		if !ok {
			log.Printf("Error: [getJSONMetricHandlerFunc] Metric info not found in context")
			http.Error(w, "Metric info not found in context", http.StatusInternalServerError)
			return
		}

		retval := &model.Metrics{
			ID:    metric.ID,
			MType: metric.MType,
		}
		switch metric.MType {
		case "counter":
			delta, err := metricStorage.GetCounter(metric.ID)
			if err != nil {
				log.Printf("Error: [getJSONMetricHandlerFunc] Counter not found: %v", err)
				http.Error(w, "Counter not found", http.StatusNotFound)
				return
			}
			retval.Delta = &delta

		case "gauge":

			value, err := metricStorage.GetGauge(metric.ID)
			if err != nil {
				log.Printf("Error: [getJSONMetricHandlerFunc] Gauge not found: %v", err)
				http.Error(w, "Gauge not found", http.StatusNotFound)
				return
			}
			retval.Value = &value

		default:
			log.Printf("Error: [getJSONMetricHandlerFunc] Unknown metric type")
			http.Error(w, "Unknown metric type", http.StatusInternalServerError)
			return

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		err := encoder.Encode(retval)
		if err != nil {
			log.Printf("Error: [getJSONMetricHandlerFunc] Result Json encode error")
			http.Error(w, "Result Json encode error", http.StatusInternalServerError)
		}
		log.Printf("Возвращаем UpdateJSON result :%v", retval)
	}
}

func MetricContextBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		metric := &model.Metrics{}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		if err := decoder.Decode(metric); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if metric.ID == "" {
			http.Error(w, "metric.ID is empty", http.StatusBadRequest)
			return
		}
		if metric.MType == "" {
			http.Error(w, "metric.MType is empty", http.StatusBadRequest)
			return
		}
		log.Printf("Incoming request Method: %v, Body: %v", r.RequestURI, metric)
		ctx := context.WithValue(r.Context(), metricContextKey{key: "metric"}, metric)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetMetricHandlerFunc(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
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

func GetAllHandler(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
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
