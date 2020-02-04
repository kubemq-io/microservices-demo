package main

import (
	"errors"
	"github.com/go-redis/redis"
	"time"
)

type RedisMetadata struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type Redis struct {
	client *redis.Client
}

func NewRedisClient(url string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	r := &Redis{
		client: client,
	}
	return r, nil
}
func (r *Redis) Set(key string, value []byte, exp time.Time) error {
	err := r.client.Set(key, value, exp.Sub(time.Now())).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Get(key string) ([]byte, error) {
	result, err := r.client.Get(key).Result()
	if err == redis.Nil {
		return nil, errors.New("key not found")
	}
	if err != nil {
		return nil, err
	}
	return []byte(result), nil

}
func (r *Redis) Del(key string) error {
	_, err := r.client.Del(key).Result()
	if err == redis.Nil {
		return errors.New("key not found")
	}
	if err != nil {
		return err
	}

	return nil

}
