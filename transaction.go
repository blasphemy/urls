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
	Clicks int
	Links  int
}

func GetUrlById(id string) (*Url, error) {
	DB := pool.Get()
	defer DB.Close()
	id = strings.Split(id, ":")[0]
	k, err := DB.Do("GET", "url:link:"+id)
	if err != nil {
		return nil, err
	}
	switch k.(type) {
	case nil:
		return nil, nil
	default:
		c, _ := DB.Do("INCR", "url:clicks:"+id)
		resp := &Url{}
		resp.id = id
		resp.Short = config.BaseURL + id
		resp.Link, _ = redis.String(k, err)
		resp.Clicks = c.(int64)
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
	k := SiteStats{}
	a, _ := GetTotalClicks()
	b, _ := GetTotalUrls()
	k.Clicks = a
	k.Links = b
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

func SetTotalClicks() {
	log.Print("Setting total number of clicks in DB...")
	i, err := GetTotalClicksFromScript()
	if err != nil {
		log.Print("error updating total clicks in db.... ", err.Error())
		return
	}
	DB := pool.Get()
	defer DB.Close()
	_, err = DB.Do("SET", "meta:total:clicks", i)
	if err != nil {
		log.Print("error updating total clicks in db.... ", err.Error())
		return
	}
	log.Printf("Total number of clicks set to %d", i)
}

func SetTotalUrls() {
	log.Print("Setting total number of urls in DB....")
	i, err := GetTotalUrlsFromScript()
	if err != nil {
		log.Print("Error updating total urls", err.Error())
		return
	}
	DB := pool.Get()
	defer DB.Close()
	_, err = DB.Do("SET", "meta:total:links", i)
	if err != nil {
		log.Print("Error updating total urls", err.Error())
		return
	}
	log.Printf("Total number of URLS set to %d", i)
}
