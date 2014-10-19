package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
	"strings"
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
	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Extensions: []string{".tmpl", ".html"},
	}))
	m.Get("/", IndexHandler)
	m.Get("/api/add/**", ApiAddURLHandler)
	m.Get("/add", WebAddHandler)
	m.Get("/:id", GetURLAndRedirect)
	log.Println("Listening on " + config.ListenAt)
	log.Fatal(http.ListenAndServe(config.ListenAt, m))
}

func IndexHandler(r render.Render) {
	r.HTML(200, "index", "")
}

func WebAddHandler(w http.ResponseWriter, r *http.Request, r2 render.Render) {
	k := r.URL.Query()["url"][0]
	if k == "" {
		http.Redirect(w, r, "/", 302)
	} else {
		new, err := GetNewUrl(k)
		if err != nil {
			r2.HTML(500, "add", err.Error())
		} else {
			r2.HTML(200, "add", config.BaseURL+new.id)
		}
	}
}

func GetURLAndRedirect(params martini.Params, w http.ResponseWriter, r *http.Request) {
	k, err := GetUrlById(params["id"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if k != nil {
		if strings.Contains(k.link, config.BaseURL) || strings.Split(k.link, ":")[0] == "/"+k.id {
			k.link = config.BaseURL
		}
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
