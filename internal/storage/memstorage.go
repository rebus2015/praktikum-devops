package storage

import (
	"fmt"
	"log"
	"strconv"
)

func (g GMetric) String() string {
	x := fmt.Sprintf("%v", float64(g.Val))
	return x
}
func (c CMetric) String() string {
	x := fmt.Sprintf("%v", int64(c.Val))
	return x
}

func (c *CMetric) TryParse(cname string, cval string) error {
	v, err := strconv.ParseInt(cval, 10, 64)
	if err != nil {
		log.Panic(err)
		return err
	}

	c = &CMetric{cname, counter(v)}
	return nil
}
func (g *GMetric) TryParse(cname string, gval string) error {
	v, err := strconv.ParseFloat(gval, 64)
	if err != nil {
		log.Panic(err)
		return err
	}
	g = &GMetric{cname, gauge(v)}
	return err
}

type gauge float64
type counter int64

type GMetric struct {
	Name string
	Val  gauge
}

type CMetric struct {
	Name string
	Val  counter
}

type MemStorage struct {
	Gauges   map[string]gauge
	Counters map[string]counter
}

type Repository interface {
	Init()
	AddGauge(g GMetric)
	AddCounter(c CMetric)
}

func (m *MemStorage) Init() {
	m.Gauges = make(map[string]gauge)
	m.Counters = make(map[string]counter)
}

func (m *MemStorage) AddGauge(g GMetric) {
	m.Gauges[g.Name] = g.Val
}

func (m *MemStorage) AddCounter(c CMetric) {
	m.Counters[c.Name] = m.Counters[c.Name] + c.Val
}
