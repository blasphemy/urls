package main

type UrlDiag struct {
	CacheHit  int64
	CacheMiss int64
	CacheLen  int
}

func GetDiagnostics() UrlDiag {
	k := UrlDiag{}
	k.CacheHit = UrlCache.Hits()
	k.CacheMiss = UrlCache.Misses()
	k.CacheLen = UrlCache.Len()
	return k
}
