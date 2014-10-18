package main

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
)

var (
	protected = []string{"list", "add", "api", "counter", "css", "img", "js"}
)

type Url struct {
	id   string
	link string
}

func GetUrlById(id string) (*Url, error) {
	DB := pool.Get()
	defer DB.Close()
	id = strings.ToLower(strings.Split(id, ":")[0])
	k, err := DB.Do("GET", "url:"+id)
	if err != nil {
		return nil, err
	}
	if k != "" {
		DB.Do("INCR", "url:"+id+":clicks")
		resp := &Url{}
		resp.id = id
		resp.link, _ = redis.String(k, err)
		return resp, nil
	} else {
		return nil, nil
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
		for strconv.FormatInt(i, 36) == k {
			i, err = GetNewCounter()
			if err != nil {
				return nil, err
			}
		}
	}
	pos := strconv.FormatInt(i, 36)
	_, err = DB.Do("SET", "url:"+pos, link)
	if err != nil {
		return nil, err
	}
	new := &Url{}
	new.id = pos
	new.link = link
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

func newPool() *redis.Pool {
	return redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", ":6379")
	}, 3)
}
