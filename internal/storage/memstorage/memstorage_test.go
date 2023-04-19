package memstorage

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
		{
			"GMetric to string int",
			GMetric{
				"Gauge1",
				424,
			},
			"424.000000",
		},
		{
			"GMetric to string negative float",
			GMetric{
				"Gauge2",
				-424.000234,
			},
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
		{
			"CMetric to string int",
			CMetric{
				"Counter1",
				424,
			},
			"424",
		},
		{
			"GMetric to string negative float",
			CMetric{
				"Counter2",
				-333,
			},
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
		{
			name: "test1 gauge",
			g: data{
				"gauge1",
				"-234.300043",
			},
			want: &GMetric{
				"gauge1", -234.300043,
			},
		},
		{
			name: "test2 gauge",
			g: data{
				"gauge2",
				"102342",
			},
			want: &GMetric{
				"gauge2",
				102342,
			},
		},
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
		{
			name: "test1",
			c: data{
				"counter1",
				"-234",
			},
			want: &CMetric{
				"counter1",
				-234,
			},
		},
		{
			name: "test2",
			c: data{
				"counter2",
				"10234234",
			},
			want: &CMetric{
				"counter2",
				10234234,
			},
		},
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
		Gauges   map[string]float64
		Counters map[string]int64
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
				map[string]float64{"g1": -32.00023},
				map[string]int64{"c1": 0},
			},
			args{
				name:     "gm",
				val:      "234",
				floatval: 234,
			},
		},
		{
			"shorten test",
			repo{
				map[string]float64{"g1": -23},
				map[string]int64{"c1": 0},
			},
			args{
				name:     "gm",
				val:      "44.1200",
				floatval: 44.1200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewStorage()
			_, err := m.SetGauge(tt.args.name, tt.args.val)
			if assert.NoError(t, err) {
				_, exists := m.Gauges[tt.args.name]
				assert.True(t, exists && m.Gauges[tt.args.name] == tt.args.floatval)
			}
		})
	}
}

func TestMemStorage_AddCounterString(t *testing.T) {
	type repo struct {
		Gauges   map[string]float64
		Counters map[string]int64
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
				map[string]float64{"g1": -32.00023},
				map[string]int64{"c1": 100},
			},
			args{name: "c1", val: "3", intval: 103},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewStorage()
			m.Counters = tt.r.Counters
			m.Gauges = tt.r.Gauges
			_, err := m.IncCounter(tt.args.name, tt.args.val)
			if assert.NoError(t, err) {
				_, exists := m.Counters[tt.args.name]
				assert.True(t, exists && m.Counters[tt.args.name] == tt.args.intval)
			}
		})
	}
}

func TestMemStorage_AddCounterInt(t *testing.T) {
	type repo struct {
		Gauges   map[string]float64
		Counters map[string]int64
	}
	type args struct {
		name   string
		val    *int64
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
				map[string]float64{"g1": -32.00023},
				map[string]int64{"c1": 100},
			},
			args{name: "c1", val: ptr(int64(3)), intval: 103},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewStorage()
			m.Counters = tt.r.Counters
			m.Gauges = tt.r.Gauges
			_, err := m.IncCounter(tt.args.name, tt.args.val)
			if assert.NoError(t, err) {
				_, exists := m.Counters[tt.args.name]
				assert.True(t, exists && m.Counters[tt.args.name] == tt.args.intval)
			}
		})
	}
}
