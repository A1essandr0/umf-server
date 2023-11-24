package repositories

import (
	"sync"

	"github.com/go-redis/redis/v8"
)

// using locks to be used more than just in mocks
type InmemoryKV struct {	
	store map[string]string
	mutex sync.RWMutex
}

func (kv *InmemoryKV) CreateKVStoreRecord(key, value string) error {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()

	kv.store[key] = value
	return nil
}

func (kv *InmemoryKV) GetKVStoreRecord(key string) (string, error) {
	kv.mutex.RLock()
	defer kv.mutex.RUnlock()

	value, ok := kv.store[key]
	if !ok {
		return "", redis.Nil 
	}
	return value, nil
}