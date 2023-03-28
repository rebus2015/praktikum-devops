package storage

import (
	"time"

	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

type Repository interface {
	AddGauge(name string, val interface{}) (float64, error)
	AddCounter(name string, val interface{}) (int64, error)
	GetCounter(name string) (int64, error)
	GetGauge(name string) (float64, error)
	GetView() ([]memstorage.MetricStr, error)
}

type SecondaryStorage interface {
	Save(ms *memstorage.MemStorage) error
	Restore(sf string) *memstorage.MemStorage
	SaveTicker(storeint time.Duration, ms *memstorage.MemStorage)
	SyncMode() bool
}
