package main

import (
	"errors"
	"fmt"
	r "github.com/dancannon/gorethink"
	"log"
)

var (
	protected = []string{"list", "add", "api", "counter", "css", "img", "js"}
)

type Url struct {
	Id     string `gorethink:"id"`
	Link   string `gorethink:"link"`
	Short  string `gorethink:"-"`
	Clicks int64  `gorethink:"clicks"`
}

type SiteStats struct {
	Clicks       int64
	Links        int64
	ClicksPerUrl float64
}

func (u *Url) setShortLink(host string) {
	u.Short = config.GetBaseUrl(host) + u.Id
}

func GetUrlById(id string, host string) (*Url, error) {
	cursor, err := r.Table("urls").Get(id).Run(session)
	if err != nil {
		return nil, err
	}
	result := Url{}
	err = cursor.One(&result)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = r.Table("urls").Update(map[string]interface{}{"id": result.Id, "clicks": r.Row.Field("clicks").Add(1)}).Exec(session)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	result.setShortLink(host)
	return &result, nil
}

func GetNewUrl(link string, host string) (*Url, error) {
	i, err := GetNewID()
	if err != nil {
		return nil, err
	}
	for _, k := range protected {
		for b62_Encode(uint64(i)) == k {
			i, err = GetNewID()
			if err != nil {
				return nil, err
			}
		}
	}
	pos := b62_Encode(uint64(i))
	result := Url{}
	result.Id = pos
	result.Clicks = 0
	result.Link = link
	result.Short = config.GetBaseUrl(host) + result.Id
	err = r.Table("urls").Insert(result).Exec(session)
	log.Println(result)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &result, nil
}

func GetNewID() (int64, error) {
	var target interface{}
	err := r.Table("meta").Get("counter").Update(map[string]interface{}{"value": r.Row.Field("value").Add(1)}).Exec(session)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	cursor, err := r.Table("meta").Get("counter").Field("value").Run(session)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	cursor.One(&target)
	if cursor.Err() != nil {
		return 0, cursor.Err()
	}
	final, ok := target.(float64)
	if !ok {
		return 0, errors.New("Cannot convert counter to float64")
	}
	return int64(final), nil
}

func GetSiteStats() SiteStats {
	k := SiteStats{}
	totalClicks, _ := GetTotalClicks()
	totalLinks, _ := GetTotalLinks()
	clicksPerUrl, _ := GetClicksPerUrl()
	k.Clicks = totalClicks
	k.Links = totalLinks
	k.ClicksPerUrl = clicksPerUrl
	return k
}

func GetTotalClicks() (int64, error) {
	var result interface{}
	cursor, err := r.Table("urls").Sum("clicks").Run(session)
	if err != nil {
		return 0, err
	}
	err = cursor.One(&result)
	if err != nil {
		return 0, err
	}
	final, ok := result.(float64)
	if !ok {
		return 0, errors.New("urls.sum(\"clicks\") is not a float64")
	}
	return int64(final), nil

}

func GetTotalLinks() (int64, error) {
	var result interface{}
	cursor, err := r.Table("urls").Count().Run(session)
	if err != nil {
		return 0, err
	}
	err = cursor.One(&result)
	if err != nil {
		return 0, err
	}
	final, ok := result.(float64)
	if !ok {
		return 0, errors.New("urls.count is not a float64")
	}
	return int64(final), nil

}

//Totally being lazy on this one...
func GetClicksPerUrl() (float64, error) {
	totalLinks, err := GetTotalLinks()
	if err != nil {
		return 0, err
	}
	totalClicks, err := GetTotalClicks()
	if err != nil {
		return 0, err
	}
	return float64(totalClicks) / float64(totalLinks), nil

}

func GetMetaValue(key string) (interface{}, error) {
	var target interface{}
	cursor, err := r.Table("meta").Get(key).Field("value").Run(session)
	if err != nil {
		return 0, err
	}
	cursor.One(&target)
	if cursor.Err() != nil {
		//TODO
		//if result is empty set it to 0
		return 0, cursor.Err()
	}
	result, ok := target.(float64)
	if !ok {
		return 0, errors.New(fmt.Sprintf("meta.%s is not a float64", key))
	}
	return result, nil
}

//Idea for unique urls?
/*
func ThisisATest() {
	r.Table("urls").AtIndex("link").Distinct().Run(session)
}
*/
