package main

import (
	"context"
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

type Cache struct {
	kube *KubeMQ
}

func NewCahce(kube *KubeMQ) *Cache {
	c := &Cache{
		kube: kube,
	}
	return c
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiry time.Time) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	cm := &CacheMessage{
		Type:   "set",
		Key:    key,
		Value:  data,
		Expiry: expiry,
	}

	err = c.kube.SendCommandToCache(ctx, cm)
	return err
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {

	cm := &CacheMessage{
		Type: "get",
		Key:  key,
	}

	data, err := c.kube.SendQueryToCache(ctx, cm)
	return data, err
}

func (c *Cache) Del(ctx context.Context, key string) error {

	cm := &CacheMessage{
		Type: "del",
		Key:  key,
	}

	err := c.kube.SendCommandToCache(ctx, cm)
	return err
}
