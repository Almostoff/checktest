package store

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"wbL0/server_sub/cache"
	"wbL0/server_sub/db"
	"wbL0/server_sub/entity"
)

type StoreService struct {
	cache cache.CacheService
	db    db.DBService
}

func InitStore(cache cache.CacheService, db db.DBService) *StoreService {
	StoreService := StoreService{
		cache: cache,
		db:    db,
	}
	return &StoreService
}

func (ss *StoreService) SaveOrderData(data []byte) error {
	od := new(entity.OrderData)
	err := od.Scan(data)
	if err != nil {
		log.Println("Wrong format")
		return err
	}
	validate := validator.New()
	err = validate.Struct(od)
	if err != nil {
		log.Println(err)
		return err
	}
	itemData := new(entity.DataItem)
	itemData.OrderData = *od
	itemData.ID = od.OrderUid
	ss.cache.AddToCache(*od)
	_, err = ss.db.SaveOrder(itemData)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func (ss *StoreService) GetFromCacheByUID(id string) entity.OrderData {
	return ss.cache.GetFromCache(id)
}

func (ss *StoreService) GetAllOrders() ([]entity.DataItem, error) {
	di, err := ss.db.GetAllOrders()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return di, err
}

func (ss *StoreService) RestoreCache() error {
	dItems, err := ss.GetAllOrders()
	if dItems == nil {
		log.Println(err)
		return err
	}
	for _, dItem := range dItems {
		ss.cache.AddToCache(dItem.OrderData)
	}
	log.Println("--CACHE IS RESTORED--")
	return err
}
