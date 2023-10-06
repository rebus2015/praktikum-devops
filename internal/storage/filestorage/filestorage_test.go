package filestorage

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/rebus2015/praktikum-devops/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
)

func TestFileStorage_Restore(t *testing.T) {
	type fields struct {
		StoreFile string
		Content   string
		Sync      bool
	}
	tests := []struct {
		want    *memstorage.MemStorage
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Empty file test",
			fields: fields{
				StoreFile: "emptyStorage.txt",
				Content:   "",
				Sync:      false,
			},
			want: &memstorage.MemStorage{
				Gauges:   map[string]float64{},
				Counters: map[string]int64{},
				Mux:      &sync.RWMutex{},
			},
			wantErr: false,
		},
		{
			name: "Bad format file test",
			fields: fields{
				StoreFile: "badStorage.txt",
				Content:   "Some uunstructured text",
				Sync:      false,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := testing.TB.TempDir(t)
			f, err := os.CreateTemp(dir, tt.fields.StoreFile)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				err := os.Remove(f.Name())
				if err != nil {
					t.Errorf("os.Remove(f.Name()) error = %v", err)
				}
				err = f.Close()
				if err != nil {
					t.Errorf("f.Close() error = %v", err)
				}
			}()

			if tt.fields.Content != "" {
				if _, err := f.WriteString(tt.fields.Content); err != nil {
					t.Errorf("FileStorage.Restore() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if err != nil {
				t.Errorf("FileStorage.Restore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fs := &FileStorage{
				StoreFile: f.Name(),
				Sync:      tt.fields.Sync,
			}
			got, err := fs.Restore(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStorage.Restore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.True(t, reflect.DeepEqual(got, tt.want))
		})
	}
}

func Test_producer_Close(t *testing.T) {
	type fields struct {
		file    *os.File
		encoder *json.Encoder
	}
	dir := testing.TB.TempDir(t)
	f, err := os.CreateTemp(dir, "1.test")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.Remove(f.Name())
		if err != nil {
			log.Fatal("Test_producer_Close:os.Remover error: %w", err)
		}
	}()

	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				file:    f,
				encoder: json.NewEncoder(f),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &producer{
				file:    tt.fields.file,
				encoder: tt.fields.encoder,
			}
			if err := p.Close(); (err != nil) != tt.wantErr {
				t.Errorf("producer.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newReader(t *testing.T) {
	type args struct {
		filename string
	}

	tests := []struct {
		want    *consumer
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test 1",
			args:    args{"1.test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := testing.TB.TempDir(t)
			f, err := os.CreateTemp(dir, tt.args.filename)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err := os.Remove(f.Name()); err != nil {
					log.Fatal("remove file error:%w", err)
				}
				if err := f.Close(); err != nil {
					log.Fatal("close file error:%w", err)
				}
			}()
			got, err := newReader(f.Name())
			if (err != nil) != tt.wantErr {
				t.Errorf("newReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.True(t, reflect.DeepEqual(got.file.Name(), f.Name())) {
				t.Errorf("newReader() = %v, want %v.", got.file.Name(), f.Name())
			}
		})
	}
}

func TestSaveTicker(t *testing.T) {
	// Create a mock storage with sample data
	ms := &memstorage.MemStorage{
		Gauges: map[string]float64{
			"metric1": 10.5,
			"metric2": 20.0,
		},
		Counters: map[string]int64{
			"metric3": 30,
			"metric4": 40,
		},
	}

	// Use a shorter ticker duration for testing
	storeInterval := 100 * time.Millisecond
	storeFile := "tempstore"
	// Create the FileStorage instance and call the SaveTicker function
	f := &FileStorage{StoreFile: storeFile}
	go f.SaveTicker(storeInterval, ms)

	// Wait for some time to allow the ticker to trigger
	time.Sleep(500 * time.Millisecond)

	// Stop the ticker
	tickerStop := make(chan bool)
	go func() {
		time.Sleep(200 * time.Millisecond)
		tickerStop <- true
	}()
	defer func() {
		err := os.Remove(storeFile)
		if err != nil {
			t.Errorf("os.Remove() error = %v", err)
		}
	}()
	assert.FileExists(t, storeFile)
}
