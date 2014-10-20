package main

import (
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
	m.Get("/view/:id", ViewHandler)
	m.Get("/:id", GetURLAndRedirect)
	log.Println("Listening on " + config.ListenAt)
	go RunJobs()
	log.Fatal(http.ListenAndServe(config.ListenAt, m))
}

func IndexHandler(r render.Render) {
	r.HTML(http.StatusOK, "index", "")
}

func ViewHandler(m martini.Params, w http.ResponseWriter, r *http.Request, r2 render.Render) {
	k, err := GetUrlById(m["id"])
	if err != nil {
		r2.HTML(http.StatusInternalServerError, "error", err.Error())
		return
	}
	if k != nil {
		r2.HTML(http.StatusOK, "view", k)
		return
	} else {
		r2.HTML(http.StatusNotFound, "error", "404 Not Found")
		return
	}
}

func WebAddHandler(w http.ResponseWriter, r *http.Request, r2 render.Render) {
	if len(r.URL.Query()["url"]) < 1 {
		r2.HTML(500, "error", "No arguments specified.")
		return
	}
	k := r.URL.Query()["url"][0]
	if k == "" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	} else {
		new, err := GetNewUrl(k)
		if err != nil {
			r2.HTML(500, "error", err.Error())
		} else {
			http.Redirect(w, r, "/view/"+new.id, http.StatusMovedPermanently)
		}
	}
}

func GetURLAndRedirect(params martini.Params, w http.ResponseWriter, r *http.Request, r2 render.Render) {
	k, err := GetUrlById(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if k != nil {
		if strings.Contains(k.Link, config.BaseURL) || strings.Split(k.Link, ":")[0] == "/"+k.id {
			k.Link = config.BaseURL
		}
		http.Redirect(w, r, k.Link, http.StatusMovedPermanently)
	} else {
		r2.HTML(http.StatusNotFound, "error", "404 Not Found")
	}
}

func ApiAddURLHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
	k, err := GetNewUrl(params["_1"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write([]byte(config.BaseURL + k.id))
	}
}
