package main

import (
	"encoding/json"
	"time"
)

type CacheMessage struct {
	Type   string
	Key    string
	Value  []byte
	Expiry time.Time
}

func (cm *CacheMessage) Data() []byte {
	data, _ := json.Marshal(cm)
	return data
}

func NewCacheMessage(data []byte) (*CacheMessage, error) {
	cm := &CacheMessage{}
	err := json.Unmarshal(data, cm)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
