package cache

import (
	"log"
	"wbL0/server_sub/entity"
)

type CacheStore map[string]entity.OrderData

type CacheService struct {
	CacheStore CacheStore
}

func CacheInit() *CacheService {
	Cs := make(CacheStore)
	CacheService := CacheService{
		CacheStore: Cs,
	}
	return &CacheService
}

func (Cservice *CacheService) AddToCache(data entity.OrderData) {
	Cservice.CacheStore[data.OrderUid] = data
	log.Println("new data in cache stored: ", data)
}

func (Cservice *CacheService) GetFromCache(order_uid string) entity.OrderData {
	return Cservice.CacheStore[order_uid]
}
