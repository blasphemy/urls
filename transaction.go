package main

import (
	"strconv"
	"strings"
)
import "github.com/garyburd/redigo/redis"

var (
	protected = []string{"list", "add", "api", "counter"}
)

type Url struct {
	id   string
	link string
}

func GetUrlById(id string) *Url {
	DB := pool.Get()
	defer DB.Close()
	k, err := DB.Do("GET", strings.ToLower(id))
	if k != "" {
		resp := &Url{}
		resp.id = id
		resp.link, _ = redis.String(k, err)
		return resp
	} else {
		return nil
	}
}

func GetNewUrl(link string) *Url {
	DB := pool.Get()
	defer DB.Close()
	i := GetNewCounter()
	for _, k := range protected {
		for strconv.FormatInt(i, 36) == k {
			i = GetNewCounter()
		}
	}
	pos := strconv.FormatInt(i, 36)
	DB.Do("SET", pos, link)
	new := &Url{}
	new.id = pos
	new.link = link
	return new
}

func GetNewCounter() int64 {
	DB := pool.Get()
	defer DB.Close()
	n, _ := DB.Do("INCR", "COUNTER")
	return n.(int64)
}

func newPool() *redis.Pool {
	return redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", ":6379")
	}, 3)
}
