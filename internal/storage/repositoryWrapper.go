package storage

import (
	"log"
	"sync"

	"github.com/rebus2015/praktikum-devops/internal/storage/filestorage"
	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

type RepositoryWrapper struct {
	memstorage  memstorage.MemStorage
	filestorage filestorage.FileStorage
	mux         sync.RWMutex
}

var _ Repository = new(RepositoryWrapper)

func NewRepositoryWrapper(mes memstorage.MemStorage, fs filestorage.FileStorage) *RepositoryWrapper {
	return &RepositoryWrapper{
		memstorage:  mes,
		filestorage: fs,
		mux:         sync.RWMutex{},
	}
}

func (rw *RepositoryWrapper) AddGauge(name string, val interface{}) (float64, error) {
	rw.mux.Lock()
	defer rw.mux.Unlock()
	retval, err := rw.memstorage.SetGauge(name, val)
	if rw.filestorage.SyncMode {
		errs := rw.filestorage.Save(&rw.memstorage)
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
	if rw.filestorage.SyncMode {
		errs := rw.filestorage.Save(&rw.memstorage)
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
