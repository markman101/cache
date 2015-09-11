package main

import (
	"cache"
	"fmt"
	"time"
)

func main() {
	table := cache.CreateCacheTable()
	item1 := cache.CreateCacheItem("1", "wangshitao")
	item2 := cache.CreateCacheItem("2", "xifengming")
	item3 := cache.CreateCacheItem("3", "peibaoqing")

	table.AddItem(item1.Key, item1, true)
	table.AddItem(item2.Key, item2, true)
	table.AddItem(item3.Key, item3, true)

	for key, value := range table.Items {
		hour, min, sec := value.TimeStamp.Clock()
		fmt.Printf("key:%s,TimeStamp[%d:%d:%d]\n", key, hour, min, sec)
	}
	table.ExpireCheck(20)
	time.Sleep(60 * time.Second)
}
