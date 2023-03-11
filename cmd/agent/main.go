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
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/rebus2015/praktikum-devops/internal/model"
)

type config struct {
	ServerAddress  string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func getConfig() *config {
	return &config{
		ServerAddress:  "127.0.0.1:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 5 * time.Second,
	}
}

type gauge float64
type counter int64

const (
	Gauge string = "gauge"
	Count string = "counter"
)

func (g gauge) String() string {
	x := fmt.Sprintf("%v", float64(g))
	return x
}
func (c counter) String() string {
	x := fmt.Sprintf("%v", int64(c))
	return x
}

type metricset struct {
	gauges    map[string]gauge
	PollCount counter
	sync.RWMutex
}

func (m *metricset) Declare() {
	m.PollCount = 0
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

func (m *metricset) Update() {

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	m.Lock()
	defer m.Unlock()
	m.PollCount++

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

func makereq(cfg *config, typename string, name string, val string) *http.Request {
	queryurl := url.URL{
		Scheme: "http",
		Host:   cfg.ServerAddress,
		Path:   "update",
	}
	metric := model.Metrics{
		ID:    name,
		MType: typename,
	}
	switch typename {
	case Gauge:
		{
			val, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Printf("convertion error %v", err)
				//log.Panicf("convertion error %v", err)
			}
			metric.Value = &val
		}
	case Count:
		{
			val, err := strconv.ParseInt(val, 0, 64)
			if err != nil {
				log.Printf("convertion error %v", err)
				//log.Panicf("convertion error %v", err)
			}
			metric.Delta = &val
		}
	}
	buf, err := json.Marshal(metric)
	if err != nil {
		log.Printf("Json Marshal error %v", err)
		//log.Panicf("Json Marshal error %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, queryurl.String(), bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("Create Request failed! with error: %v", err)
		//log.Panicf("Create Request failed! with error: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	return req
}

func sendreq(r *http.Request, c *http.Client) error {
	response, err := c.Do(r)
	if err != nil {
		//log.Panicf("Client request %v failed with error: %v", r.RequestURI, err)
		log.Printf("Client request %v failed with error: %v", r.RequestURI, err)
		return err
	}
	defer response.Body.Close()
	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {

		log.Printf("Read reaponce body error: %v", err)
		return err
		//log.Panic(err1)
	}
	return nil
}

// const hostip string = "127.0.0.1:8080"
// const pollinterval time.Duration = 2 * time.Second
// const reportintelval time.Duration = 10 * time.Second

func main() {

	m := metricset{}
	m.Declare()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	cfg := getConfig()
	updticker := time.NewTicker(cfg.PollInterval)
	sndticker := time.NewTicker(cfg.ReportInterval)

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
				m.RLock()
				//отправляем статистику для gauge
				for g, v := range m.gauges {
					_ = sendreq(
						makereq(cfg, Gauge, g, v.String()),
						client)
					fmt.Printf("%v Send Statistic", s)
					fmt.Println("")
				}
				m.RUnlock()
				//отправляем статистику counter
				_ =
					sendreq(
						makereq(cfg, Count, "PollCount",
							m.PollCount.String()), client)
				fmt.Printf("%v Send Statistic", s)
				fmt.Println("")
				m.Lock()
				m.PollCount = 0
				m.Unlock()
			}
		case q := <-sigChan:
			{
				fmt.Printf("q: %v\n", q)
				//TODO корректно завершить обработку
			}
		}
	}
}
