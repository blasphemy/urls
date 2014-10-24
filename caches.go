package main

import (
	"github.com/blasphemy/cache"
	"time"
)

var (
	StatsCache *cache.Cache
)

func initCaches() {
	SCOP := cache.CacheOptions{}
	SCOP.ExpirationTime = time.Minute * 5
	SCOP.Upper = 20
	StatsCache = cache.NewCache(SCOP)
}
