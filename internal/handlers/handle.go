// Package handlers создает экземпляр роутера и описывает все доступные эндпоинты
// содержит реализацию необходимого middleware.
package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/rebus2015/praktikum-devops/internal/config"
	"github.com/rebus2015/praktikum-devops/internal/model"
	"github.com/rebus2015/praktikum-devops/internal/signer"
	"github.com/rebus2015/praktikum-devops/internal/storage"
	"github.com/rebus2015/praktikum-devops/internal/storage/dbstorage"
)

// Шаблон html-страницы для вывода всех имеющихся метрик.
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

type multipleMetricsContextKey struct{}

type singleMetricContextKey struct{}

type bodyContextKey struct{}

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
	counter                    string = "counter"
	gauge                      string = "gauge"
	compressed                 string = `gzip`
	keyCT                      string = "Content-Type"
	keyValueJSON               string = "application/json"
	missingContextMessage      string = "Metric info not found in context"
	missingContextMessageLong  string = "Error: [updateJSONMetricHandlerFunc] Metric info not found in context ('500')"
	unkMTMessage               string = "Unknown metric type"
	resJSONSignErrorMessage    string = "Error: [updateJSONMetricHandlerFunc] Result Json Sign data error :%v"
	httpJSONSignErrorMessage   string = "Result Json Sign error"
	retUpdateJSONResultMessage string = "Возвращаем UpdateJSON result :%v"
)

// NewRouter инициализация роутера с помощью библиотеки chi и описание доступных эндпоинтов.
func NewRouter(
	metricStorage storage.Repository,
	postgreStorage dbstorage.SQLStorage,
	cfg config.Config,
) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(gzip.BestSpeed, contentTypes...))
	r.Mount("/debug", middleware.Profiler())
	r.Get("/", GetAllHandler(metricStorage))

	r.Get("/ping", GetDBConnState(postgreStorage))

	r.Route("/update", func(r chi.Router) {
		r.With(gzipMiddleware).
			With(rsaMiddleware(cfg.CryptoKey)).
			With(MiddlewareGeneratorSingleJSON(cfg.Key)).
			Post("/", UpdateJSONMetricHandlerFunc(metricStorage, cfg.Key))
		r.Route("/{mtype}/{name}/{val}", func(r chi.Router) {
			r.Post("/", UpdateMetricHandlerFunc(metricStorage))
		})
	})

	r.Route("/updates", func(r chi.Router) {
		r.With(gzipMiddleware).
			With(rsaMiddleware(cfg.CryptoKey)).
			With(MiddlewareGeneratorMultipleJSON(cfg.Key)).
			Post("/", UpdateJSONMultipleMetricHandlerFunc(metricStorage, cfg.Key))
	})

	r.Route("/value", func(r chi.Router) {
		r.With(gzipMiddleware).
			With(rsaMiddleware(cfg.CryptoKey)).
			With(MiddlewareGeneratorSingleJSON(cfg.Key)).
			Post("/", GetJSONMetricHandlerFunc(metricStorage, cfg.Key))
		r.Route("/{mtype}/{name}", func(r chi.Router) {
			r.Get("/", GetMetricHandlerFunc(metricStorage))
		})
	})
	return r
}

// GetDBConnState реализует пинг состояния БД.
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

// UpdateJSONMultipleMetricHandlerFunc обрабатывает обновления значений метрик, которыые приходят в виде массивов JSON.
func UpdateJSONMultipleMetricHandlerFunc(
	metricStorage storage.Repository,
	key string,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics, ok := r.Context().Value(multipleMetricsContextKey{}).([]*model.Metrics)
		if !ok {
			log.Printf(
				"Error: [UpdateJSONMultipleMetricHandlerFunc] Metric info not found in context status-'500'",
			)
			http.Error(w, missingContextMessage, http.StatusInternalServerError)
			return
		}
		err := metricStorage.AddMetrics(metrics)
		if err != nil {
			log.Printf("Error: [UpdateJSONMultipleMetricHandlerFunc] Add multiple metrics error: %v", err)
			http.Error(
				w,
				fmt.Sprintf("Add multiple metrics error: %v", err),
				http.StatusInternalServerError,
			)
			return
		}
		l := len(metrics)
		retval := make([]model.Metrics, l)
		for i, m := range metrics {
			if key != "" {
				hashObject := signer.NewHashObject(key)
				sssignErr := hashObject.Sign(metrics[i])
				if sssignErr != nil {
					log.Printf(
						resJSONSignErrorMessage,
						err,
					)
					http.Error(w, httpJSONSignErrorMessage, http.StatusInternalServerError)
				}
			}
			retval[i] = model.Metrics{
				ID:    m.ID,
				MType: m.MType,
				Value: m.Value,
				Delta: m.Delta,
				Hash:  m.Hash,
			}
		}

		w.Header().Set(keyCT, keyValueJSON)
		w.WriteHeader(http.StatusOK)
		log.Printf("Try to encode :%v to metric[]", retval)

		encoder := json.NewEncoder(w)
		err = encoder.Encode(retval)
		if err != nil {
			log.Printf("Error: [updateJSONMetricHandlerFunc] Result Json encode error :%v", err)
			http.Error(w, "sjosn encode error", http.StatusInternalServerError)
		}
		log.Printf(retUpdateJSONResultMessage, retval)
	}
}

// UpdateJSONMetricHandlerFunc обрабатывает обновления одиночных значений метрик в формате JSON.
func UpdateJSONMetricHandlerFunc(
	metricStorage storage.Repository,
	key string,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metric, ok := r.Context().Value(singleMetricContextKey{}).(*model.Metrics)
		if !ok {
			log.Printf(
				missingContextMessageLong,
			)
			http.Error(w, missingContextMessage, http.StatusInternalServerError)
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
					http.Error(w, "error: counter not found", http.StatusBadRequest)
					return
				}

				delta, err := metricStorage.AddCounter(metric.ID, metric.Delta)
				if err != nil {
					log.Printf("Error: [updateJSONMetricHandlerFunc] Update counter error: %v", err)
					http.Error(
						w,
						fmt.Sprintf("error: update counter failed: %v", err),
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
					http.Error(w, "error: gauge not found", http.StatusBadRequest)
					return
				}

				value, err := metricStorage.AddGauge(metric.ID, metric.Value)
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
				http.Error(w, unkMTMessage, http.StatusInternalServerError)
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

		w.Header().Set(keyCT, keyValueJSON)
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

// UpdateMetricHandlerFunc обрабатывает обновления одиночных значений метрик из тела URL.
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
			_, err = metricStorage.AddGauge(name, val)
		case counter:
			_, err = metricStorage.AddCounter(name, val)
		default:
			{
				http.Error(w, unkMTMessage, http.StatusNotImplemented)
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

// GetJSONMetricHandlerFunc возвращает значение метрики по ключу в формате JSON.
func GetJSONMetricHandlerFunc(
	metricStorage storage.Repository,
	key string,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metric, ok := r.Context().Value(singleMetricContextKey{}).(*model.Metrics)
		if !ok {
			log.Printf("Error: [getJSONMetricHandlerFunc] Metric info not found in context")
			http.Error(w, missingContextMessage, http.StatusInternalServerError)
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
			http.Error(w, unkMTMessage, http.StatusInternalServerError)
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

		w.Header().Set(keyCT, keyValueJSON)
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

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		var buf []byte
		var err error
		if r.Header.Get(`Content-Encoding`) == compressed {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Printf("Failed to create gzip reader: %v", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			buf, err = io.ReadAll(gz)
			if err != nil {
				log.Printf("Failed to read bytes from gzip reader: %v", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func() {
				if err := gz.Close(); err != nil {
					log.Errorf("error occured when closing gzip: %v", err)
				}
			}()
		} else {
			buf, err = io.ReadAll(r.Body)

			if err != nil {
				log.Printf("Failed to read bytes from request body in gzip encoder: %v", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		ctx := context.WithValue(r.Context(), bodyContextKey{}, buf)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func rsaMiddleware(key *rsa.PrivateKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			content, ok := r.Context().Value(bodyContextKey{}).([]byte)
			if !ok {
				log.Printf(
					missingContextMessageLong,
				)
				http.Error(w, missingContextMessage, http.StatusInternalServerError)
				return
			}
			var nextData []byte
			var err error
			if key != nil {
				nextData, err = signer.DecryptMessage(key, content)
				if err != nil {
					log.Printf(
						missingContextMessageLong,
					)
					http.Error(w, missingContextMessage, http.StatusInternalServerError)
					return
				}
			} else {
				nextData = content
			}
			ctx := context.WithValue(r.Context(), bodyContextKey{}, nextData)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// MiddlewareGeneratorSingleJSON промежуточная функция обработки одиночной метрики в формате JSON.
func MiddlewareGeneratorSingleJSON(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var reader io.Reader
			content, ok := r.Context().Value(bodyContextKey{}).([]byte)
			if !ok {
				log.Printf(
					missingContextMessageLong,
				)
				http.Error(w, missingContextMessage, http.StatusInternalServerError)
				return
			}
			reader = bytes.NewReader(content)

			metric := &model.Metrics{}
			decoder := json.NewDecoder(reader)

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
			ctx := context.WithValue(r.Context(), singleMetricContextKey{}, metric)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// MiddlewareGeneratorMultipleJSON промежуточная функция обработки vfccbdf метрик в формате JSON.
func MiddlewareGeneratorMultipleJSON(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var reader io.Reader
			content, ok := r.Context().Value(bodyContextKey{}).([]byte)
			if !ok {
				log.Printf(
					"Error: [MiddlewareGeneratorMultipleJSON] Metric info not found in context status-'500'",
				)
				http.Error(w, missingContextMessage, http.StatusInternalServerError)
				return
			}
			reader = bytes.NewReader(content)
			log.Println("Incoming request Updates, before decoder")
			defer func() {
				if err := r.Body.Close(); err != nil {
					log.Printf("error on closing req body: %v", err)
				}
			}()
			var metrics []*model.Metrics
			bodyBytes, _ := io.ReadAll(reader)
			err := json.Unmarshal(bodyBytes, &metrics)
			if err != nil {
				log.Printf("Failed to Decode incoming metricList %v, error: %v", string(bodyBytes), err)
				http.Error(w, fmt.Sprintf("Failed to Decode incoming metricList %v", err), http.StatusBadRequest)
				return
			}
			log.Printf("Incoming request Method: %v, Body: %v", r.RequestURI, string(bodyBytes))
			for i := range metrics {
				if key != "" {
					pass, err := checkMetric(metrics[i], key)
					if err != nil || !pass {
						http.Error(
							w,
							fmt.Sprintf("check Metric Error:%v", err),
							http.StatusBadRequest,
						)
						return
					}
				}
			}

			ctx := context.WithValue(r.Context(), multipleMetricsContextKey{}, metrics)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetMetricHandlerFunc - функция обработки метрик в URL.
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
			val = fmt.Sprintf("%v", g)

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
			_, err := w.Write([]byte(unkMTMessage))
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

// GetAllHandler возвращает значения всех метрик в виде html-страницы.
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

// checkMetric внутренняя функция проверки целостности метрики.
func checkMetric(metric *model.Metrics, key string) (bool, error) {
	if metric.ID == "" {
		return false, fmt.Errorf("metric.ID is empty /n Body: %v", metric)
	}
	if metric.MType == "" {
		return false, fmt.Errorf("metric.MType is empty /n Body: %v", metric)
	}
	if metric.Hash != "" {
		hashObject := signer.NewHashObject(key)
		passed, err := hashObject.Verify(metric)
		if err != nil {
			log.Printf(
				"Incoming Metric verification error: \nBody: %v, \n error: %v",
				metric,
				err)
			return false, fmt.Errorf("incoming Metric verification error: \nBody: %v, \n error: %w",
				metric,
				err)
		}
		if !passed {
			log.Printf(
				"Error: Incoming Metric could not pass signature verification: \nBody: %v",
				metric)

			return false, fmt.Errorf("error: Incoming Metric could not pass signature verification: \nBody: %v",
				metric)
		}
	}
	return true, nil
}
