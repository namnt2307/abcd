package live_event

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
	LocalData 	map[string]LocalDataStruct
	Lock      	sync.RWMutex
}


// func init() {
// 	go initWorkerCleanCache()
// }

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