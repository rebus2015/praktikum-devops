package storage

import (
	"fmt"
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
	gauges    map[string]gauge
	PollCount counter
}
