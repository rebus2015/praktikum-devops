package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"

	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/rebus2015/praktikum-devops/internal/model"
)

type config struct {
	ServerAddress  string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"PUSH_TIMEOUT"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

func getConfig() (*config, error) {
	conf := config{}
	flag.StringVar(&conf.ServerAddress, "a", "127.0.0.1:8080", "Server address")
	flag.DurationVar(&conf.ReportInterval, "r", time.Second*10, "Interval before push metrics to server")
	flag.DurationVar(&conf.ReportInterval, "p", time.Second*2, "Interval between metrics reads from runtime")

	flag.Parse()
	err := env.Parse(&conf)

	return &conf, err
}

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
	m.Lock()
	defer m.Unlock()
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

func (m *metricset) flushCounter(c string) {
	m.Lock()
	defer m.Unlock()
	m.counters[c] = 0
}

func (m *metricset) updateSend(cfg *config) error {

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	client := &http.Client{}
	//m.Lock()

	//отправляем статистику для gauge
	for g := range m.gauges {
		gmetric := m.Get(Gauge, g)
		err := sendreq(
			request(ctx, gmetric, cfg), client)
		if err != nil {
			log.Printf("Error send gauge Statistic: %v\n", err)
			return err
		}
	}

	//отправляем статистику counter
	for c := range m.counters {
		cmetric := m.Get(Count, c)
		err := sendreq(request(ctx, cmetric, cfg), client)
		if err != nil {
			log.Printf("Error send counter Statistic: %v\n", err)
			return err
		}
		m.flushCounter(c)
	}
	return nil
}

func (m *metricset) Get(mtype string, name string) *model.Metrics {
	metric := model.Metrics{
		ID:    name,
		MType: mtype,
	}
	m.RLock()
	defer m.RUnlock()
	switch mtype {
	case Gauge:
		{
			if v, ok := m.gauges[name]; ok {
				metric.Value = Ptr(float64(v))
				break
			}
			log.Printf("Client '%v': no such gauge metric", name)
		}
	case Count:
		{
			if v, ok := m.counters[name]; ok {
				metric.Delta = Ptr(int64(v))
				break
			}
			log.Printf("Client '%v': no such counter metric", name)
		}
	}
	return &metric
}

func request(ctx context.Context, metric *model.Metrics, cfg *config) *http.Request {

	queryurl := url.URL{
		Scheme: "http",
		Host:   cfg.ServerAddress,
		Path:   "update",
	}
	data, err := json.Marshal(metric)
	if err != nil {
		log.Printf("Error request '%v'\n", err)
		log.Panicf("Error request '%v'\n", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, queryurl.String(), bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Create Request failed! with error: %v\n", err)
		log.Panicf("Create Request failed! with error: %v\n", err)
	}
	req.Header.Add("Content-type", "application/json")

	return req
}

func sendreq(r *http.Request, c *http.Client) error {
	response, err := c.Do(r)
	if err != nil {
		log.Printf("Send request error: %v", err)
		return err
	}
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Read response body error: %v", err)
		return err
	}

	log.Printf("Client request for update metric %s\n", b)
	fmt.Println()
	return nil
}

func main() {

	m := metricset{}
	m.Declare()
	cfg, err := getConfig()
	if err != nil {
		log.Panicf("Error reading configuration from env variables: %v", err)
		return
	}
	log.Printf("agent started on %v", cfg.ServerAddress)

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	updticker := time.NewTicker(cfg.PollInterval)
	sndticker := time.NewTicker(cfg.ReportInterval)

	defer updticker.Stop()
	defer sndticker.Stop()

	for {
		select {
		case <-updticker.C:
			{
				m.Update()
			}
		case <-sndticker.C:
			{
				err := m.updateSend(cfg)
				if err != nil {
					log.Printf("Error send metrics: %v\n", err)
				}
			}
		case q := <-sigChan:
			{
				log.Printf("Signal notification: %v\n", q)
				os.Exit(0)
				//TODO корректно завершить обработку
			}
		}
	}
}
