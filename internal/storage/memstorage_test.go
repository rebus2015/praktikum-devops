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
		name     string
		val      string
		floatval float64
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
			args{name: "gm", val: "234", floatval: 234},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newStorage()
			m.AddGauge(tt.args.name, tt.args.val)
			_, exists := m.Gauges[tt.args.name]
			assert.True(t, exists && m.Gauges[tt.args.name] == Gauge(tt.args.floatval))
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	type repo struct {
		Gauges   map[string]Gauge
		Counters map[string]Counter
	}
	type args struct {
		name   string
		val    string
		intval int64
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
			args{name: "c1", val: "3", intval: 103},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newStorage()
			m.Counters = tt.r.Counters
			m.Gauges = tt.r.Gauges
			m.AddCounter(tt.args.name, tt.args.val)
			_, exists := m.Counters[tt.args.name]
			assert.True(t, exists && m.Counters[tt.args.name] == Counter(tt.args.intval))
		})
	}
}
