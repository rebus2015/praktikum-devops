package storage

import (
	"log"
	"sync"

	"github.com/rebus2015/praktikum-devops/internal/model"
	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

type RepositoryWrapper struct {
	memstorage       memstorage.MemStorage
	secondarystorage SecondaryStorage
	mux              sync.RWMutex
}

var _ Repository = new(RepositoryWrapper)

func NewRepositoryWrapper(mes memstorage.MemStorage, sec SecondaryStorage) *RepositoryWrapper {
	return &RepositoryWrapper{
		memstorage:       mes,
		secondarystorage: sec,
		mux:              sync.RWMutex{},
	}
}

func (rw *RepositoryWrapper) AddGauge(name string, val interface{}) (float64, error) {
	rw.mux.Lock()
	defer rw.mux.Unlock()
	retval, err := rw.memstorage.SetGauge(name, val)
	if rw.secondarystorage.SyncMode() {
		errs := rw.secondarystorage.Save(&rw.memstorage)
		if errs != nil {
			log.Printf("FileStorage Save error: %v", err)
		}
	}
	return retval, err
}

func (rw *RepositoryWrapper) AddCounter(name string, val interface{}) (int64, error) {
	rw.mux.Lock()
	defer rw.mux.Unlock()
	retval, err := rw.memstorage.IncCounter(name, val)
	if rw.secondarystorage.SyncMode() {
		errs := rw.secondarystorage.Save(&rw.memstorage)
		if errs != nil {
			log.Printf("FileStorage Save error: %v", err)
		}
	}
	return retval, err
}

func (rw *RepositoryWrapper) GetCounter(name string) (int64, error) {
	rw.mux.RLock()
	defer rw.mux.RUnlock()
	return rw.memstorage.GetCounter(name)
}

func (rw *RepositoryWrapper) GetGauge(name string) (float64, error) {
	rw.mux.RLock()
	defer rw.mux.RUnlock()
	return rw.memstorage.GetGauge(name)
}

func (rw *RepositoryWrapper) GetView() ([]memstorage.MetricStr, error) {
	return rw.memstorage.GetView()
}

func (rw *RepositoryWrapper) AddMetrics(m []*model.Metrics) error {
	rw.mux.Lock()
	defer rw.mux.Unlock()
	err := rw.memstorage.AddMetrics(m)
	if rw.secondarystorage.SyncMode() {
		errs := rw.secondarystorage.Save(&rw.memstorage)
		if errs != nil {
			log.Printf("FileStorage Save error: %v", err)
		}
	}
	return nil
}
