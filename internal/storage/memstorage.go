package storage

import (
	"fmt"
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
		return err
	}
	c.Name = cname
	c.Val = Counter(v)
	return nil
}
func (g *GMetric) TryParse(gname string, gval string) error {
	v, err := strconv.ParseFloat(gval, 64)
	if err != nil {
		return err
	}
	g.Name = gname
	g.Val = Gauge(v)
	return nil
}

type Gauge float64
type Counter int64

type GMetric struct {
	Name string
	Val  Gauge
}

type CMetric struct {
	Name string
	Val  Counter
}

type MemStorage struct {
	Gauges   map[string]Gauge
	Counters map[string]Counter
}

type Repository interface {
	Init()
	AddGauge(g GMetric)
	AddCounter(c CMetric)
}

func (m *MemStorage) Init() {
	m.Gauges = make(map[string]Gauge)
	m.Counters = make(map[string]Counter)
}

func (m *MemStorage) AddGauge(g GMetric) {
	if m == nil {
		return
	}
	if m.Gauges == nil {
		return
	}
	m.Gauges[g.Name] = g.Val
}

func (m *MemStorage) AddCounter(c CMetric) {
	if m == nil {
		return
	}
	if m.Counters == nil {
		return
	}
	if _, ok := m.Counters[c.Name]; !ok {
		m.Counters[c.Name] = c.Val
		return
	}
	m.Counters[c.Name] = m.Counters[c.Name] + c.Val
}
