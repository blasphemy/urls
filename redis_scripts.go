package main

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

const (
	redis_lua_get_set_click_sum = `
		local sum = 0
		local matches = redis.call('KEYS', 'url:clicks:*')
		for _,key in ipairs(matches) do
		    local val = redis.call('GET', key)
		    sum = sum + tonumber(val)
		end
		redis.call('SET', 'meta:total:clicks', sum)
		return sum`
	redis_lua_get_set_link_count = `
		local num = table.getn(redis.call('keys', 'url:link:*'))
		redis.call('SET', 'meta:total:links', num)
		return num`
)

func SetGetTotalClicksFromScript() (int, error) {
	db := pool.Get()
	s := redis.NewScript(0, redis_lua_get_set_click_sum)
	a, err := s.Do(db)
	if err != nil {
		return 0, err
	}
	b, err := redis.Int(a, err)
	if err != nil {
		return 0, err
	}
	log.Printf("Number of clicks set to %d", b)
	return b, nil
}

func SetGetTotalUrlsFromScript() (int, error) {
	db := pool.Get()
	s := redis.NewScript(0, redis_lua_get_set_link_count)
	a, err := s.Do(db)
	if err != nil {
		return 0, err
	}
	b, err := redis.Int(a, err)
	if err != nil {
		return 0, err
	}
	log.Printf("Number of links set to %d", b)
	return b, nil
}
