package main

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"strings"
)

var (
	protected = []string{"list", "add", "api", "counter", "css", "img", "js"}
)

type Url struct {
	id     string
	Link   string
	Short  string
	Clicks int64
}

type SiteStats struct {
	Clicks       int
	Links        int
	ClicksPerUrl float64
}

func GetUrlById(id string) (*Url, error) {
	DB := pool.Get()
	defer DB.Close()
	cr := UrlCache.Get(id)
	if cr != nil {
		log.Print("UrlCache: Cache HIT!")
		log.Print("Updating click count in goroutine")
		go UpdateClickCount(id)
		return cr.(*Url), nil
	}
	log.Print("UrlCache: Cache Miss, retrieving from DB")
	id = strings.Split(id, ":")[0]
	k, err := DB.Do("GET", "url:link:"+id)
	if err != nil {
		return nil, err
	}
	switch k.(type) {
	case nil:
		return nil, nil
	default:
		c, _ := UpdateClickCount(id)
		resp := &Url{}
		resp.id = id
		resp.Short = config.GetBaseUrl() + id
		resp.Link, _ = redis.String(k, err)
		resp.Clicks = int64(c)
		UrlCache.Set(id, resp)
		return resp, nil
	}
}

func GetNewUrl(link string) (*Url, error) {
	DB := pool.Get()
	defer DB.Close()
	i, err := GetNewCounter()
	if err != nil {
		return nil, err
	}
	for _, k := range protected {
		for b62_Encode(uint64(i)) == k {
			i, err = GetNewCounter()
			if err != nil {
				return nil, err
			}
		}
	}
	pos := b62_Encode(uint64(i))
	_, err = DB.Do("SET", "url:link:"+pos, link)
	_, err = DB.Do("SET", "url:clicks:"+pos, 0)
	if err != nil {
		return nil, err
	}
	new := &Url{}
	new.id = pos
	new.Link = link
	new.Clicks = 0
	new.Short = config.GetBaseUrl() + new.id
	UrlCache.Set(new.id, new)
	log.Printf("Shortened %s to %s", new.Link, config.GetBaseUrl()+new.id)
	return new, nil
}

func GetNewCounter() (int64, error) {
	DB := pool.Get()
	defer DB.Close()
	n, err := DB.Do("INCR", "meta:COUNTER")
	if err != nil {
		return 0, err
	}
	return n.(int64), nil
}

func GetSiteStats() SiteStats {
	cc := StatsCache.Get("Stats")
	if cc != nil {
		log.Print("Cache: Site Stats HIT")
		return cc.(SiteStats)
	} else {
		log.Print("Cache: Site Stats MISS")
	}
	k := SiteStats{}
	a, _ := GetTotalClicks()
	b, _ := GetTotalUrls()
	c, _ := GetClicksPerUrl()
	k.Clicks = a
	k.Links = b
	k.ClicksPerUrl = c
	StatsCache.Set("Stats", k)
	return k
}

func newPool() *redis.Pool {
	return redis.NewPool(func() (redis.Conn, error) {
		conn, err := redis.Dial("tcp", config.DBAddress)
		if err != nil {
			return nil, err
		}
		_, err = conn.Do("AUTH", config.DBPassword)
		if err != nil {
			log.Print(err.Error())
			return nil, err
		}
		return conn, nil

	}, 3)
}

func GetTotalUrls() (int, error) {
	db := pool.Get()
	defer db.Close()
	k, err := db.Do("GET", "meta:total:links")
	if err != nil {
		return 0, err
	}
	l, err := redis.Int(k, err)
	if err != nil {
		return 0, err
	}
	return l, err
}

func GetTotalClicks() (int, error) {
	DB := pool.Get()
	defer DB.Close()
	k, err := DB.Do("GET", "meta:total:clicks")
	if err != nil {
		return 0, nil
	}
	j, err := redis.Int(k, err)
	if err != nil {
		return 0, nil
	}
	return j, nil
}

func UpdateClickCount(id string) (int, error) {
	DB := pool.Get()
	defer DB.Close()
	k, err := DB.Do("INCR", "url:clicks:"+id)
	i, err := redis.Int(k, err)
	return i, err
}
