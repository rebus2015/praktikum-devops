package storage

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/rebus2015/praktikum-devops/internal/model"
)

func Ptr[T any](v T) *T {
	return &v
}
func (g GMetric) String() string {
	x := fmt.Sprintf("%v", float64(g.Val))
	return x
}

func (c CMetric) String() string {
	x := fmt.Sprintf("%v", int64(c.Val))
	return x
}

func (g *GMetric) Metric() *model.Metrics {
	m := model.Metrics{
		ID:    g.Name,
		MType: "gauge",
		Delta: new(int64),
		Value: Ptr(float64(g.Val)),
	}
	return &m
}

func (c *CMetric) Metric() *model.Metrics {
	m := model.Metrics{
		ID:    c.Name,
		MType: "counter",
		Delta: Ptr(int64(c.Val)),
		Value: new(float64),
	}
	return &m
}

func (c *CMetric) TryParse(cname string, cval string) error {
	v, err := strconv.ParseInt(cval, 10, 64)
	if err != nil {
		return err
	}
	c.Name = cname
	c.Val = v
	return nil
}
func (g *GMetric) TryParse(gname string, gval string) error {
	v, err := strconv.ParseFloat(gval, 64)
	if err != nil {
		return err
	}
	g.Name = gname
	g.Val = v
	return nil
}

type MetricStr struct {
	Name string
	Val  string
}

type GMetric struct {
	Name string
	Val  float64
}

type CMetric struct {
	Name string
	Val  int64
}

type memStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
	sync.RWMutex
}

func CreateRepository() Repository {
	return newStorage()
}

func newStorage() *memStorage {
	return &memStorage{
		map[string]float64{},
		map[string]int64{},
		sync.RWMutex{},
	}
}

type Repository interface {
	AddGauge(name string, val interface{}) (model.Metrics, error)
	AddCounter(name string, val interface{}) (model.Metrics, error)
	GetCounter(name string) (int64, error)
	GetGauge(name string) (float64, error)
	FillMetric(data *model.Metrics) error
	GetView() ([]MetricStr, error)
}

func (m *memStorage) AddGauge(name string, val interface{}) (model.Metrics, error) {

	g := GMetric{}
	switch v := val.(type) {
	case string:
		{
			err := g.TryParse(name, v)
			if err != nil {
				return model.Metrics{}, err
			}
		}
	case *float64:
		{
			g.Name = name
			g.Val = *v
		}
	default:
		return model.Metrics{}, errors.New("unexpected gauge value")

	}

	m.Lock()
	defer m.Unlock()
	m.Gauges[g.Name] = g.Val
	return *g.Metric(), nil
}

func (m *memStorage) AddCounter(name string, val interface{}) (model.Metrics, error) {
	c := CMetric{}
	switch v := val.(type) {
	case string:
		{
			err := c.TryParse(name, v)
			if err != nil {
				return model.Metrics{}, err
			}
		}
	case *int64:
		{
			c.Name = name
			c.Val = *v
		}
	default:
		return model.Metrics{}, errors.New("unexpected counter value")

	}
	m.Lock()
	defer m.Unlock()
	if _, ok := m.Counters[c.Name]; !ok {
		m.Counters[c.Name] = c.Val
	} else {
		m.Counters[c.Name] = m.Counters[c.Name] + c.Val
		c.Val = m.Counters[c.Name]
	}
	return *c.Metric(), nil
}

// TODO: change retrurn value to native, use convert after call when necessery
func (m *memStorage) GetCounter(name string) (int64, error) {
	m.RLock()
	defer m.RUnlock()
	if _, ok := m.Counters[name]; !ok {
		return 0, fmt.Errorf("%v: Counter with name is not found", name)
	}
	return m.Counters[name], nil
}

// TODO: change retrurn value to native, use convert after call when necessery
func (m *memStorage) GetGauge(name string) (float64, error) {
	m.RLock()
	defer m.RUnlock()
	if _, ok := m.Gauges[name]; !ok {
		return 0, fmt.Errorf("%v: Gauge with name is not found", name)
	}
	return m.Gauges[name], nil
}

func (m *memStorage) FillMetric(data *model.Metrics) error {
	m.RLock()
	defer m.RUnlock()
	switch data.MType {
	case "counter":
		{
			if v, ok := m.Counters[data.ID]; ok {
				data.Delta = Ptr(int64(v))
				break
			}
			return fmt.Errorf("%v: Counter with name is not found", data.ID)

		}
	case "gauge":
		{
			if v, ok := m.Gauges[data.ID]; ok {
				data.Value = Ptr(float64(v))
				break
			}
			return fmt.Errorf("%v: Gauge with name is not found", data.ID)

		}
	default:
		{
			return nil
		}
	}

	return nil
}

func (m *memStorage) GetView() ([]MetricStr, error) {
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
