package storage

import (
	"github.com/rebus2015/praktikum-devops/internal/config"
	"github.com/rebus2015/praktikum-devops/internal/storage/filestorage"
	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

type Repository interface {
	AddGauge(name string, val interface{}) (float64, error)
	AddCounter(name string, val interface{}) (int64, error)
	GetCounter(name string) (int64, error)
	GetGauge(name string) (float64, error)
	GetView() ([]memstorage.MetricStr, error)
	Storage() *memstorage.MemStorage
}

func InitStorage(cfg *config.Config) *Repository{
	if(cfg.StoreFile=="") {
		repo := CreateMemoryRepository()
		return &repo
	} else {
		repo:= CreateFileRepository(cfg)
		return &repo
	}
}

func CreateMemoryRepository() Repository {
	repo := memstorage.NewStorage()
	return repo
}

func CreateFileRepository(cfg *config.Config) Repository {
	repo := filestorage.NewStorage(cfg)
	if(cfg.Restore){repo.Restore(cfg.StoreFile)}
	if(cfg.StoreInterval!=0){
		go repo.SaveTicker(cfg.StoreInterval)
	}
	return repo

}

