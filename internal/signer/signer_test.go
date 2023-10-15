// Пакет signer выполняет функцию проверки целостнсти данных при обмене метриками между клиентом и сервисом
// Выполняет функции подписи данных в структуре данных и их верификацию генерируя SHA256 HMAC Hash.
package signer

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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

func Test_hash(t *testing.T) {
	type args struct {
		src string
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"1st positive",
			args{
				src: "srcStringforTest",
				key: "superSecretKey",
			},
			"f32c0c974fcd64152114df2098df1b77694669fc25338e8b0c8b88122da98bf4",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hash(tt.args.src, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("hash() = %v, want %v", got, tt.want)
			}
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_srcString(t *testing.T) {
	type args struct {
		model *model.Metrics
	}
	tests := []struct {
		name    string
		args    args
		want    string
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
			},
			"TotalMemory:gauge:7268679680.000000",
			false,
		},
		{
			"1st negative",
			args{
				&model.Metrics{
					ID:    "TotalMemory",
					MType: "gauge",
					Delta: nil,
					Value: nil,
					Hash:  "wjhfwbwih388",
				},
			},
			"",
			true,
		},
		{
			"2nd positive",
			args{
				&model.Metrics{
					ID:    "TotalMemory",
					MType: "counter",
					Delta: ptr(int64(1)),
					Value: nil,
					Hash:  "wjhf383wih388",
				},
			},
			"TotalMemory:counter:1",
			false,
		},
		{
			"2nd negative",
			args{
				&model.Metrics{
					ID:    "TotalMemory",
					MType: "counter",
					Delta: nil,
					Value: nil,
					Hash:  "wjhf383wih388",
				},
			},
			"",
			true,
		},
		{
			"3nd negative",
			args{
				&model.Metrics{
					ID:    "TotalMemory",
					MType: "unk",
					Delta: nil,
					Value: nil,
					Hash:  "wjhf383wih388",
				},
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srcString(tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("srcString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("srcString() = %v, want %v", got, tt.want)
			}
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestDecryptMessage(t *testing.T) {
	// Generate a key pair for testing
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err, "Failed to generate key pair")

	// Encrypt a test message using the public key
	message := []byte("test message")
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &key.PublicKey, message, []byte(""))
	assert.NoError(t, err, "Failed to encrypt test message")

	// Decrypt the encrypted message
	decrypted, err := DecryptMessage(key, encrypted)

	// Assertions
	assert.NoError(t, err, "DecryptMessage should not return an error")
	assert.Equal(t, message, decrypted, "Decrypted message should match the original message")
}
