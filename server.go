package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
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
	router := gin.Default()
	router.LoadHTMLTemplates("templates/*")
	router.GET("/*params", MyCustomRouter)
	router.GET("/", IndexHandler)
	router.GET("/add/", WebAddHandler)
	router.GET("/view/:id/", ViewHandler)
	router.GET("/:id/", GetURLAndRedirect)
	router.GET("/api/add/*url", ApiAddURLHandler)
	log.Println("Listening on " + config.ListenAt)
	log.Fatal(http.ListenAndServe(config.ListenAt, router))
}

func IndexHandler(c *gin.Context) {
	log.Println(c.Params)
	c.HTML(http.StatusOK, "index", "")
}

func ViewHandler(c *gin.Context) {
	k, err := GetUrlById(c.Params.ByName("id"))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error", err.Error())
		return
	}
	if k != nil {
		c.HTML(http.StatusOK, "view", k)
		return
	} else {
		c.HTML(http.StatusNotFound, "error", "404 Not Found")
		return
	}
}

func WebAddHandler(c *gin.Context) {
	if len(c.Request.URL.Query()["url"]) < 1 {
		c.HTML(500, "error", "No arguments specified.")
		return
	}
	k := c.Request.URL.Query()["url"][0]
	if k == "" {
		c.Redirect(http.StatusMovedPermanently, "/")
	} else {
		new, err := GetNewUrl(k)
		if err != nil {
			c.HTML(500, "error", err.Error())
		} else {
			c.Redirect(http.StatusMovedPermanently, "/view/"+new.id)
		}
	}
}

func GetURLAndRedirect(c *gin.Context) {
	k, err := GetUrlById(c.Params.ByName("id"))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error", err.Error())
		return
	}
	if k != nil {
		if strings.Contains(k.Link, config.BaseURL) || strings.Split(k.Link, ":")[0] == "/"+k.id {
			k.Link = config.BaseURL
		}
		c.Redirect(http.StatusMovedPermanently, k.Link)
	} else {
		c.HTML(http.StatusNotFound, "error", "404 Not Found")
	}
}

func ApiAddURLHandler(c *gin.Context) {
	k, err := GetNewUrl(c.Params.ByName("url"))
	if err != nil {
		c.Fail(http.StatusInternalServerError, err)
	} else {
		c.Writer.Write([]byte(config.BaseURL + k.id))
	}
}
