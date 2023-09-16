package filestorage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"sync"
	"testing"

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
		name    string
		fields  fields
		want    *memstorage.MemStorage
		wantErr bool
	}{
		{
			"Empty file test",
			fields{
				StoreFile: "emptyStorage.txt",
				Content:   "",
				Sync:      false,
			},
			&memstorage.MemStorage{
				Gauges:   map[string]float64{},
				Counters: map[string]int64{},
				Mux:      &sync.RWMutex{},
			},
			false,
		},
		{
			"Wrong format file test",
			fields{
				StoreFile: "badStorage.txt",
				Content:   "Some uunstructured text",
				Sync:      false,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ef, err := os.Create(tt.fields.StoreFile)
			if tt.fields.Content != "" {
				if _, err := ef.WriteString(tt.fields.Content); err != nil {
					t.Errorf("FileStorage.Restore() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			defer func() {
				if _, fileErr := os.Stat(tt.fields.StoreFile); errors.Is(fileErr, os.ErrNotExist) {
					return
				}
				ef.Close()
				os.Remove(tt.fields.StoreFile)
			}()
			if err != nil {
				t.Errorf("FileStorage.Restore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			f := &FileStorage{
				StoreFile: tt.fields.StoreFile,
				Sync:      tt.fields.Sync,
			}
			got, err := f.Restore()
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
	f, _ := os.Create("1.test")
	defer os.Remove("1.test")
	tests := []struct {
		name    string
		fields  fields
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

	f, _ := os.Create("1.test")
	tests := []struct {
		name    string
		args    args
		want    *consumer
		wantErr bool
	}{
		{
			name: "test1",
			args: args{"1.test"},
			want: &consumer{
				file:    f,
				scanner: bufio.NewScanner(f),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newReader(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("newReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.True(t, reflect.DeepEqual(got.file.Name(), tt.want.file.Name())) {
				t.Errorf("newReader() = %v, want %v.", got, tt.want)
			}
		})
	}
}
