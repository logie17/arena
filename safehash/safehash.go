package safehash

import (
	"sync"
)

type safeMap struct {
	myHash map[string]int
	mutex  *sync.RWMutex
}

func (sf *safeMap) Insert(key string, val int) {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()
	sf.myHash[key] = val
}

func NewSafeMap() *safeMap {
	return &safeMap{make(map[string]int), new(sync.RWMutex)}
}

func (sf *safeMap) Find(key string) int {
	sf.mutex.RLock()
	defer sf.mutex.RUnlock()
	return sf.myHash[key]
}
