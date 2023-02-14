package storage

import (
	"fmt"
)



func (g GMetric) String() string {
	x := fmt.Sprintf("%v", float64(g.Val))
	return x
}
func (c CMetric) String() string {
	x := fmt.Sprintf("%v", int64(c.Val))
	return x
}

type gauge float64
type counter int64

type GMetric struct 
{
	Name string
	Val gauge 
}

type CMetric struct
{
	Name string
	Val counter
}

type MemStorage struct {
	Gauges    map[string]gauge
	Counters  map[string]counter
}


type Repository interface{
	Init()
	AddGauge(g GMetric)
	AddCounter(c CMetric)
} 

func (m *MemStorage) Init (){
    m.Gauges=make(map[string]gauge)
	m.Counters=make(map[string]counter)
}

func (m *MemStorage) AddGauge(g GMetric){
    m.Gauges[g.Name]=g.Val
}

func (m *MemStorage) AddCounter(c CMetric){
    m.Counters[c.Name]=c.Val
}

