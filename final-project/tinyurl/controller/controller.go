package controller

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/iris/v12"
	"log"
	"tinyurl/tinyencoder"
)

var db *gorm.DB
var client *redis.Client

type UrlPair struct {
	ID 	int 	`gorm:"primary_key;auto_increment"`
	Url string
}

func GetUrl(ctx iris.Context) {
	turl := ctx.Params().Get("url")
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}
	val, err := client.Get(turl).Result()
	if err != redis.Nil {
		log.Println(err)
	}
	if val != "" {
		_, _ = ctx.JSON(iris.Map{
			"url": val,
		})
		return
	}
	// cache miss
	id := tinyencoder.Decode(turl)
	url := UrlPair{}
	db.First(&url, "id = ?", id)
	if url.ID != 0 {
		_, _ = ctx.JSON(iris.Map{
			"url": url.Url,
		})
		client.Set(turl, url.Url, 300000000000)
	} else {
		ctx.StatusCode(404)
	}
}

func AddUrl(ctx iris.Context) {
	urlString := ctx.URLParam("url")
	url := UrlPair{}
	db.First(&url, "url = ?", urlString)
	if url.ID == 0 {
		url.Url = urlString
		db.Create(&url)
		db.First(&url, "url = ?", urlString)
		ctx.StatusCode(200)
		_, _ = ctx.JSON(iris.Map{
			"message": "ok",
			"tiny_url": tinyencoder.Encode(url.ID),
		})
	} else {
		ctx.StatusCode(400)
		_, _ = ctx.JSON(iris.Map{
			"message": "the url has been added",
			"tiny_url": tinyencoder.Encode(url.ID),
		})
	}
}

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:tiny@tcp(localhost:3306)/tinyurl")
	if err != nil {
		log.Fatal(err)
	}
	if !db.HasTable(&UrlPair{}) {
		db.CreateTable(&UrlPair{})
	}
}