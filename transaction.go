package main

import (
	"errors"
	r "github.com/dancannon/gorethink"
	"github.com/garyburd/redigo/redis"
	"log"
	//"strings"
)

var (
	protected = []string{"list", "add", "api", "counter", "css", "img", "js"}
)

type Url struct {
	Id     string `gorethink:"id"`
	Link   string `gorethink:"link"`
	Short  string `gorethink:"-"`
	Clicks int64  `gorethink:"clicks"`
}

type SiteStats struct {
	Clicks       int
	Links        int
	ClicksPerUrl float64
}

/*
func GetUrlById(id string, host string) (*Url, error) {
	DB := pool.Get()
	defer DB.Close()
	cr := UrlCache.Get(id)
	if cr != nil {
		log.Print("UrlCache: Cache HIT!")
		log.Print("Updating click count in goroutine")
		go UpdateClickCount(id)
		//return cr.(*Url), nil
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
		resp.Id = id
		resp.Short = config.GetBaseUrl(host) + id
		resp.Link, _ = redis.String(k, err)
		resp.Clicks = int64(c)
		UrlCache.Set(id, resp)
		return resp, nil
	}
}
*/

func GetUrlById(id string, host string) (*Url, error) {
	cursor, err := r.Table("urls").Get(id).Run(session)
	if err != nil {
		return nil, err
	}
	result := Url{}
	err = cursor.One(&result)
	if err != nil {
		return nil, err
	}
	err = r.Table("urls").Update(map[string]interface{}{"clicks": r.Row.Field("clicks").Add(1)}).Exec(session)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetNewUrl(link string, host string) (*Url, error) {
	i, err := GetNewID()
	if err != nil {
		return nil, err
	}
	for _, k := range protected {
		for b62_Encode(uint64(i)) == k {
			i, err = GetNewID()
			if err != nil {
				return nil, err
			}
		}
	}
	pos := b62_Encode(uint64(i))
	result := Url{}
	result.Id = pos
	result.Clicks = 0
	result.Link = link
	result.Short = config.GetBaseUrl(host) + result.Id
	err = r.Table("urls").Insert(result).Exec(session)
	log.Println(result)
	if err != nil {
		return nil, err
	}
	UrlCache.Set(result.Id, result)
	return &result, nil
}

func GetNewID() (int64, error) {
	var target interface{}
	err := r.Table("meta").Get("counter").Update(map[string]interface{}{"value": r.Row.Field("value").Add(1)}).Exec(session)
	if err != nil {
		return 0, err
	}
	cursor, err := r.Table("meta").Get("counter").Field("value").Run(session)
	if err != nil {
		return 0, err
	}
	cursor.One(&target)
	if cursor.Err() != nil {
		return 0, cursor.Err()
	}
	final, ok := target.(float64)
	if !ok {
		return 0, errors.New("Cannot convert counter to float64")
	}
	return int64(final), nil
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
	var target interface{}
	cursor, err := r.Table("meta").Get("total_links").Field("value").Run(session)
	if err != nil {
		return 0, err
	}
	cursor.One(&target)
	if cursor.Err() != nil {
		return 0, cursor.Err()
	}
	result, ok := target.(float64)
	if !ok {
		return 0, errors.New("meta.total_links is not a float64")
	}
	return int(result), nil
}

func GetTotalClicks() (int, error) {
	var target interface{}
	cursor, err := r.Table("meta").Get("total_clicks").Field("value").Run(session)
	if err != nil {
		return 0, err
	}
	cursor.One(&target)
	if cursor.Err() != nil {
		return 0, cursor.Err()
	}
	result, ok := target.(float64)
	if !ok {
		return 0, errors.New("meta.total_clicks is not a float64")
	}
	return int(result), nil
}

//No longer needed?
func UpdateClickCount(id string) (int, error) {
	DB := pool.Get()
	defer DB.Close()
	k, err := DB.Do("INCR", "url:clicks:"+id)
	i, err := redis.Int(k, err)
	return i, err
}
