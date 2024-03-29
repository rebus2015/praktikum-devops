// Package signer выполняет функцию проверки целостнсти данных при обмене метриками между клиентом и сервисом
// выполняет функции подписи данных в структуре данных и их верификацию генерируя SHA256 HMAC Hash.
package signer

import (
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/rebus2015/praktikum-devops/internal/model"
)

// HashObject подпись.
type HashObject struct {
	key string
}

// NewHashObject creation.
func NewHashObject(key string) HashObject {
	h := HashObject{key: key}
	return h
}

// Sign формирование подписи для метрики.
func (s *HashObject) Sign(m *model.Metrics) error {
	src, err := srcString(m)
	if err != nil {
		return err
	}
	h, err := hash(src, s.key)
	if err != nil {
		return err
	}
	m.Hash = h
	return nil
}

// Verify проверка целостности пришедших данных.
func (s *HashObject) Verify(m *model.Metrics) (bool, error) {
	src, err := srcString(m)
	if err != nil {
		return false, err
	}
	h, err := hash(src, s.key)
	if err != nil {
		return false, err
	}
	return m.Hash == h, nil
}

// hash формирование  hash shá56 от указанной строки с ключом key.
func hash(src string, key string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))
	_, err := h.Write([]byte(src))
	if err != nil {
		log.Printf("Sign Meric failed, new hamc writer error:%v", err)
		return "", fmt.Errorf("Sign Meric failed, new hamc writer error:%w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// srcString получение текстового предстваления метрики для создания хэша.
func srcString(model *model.Metrics) (string, error) {
	switch model.MType {
	case "gauge":
		if model.Value == nil {
			log.Printf("metric of type '%v' error trying to Sign metric, model.Value == nil", model.MType)
			return "", fmt.Errorf("gauge '%v' error trying to Sign metric, model.Value == nil", model.MType)
		}
		return fmt.Sprintf("%s:%v:%f", model.ID, model.MType, *model.Value), nil
	case "counter":
		if model.Delta == nil {
			log.Printf("metric of type '%v' error trying to Sign metric, model.Delta == nil", model.MType)
			return "", fmt.Errorf("counter '%v' error trying to Sign metric, model.Delta == nil", model.MType)
		}
		return fmt.Sprintf("%s:%v:%d", model.ID, model.MType, *model.Delta), nil
	default:
		log.Printf("unknown metric type exception '%v' trying to Sign metric", model.MType)
		return "", fmt.Errorf("unknown metric type exception '%v' trying to Sign metric", model.MType)
	}
}

func DecryptMessage(key *rsa.PrivateKey, msg []byte) ([]byte, error) {
	size := key.PublicKey.Size()
	if len(msg)%size != 0 {
		return nil, errors.New("message length error")
	}
	hash := sha256.New()
	dectipted := make([]byte, 0)
	for i := 0; i < len(msg); i += size {
		data, err := rsa.DecryptOAEP(hash, nil, key, msg[i:i+size], []byte(""))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt message body err: %w", err)
		}
		dectipted = append(dectipted, data...)
	}
	return dectipted, nil
}
