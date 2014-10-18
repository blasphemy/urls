package main

import (
	"strconv"
	"strings"
)

var (
	protected = []string{"list", "add", "api", "counter"}
)

type Url struct {
	id   string
	link string
}

func GetUrlById(id string) *Url {
	k, _ := DB.Do("GET", strings.ToLower(id))
	if k != "" {
		resp := &Url{}
		resp.id = id
		resp.link = k.(string)
		return resp
	} else {
		return nil
	}
}

func GetNewUrl(link string) *Url {
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
	n, _ := DB.Do("INCR", "COUNTER")
	return n.(int64)
}
