package main

import "github.com/go-martini/martini"
import "strconv"
import "net/http"
import "fmt"
import "strings"
import "log"


var (
	urlmap    = make(map[string]string)
	counter   int64
	protected = []string{"list", "add"}
	listenstring = ":5596"
)

func main() {
	m := martini.Classic()
	m.Get("/api/add/**", AddURL)
	m.Get("/list", ListURLS)
	m.Get("/:id", GetURLById)
	log.Println("Listening on " + listenstring)
	log.Fatal(http.ListenAndServe(listenstring, m))
}

func AddURL(params martini.Params) string {
	counter++
	//Check to make sure it doesn't match "list or add"
	for _, k := range protected {
		for strconv.FormatInt(counter, 36) == k {
			counter++
		}
	}
	pos := strconv.FormatInt(counter, 36)
	urlmap[pos] = params["_1"]
	return pos
}

func GetURLById(params martini.Params, w http.ResponseWriter, r *http.Request) {
	k := urlmap[strings.ToLower(params["id"])]
	if k != "" {
		http.Redirect(w, r, k, http.StatusMovedPermanently)
	} else {
		http.Error(w, fmt.Sprintf("/%s not found", params["id"]), 404)
	}
}

func ListURLS() string {
	return fmt.Sprint(urlmap)
}
