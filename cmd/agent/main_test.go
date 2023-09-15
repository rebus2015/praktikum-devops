package main

import (
	"reflect"
	"sync"
	"testing"

	"github.com/rebus2015/praktikum-devops/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_metricset_gatherJSONMetrics(t *testing.T) {
	type fields struct {
		gauges   map[string]gauge
		counters map[string]counter
		mux      *sync.RWMutex
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Metrics
		wantErr bool
	}{
		{
			"1st test",
			fields{
				map[string]gauge{"Test": gauge((float64(23.665)))},
				map[string]counter{},
				&sync.RWMutex{},
			},
			args{"Test"},
			[]*model.Metrics{{
				ID:    "Test",
				MType: "gauge",
				Delta: nil,
				Value: ptr(float64(23.665)),
				Hash:  "5443e7ef79efe747527ba0d2cf57c1f83015728b93304324fa50087785aedc40",
			},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &metricset{
				gauges:   tt.fields.gauges,
				counters: tt.fields.counters,
				mux:      sync.RWMutex{},
			}
			got, err := m.gatherJSONMetrics(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("metricset.gatherJSONMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(&got, &tt.want) {
				t.Errorf("metricset.gatherJSONMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metricset_Declare(t *testing.T) {
	m := metricset{}
	m.Declare()
	assert.NotNil(t, m.counters)
	assert.NotNil(t, m.gauges)
}

func Test_metricset_flushCounter(t *testing.T) {
	type fields struct {
		gauges   map[string]gauge
		counters map[string]counter
		mux      *sync.RWMutex
	}
	type args struct {
		c string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{{
		"1st test",
		fields{
			map[string]gauge{},
			map[string]counter{"Test": counter((int64(23)))},
			&sync.RWMutex{},
		},
		args{"Test"},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &metricset{
				gauges:   tt.fields.gauges,
				counters: tt.fields.counters,
				mux:      sync.RWMutex{},
			}
			m.flushCounter(tt.args.c)
			assert.Equal(t, tt.fields.counters[string(tt.args.c)], counter(0))
		})
	}
}
