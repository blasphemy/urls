package main

import "github.com/go-martini/martini"
import "github.com/garyburd/redigo/redis"
import "net/http"
import "fmt"
import "log"

var (
	listenstring = ":5596"
	DB           redis.Conn
)

func main() {
	DB, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err.Error())
	} else {
		k, _ := DB.Do("INCR", "COUNTER")
		log.Print(k)
	}
	m := martini.Classic()
	m.Get("/api/add/**", ApiAddURLHandler)
	m.Get("/:id", GetURLAndRedirect)
	log.Println("Listening on " + listenstring)
	log.Fatal(http.ListenAndServe(listenstring, m))
}

func GetURLAndRedirect(params martini.Params, w http.ResponseWriter, r *http.Request) {
	k := GetUrlById(params["id"])
	if k != nil {
		http.Redirect(w, r, k.link, http.StatusMovedPermanently)
	} else {
		http.Error(w, fmt.Sprintf("/%s not found", params["id"]), 404)
	}
}

func ApiAddURLHandler(params martini.Params) string {
	k := GetNewUrl(params["_1"])
	return k.id
}
