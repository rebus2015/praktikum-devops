package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/rebus2015/praktikum-devops/internal/model"
)

type Sign struct {
model *model.Metrics
key string 
}

type Signer interface{
	Sign()  error
	Ensure() (bool,error)
}

func  NewSigner (m *model.Metrics, key string) Signer{
 return &Sign{
model: m,
key: key,
 }
}

func (s *Sign) Sign() error{
	var src string
	retVal:= &model.Metrics{
		ID:    s.model.ID,
		MType: s.model.MType,
	}

	switch s.model.MType {
	case "gauge":
		if s.model.Value == nil {
			log.Printf("metric of type '%v' error trying to Sign metric, model.Value == nil", s.model.MType)
			return fmt.Errorf("metric of type '%v' error trying to Sign metric, model.Value == nil", s.model.MType)
		}
		retVal.Value = s.model.Value
		src = fmt.Sprintf("%s:gauge:%f", s.model.ID, *s.model.Value)
	case "counter":
		if s.model.Delta == nil {
			log.Printf("metric of type '%v' error trying to Sign metric, model.Delta == nil", s.model.MType)
			return fmt.Errorf("metric of type '%v' error trying to Sign metric, model.Delta == nil", s.model.MType)
		}
		retVal.Delta=s.model.Delta
		src = fmt.Sprintf("%s:gauge:%d", s.model.ID, *s.model.Delta)
	default:
		log.Printf("unknown metric type exception '%v' trying to Sign metric", s.model.MType)
		return fmt.Errorf("unknown metric type exception '%v' trying to Sign metric", s.model.MType)
	}
	h := hmac.New(sha256.New, []byte(s.key))
	_, err := h.Write([]byte(src))
	if err != nil {
		log.Printf("Sign Meric failed, new hamc writer error:%v", err)
		return fmt.Errorf("Sign Meric failed, new hamc writer error:%w", err)
	}
	s.model.Hash = hex.EncodeToString(h.Sum(nil))
	
	return nil
}

func (s *Sign) Ensure() (bool,error){
	return true,nil
}