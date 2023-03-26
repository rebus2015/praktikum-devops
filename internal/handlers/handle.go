package handlers

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/rebus2015/praktikum-devops/internal/config"
	"github.com/rebus2015/praktikum-devops/internal/model"
	"github.com/rebus2015/praktikum-devops/internal/signer"
	"github.com/rebus2015/praktikum-devops/internal/storage"
	"github.com/rebus2015/praktikum-devops/internal/storage/dbstorage"
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

type metricContextKey struct{}

var contentTypes = []string{
	"application/javascript",
	"application/json",
	"text/css",
	"text/html",
	"text/html; charset=utf-8",
	"text/plain",
	"text/xml",
}

const (
	counter string = "counter"
	gauge   string = "gauge"
)

func NewRouter(
	metricStorage *storage.Repository,
	postgreStorage dbstorage.SQLStorage,
	cfg config.Config,
) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(gzip.BestSpeed, contentTypes...))

	r.Get("/", GetAllHandler(*metricStorage))
	r.Get("/ping", GetDBConnState(postgreStorage))
	r.Route("/update", func(r chi.Router) {
		r.With(MiddlewareGeneratorJSON(cfg.Key)).
			Post("/", UpdateJSONMetricHandlerFunc(*metricStorage, cfg.Key))
		r.Route("/{mtype}/{name}/{val}", func(r chi.Router) {
			r.Post("/", UpdateMetricHandlerFunc(*metricStorage))
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.With(MiddlewareGeneratorJSON(cfg.Key)).
			Post("/", GetJSONMetricHandlerFunc(*metricStorage, cfg.Key))
		r.Route("/{mtype}/{name}", func(r chi.Router) {
			r.Get("/", GetMetricHandlerFunc(*metricStorage))
		})
	})

	return r
}

func GetDBConnState(
	sqlStorage dbstorage.SQLStorage,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// При успешной проверке хендлер должен вернуть HTTP-статус 200 OK, при неуспешной — 500 Internal Server Error.
		if err := sqlStorage.Ping(r.Context()); err != nil {
			log.Printf("Cannot ping database because %s", err)
			http.Error(
				w,
				fmt.Sprintf("Cannot ping database %v", err),
				http.StatusInternalServerError,
			)
			return
		}
		// устанавливаем статус-код 200
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateJSONMetricHandlerFunc(
	metricStorage storage.Repository,
	key string,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metric, ok := r.Context().Value(metricContextKey{}).(*model.Metrics)
		if !ok {
			log.Printf(
				"Error: [updateJSONMetricHandlerFunc] Metric info not found in context status-'500'",
			)
			http.Error(w, "Metric info not found in context", http.StatusInternalServerError)
			return
		}

		retval := &model.Metrics{
			ID:    metric.ID,
			MType: metric.MType,
		}

		switch metric.MType {
		case counter:
			{
				if metric.Delta == nil {
					log.Printf("Error: [updateJSONMetricHandlerFunc] Counter not found status- 400")
					http.Error(w, "Counter not found", http.StatusBadRequest)
					return
				}

				delta, err := metricStorage.IncCounter(metric.ID, metric.Delta)
				if err != nil {
					log.Printf("Error: [updateJSONMetricHandlerFunc] Update counter error: %v", err)
					http.Error(
						w,
						fmt.Sprintf("Update counter error: %v", err),
						http.StatusInternalServerError,
					)
					return
				}
				retval.Delta = &delta
			}
		case gauge:
			{
				if metric.Value == nil {
					log.Printf("Error: [updateJSONMetricHandlerFunc] gauge not found status- 400")
					http.Error(w, "gauge not found", http.StatusBadRequest)
					return
				}

				value, err := metricStorage.SetGauge(metric.ID, metric.Value)
				if err != nil {
					log.Printf("Error: [updateJSONMetricHandlerFunc] Update gauge error: %v", err)
					http.Error(
						w,
						fmt.Sprintf("Update counter error: %v", err),
						http.StatusInternalServerError,
					)
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

		if key != "" {
			hashObject := signer.NewHashObject(key)
			err := hashObject.Sign(retval)
			if err != nil {
				log.Printf(
					"Error: [updateJSONMetricHandlerFunc] Result Json Sign data error :%v",
					err,
				)
				http.Error(w, "Result Json Sign error", http.StatusInternalServerError)
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

func UpdateMetricHandlerFunc(
	metricStorage storage.Repository,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mtype := chi.URLParam(r, "mtype")
		name := chi.URLParam(r, "name")
		val := chi.URLParam(r, "val")
		var err error
		switch mtype {
		case gauge:
			_, err = metricStorage.SetGauge(name, val)
		case counter:
			_, err = metricStorage.IncCounter(name, val)
		default:
			{
				http.Error(w, "Unknown metric Type", http.StatusNotImplemented)
				return
			}
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // 400
		}
		// устанавливаем статус-код 200
		w.WriteHeader(http.StatusOK)
	}
}

func GetJSONMetricHandlerFunc(
	metricStorage storage.Repository,
	key string,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metric, ok := r.Context().Value(metricContextKey{}).(*model.Metrics)
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
		case counter:
			delta, err := metricStorage.GetCounter(metric.ID)
			if err != nil {
				log.Printf("Error: [getJSONMetricHandlerFunc] Counter not found: %v", err)
				http.Error(w, "Counter not found", http.StatusNotFound)
				return
			}
			retval.Delta = &delta

		case gauge:

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

		if key != "" {
			hashObject := signer.NewHashObject(key)
			err := hashObject.Sign(retval)
			if err != nil {
				log.Printf(
					"Error: [updateJSONMetricHandlerFunc] Result Json Sign data error :%v",
					err,
				)
				http.Error(w, "Result Json Sign error", http.StatusInternalServerError)
			}
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

func MiddlewareGeneratorJSON(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var reader io.Reader
			if r.Header.Get(`Content-Encoding`) == `gzip` {
				gz, err := gzip.NewReader(r.Body)
				if err != nil {
					log.Printf("Failed to create gzip reader: %v", err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				reader = gz
				defer gz.Close()
			} else {
				reader = r.Body
			}

			metric := &model.Metrics{}
			decoder := json.NewDecoder(reader)
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

			if key != "" && metric.Hash != "" {
				hashObject := signer.NewHashObject(key)
				passed, err := hashObject.Verify(metric)
				if err != nil {
					log.Printf(
						"Incoming Metric could not pass signature verification: %v, \nBody: %v, \n error: %v",
						r.RequestURI,
						metric,
						err,
					)
					http.Error(
						w,
						fmt.Sprintf(
							"Incoming Metric could not pass signature verification, error:%v",
							err,
						),
						http.StatusBadRequest,
					)
				}
				if !passed {
					log.Printf(
						"Error: Incoming Metric could not pass signature verification: %v, Body: %v",
						r.RequestURI,
						metric,
					)
					http.Error(
						w,
						"Incoming Metric could not pass signature verification",
						http.StatusBadRequest,
					)
					return
				}
			}
			ctx := context.WithValue(r.Context(), metricContextKey{}, metric)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetMetricHandlerFunc(
	metricStorage storage.Repository,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		mtype := chi.URLParam(r, "mtype")

		var val string

		switch mtype {
		case gauge:

			g, err := metricStorage.GetGauge(name)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				_, err = w.Write([]byte(err.Error()))
				if err != nil {
					log.Printf("GetMetricHandlerFunc gauge writer.Write error:%v", err)
				}
				return
			}
			val = fmt.Sprintf("%.3f", g)

		case counter:

			c, err := metricStorage.GetCounter(name)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				_, err = w.Write([]byte(err.Error()))
				if err != nil {
					log.Printf("GetMetricHandlerFunc counter writer.Write error:%v", err)
				}
				return
			}
			val = fmt.Sprintf("%v", c)

		default:

			w.WriteHeader(http.StatusNotImplemented)
			_, err := w.Write([]byte("Unknown metric type"))
			if err != nil {
				log.Printf("GetMetricHandlerFunc unknown type writer.Write error:%v", err)
			}
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(val))
		if err != nil {
			log.Printf("GetMetricHandlerFunc metric success writer.Write error:%v", err)
		}
	}
}

func GetAllHandler(metricStorage storage.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		metrics, err := metricStorage.GetView()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
			return
		}
		template, err := template.New("metrics").Parse(templ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // 400
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		err = template.ExecuteTemplate(w, "metrics", metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
			return
		}
	}
}
