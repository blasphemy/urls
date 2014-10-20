package main

import (
	"github.com/garyburd/redigo/redis"
)

const (
	redis_lua_get_click_sum = `
local sum = 0
local matches = redis.call('KEYS', 'url:clicks:*')

for _,key in ipairs(matches) do
    local val = redis.call('GET', key)
    sum = sum + tonumber(val)
end

return sum`

	redis_lua_get_link_count = `return table.getn(redis.call('keys', 'url:link:*'))`
)

func GetTotalClicksFromScript() (int, error) {
	db := pool.Get()
	defer db.Close()
	i, err := db.Do("EVAL", redis_lua_get_click_sum, 0)
	if err != nil {
		return 0, nil
	}
	k, err := redis.Int(i, err)
	if err != nil {
		return 0, err
	}
	return k, nil
}

func GetTotalUrlsFromScript() (int, error) {
	db := pool.Get()
	defer db.Close()
	i, err := db.Do("EVAL", redis_lua_get_link_count, 0)
	if err != nil {
		return 0, nil
	}
	k, err := redis.Int(i, err)
	if err != nil {
		return 0, err
	}
	return k, nil
}
