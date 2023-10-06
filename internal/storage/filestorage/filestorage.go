// Package filestorage реализует механизм хранения метрик в текстовом файле.
package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/rebus2015/praktikum-devops/internal/config"
	"github.com/rebus2015/praktikum-devops/internal/storage"
	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

type FileStorage struct {
	StoreFile string
	Sync      bool
}

var _ storage.SecondaryStorage = new(FileStorage)

func NewStorage(ctx context.Context, c *config.Config) *FileStorage {
	return &FileStorage{
		c.StoreFile,
		c.StoreInterval == 0,
	}
}

func (f *FileStorage) SyncMode() bool {
	return f.Sync
}

func (f *FileStorage) Save(ctx context.Context, ms *memstorage.MemStorage) error {
	writer, err := newWriter(f.StoreFile)
	if err != nil {
		log.Printf("error FileStorage save metrics to file '%s' error: %v", f.StoreFile, err)
		return fmt.Errorf("error FileStorage save metrics to file '%s' error:%w", f.StoreFile, err)
	}

	err = writer.encoder.Encode(ms)
	if err != nil {
		log.Printf("writer.encoder.Encode from file: %v, error: %v", f.StoreFile, err)
		return fmt.Errorf("error encode from file '%s' error: %w", f.StoreFile, err)
	}
	return nil
}

func (f *FileStorage) Restore(ctx context.Context) (*memstorage.MemStorage, error) {
	reader, err := newReader(f.StoreFile)
	if err != nil {
		log.Printf("Restore metrics from file '%s' reader error: %v", f.StoreFile, err)
		return nil, fmt.Errorf("error restore metrics from file '%s' reader error: %w", f.StoreFile, err)
	}

	checkFile, err := os.Stat(f.StoreFile)
	if err != nil {
		log.Printf("Restore metrics from file '%s' Stat error: %v", f.StoreFile, err)
		return nil, fmt.Errorf("Restore metrics from file '%s' reader error: %w", f.StoreFile, err)
	}

	size := checkFile.Size()
	if size == 0 {
		log.Printf("Metrics store file '%s' is emmpty", f.StoreFile)
		return memstorage.NewStorage(), nil
	}
	restored, err := reader.readStorage()
	if err != nil {
		log.Printf("Restore metrics from file '%s' ReadStorage error: %v", f.StoreFile, err)
		return nil, fmt.Errorf("Restore metrics from file '%s' ReadStorage error: %w", f.StoreFile, err)
	}
	return restored, nil
}

func (f *FileStorage) SaveTicker(storeint time.Duration, ms *memstorage.MemStorage) {
	ticker := time.NewTicker(storeint)
	for range ticker.C {
		errs := f.Save(context.Background(), ms)
		if errs != nil {
			log.Printf("FileStorage Save error: %v", errs)
		}
	}
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func newWriter(filename string) (*producer, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("newWriter error:%w", err)
	}

	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) Close() error {
	err := p.file.Close()
	if err != nil {
		return fmt.Errorf("error on close :%w", err)
	}
	return nil
}

type consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func newReader(filename string) (*consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("error on os.OpenFile :%w", err)
	}

	return &consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (r *consumer) readStorage() (*memstorage.MemStorage, error) {
	if !r.scanner.Scan() {
		return nil, fmt.Errorf("error r.scanner.Scan() :%w", r.scanner.Err())
	}

	data := r.scanner.Bytes()

	ms := memstorage.MemStorage{}
	err := json.Unmarshal(data, &ms)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error:%w", r.scanner.Err())
	}
	ms.Mux = &sync.RWMutex{}
	return &ms, nil
}
