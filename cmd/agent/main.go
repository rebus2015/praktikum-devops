package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/rebus2015/praktikum-devops/internal/model"
)

type gauge float64
type counter int64

func (g gauge) String() string {
	x := fmt.Sprintf("%v", float64(g))
	return x
}
func (c counter) String() string {
	x := fmt.Sprintf("%v", int64(c))
	return x
}

type metricset struct {
	gauges   map[string]gauge
	counters map[string]counter
	sync.RWMutex
}

func (m *metricset) Declare() {
	m.counters = map[string]counter{
		"PollCount": 0,
	}
	m.gauges = map[string]gauge{
		"Alloc":         0,
		"BuckHashSys":   0,
		"Frees":         0,
		"GCCPUFraction": 0,
		"GCSys":         0,
		"HeapAlloc":     0,
		"HeapIdle":      0,
		"HeapInuse":     0,
		"HeapObjects":   0,
		"HeapReleased":  0,
		"HeapSys":       0,
		"LastGC":        0,
		"Lookups":       0,
		"MCacheInuse":   0,
		"MCacheSys":     0,
		"MSpanInuse":    0,
		"MSpanSys":      0,
		"Mallocs":       0,
		"NextGC":        0,
		"NumForcedGC":   0,
		"NumGC":         0,
		"OtherSys":      0,
		"PauseTotalNs":  0,
		"StackInuse":    0,
		"StackSys":      0,
		"Sys":           0,
		"TotalAlloc":    0,
		"RandomValue":   0,
	}
}

const (
	Gauge string = "gauge"
	Count string = "counter"
)

func (m *metricset) Update() {

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	m.Lock()
	defer m.Unlock()
	m.counters["PollCount"]++

	m.gauges["Alloc"] = gauge(ms.Alloc)
	m.gauges["BuckHashSys"] = gauge(ms.BuckHashSys)
	m.gauges["Frees"] = gauge(ms.Frees)
	m.gauges["GCCPUFraction"] = gauge(ms.GCCPUFraction)
	m.gauges["GCSys"] = gauge(ms.GCSys)
	m.gauges["HeapAlloc"] = gauge(ms.HeapAlloc)
	m.gauges["HeapIdle"] = gauge(ms.HeapIdle)
	m.gauges["HeapInuse"] = gauge(ms.HeapInuse)
	m.gauges["HeapObjects"] = gauge(ms.HeapObjects)
	m.gauges["HeapReleased"] = gauge(ms.HeapReleased)
	m.gauges["HeapSys"] = gauge(ms.HeapSys)
	m.gauges["LastGC"] = gauge(ms.LastGC)
	m.gauges["Lookups"] = gauge(ms.Lookups)
	m.gauges["MCacheInuse"] = gauge(ms.MCacheInuse)
	m.gauges["MCacheSys"] = gauge(ms.MCacheSys)
	m.gauges["MSpanInuse"] = gauge(ms.MSpanInuse)
	m.gauges["MSpanSys"] = gauge(ms.MSpanSys)
	m.gauges["Mallocs"] = gauge(ms.Mallocs)
	m.gauges["NextGC"] = gauge(ms.NextGC)
	m.gauges["NumForcedGC"] = gauge(ms.NumForcedGC)
	m.gauges["NumGC"] = gauge(ms.NumGC)
	m.gauges["OtherSys"] = gauge(ms.OtherSys)
	m.gauges["PauseTotalNs"] = gauge(ms.PauseTotalNs)
	m.gauges["StackInuse"] = gauge(ms.StackInuse)
	m.gauges["StackSys"] = gauge(ms.StackSys)
	m.gauges["Sys"] = gauge(ms.Sys)
	m.gauges["TotalAlloc"] = gauge(ms.TotalAlloc)
	m.gauges["RandomValue"] = gauge(rand.Float32())
}
func Ptr[T any](v T) *T {
	return &v
}

func (m *metricset) Get(mtype string, name string) *model.Metrics {
	m.RLock()
	defer m.RUnlock()
	metric := model.Metrics{
		ID:    name,
		MType: mtype,
	}
	switch mtype {
	case Gauge:
		{
			if v, ok := m.gauges[name]; ok {
				metric.Value = Ptr(float64(v))
				break
			}
			log.Panicf("%v: no such gauge metric", name)
		}
	case Count:
		{
			if v, ok := m.counters[name]; ok {
				metric.Delta = Ptr(int64(v))
				break
			}
			log.Panicf("%v: no such counter metric", name)
		}
	}
	return &metric
}
func request(metric *model.Metrics, cfg *config) *http.Request {

	queryurl := url.URL{
		Scheme: "http",
		Host:   cfg.ServerAddress,
		Path:   "update",
	}
	data, err := json.Marshal(metric)
	if err != nil {
		log.Panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, queryurl.String(), bytes.NewBuffer(data))
	if err != nil {
		log.Panicf("Create Request failed! with error: %v", err)
	}
	req.Header.Add("content-type", "application/json")

	return req
}
func makereq(typename string, name string, val string, cfg *config) *http.Request {
	path, err := url.JoinPath(
		"update",
		typename,
		name,
		val)
	if err != nil {
		log.Panicf("Url JoinPath failed! with error: %v", err)
	}
	queryurl := url.URL{
		Scheme: "http",
		Host:   cfg.ServerAddress,
		Path:   path,
	}
	req, err := http.NewRequest(http.MethodPost, queryurl.String(), nil)
	if err != nil {
		log.Panicf("Create Request failed! with error: %v", err)
	}
	req.Header.Add("Content-Type", "text/plain")
	return req
}

func sendreq(r *http.Request, c *http.Client) {
	response, err := c.Do(r)
	if err != nil {
		log.Panicf("Client request %v failed with error: %v", r.RequestURI, err)
	}
	defer response.Body.Close()
	_, err1 := io.Copy(io.Discard, response.Body)
	if err != nil {
		log.Panic(err1)
	}
}

func main() {

	m := metricset{}
	m.Declare()
	cfg := getConfig()
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	updticker := time.NewTicker(cfg.PollInterval)
	sndticker := time.NewTicker(cfg.ReportInternal)

	defer updticker.Stop()
	defer sndticker.Stop()

	for {
		select {
		case t := <-updticker.C:
			{
				m.Update()
				fmt.Printf("%v Updateed metrics", t)
				fmt.Println("")
			}
		case s := <-sndticker.C:
			{

				client := &http.Client{}

				//отправляем статистику для gauge
				for g, v := range m.gauges {
					sendreq(request(m.Get(Gauge, g), cfg), client)
					fmt.Printf("%v %v Send Statistic", s, makereq(reflect.TypeOf(v).Name(), g, v.String(), cfg).URL)
					fmt.Println("")
				}

				//отправляем статистику counter
				for c, v := range m.counters {
					sendreq(request(m.Get(Count, c), cfg), client)
					fmt.Printf("%v %v Send Statistic", s, makereq(reflect.TypeOf(v).Name(), c, v.String(), cfg).URL)
					fmt.Println("")
					m.Lock()
					m.counters[c] = 0
					m.Unlock()
				}

			}
		case q := <-sigChan:
			{
				fmt.Printf("q: %v\n", q)
				//TODO корректно завершить обработку
			}
		}
	}
}
