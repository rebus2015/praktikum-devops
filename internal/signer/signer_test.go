// Пакет signer выполняет функцию проверки целостнсти данных при обмене метриками между клиентом и сервисом
// Выполняет функции подписи данных в структуре данных и их верификацию генерируя SHA256 HMAC Hash
package signer

import (
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
				"wjhfwbwih388",
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
	type fields struct {
		key string
	}
	type args struct {
		m *model.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HashObject{
				key: tt.fields.key,
			}
			got, err := s.Verify(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashObject.Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HashObject.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}
