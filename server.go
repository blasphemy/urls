package main

import "github.com/go-martini/martini"
import "net/http"
import "fmt"
import "log"

var (
	listenstring = ":5596"
)

func main() {
	m := martini.Classic()
	m.Get("/api/add/**", ApiAddURLHandler)
	m.Get("/list", ListURLS)
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

func ListURLS() string {
	return fmt.Sprint(urlmap)
}
