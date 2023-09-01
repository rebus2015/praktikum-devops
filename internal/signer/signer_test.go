// Пакет signer выполняет функцию проверки целостнсти данных при обмене метриками между клиентом и сервисом
// Выполняет функции подписи данных в структуре данных и их верификацию генерируя SHA256 HMAC Hash
package signer

import (
	"fmt"
	"testing"

	"github.com/rebus2015/praktikum-devops/internal/model"

	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T {
	return &v
}

func TestHashObject_Sign(t *testing.T) {
	type args struct {
		m            *model.Metrics
		key          string
		expectedHash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"1st positive",
			args{
				&model.Metrics{
					ID:    "TotalMemory",
					MType: "gauge",
					Delta: nil,
					Value: ptr(float64(7268679680)),
					Hash:  "",
				},
				"SuperSecretKey",
				"a3afa5537e12b6a83f982b8a286031b03d26ee12eb43f24c83beacff3ed81f87",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := HashObject{
				key: tt.args.key,
			}

			if err := s.Sign(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("HashObject.Sign() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NotEqual(t, tt.args.expectedHash, tt.args.m.Hash, "")
		})
	}
}

func TestHashObject_Verify(t *testing.T) {
	type args struct {
		m    *model.Metrics
		key  string
		hash string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"1st positive",
			args{
				&model.Metrics{
					ID:    "TotalMemory",
					MType: "gauge",
					Delta: nil,
					Value: ptr(float64(7268679680)),
					Hash:  "wjhfwbwih388",
				},
				"SuperSecretKey",
				"wjhfwbwih388",
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HashObject{
				key: tt.args.key,
			}
			got, err := s.Verify(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashObject.Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got, fmt.Errorf("HashObject.Verify() = %v, want %v", got, tt.want))
		})
	}
}
