package content

import (
	"errors"
	"sync"
	"time"
)

type LocalDataStruct struct {
	Data interface{}
	TTL  int64
}

// var (
// 	LocalData = make(map[string]LocalDataStruct)
// 	lock      = sync.RWMutex{}
// )

var LocalCache LocalModelStruct = LocalModelStruct{LocalData: make(map[string]LocalDataStruct)}

type LocalModelStruct struct {
	LocalData map[string]LocalDataStruct
	Lock      sync.RWMutex
}

func init() {
	go LocalCache.initWorkerCleanCache()
}

func (this *LocalModelStruct) GetValue(key string) (interface{}, error) {
	this.Lock.Lock()
	defer this.Lock.Unlock()

	if val, ok := this.LocalData[key]; ok {
		currentTime := time.Now().Unix()
		if val.TTL <= currentTime {
			delete(this.LocalData, key)
			return "", errors.New("Not exist data local: " + key)
		}
		return val.Data, nil
	}
	return "", errors.New("Not exist data local: " + key)
}

func (this *LocalModelStruct) SetValue(key string, value interface{}, ttl int64) error {
	this.Lock.Lock()
	defer this.Lock.Unlock()

	currentTime := time.Now().Unix()
	var dataLocal LocalDataStruct
	dataLocal.Data = value
	dataLocal.TTL = currentTime + ttl

	this.LocalData[key] = dataLocal
	return nil
}
func (this *LocalModelStruct) initWorkerCleanCache() {
	for {
		time.Sleep(1 * time.Second)
		currentTime := time.Now().Unix()
		for k, v := range this.LocalData {
			if v.TTL <= currentTime {
				this.Lock.Lock()
				delete(this.LocalData, k)
				this.Lock.Unlock()
			}
		}
	}
}
