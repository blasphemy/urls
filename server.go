package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"log"
	"net/http"
)

var (
	pool *redis.Pool
)

func main() {
	err := MakeConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	pool = newPool()
	m := martini.Classic()
	m.Get("/api/add/**", ApiAddURLHandler)
	m.Get("/add", WebAddHandler)
	m.Get("/:id", GetURLAndRedirect)
	log.Println("Listening on " + config.ListenAt)
	log.Fatal(http.ListenAndServe(config.ListenAt, m))
}

func WebAddHandler(w http.ResponseWriter, r *http.Request) {
	k := r.URL.Query()["url"][0]
	if k == "" {
		http.Redirect(w, r, "/", 302)
	} else {
		w.Write([]byte(k))
	}
}

func GetURLAndRedirect(params martini.Params, w http.ResponseWriter, r *http.Request) {
	k, err := GetUrlById(params["id"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if k != nil {
		http.Redirect(w, r, k.link, http.StatusMovedPermanently)
	} else {
		http.Error(w, fmt.Sprintf("/%s not found", params["id"]), 404)
	}
}

func ApiAddURLHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
	k, err := GetNewUrl(params["_1"])
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		w.Write([]byte(config.BaseURL + k.id))
	}
}
