package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
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

	"github.com/go-deeper/chunks"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/shirou/gopsutil/v3/mem"
	log "github.com/sirupsen/logrus"

	"github.com/rebus2015/praktikum-devops/internal/agent"
	"github.com/rebus2015/praktikum-devops/internal/model"
	"github.com/rebus2015/praktikum-devops/internal/signer"
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

func (m *metricset) gatherJSONMetrics(key string) ([]*model.Metrics, error) {
	metricList := []*model.Metrics{}
	for g := range m.gauges {
		gmetric := m.get(Gauge, g)
		if key != "" {
			hashObject := signer.NewHashObject(key)
			err := hashObject.Sign(gmetric)
			if err != nil {
				log.Printf("Error Sign gauges Statistic: %v,\n Gauge:%v", err, gmetric)
				return nil, fmt.Errorf("sign gauge error:%w", err)
			}
		}
		metricList = append(metricList, gmetric)
	}

	// собираем статистику counter
	for c := range m.counters {
		cmetric := m.get(Count, c)
		if key != "" {
			hashObject := signer.NewHashObject(key)
			err := hashObject.Sign(cmetric)
			if err != nil {
				log.Printf("Error Sign counter Statistic: %v,\n Counter: %v", err, cmetric)
				return nil, fmt.Errorf("sign counter stat error:%w", err)
			}
		}
		metricList = append(metricList, cmetric)
		m.flushCounter(c)
	}
	return metricList, nil
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

func encrypt(d []byte, pubKey *rsa.PublicKey) ([]byte, error) {
	rnd := rand.Reader
	hash := sha256.New()
	size := pubKey.Size() - 2*hash.Size() - 2
	encripted := make([]byte, 0)
	slices := chunks.Split(d, size)
	for _, slice := range slices {
		data, err := rsa.EncryptOAEP(hash, rnd, pubKey, slice, []byte(""))
		if err != nil {
			return nil, fmt.Errorf("message encript error: %w", err)
		}
		encripted = append(encripted, data...)
	}
	fmt.Fprintln(os.Stdout, len(d), len(encripted))
	return encripted, nil
}

func request(ctx context.Context, metrics []model.Metrics, cfg *agent.Config) *retryablehttp.Request {
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

	buf := bytes.NewBuffer(data)
	if cfg.CryptoKey != nil {
		d, err1 := encrypt(data, cfg.CryptoKey)
		if err1 != nil {
			log.Printf("Create Request failed! with error: %s\n", err)
			log.Panicf("Create Request failed! with error: %s\n", err)
		}
		buf = bytes.NewBuffer(d)
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, queryurl.String(), buf)
	if err != nil {
		log.Printf("Create Request failed! with error: %v\n", err)
		log.Panicf("Create Request failed! with error: %v\n", err)
	}
	req.Header.Add("Content-type", "application/json")

	return req
}

func sendreq(ctx context.Context, args agent.Args) error {
	r := request(ctx, args.Metrics, args.Config)
	response, err := args.Client.Do(r)
	if err != nil {
		log.Printf("Send request error: %v", err)
		return fmt.Errorf("send request error:%w", err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("response.Body.Close error: %v", err)
		}
	}()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Read response body error: %v", err)
		return fmt.Errorf("responce read error:%w", err)
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

func valuer(m []*model.Metrics) []model.Metrics {
	var mm = make([]model.Metrics, len(m))
	for _, s := range m {
		mm = append(mm, *s)
	}
	return mm
}

func (m *metricset) updateSendMultiple(ctx context.Context, cfg *agent.Config) error {
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.RetryWaitMax = time.Duration(5 * time.Second)
	metricList, err := m.gatherJSONMetrics(cfg.Key)
	if err != nil {
		log.Printf("Error send metricList Statistic: %v,\n Values: %v", err, metricList)
		return err
	}
	jobs := []agent.Job{}
	length := len(metricList)
	for i := 0; i < length; i += 2 {
		var section []*model.Metrics
		if i > length-4 {
			section = metricList[i:]
		} else {
			section = metricList[i : i+4]
		}
		jobs = append(jobs, agent.Job{
			Descriptor: i,
			ExecFn:     sendreq,
			Args: agent.Args{
				Client:  client,
				Metrics: valuer(section),
				Config:  cfg,
			},
		})
	}
	wp := agent.New(cfg.RateLimit)

	go wp.GenerateFrom(jobs)
	go wp.Run(ctx)

	for {
		select {
		case r, ok := <-wp.ErrCh():
			if !ok {
				continue
			}
			if r.Err != nil {
				log.Printf("unexpected error: %v from worker on Job %v", r.Err, r.Descriptor)
			}
			log.Printf("worker processed Job %v", r.Descriptor)

		case <-wp.Done:
			log.Printf("worker FINISHED")
			return nil
		}
	}
}

func (m *metricset) sndWorker(ctx context.Context, cfg *agent.Config, errCh chan<- error) {
	ticker := time.NewTicker(cfg.ReportInterval)
	defer close(errCh)
	for {
		select {
		case <-ticker.C:
			err := m.updateSendMultiple(ctx, cfg)
			if err != nil {
				errCh <- fmt.Errorf("error send metrics: %w", err)
			}
		case <-ctx.Done():
			log.Println("sndWorker stopped")
			return
		}
	}
}

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n\n", buildCommit)

	m := metricset{}
	m.Declare()
	cfg, err := agent.GetConfig()
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

	go m.updWorkerRuntime(ctx, cfg.ReportInterval)
	go m.updWorkerPs(ctx, cfg.ReportInterval)
	go m.sndWorker(ctx, cfg, errCh)

	for {
		select {
		case q := <-sigChan:
			cancel()
			log.Printf("Signal notification: %v\n", q)
			syscall.Exit(0)

		case err := <-errCh:
			if err != nil {
				log.Println(err)
			}
		}
	}
}
