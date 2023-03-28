package memstorage

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

func ptr[T any](v T) *T {
	return &v
}

func (g GMetric) String() string {
	x := fmt.Sprintf("%g", g.Val)
	return x
}

func (c CMetric) String() string {
	x := fmt.Sprintf("%v", c.Val)
	return x
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

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func NewStorage() *MemStorage {
	return &MemStorage{
		map[string]float64{},
		map[string]int64{},
	}
}

func (m *MemStorage) SetGauge(name string, val interface{}) (float64, error) {
	g := GMetric{}
	switch v := val.(type) {
	case string:
		{
			err := g.TryParse(name, v)
			if err != nil {
				return 0, err
			}
		}
	case *float64:
		{
			g.Name = name
			g.Val = *v
		}
	default:
		return 0, errors.New("unexpected gauge value")
	}

	m.Gauges[g.Name] = g.Val
	return m.Gauges[g.Name], nil
}

func (m *MemStorage) IncCounter(name string, val interface{}) (int64, error) {
	c := CMetric{}
	switch v := val.(type) {
	case string:
		{
			err := c.TryParse(name, v)
			if err != nil {
				return 0, err
			}
		}
	case *int64:
		{
			c.Name = name
			c.Val = *v
		}
	default:
		return 0, errors.New("unexpected counter value")
	}
	if _, ok := m.Counters[c.Name]; !ok {
		m.Counters[c.Name] = c.Val
	} else {
		m.Counters[c.Name] += c.Val
		c.Val = m.Counters[c.Name]
	}
	return m.Counters[c.Name], nil
}

func (m *MemStorage) GetCounter(name string) (int64, error) {
	if _, ok := m.Counters[name]; !ok {
		return 0, fmt.Errorf("counter with name '%v' is not found", name)
	}
	return m.Counters[name], nil
}

func (m *MemStorage) GetGauge(name string) (float64, error) {
	if _, ok := m.Gauges[name]; !ok {
		return 0, fmt.Errorf("cauge with name '%v' is not found", name)
	}
	return m.Gauges[name], nil
}

func (m *MemStorage) GetView() ([]MetricStr, error) {
	view := []MetricStr{}
	keys := []string{}
	for k := range m.Gauges {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for key, val := range m.Counters {
		view = append(view, MetricStr{key, fmt.Sprintf("%v", val)})
	}
	for _, key := range keys {
		view = append(view, MetricStr{key, fmt.Sprintf("%v", m.Gauges[key])})
	}

	return view, nil
}
