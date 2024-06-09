package utils

import (
	"sync"
	"time"
)

type ExpiringSet struct {
	//data     map[string]time.Time
	data *sync.Map

	//mutex    sync.RWMutex
	expireIn time.Duration
}

func NewExpiringSet(expireIn time.Duration) *ExpiringSet {
	set := &ExpiringSet{
		data:     &sync.Map{},
		expireIn: expireIn,
	}
	go set.startCleanup()
	return set
}

func (s *ExpiringSet) Add(key string) {
	//s.mutex.Lock()
	//defer s.mutex.Unlock()
	//s.data[key] = time.Now().Add(s.expireIn)
	s.data.Store(key, time.Now().Add(s.expireIn))

}

// Modify the function to support a string slice and add each element to the data structure
func (s *ExpiringSet) AddSlice(keys []string) {
	//s.mutex.Lock()
	//defer s.mutex.Unlock()
	//里面的值就是用作key 而不是index
	for _, key := range keys {
		//s.data[key] = time.Now().Add(s.expireIn)
		//s.data[key] = time.Now().Add(s.expireIn)
		s.data.Store(key, time.Now().Add(s.expireIn))
	}
}
func (s *ExpiringSet) Contains(key string) bool {
	//s.mutex.RLock()
	//
	//defer s.mutex.RUnlock()
	expiration, ok := s.data.Load(key)
	if !ok {
		return false
	}
	return time.Now().Before(expiration.(time.Time))
}

func (s *ExpiringSet) startCleanup() {
	ticker := time.NewTicker(s.expireIn)
	for range ticker.C {
		//s.mutex.Lock()
		//for key, expiration := range s.data {
		//	if time.Now().After(expiration) {
		//		delete(s.data, key)
		//	}
		//}
		//s.mutex.Unlock()

		// 使用 Range 方法遍历 sync.Map
		s.data.Range(func(key, value any) bool {
			//fmt.Printf("Key: %v, Value: %v\n", key, value)
			if time.Now().After(value.(time.Time)) {
				s.data.Delete(key)
			}
			return true // 返回 true 继续遍历，返回 false 停止遍历
		})
	}
}
