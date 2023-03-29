package filestorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rebus2015/praktikum-devops/internal/config"
	"github.com/rebus2015/praktikum-devops/internal/storage"
	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
)

type FileStorage struct {
	StoreFile string
	Sync      bool
}

var _ storage.SecondaryStorage = new(FileStorage)

func NewStorage(c *config.Config) *FileStorage {
	return &FileStorage{
		c.StoreFile,
		c.StoreInterval == 0,
	}
}

func (f *FileStorage) SyncMode() bool {
	return f.Sync
}

func (f *FileStorage) Save(ms *memstorage.MemStorage) error {
	writer, err := newWriter(f.StoreFile)
	if err != nil {
		log.Printf("Save metrics to file '%s' error: %v", f.StoreFile, err)
		return fmt.Errorf("Save metrics to file '%s' error:%w", f.StoreFile, err)
	}

	err = writer.encoder.Encode(ms)
	if err != nil {
		log.Printf("Save metrics to file '%s' Encode error: %v", f.StoreFile, err)
		return err
	}
	return nil
}

func (f *FileStorage) Restore() (*memstorage.MemStorage, error) {
	reader, err := newReader(f.StoreFile)
	if err != nil {
		log.Printf("Restore metrics from file '%s' reader error: %v", f.StoreFile, err)
		return nil, fmt.Errorf("Restore metrics from file '%s' reader error: %w", f.StoreFile, err)
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
		errs := f.Save(ms)
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

func newReader(filename string) (*consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0o777)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (r *consumer) readStorage() (*memstorage.MemStorage, error) {
	if !r.scanner.Scan() {
		return nil, r.scanner.Err()
	}

	data := r.scanner.Bytes()

	ms := memstorage.MemStorage{}
	err := json.Unmarshal(data, &ms)
	if err != nil {
		return nil, err
	}

	return &ms, nil
}
