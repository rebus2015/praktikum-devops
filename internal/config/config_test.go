// Пакет config выполняет функцию параметризации сервиса сбора метрик
// Поддерживает задание параметров запуска через переменные окружения и параметры командной строки.
package config

import (
	"crypto/rsa"
	"fmt"
	"os"
	"reflect"
	"strings"
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

func TestConfig_UnmarshalJSON(t *testing.T) {
	type fields struct {
		ServerAddress    string
		StoreInterval    time.Duration
		StoreFile        string
		Restore          bool
		Key              string
		ConnectionString string
		CryptoKeyFile    string
		confFile         string
		CryptoKey        *rsa.PrivateKey
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "positive test 1",
			args: args{
				data: []byte(`{
				"address": "localhost:8080", 
				"restore": true, 
				"store_interval": "1s", 
				"store_file": "/path/to/file.db", 
				"database_dsn": "", 
				"crypto_key": "/path/to/key.pem" 
			}`),
			},
			fields: fields{
				ServerAddress:    "localhost:8080",
				StoreInterval:    time.Second * 1,
				StoreFile:        "/path/to/file.db",
				Restore:          true,
				ConnectionString: "",
				CryptoKeyFile:    "/path/to/key.pem",
			},
			wantErr: false,
		},
		{
			name: "negative test 1",
			args: args{
				data: []byte(`{
				"address": "localhost:8080", 
				"restore": true, 
				"store_interval": "day", 
				"store_file": "/path/to/file.db", 
				"database_dsn": "", 
				"crypto_key": "/path/to/key.pem" 
			}`),
			},
			fields: fields{
				ServerAddress:    "localhost:8080",
				StoreInterval:    time.Second * 1,
				StoreFile:        "/path/to/file.db",
				Restore:          true,
				ConnectionString: "",
				CryptoKeyFile:    "/path/to/key.pem",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				ServerAddress:    tt.fields.ServerAddress,
				StoreInterval:    tt.fields.StoreInterval,
				StoreFile:        tt.fields.StoreFile,
				Restore:          tt.fields.Restore,
				Key:              tt.fields.Key,
				ConnectionString: tt.fields.ConnectionString,
				CryptoKeyFile:    tt.fields.CryptoKeyFile,
				confFile:         tt.fields.confFile,
				CryptoKey:        tt.fields.CryptoKey,
			}
			var err error
			if err = c.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Config.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.fields.ServerAddress, c.ServerAddress)
			assert.Equal(t, tt.fields.StoreInterval, c.StoreInterval)
			assert.Equal(t, tt.fields.StoreFile, c.StoreFile)
			assert.Equal(t, tt.fields.Restore, c.Restore)
			assert.Equal(t, tt.fields.Key, c.Key)
			assert.Equal(t, tt.fields.ConnectionString, c.ConnectionString)
			assert.Equal(t, tt.fields.CryptoKeyFile, c.CryptoKeyFile)
			assert.Equal(t, tt.fields.confFile, c.confFile)
			assert.Equal(t, tt.fields.CryptoKey, c.CryptoKey)
		})
	}
}

func TestConfig_parseConfigFile(t *testing.T) {
	type fields struct {
		ServerAddress    string
		StoreInterval    time.Duration
		StoreFile        string
		Restore          bool
		Key              string
		ConnectionString string
		CryptoKeyFile    string
		confFile         string
		CryptoKey        *rsa.PrivateKey
	}
	tests := []struct {
		name    string
		f       fields
		wantf   fields
		data    string
		wantErr bool
	}{
		{
			name:    "negative test: no file",
			wantErr: false,
		},
		{
			name:    "negative test: bad file",
			wantErr: true,
		},
		{
			name: "positive test 1",
			f: fields{
				ServerAddress:    "localhost:8080",
				StoreInterval:    time.Second * 1,
				StoreFile:        "/path/to/file.db",
				Restore:          true,
				ConnectionString: "",
				CryptoKeyFile:    "",
				confFile:         "",
			},
			wantf: fields{
				ServerAddress:    "localhost:8080",
				StoreInterval:    time.Second * 1,
				StoreFile:        "/path/to/file.db",
				Restore:          true,
				ConnectionString: "new value",
				CryptoKeyFile:    "/path/to/key.pem",
				confFile:         "",
			},
			data: `{
				"address": "localhost:8080", 
				"restore": true, 
				"store_interval": "1s", 
				"store_file": "/path/to/file.db", 
				"database_dsn": "new value", 
				"crypto_key": "/path/to/key.pem" 
			}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile, errt := os.CreateTemp("", "conf.json")
			if errt != nil {
				fmt.Println("Error creating temporary file:", errt)
				return
			}

			defer os.Remove(tempFile.Name())
			_, errw := tempFile.Write([]byte(tt.data))
			if errw != nil {
				fmt.Println("Error writing to temporary file:", errw)
				return
			}

			c := &Config{
				ServerAddress:    tt.f.ServerAddress,
				StoreInterval:    tt.f.StoreInterval,
				StoreFile:        tt.f.StoreFile,
				Restore:          tt.f.Restore,
				Key:              tt.f.Key,
				ConnectionString: tt.f.ConnectionString,
				CryptoKeyFile:    tt.f.CryptoKeyFile,
				confFile:         tt.f.confFile,
				CryptoKey:        tt.f.CryptoKey,
			}
			if strings.Contains(tt.name, "no file") {
				c.confFile = ""
			}
			if tt.wantErr && strings.Contains(tt.name, "bad file") {
				c.confFile = "some/bad/file.json"
			}
			if !tt.wantErr {
				c.confFile = tempFile.Name()
			}
			if strings.Contains(tt.name, "no file") {
				c.confFile = ""
			}
			var err error
			if err = c.parseConfigFile(); (err != nil) != tt.wantErr {
				t.Errorf("Config.parseConfigFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantf.ServerAddress, c.ServerAddress)
				assert.Equal(t, tt.wantf.StoreInterval, c.StoreInterval)
				assert.Equal(t, tt.wantf.StoreFile, c.StoreFile)
				assert.Equal(t, tt.wantf.Restore, c.Restore)
				assert.Equal(t, tt.wantf.Key, c.Key)
				assert.Equal(t, tt.wantf.ConnectionString, c.ConnectionString)
				assert.Equal(t, tt.wantf.CryptoKeyFile, c.CryptoKeyFile)
				assert.Equal(t, tt.wantf.CryptoKey, c.CryptoKey)
			}

		})
	}
}

const rsaKey string = `-----BEGIN rsa private key-----
MIICXgIBAAKBgQC3swOLvQsBKQcb6o/medAanjuv3PMgK1jYvSD31IcMvU9N0nOv
qYgytSloHlu1GyKURBs5Vg3xVRsVaVgeRQnBtUEtaJNhT3+Qqe0Pynv+ifx2YidV
LyqgTzG+hHGp9uH2JAupkAqLzYwpw/rulOt2dq1J6hThT937FvpoH/KbBwIDAQAB
AoGBAIeeof+IkZdJsvXpNlPxmrIMIAS2GsilN/LLrotJXGsLWIEb3kzR3LuTA/7a
atpKLj1ICtFJtwF004n7PBMc5RXbeYQ5Sdb/6pYvlRWtkm4jaUZaB63SnbzL6RPp
NwCJjVGdJOdUG6yXeuHrvuTsbo/hlgoYWdr80Lb7Jf9ex6gBAkEA8tktzAU2+RkK
ue6elDUiytWk/XbEK3upEDp1L1FIJE2ioYRtrOxEEcBVbYXL+YkipL9rLf3C4k74
0STvprij1QJBAMGlzIyQyccnbBBF3TgYqKsMhdzTvemMWDX2Giq9nxxTgVaOqPYQ
u0W+/yN9gzlrJrKTOer4CEwkaCLNX2PbHWsCQQDZ6QVOODOu68iTNMo5JUD2DyVA
hyzZ89mthTcX4XDBmqRfGIytiUg/QX2mjFOOs35RpK4RE86m8cQVL3aX/MCNAkEA
qN2qeFGyg6cPB0nFVau7Oh4bhaxoCgfGzJelzeu5mnv/Z7nUAXApvvKFjy9ehW25
OzRD53EP20ZMQT0SmAN1rQJAQp949K4gjYLoIjmNKKulRix6HLDCVMsp132w/Hjv
kVpKJJ4IyDiPV8+6A3VAlxvSpjVusaR7Jq8VDiHBKEf5jQ==
-----END rsa private key-----`

func TestConfig_getCryptoKey(t *testing.T) {
	type fields struct {
		ServerAddress    string
		StoreInterval    time.Duration
		StoreFile        string
		Restore          bool
		Key              string
		ConnectionString string
		CryptoKeyFile    string
		confFile         string
		CryptoKey        *rsa.PrivateKey
	}
	tests := []struct {
		name    string
		data    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "negative: no file set",
			fields:  fields{},
			wantErr: false,
		},
		{
			name: "negative: file not exists",
			fields: fields{
				CryptoKeyFile: "/file/not/exist",
			},
			wantErr: true,
		},
		{
			name: "negative: bad file",
			fields: fields{
				CryptoKeyFile: "/file/not/exist",
			},
			data:    "some not valid data",
			wantErr: true,
		},
		{
			name: "positive",
			data: rsaKey,
			fields: fields{
				CryptoKeyFile: "/file/not/exist",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			var err error
			tempFile, errt := os.CreateTemp("", "key.pem")
			if errt != nil {
				fmt.Println("Error creating temporary file:", errt)
				return
			}

			defer os.Remove(tempFile.Name())
			_, errw := tempFile.Write([]byte(tt.data))
			if errw != nil {
				fmt.Println("Error writing to temporary file:", errw)
				return
			}

			if tt.wantErr && strings.Contains(tt.name, "bad file") {
				tt.fields.CryptoKeyFile = tempFile.Name()
			}

			if !tt.wantErr {
				tt.fields.CryptoKeyFile = tempFile.Name()
			}

			if !tt.wantErr && strings.Contains(tt.name, "no file") {
				tt.fields.CryptoKeyFile = ""
			}

			c := &Config{
				CryptoKeyFile: tt.fields.CryptoKeyFile,
				CryptoKey:     tt.fields.CryptoKey,
			}
			if err = c.getCryptoKey(); (err != nil) != tt.wantErr {
				t.Errorf("Config.getCryptoKey() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if tt.fields.CryptoKeyFile == "" {
				assert.True(t, c.CryptoKey == nil)
				return
			}
			assert.True(t, c.CryptoKey != nil)
		})
	}
}
