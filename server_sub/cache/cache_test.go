package cache

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"wbL0/server_sub/entity"
)

func TestCacheInit(t *testing.T) {
	cache := CacheInit()
	require.NotNil(t, cache)
}

func TestCacheService_AddToCache(t *testing.T) {
	cache := CacheInit()
	require.NotNil(t, cache)
	fakeCache := make(CacheStore)
	text, err := os.ReadFile("../d4test/t4data.json")
	if err != nil {
		t.Fatal("error reading testdata file")
	}
	od := new(entity.OrderData)
	err = od.Scan(text)
	fakeCache[od.OrderUid] = *od
	cache.AddToCache(*od)
	realC := cache.CacheStore
	require.Equal(t, realC, fakeCache)
}

func TestCacheService_GetFromCache(t *testing.T) {
	cache := CacheInit()
	text, err := os.ReadFile("../d4test/t4data.json")
	if err != nil {
		t.Fatal("error reading testdata file")
	}
	od := new(entity.OrderData)
	err = od.Scan(text)
	cache.CacheStore[od.OrderUid] = *od
	fakeCache := make(CacheStore)
	fakeCache[od.OrderUid] = *od
	realOr := cache.GetFromCache(od.OrderUid)
	fakeOr := fakeCache[od.OrderUid]
	require.Equal(t, realOr, fakeOr)
}
