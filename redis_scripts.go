package main

import (
	"errors"
	r "github.com/dancannon/gorethink"
	"github.com/garyburd/redigo/redis"
	"log"
)

const (
	redis_lua_get_set_click_sum = `
		local sum = 0
		local matches = redis.call('KEYS', 'url:clicks:*')r 
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

func SetGetTotalLinks() (int64, error) {
	cursor, err := r.Table("urls").Count().Run(session)
	if err != nil {
		return 0, err
	}
	var result interface{}
	err = cursor.One(&result)
	if err != nil {
		return 0, err
	}
	result2, ok := result.(float64)
	if !ok {
		return 0, errors.New("urls.count() is not a float64")
	}
	err = r.Table("meta").Get("total_links").Update(map[string]interface{}{"value": result2}).Exec(session)
	if err != nil {
		return 0, err
	}
	return int64(result2), nil
}

/*
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
*/

func SetGetClicksPerUrl() (float64, error) {
	db := pool.Get()
	defer db.Close()
	a, err := db.Do("GET", "meta:total:clicks")
	b, err := db.Do("GET", "meta:total:links")
	c, err := redis.Int(a, err)
	d, err := redis.Int(b, err)
	if err != nil {
		return 0, err
	}
	x := float64(c) / float64(d)
	_, err = db.Do("SET", "meta:clicksperurl", x)
	if err != nil {
		log.Print(err.Error())
		return 0, err
	} else {
		log.Print("Clicks per Url set to ", x)
		return x, nil
	}
}

func GetClicksPerUrl() (float64, error) {
	db := pool.Get()
	defer db.Close()
	a, err := db.Do("GET", "meta:clicksperurl")
	b, err := redis.Float64(a, err)
	if err != nil {
		return 0, err
	} else {
		return b, nil
	}
}
