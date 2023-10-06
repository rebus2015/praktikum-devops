package memstorage

import (
	"sync"
	"testing"

	"github.com/rebus2015/praktikum-devops/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestGMetric_String(t *testing.T) {
	tests := []struct {
		name string
		want string
		g    GMetric
	}{
		{
			name: "GMetric to string int",
			g: GMetric{
				"Gauge1",
				424,
			},
			want: "424.000000",
		},
		{
			name: "GMetric to string negative float",
			g: GMetric{
				"Gauge2",
				-424.000234,
			},
			want: "-424.000234",
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
		want string
		c    CMetric
	}{
		{
			name: "CMetric to string int",
			c: CMetric{
				"Counter1",
				424,
			},
			want: "424",
		},
		{
			name: "GMetric to string negative float",
			c: CMetric{
				"Counter2",
				-333,
			},
			want: "-333",
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
		want *GMetric
		g    data
		name string
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
		want *CMetric
		name string
		c    data
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
		val    *int64
		name   string
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

func TestMemStorage_GetGauge(t *testing.T) {
	type fields struct {
		Gauges   map[string]float64
		Counters map[string]int64
		Mux      *sync.RWMutex
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			"the only positive test",
			fields{
				map[string]float64{"g1": -32.00023},
				map[string]int64{"c1": 100},
				&sync.RWMutex{},
			},
			args{name: "g1"},
			float64(-32.00023),
			false,
		},
		{
			"the only negative test",
			fields{
				map[string]float64{"g1": -32.00023},
				map[string]int64{"c1": 100},
				&sync.RWMutex{},
			},
			args{name: "c3"},
			float64(22),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
				Mux:      tt.fields.Mux,
			}
			got, err := m.GetGauge(tt.args.name)
			if err != nil {
				if tt.wantErr {
					assert.Equal(t, err != nil, tt.wantErr)
					return
				}
				t.Errorf("MemStorage.GetGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemStorage.GetGauge() = %v, want %v", got, tt.want)
			}
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	type fields struct {
		Gauges   map[string]float64
		Counters map[string]int64
		Mux      *sync.RWMutex
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			"the only positive test",
			fields{
				map[string]float64{"g1": -32.00023},
				map[string]int64{"c1": 100},
				&sync.RWMutex{},
			},
			args{name: "c1"},
			int64(100),
			false,
		},
		{
			"the only negative test",
			fields{
				map[string]float64{"g1": -32.00023},
				map[string]int64{"c1": 100},
				&sync.RWMutex{},
			},
			args{name: "c3"},
			int64(22),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
				Mux:      tt.fields.Mux,
			}
			got, err := m.GetCounter(tt.args.name)
			if err != nil {
				if tt.wantErr {
					assert.Equal(t, err != nil, tt.wantErr)
					return
				}
				t.Errorf("MemStorage.GetGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemStorage.GetGauge() = %v, want %v", got, tt.want)
			}
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMemStorage_AddMetrics(t *testing.T) {
	type fields struct {
		Gauges   map[string]float64
		Counters map[string]int64
		Mux      *sync.RWMutex
	}
	type args struct {
		metrics []*model.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			"the only positive test",
			fields{
				map[string]float64{"g1": -32.00023},
				map[string]int64{"c1": 100},
				&sync.RWMutex{},
			},
			args{
				[]*model.Metrics{
					{
						ID:    "TotalMemory",
						MType: "gauge",
						Delta: nil,
						Value: ptr(float64(7268679680)),
						Hash:  "",
					},
					{
						ID:    "CPU",
						MType: "counter",
						Delta: ptr(int64(25)),
						Value: nil,
						Hash:  "",
					},
				},
			},
			false,
			false,
		}, {
			"Positive test 1 unknown type",
			fields{
				map[string]float64{"g1": -32.00023},
				map[string]int64{"c1": 100},
				&sync.RWMutex{},
			},
			args{
				[]*model.Metrics{
					{
						ID:    "mm2",
						MType: "gauge",
						Delta: nil,
						Value: ptr(float64(-0.122)),
						Hash:  "",
					},
					{
						ID:    "SomeMetric",
						MType: "unknown",
						Delta: ptr(int64(25)),
						Value: nil,
						Hash:  "",
					},
				},
			},
			false,
			false,
		}, {
			"Negative test 2  metrics type mismatch",
			fields{
				map[string]float64{"g1": -32.00023},
				map[string]int64{"c1": 100},
				&sync.RWMutex{},
			},
			args{
				[]*model.Metrics{
					{
						ID:    "mm2",
						MType: "counter",
						Delta: nil,
						Value: ptr(float64(-0.122)),
						Hash:  "",
					},
				},
			},
			false,
			true,
		},
		{
			"Negative test 2  nil",
			fields{
				map[string]float64{"g1": -32.00023},
				map[string]int64{"c1": 100},
				&sync.RWMutex{},
			},
			args{
				[]*model.Metrics{
					{
						ID:    "mm2",
						MType: "gauge",
						Delta: nil,
						Value: nil,
						Hash:  "",
					},
				},
			},
			false,
			true,
		},
		{
			"Negative test 2  nil",
			fields{
				map[string]float64{"g1": -32.00023},
				map[string]int64{"c1": 100},
				&sync.RWMutex{},
			},
			args{
				[]*model.Metrics{
					{
						ID:    "mm2",
						MType: "counter",
						Delta: nil,
						Value: nil,
						Hash:  "",
					},
				},
			},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
				Mux:      tt.fields.Mux,
			}
			funcErr := m.AddMetrics(tt.args.metrics)
			if (funcErr != nil) != tt.wantErr {
				t.Errorf("MemStorage.AddMetrics() error = %v, wantErr %v", funcErr, tt.wantErr)
			}
			if tt.wantErr {
				assert.ErrorContains(t, funcErr, "[updateJSONMetricHandlerFunc] ")
				return
			}

			for _, metric := range tt.args.metrics {
				switch metric.MType {
				case "gauge":
					val, err := m.GetGauge(metric.ID)
					assert.True(t, err == nil)
					assert.Equal(t, val, *metric.Value)
				case "counter":
					val, err := m.GetCounter(metric.ID)
					assert.True(t, err == nil)
					assert.Equal(t, val, *metric.Delta)
				}
			}
		})
	}
}
