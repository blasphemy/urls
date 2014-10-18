package main

import "github.com/go-martini/martini"
import "github.com/garyburd/redigo/redis"
import "net/http"
import "fmt"
import "log"

var (
	listenstring = ":5596"
	pool         *redis.Pool
)

func main() {
	pool = newPool()
	m := martini.Classic()
	m.Get("/api/add/**", ApiAddURLHandler)
	m.Get("/:id", GetURLAndRedirect)
	log.Println("Listening on " + listenstring)
	log.Fatal(http.ListenAndServe(listenstring, m))
}

func GetURLAndRedirect(params martini.Params, w http.ResponseWriter, r *http.Request) {
	k, err := GetUrlById(params["id"])
	if err != nil {
		http.Error(w, err.Error(), 500)
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
		w.Write([]byte(k.id))
	}
}
