package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/rebus2015/praktikum-devops/internal/model"
)

type HashObject struct {
	key string
}

type Signer interface {
	Sign(m *model.Metrics) error
	Verify(m *model.Metrics) (bool, error)
}

func NewHashObject(key string) Signer {
	h := HashObject{key: key}
	return &h
}

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

func hash(src string, key string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))
	_, err := h.Write([]byte(src))
	if err != nil {
		log.Printf("Sign Meric failed, new hamc writer error:%v", err)
		return "", fmt.Errorf("Sign Meric failed, new hamc writer error:%w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func srcString(model *model.Metrics) (string, error) {
	switch model.MType {
	case "gauge":
		if model.Value == nil {
			log.Printf("metric of type '%v' error trying to Sign metric, model.Value == nil", model.MType)
			return "", fmt.Errorf("metric of type '%v' error trying to Sign metric, model.Value == nil", model.MType)
		}
		return fmt.Sprintf("%s:%v:%f", model.ID, model.MType, *model.Value), nil
	case "counter":
		if model.Delta == nil {
			log.Printf("metric of type '%v' error trying to Sign metric, model.Delta == nil", model.MType)
			return "", fmt.Errorf("metric of type '%v' error trying to Sign metric, model.Delta == nil", model.MType)
		}
		return fmt.Sprintf("%s:%v:%d", model.ID, model.MType, *model.Delta), nil
	default:
		log.Printf("unknown metric type exception '%v' trying to Sign metric", model.MType)
		return "", fmt.Errorf("unknown metric type exception '%v' trying to Sign metric", model.MType)
	}
}
