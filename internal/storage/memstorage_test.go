package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGMetric_String(t *testing.T) {
	tests := []struct {
		name string
		g    GMetric
		want string
	}{
		{"GMetric to string int",
			GMetric{"Gauge1", 424},
			"424",
		},
		{"GMetric to string negative float",
			GMetric{"Gauge2", -424.000234},
			"-424.000234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.g
			got := g.String()
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("GMetric.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestCMetric_String(t *testing.T) {
	tests := []struct {
		name string
		c    CMetric
		want string
	}{
		{"CMetric to string int",
			CMetric{"Counter1", 424},
			"424",
		},
		{"GMetric to string negative float",
			CMetric{"Counter2", -333},
			"-333",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.c
			got := c.String()
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("CMetric.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGMetric_TryParse(t *testing.T) {
	type data struct {
		name string
		val  string
	}
	tests := []struct {
		name string
		g    data
		want *GMetric
	}{
		{name: "test1 gauge", g: data{"gauge1", "-234.300043"}, want: &GMetric{"gauge1", Gauge(-234.300043)}},
		{name: "test2 gauge", g: data{"gauge2", "102342"}, want: &GMetric{"gauge2", Gauge(102342)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := new(GMetric)
			d := tt.g
			err := g.TryParse(d.name, d.val)
			assert.NoError(t, err)
			assert.Equal(t, g, tt.want)
		})
	}
}

func TestCMetric_TryParse(t *testing.T) {
	type data struct {
		name string
		val  string
	}
	tests := []struct {
		name string
		c    data
		want *CMetric
	}{
		{name: "test1", c: data{"counter1", "-234"}, want: &CMetric{"counter1", Counter(-234)}},
		{name: "test2", c: data{"counter2", "10234234"}, want: &CMetric{"counter2", Counter(10234234)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := new(CMetric)
			d := tt.c
			err := c.TryParse(d.name, d.val)
			assert.NoError(t, err)
			assert.Equal(t, c, tt.want)
		})
	}

}

func TestMemStorage_AddGauge(t *testing.T) {
	type repo struct {
		Gauges   map[string]Gauge
		Counters map[string]Counter
	}
	type args struct {
		g GMetric
	}
	tests := []struct {
		name string
		r    repo
		args args
	}{
		{
			"the only test",
			repo{
				map[string]Gauge{"g1": Gauge(-32.00023)},
				map[string]Counter{"c1": Counter(0)},
			},
			args{GMetric{Name: "gm", Val: Gauge(234)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Gauges:   tt.r.Gauges,
				Counters: tt.r.Counters,
			}
			m.AddGauge(tt.args.g)
			_, exists := m.Gauges[tt.args.g.Name]
			assert.True(t, exists && m.Gauges[tt.args.g.Name] == tt.args.g.Val)
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	type repo struct {
		Gauges   map[string]Gauge
		Counters map[string]Counter
	}
	type args struct {
		c CMetric
	}
	tests := []struct {
		name string
		r    repo
		args args
	}{
		{
			"the only test",
			repo{
				map[string]Gauge{"g1": Gauge(-32.00023)},
				map[string]Counter{"c1": Counter(100)},
			},
			args{CMetric{Name: "c1", Val: Counter(3)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Gauges:   tt.r.Gauges,
				Counters: tt.r.Counters,
			}
			c := tt.args.c
			m.AddCounter(c)
			_, exists := m.Counters[tt.args.c.Name]
			assert.True(t, exists && m.Counters[tt.args.c.Name] == Counter(103))
		})
	}
}
