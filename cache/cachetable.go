package cache

import (
	"fmt"
	"sync"
	"time"
)

const (
	ITEM_EXIST_NO_UPDATE = iota
	ITEM_EXIST_UPDATE
	ADD_NEW_ITEM_SUCC
	TYPE_NUM
)

//set other name
type CacheMap map[interface{}]*CacheItem

type CacheTable struct {
	sync.RWMutex
	Items    CacheMap
	LifeSpan time.Duration
	//expire   time.Timer
}

func CreateCacheTable() *CacheTable {
	table := make(CacheMap)
	return &CacheTable{Items: table, LifeSpan: 0}
}

func (table *CacheTable) AddItem(key interface{}, item *CacheItem, update bool) (ret int) {
	ret = TYPE_NUM
	table.RWMutex.Lock()
	defer table.RWMutex.Unlock()

	value, ok := table.Items[key]
	if ok == true {
		//the value has exist
		if !update {
			//the update flag is false.so not update the value.
			value.TimeStamp = time.Now()
			ret = ITEM_EXIST_NO_UPDATE
		} else {
			item.TimeStamp = time.Now()
			table.Items[key] = item
			ret = ITEM_EXIST_UPDATE
		}
	} else {
		//the value no exist,so added it
		item.TimeStamp = time.Now()
		table.Items[key] = item
		ret = ADD_NEW_ITEM_SUCC
	}

	return
}

func (table *CacheTable) RmItem(key interface{}) (ret bool) {
	ret = false
	table.RWMutex.Lock()
	defer table.RWMutex.Unlock()

	_, ok := table.Items[key]
	if ok == true {
		delete(table.Items, key)
		ret = true
	}

	return
}

func (table *CacheTable) GetItem(key interface{}) (ret bool, item CacheItem) {
	ret = false
	table.RWMutex.RLock()
	defer table.RWMutex.RUnlock()

	value, ok := table.Items[key]
	if ok == true {
		ret = true
		item = *value
	}
	return
}

//timeou is second
func (table *CacheTable) ExpireCheck(timeout time.Duration) {
	tmp_items := make(CacheMap)

	//get all expire item
	curTime := time.Now()
	table.RWMutex.RLock()
	for key, value := range table.Items {
		//curTime is Nanosecond.1 Second=10的9次方Nanosecond
		if curTime.Sub(value.TimeStamp)/1000/1000/1000 > timeout {
			hour, min, sec := curTime.Clock()
			fmt.Printf("key:%s,curTime[%d:%d:%d] dif[%d]\n", key, hour, min, sec, curTime.Sub(value.TimeStamp))
			tmp_items[key] = value
		}
	}
	table.RWMutex.RUnlock()

	fmt.Printf("get invalid cache num[%d]\n", len(tmp_items))
	//delete itme in WRITE_LOCK
	table.RWMutex.Lock()
	for key, _ := range tmp_items {
		delete(table.Items, key)
	}
	table.RWMutex.Unlock()

	//set timer
	time.AfterFunc(3*time.Second, func() {
		table.ExpireCheck(timeout)
	})
}
