package filestorage

import (
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
				Sync:      false,
			},
			&memstorage.MemStorage{
				Gauges:   map[string]float64{},
				Counters: map[string]int64{},
				Mux:      &sync.RWMutex{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ef, err := os.Create(tt.fields.StoreFile)
			defer func() {
				if _, fileErr := os.Stat("/path/to/whatever"); errors.Is(fileErr, os.ErrNotExist) {
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
