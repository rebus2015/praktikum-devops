package storage

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/rebus2015/praktikum-devops/internal/model"
	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

type RepositoryWrapper struct {
	memstorage       *memstorage.MemStorage
	secondarystorage SecondaryStorage
}

const fsSaveErrorMsg string = "FileStorage Save error: %v"

var _ Repository = new(RepositoryWrapper)

func NewRepositoryWrapper(mes *memstorage.MemStorage, sec SecondaryStorage) *RepositoryWrapper {
	return &RepositoryWrapper{
		memstorage:       mes,
		secondarystorage: sec,
	}
}
func (rw *RepositoryWrapper) AddGauge(name string, val interface{}) (float64, error) {
	retval, err := rw.memstorage.SetGauge(name, val)
	if rw.secondarystorage != nil {
		if rw.secondarystorage.SyncMode() {
			errs := rw.secondarystorage.Save(context.Background(), rw.memstorage)
			if errs != nil {
				log.Printf("FileStorage Save error: %v", err)
			}
		}
	}
	if err != nil {
		return 0, fmt.Errorf("AddGauge error:%w", err)
	}
	return retval, nil
}

func (rw *RepositoryWrapper) AddCounter(name string, val interface{}) (int64, error) {
	retval, err := rw.memstorage.IncCounter(name, val)
	if rw.secondarystorage != nil {
		if rw.secondarystorage.SyncMode() {
			errs := rw.secondarystorage.Save(context.Background(), rw.memstorage)
			if errs != nil {
				log.Printf("FileStorage Save error: %v", err)
			}
		}
	}
	if err != nil {
		return 0, fmt.Errorf("AddCounter error:%w", err)
	}
	return retval, nil
}

// Func (rw *RepositoryWrapper) AddGauge(name string, val interface{}) (float64, error) {
// 	retval, err := rw.memstorage.SetGauge(name, val)
// 	if err != nil {
// 		return retval, fmt.Errorf("addGauge error:%w", err)
// 	}
// 	if rw.secondarystorage != nil {
// 		if rw.secondarystorage.SyncMode() {
// 			err = rw.secondarystorage.Save(context.Background(), rw.memstorage)
// 			if err != nil {
// 				log.Printf(fsSaveErrorMsg, err)
// 				return 0, fmt.Errorf("fs save error:%w", err)
// 			}
// 		}
// 	}
// 	return retval, nil
// }.

// Func (rw *RepositoryWrapper) AddCounter(name string, val interface{}) (int64, error) {
// 	retval, err := rw.memstorage.IncCounter(name, val)
// 	if err != nil {
// 		return 0, fmt.Errorf("addCounter error:%w", err)
// 	}
// 	if rw.secondarystorage != nil {
// 		if rw.secondarystorage.SyncMode() {
// 			errs := rw.secondarystorage.Save(context.Background(), rw.memstorage)
// 			if errs != nil {
// 				log.Printf(fsSaveErrorMsg, err)
// 				return 0, fmt.Errorf("secondarystorage Save error: %w", errs)
// 			}
// 		}
// 	}
// 	return retval, nil
// }.

func (rw *RepositoryWrapper) GetCounter(name string) (int64, error) {
	result, err := rw.memstorage.GetCounter(name)
	if err != nil {
		return 0, fmt.Errorf("GetCounter error: %w", err)
	}
	return result, nil
}

func (rw *RepositoryWrapper) GetGauge(name string) (float64, error) {
	result, err := rw.memstorage.GetGauge(name)
	if err != nil {
		return 0, fmt.Errorf("GetGauge error: %w", err)
	}
	return result, nil
}

func (rw *RepositoryWrapper) GetView() ([]memstorage.MetricStr, error) {
	result, err := rw.memstorage.GetView()
	if err != nil {
		return nil, fmt.Errorf("GetView error: %w", err)
	}
	return result, nil
}

func (rw *RepositoryWrapper) AddMetrics(m []*model.Metrics) error {
	err := rw.memstorage.AddMetrics(m)
	if rw.secondarystorage != nil {
		if rw.secondarystorage.SyncMode() {
			errs := rw.secondarystorage.Save(context.Background(), rw.memstorage)
			if errs != nil {
				log.Printf(fsSaveErrorMsg, err)
			}
		}
	}
	return nil
}
