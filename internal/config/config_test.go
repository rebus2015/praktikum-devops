// Пакет config выполняет функцию параметризации сервиса сбора метрик
// Поддерживает задание параметров запуска через переменные окружения и параметры командной строки.
package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		{
			"test 1",
			&Config{
				ServerAddress:    "127.0.0.1:8080",
				StoreInterval:    time.Second * 30,
				StoreFile:        "",
				Restore:          false,
				Key:              "",
				ConnectionString: "",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.True(t, reflect.DeepEqual(got, tt.want))
		})
	}
}
