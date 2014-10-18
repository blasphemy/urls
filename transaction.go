package main

import (
	"strconv"
	"strings"
)

var (
	protected = []string{"list", "add", "api"}
	urlmap    = make(map[string]string)
	counter   int64
)

type Url struct {
	id   string
	link string
}

func GetUrlById(id string) *Url {
	k := urlmap[strings.ToLower(id)]
	if k != "" {
		resp := &Url{}
		resp.id = id
		resp.link = k
		return resp
	} else {
		return nil
	}

}

func GetNewUrl(link string) *Url {
	counter++
	i := counter
	for _, k := range protected {
		for strconv.FormatInt(i, 36) == k {
			counter++
			i = counter
		}
	}
	pos := strconv.FormatInt(i, 36)
	urlmap[pos] = link
	new := &Url{}
	new.id = pos
	new.link = link
	return new
}
