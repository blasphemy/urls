package main

import (
	"log"
	"net/http"
	"strings"

	rdb "github.com/dancannon/gorethink"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var (
	session *rdb.Session
)

type PageData struct {
	CurrentUrl *Url
	Stats      SiteStats
	Message    string
}

func GetNewPageData() PageData {
	k := PageData{}
	k.Stats = GetSiteStats()
	return k
}

func main() {
	err := MakeConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	session, err = rdb.Connect(rdb.ConnectOpts{
		Address:  config.RethinkConnectionString,
		Database: "urls",
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	rdb.Db("urls").TableCreate("meta").Exec(session)
	rdb.Db("urls").TableCreate("urls").Exec(session)
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Extensions: []string{".tmpl", ".html"},
	}))
	m.Get("/", IndexHandler)
	m.Get("/api/add", ApiAddURLHandler)
	m.Get("/add", WebAddHandler)
	m.Get("/view/:id", ViewHandler)
	m.Get("/:id", GetURLAndRedirect)
	log.Println("Listening on " + config.ListenAt)
	log.Fatal(http.ListenAndServe(config.ListenAt, m))
}

func IndexHandler(r render.Render) {
	pd := GetNewPageData()
	r.HTML(http.StatusOK, "index", pd)
}

func ViewHandler(m martini.Params, w http.ResponseWriter, r *http.Request, r2 render.Render) {
	k, err := GetUrlById(m["id"], r.Host)
	if err != nil {
		if err == rdb.ErrEmptyResult {
			pd := GetNewPageData()
			pd.Message = "404 Not Found"
			r2.HTML(http.StatusNotFound, "error", pd)
		} else {
			pd := GetNewPageData()
			pd.Message = err.Error()
			r2.HTML(http.StatusInternalServerError, "error", pd)
		}
		return
	}
	if k != nil {
		pd := GetNewPageData()
		pd.CurrentUrl = k
		r2.HTML(http.StatusOK, "view", pd)
		return
	} else {
		pd := GetNewPageData()
		pd.Message = "404 Not Found"
		r2.HTML(http.StatusNotFound, "error", pd)
		return
	}
}

func WebAddHandler(w http.ResponseWriter, r *http.Request, r2 render.Render) {
	if len(r.URL.Query()["url"]) < 1 {
		pd := GetNewPageData()
		pd.Message = "No arguments specified."
		r2.HTML(500, "error", pd)
		return
	}
	k := r.URL.Query()["url"][0]
	if k == "" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	} else {
		user := UserData{}
		user.IpAddress = r.RemoteAddr
		user.UserAgent = r.UserAgent()
		k = UrlPreprocessor(k)
		new, err := GetNewUrl(k, r.Host, user)
		if err != nil {
			pd := GetNewPageData()
			pd.Message = err.Error()
			r2.HTML(500, "error", pd)
		} else {
			http.Redirect(w, r, "/view/"+new.Id, http.StatusMovedPermanently)
		}
	}
}

func ApiAddURLHandler(r *http.Request) string {
	if len(r.URL.Query()["url"]) < 1 {
		return "Error, no arguments specified"
	}
	k := r.URL.Query()["url"][0]
	if k == "" {
		return "Error, no url specified"
	} else {
		user := UserData{}
		user.IpAddress = r.RemoteAddr
		user.UserAgent = r.UserAgent()
		k := UrlPreprocessor(k)
		new, err := GetNewUrl(k, r.Host, user)
		if err != nil {
			return err.Error()
		} else {
			return config.GetBaseUrl(r.Host) + new.Id
		}
	}
}

func GetURLAndRedirect(params martini.Params, w http.ResponseWriter, r *http.Request, r2 render.Render) {
	k, err := GetUrlById(params["id"], r.Host)
	if err != nil {
		if err == rdb.ErrEmptyResult {
			pd := GetNewPageData()
			pd.Message = "404 Not Found"
			r2.HTML(http.StatusNotFound, "error", pd)
		} else {
			pd := GetNewPageData()
			pd.Message = err.Error()
			r2.HTML(http.StatusInternalServerError, "error", pd)
		}
		return
	}
	if k != nil {
		if strings.Contains(k.Link, config.HostName) || strings.Split(k.Link, ":")[0] == "/"+k.Id {
			k.Link = config.GetBaseUrl(r.Host)
		}
		http.Redirect(w, r, k.Link, http.StatusMovedPermanently)
	} else {
		pd := GetNewPageData()
		pd.Message = "404 Not Found"
		r2.HTML(http.StatusNotFound, "error", pd)
	}
}

func UrlPreprocessor(url string) string {
	if !strings.HasPrefix(url, "http:/") && !strings.HasPrefix(url, "https:/") {
		url = "http://" + url
	}
	return url
}
