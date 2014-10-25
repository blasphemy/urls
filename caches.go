package main

import (
	"github.com/blasphemy/cache"
	"time"
)

var (
	StatsCache *cache.Cache
	UrlCache   *cache.Cache
)

func initCaches() {
	SCOP := cache.CacheOptions{}
	SCOP.ExpirationTime = time.Minute * 5
	SCOP.Upper = 20
	SCOP.BurnStrategy = cache.BurnStrategyOldest
	StatsCache = cache.NewCache(SCOP)
	StatsCache.Start()
	UCOP := cache.CacheOptions{}
	UCOP.ExpirationTime = time.Hour * 1
	UCOP.Upper = 1000
	UCOP.MaxEntries = 0
	UCOP.BurnStrategy = cache.BurnStrategyOldestLRU
	UrlCache = cache.NewCache(UCOP)
	UrlCache.Start()
}