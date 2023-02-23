package storage

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
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
type MetricStr struct {
	Name string
	Val  string
}

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
	sync.RWMutex
}

type Repository interface {
	Init()
	AddGauge(name string, val string) error
	AddCounter(name string, val string) error
	GetCounter(name string) (string, error)
	GetGauge(name string) (string, error)
	GetView() ([]MetricStr, error)
}

func (m *MemStorage) Init() {
	m.Gauges = make(map[string]Gauge)
	m.Counters = make(map[string]Counter)
}

func (m *MemStorage) AddGauge(name string, val string) error {
	if m == nil {
		return errors.New("m *MemStorage is nil. Init it first")
	}
	if m.Gauges == nil {
		return errors.New("m.Gauges map[string]gauge is nil")
	}
	g := GMetric{}
	err := g.TryParse(name, val)
	if err != nil {
		return err
	}
	m.Lock()
	defer m.Unlock()
	m.Gauges[g.Name] = g.Val
	return nil
}

func (m *MemStorage) AddCounter(name string, val string) error {
	if m == nil {
		return errors.New("m *MemStorage is nil. Init it first")
	}
	if m.Counters == nil {
		return errors.New("m.Counters map[string]counter is nil")
	}
	c := CMetric{}
	err := c.TryParse(name, val)
	if err != nil {
		return err
	}
	m.Lock()
	defer m.Unlock()
	if _, ok := m.Counters[c.Name]; !ok {
		m.Counters[c.Name] = c.Val
		return nil
	}
	m.Counters[c.Name] = m.Counters[c.Name] + c.Val
	return nil
}

func (m *MemStorage) GetCounter(name string) (string, error) {
	if m == nil {
		return "", errors.New("m *MemStorage is nil. Init it first")
	}
	if m.Counters == nil {
		return "", errors.New("m.Counters map[string]counter is nil")
	}
	m.RLock()
	defer m.RUnlock()
	if _, ok := m.Counters[name]; !ok {
		return "", fmt.Errorf("Counter with name %v is not found", name)
	}
	return fmt.Sprintf("%v", int64(m.Counters[name])), nil
}

func (m *MemStorage) GetGauge(name string) (string, error) {
	if m == nil {
		return "", errors.New("m *MemStorage is nil. Init it first")
	}
	if m.Gauges == nil {
		return "", errors.New("m.Counters map[string]counter is nil")
	}
	m.RLock()
	defer m.RUnlock()
	if _, ok := m.Gauges[name]; !ok {
		return "", fmt.Errorf("Gauge with name %v is not found", name)
	}
	return fmt.Sprintf("%v", float64(m.Gauges[name])), nil
}

func (m *MemStorage) GetView() ([]MetricStr, error) {
	m.RLock()
	defer m.RUnlock()
	view := []MetricStr{}
	for key, val := range m.Counters {
		view = append(view, MetricStr{key, fmt.Sprintf("%v", val)})

	}
	for key, val := range m.Gauges {
		view = append(view, MetricStr{key, fmt.Sprintf("%v", val)})
	}
	return view, nil
}
