package storage

import (
	log "github.com/sirupsen/logrus"

	"github.com/rebus2015/praktikum-devops/internal/model"
	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

type RepositoryWrapper struct {
	memstorage       memstorage.MemStorage
	secondarystorage SecondaryStorage
}

var _ Repository = new(RepositoryWrapper)

func NewRepositoryWrapper(mes memstorage.MemStorage, sec SecondaryStorage) *RepositoryWrapper {
	return &RepositoryWrapper{
		memstorage:       mes,
		secondarystorage: sec,
	}
}

func (rw *RepositoryWrapper) AddGauge(name string, val interface{}) (float64, error) {
	retval, err := rw.memstorage.SetGauge(name, val)
	if rw.secondarystorage != nil {
		if rw.secondarystorage.SyncMode() {
			errs := rw.secondarystorage.Save(&rw.memstorage)
			if errs != nil {
				log.Printf("FileStorage Save error: %v", err)
			}
		}
	}
	return retval, err
}

func (rw *RepositoryWrapper) AddCounter(name string, val interface{}) (int64, error) {
	retval, err := rw.memstorage.IncCounter(name, val)
	if rw.secondarystorage != nil {
		if rw.secondarystorage.SyncMode() {
			errs := rw.secondarystorage.Save(&rw.memstorage)
			if errs != nil {
				log.Printf("FileStorage Save error: %v", err)
			}
		}
	}
	return retval, err
}

func (rw *RepositoryWrapper) GetCounter(name string) (int64, error) {
	return rw.memstorage.GetCounter(name)
}

func (rw *RepositoryWrapper) GetGauge(name string) (float64, error) {
	return rw.memstorage.GetGauge(name)
}

func (rw *RepositoryWrapper) GetView() ([]memstorage.MetricStr, error) {
	return rw.memstorage.GetView()
}

func (rw *RepositoryWrapper) AddMetrics(m []*model.Metrics) error {
	err := rw.memstorage.AddMetrics(m)
	if rw.secondarystorage != nil {
		if rw.secondarystorage.SyncMode() {
			errs := rw.secondarystorage.Save(&rw.memstorage)
			if errs != nil {
				log.Printf("FileStorage Save error: %v", err)
			}
		}
	}
	return nil
}
