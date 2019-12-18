package controller

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/iris/v12"
	"log"
	"tinyurl/tinyencoder"
)

var db *gorm.DB

type UrlPair struct {
	ID 	int 	`gorm:"primary_key;auto_increment"`
	Url string
}

func GetUrl(ctx iris.Context) {
	id := tinyencoder.Decode(ctx.Params().Get("url"))
	url := UrlPair{}
	db.First(&url, "id = ?", id)
	if url.ID != 0 {
		_, _ = ctx.JSON(iris.Map{
			"url": url.Url,
		})
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
	db, err = gorm.Open("mysql", "tiny:tiny@tcp(mysql:3306)/tinyurl")
	if err != nil {
		log.Fatal(err)
	}
	if !db.HasTable(&UrlPair{}) {
		db.CreateTable(&UrlPair{})
	}
}