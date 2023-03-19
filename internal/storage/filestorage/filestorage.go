package filestorage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/rebus2015/praktikum-devops/internal/config"
	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

type FileStorage struct {
	*memstorage.MemStorage
	StoreFile string
	SyncMode  bool
}

func NewStorage(c *config.Config) *FileStorage {
	return &FileStorage{
		&memstorage.MemStorage{
			Gauges:   map[string]float64{},
			Counters: map[string]int64{},
			Mux:      sync.RWMutex{},
		},
		c.StoreFile,
		c.StoreInterval == 0,
	}
}

func (f *FileStorage) AddGauge(name string, val interface{}) (float64, error) {
	retval, err := f.MemStorage.SetGauge(name, val)
	if f.SyncMode {
		errs := f.Save()
		if errs != nil {
			log.Printf("FileStorage Save error: %v", err)
		}
	}
	return retval, err
}

func (f *FileStorage) AddCounter(name string, val interface{}) (int64, error) {
	retval, err := f.MemStorage.IncCounter(name, val)
	if f.SyncMode {
		errs := f.Save()
		if errs != nil {
			log.Printf("FileStorage Save error: %v", errs)
		}
	}
	return retval, err
}

// func (f *FileStorage) GetCounter(name string) (int64, error) {
// 	return f.MemStorage.GetCounter(name)
// }

// func (f *FileStorage) GetGauge(name string) (float64, error) {
// 	return f.MemStorage.GetGauge(name)
// }

// func (f *FileStorage) GetView() ([]memstorage.MetricStr, error) {
// 	return f.MemStorage.GetView()
// }

// func (f *FileStorage) Storage() *memstorage.MemStorage {
// 	return f.MemStorage.Storage()
// }

func (f *FileStorage) Save() error {
	writer, err := NewWriter(f.StoreFile)
	if err != nil {
		log.Printf("Save metrics to file '%s' error: %v", f.StoreFile, err)
		log.Fatal(err)
	}

	err = writer.encoder.Encode(f.MemStorage)
	if err != nil {
		log.Printf("Save metrics to file '%s' Encode error: %v", f.StoreFile, err)
		return err
	}
	return nil
}

func (f *FileStorage) Restore(sf string) {
	reader, err := NewReader(sf)
	if err != nil {
		log.Printf("Restore metrics from file '%s' reader error: %v", sf, err)
		log.Fatal(err)
	}

	checkFile, err := os.Stat(sf)
	if err != nil {
		log.Printf("Restore metrics from file '%s' Stat error: %v", sf, err)
		log.Fatal(err)
	}

	size := checkFile.Size()

	if size == 0 {
		errs := f.Save()
		if errs != nil {
			log.Printf("Restore Save error: %v", errs)
		}
	}

	restored, err := reader.ReadStorage()
	if err != nil {
		log.Printf("Restore metrics from file '%s' ReadStorage error: %v", sf, err)
		log.Fatal(err)
	}
	f.MemStorage = restored.MemStorage
}

func (f *FileStorage) SaveTicker(storeint time.Duration) {
	ticker := time.NewTicker(storeint)
	for range ticker.C {
		errs := f.Save()
		if errs != nil {
			log.Printf("FileStorage Save error: %v", errs)
		}
	}
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewWriter(filename string) (*producer, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) Close() error {
	return p.file.Close()
}

type consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewReader(filename string) (*consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0o777)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (r *consumer) ReadStorage() (*FileStorage, error) {
	if !r.scanner.Scan() {
		return nil, r.scanner.Err()
	}

	data := r.scanner.Bytes()

	fs := FileStorage{}
	err := json.Unmarshal(data, &fs)
	if err != nil {
		return nil, err
	}

	return &fs, nil
}
