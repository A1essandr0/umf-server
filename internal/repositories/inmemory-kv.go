package repositories

import "github.com/go-redis/redis/v8"

// not suitable for concurrency, therefore for anything but mocks
type InmemoryKV struct {	
	store map[string]string
}

func (kv *InmemoryKV) CreateKVStoreRecord(key, value string) error {
	kv.store[key] = value
	return nil
}

func (kv *InmemoryKV) GetKVStoreRecord(key string) (string, error) {
	value, ok := kv.store[key]
	if !ok {
		return "", redis.Nil 
	}
	return value, nil
}