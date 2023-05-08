package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"github.com/shirou/gopsutil/v3/mem"
	log "github.com/sirupsen/logrus"

	"github.com/rebus2015/praktikum-devops/internal/model"
	"github.com/rebus2015/praktikum-devops/internal/signer"
)

type config struct {
	ServerAddress  string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"PUSH_TIMEOUT"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
}

func getConfig() (*config, error) {
	conf := config{}
	flag.StringVar(&conf.ServerAddress, "a", "127.0.0.1:8080", "Server address")
	flag.DurationVar(&conf.ReportInterval, "r", time.Second*11, "Interval before push metrics to server")
	flag.DurationVar(&conf.PollInterval, "p", time.Second*5, "Interval between metrics reads from runtime")
	flag.StringVar(&conf.Key, "k", "", "Key to sign up data with SHA256 algorythm")

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
	mux      sync.RWMutex
}

func (m *metricset) Declare() {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.counters = map[string]counter{
		"PollCount": 0,
	}

	m.gauges = map[string]gauge{
		"Alloc":           0,
		"BuckHashSys":     0,
		"Frees":           0,
		"GCCPUFraction":   0,
		"GCSys":           0,
		"HeapAlloc":       0,
		"HeapIdle":        0,
		"HeapInuse":       0,
		"HeapObjects":     0,
		"HeapReleased":    0,
		"HeapSys":         0,
		"LastGC":          0,
		"Lookups":         0,
		"MCacheInuse":     0,
		"MCacheSys":       0,
		"MSpanInuse":      0,
		"MSpanSys":        0,
		"Mallocs":         0,
		"NextGC":          0,
		"NumForcedGC":     0,
		"NumGC":           0,
		"OtherSys":        0,
		"PauseTotalNs":    0,
		"StackInuse":      0,
		"StackSys":        0,
		"Sys":             0,
		"TotalAlloc":      0,
		"RandomValue":     0,
		"TotalMemory":     0,
		"FreeMemory":      0,
		"CPUutilization1": 0,
	}
}

const (
	Gauge string = "gauge"
	Count string = "counter"
)

func (m *metricset) updatePs() {
	ms, _ := mem.VirtualMemory()
	m.mux.Lock()
	defer m.mux.Unlock()

	m.gauges["TotalMemory"] = gauge(ms.Total)
	m.gauges["FreeMemory"] = gauge(ms.Free)
	m.gauges["CPUutilization1"] = gauge(ms.UsedPercent)
}

func (m *metricset) updateRuntime() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	m.mux.Lock()
	defer m.mux.Unlock()
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
	m.gauges["RandomValue"] = gauge(newFloat64())
}

func intn(max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		panic(err)
	}
	return nBig.Int64()
}

func newFloat64() float64 {
	return float64(intn(1<<53)) / (1 << 53)
}

func ptr[T any](v T) *T {
	return &v
}

func (m *metricset) flushCounter(c string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.counters[c] = 0
}

func (m *metricset) updateSendMultiple(cfg *config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	client := &http.Client{}
	metricList := []*model.Metrics{}
	// собираем статистику для gauge
	for g := range m.gauges {
		gmetric := m.get(Gauge, g)
		if cfg.Key != "" {
			hashObject := signer.NewHashObject(cfg.Key)
			err := hashObject.Sign(gmetric)
			if err != nil {
				log.Printf("Error Sign gauges Statistic: %v,\n Gauge:%v", err, gmetric)
				return err
			}
		}
		metricList = append(metricList, gmetric)
	}

	// собираем статистику counter
	for c := range m.counters {
		cmetric := m.get(Count, c)
		if cfg.Key != "" {
			hashObject := signer.NewHashObject(cfg.Key)
			err := hashObject.Sign(cmetric)
			if err != nil {
				log.Printf("Error Sign counter Statistic: %v,\n Counter: %v", err, cmetric)
				return err
			}
		}
		metricList = append(metricList, cmetric)
		m.flushCounter(c)
	}

	err := sendreq(request(ctx, metricList, cfg), client)
	if err != nil {
		log.Printf("Error send metricList Statistic: %v,\n Values: %v", err, metricList)
		return err
	}

	return nil
}

func (m *metricset) get(mtype string, name string) *model.Metrics {
	metric := model.Metrics{
		ID:    name,
		MType: mtype,
	}
	m.mux.RLock()
	defer m.mux.RUnlock()
	switch mtype {
	case Gauge:
		{
			if v, ok := m.gauges[name]; ok {
				metric.Value = ptr(float64(v))
				break
			}
			log.Printf("Client '%v': no such gauge metric", name)
		}
	case Count:
		{
			if v, ok := m.counters[name]; ok {
				metric.Delta = ptr(int64(v))
				break
			}
			log.Printf("Client '%v': no such counter metric", name)
		}
	}
	return &metric
}

func request(ctx context.Context, metrics []*model.Metrics, cfg *config) *http.Request {
	queryurl := url.URL{
		Scheme: "http",
		Host:   cfg.ServerAddress,
		Path:   "updates",
	}
	data, err := json.Marshal(metrics)
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

func (m *metricset) updWorkerRuntime(ctx context.Context, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	for {
		select {
		case <-ticker.C:
			m.updateRuntime()
		case <-ctx.Done():
			log.Println("updWorker stopped")
		}
	}
}

func (m *metricset) updWorkerPs(ctx context.Context, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	for {
		select {
		case <-ticker.C:
			m.updatePs()
		case <-ctx.Done():
			log.Println("updWorker stopped")
		}
	}
}

func (m *metricset) sndWorker(ctx context.Context, cfg *config, errCh chan<- error) {
	ticker := time.NewTicker(cfg.ReportInterval)
	defer close(errCh)
	for {
		select {
		case <-ticker.C:
			err := m.updateSendMultiple(cfg)
			if err != nil {
				errCh <- fmt.Errorf("error send metrics: %w", err)
				// log.Printf("Error send metrics: %v\n", err)
				// return
			}
		case <-ctx.Done():
			log.Println("sndWorker stopped")
			return
		}
	}
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
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	errCh := make(chan error) // создаём канал, из которого будем ждать ошибку

	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	// wg := &sync.WaitGroup{}
	// errCh := make(chan error)
	// stopCh := make(chan struct{})
	go m.updWorkerRuntime(ctx, cfg.ReportInterval)
	go m.updWorkerPs(ctx, cfg.ReportInterval)
	go m.sndWorker(ctx, cfg, errCh)

	for {
		select {
		case q := <-sigChan:
			cancel()
			log.Printf("Signal notification: %v\n", q)
			os.Exit(0)

		case err := <-errCh:
			if err != nil {
				log.Println(err)
				// return
			}
		}
	}
}
