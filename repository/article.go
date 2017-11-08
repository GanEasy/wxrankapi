package repository

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/GanEasy/wxrank/orm"
	"github.com/yizenghui/sda/wechat"
)

func init() {
	orm.DB().AutoMigrate(&orm.Tag{}, &orm.Article{}, &orm.Media{}, &orm.Post{})

}

// Find wechat.Article
func Find(url string) (article wechat.Article, err error) {

	article, err = wechat.Find(url)
	return
}

func update(url string) (article wechat.Article, err error) {
	article, err = wechat.Find(url)
	return
}

// Insert wechat.Article
func Insert(url string) (article wechat.Article, err error) {

	article, err = wechat.Find(url)
	return
}

//GetArticle 获取文章列表
func GetArticle(limit, offset, tag int, order string) (articles []orm.Article, err error) {
	var a orm.Article
	// var articles []orm.Article

	if order == "id" {
		order = "id DESC"
	} else {
		order = "rank DESC"
	}

	articles = a.GetArticle(limit, offset, tag, order)

	// orm.DB().Offset(offset).Limit(limit).Order("rank DESC").Find(&articles)
	for key, article := range articles {
		// articles[key].Cover = "http://localhost:8004/image/" + base64.URLEncoding.EncodeToString([]byte(article.Cover))
		articles[key].Cover = "http://pic3.readfollow.com/" + base64.URLEncoding.EncodeToString([]byte(article.Cover))
		article.URL = strings.Replace(article.URL, `http://`, `https://`, -1)
		articles[key].URL = strings.Replace(article.URL, `#rd`, "&scene=27#wechat_redirect", 1)

		article.Title = strings.Replace(article.Title, `\x26quot;`, `"`, -1)
		article.Title = strings.Replace(article.Title, `\x26amp;`, `&`, -1)
		article.Title = strings.Replace(article.Title, `\x0a`, "\n", -1)
		articles[key].Title = article.Title

		article.Intro = strings.Replace(article.Intro, `\x0a`, "\n", -1)
		article.Intro = strings.Replace(article.Intro, `\x26quot;`, `"`, -1)
		article.Intro = strings.Replace(article.Intro, `\x26amp;`, `&`, -1)
		articles[key].Intro = article.Intro

	}

	return
}

//View ..
func View(id int) (a orm.Article, err error) {

	// var a orm.Article
	a.GetArticleByID(id)

	if a.Title == "" {
		err = errors.New("内容异常")
		return
	}

	if a.View < 1 {
		a.View = 1
	} else {
		a.View++
	}

	if a.ID != 0 {
		a.Rank = Rank(int(a.View), 0, a.PubAt.Unix())
		a.Save()
	}

	a.Cover = "http://pic3.readfollow.com/" + base64.URLEncoding.EncodeToString([]byte(a.Cover))

	return
}

// Post ... url
func Post(url string) (err error) {
	var post orm.Post
	post.GetPostByURL(url)
	// if post.State == 0 { // 检查提交状态
	var a orm.Article
	article, err := wechat.Find(url)
	if err == nil {

		if article.URL == "" {
			return errors.New("不支持该链接！")
		}

		media, err := GetMediaByAppID(article.AppID)
		if err != nil {
			return errors.New("公众号信息出错！")
		}
		// 如果公众号是新收录的
		if media.State == 0 {
			media.AppName = article.AppName
			media.Cover = article.RoundHead
			media.State = 1

			// 公众号ID作为一个标签
			var tag orm.Tag
			tag.GetTagByName(article.AppID)
			if tag.Type == "" {
				tag.Type = "wxid"
				tag.Title = article.AppName
				// tag.IsShow = 0
				tag.Save()
			}

			media.Tags = append(media.Tags, int64(tag.ID))

			media.Save()
		}

		post.ArticleURL = article.URL
		post.State = 1
		post.Save()
		a.GetArticleByURL(article.URL)
		a.MediaID = media.ID
		a.Title = article.Title
		a.Intro = article.Intro
		a.Cover = article.Cover
		a.Author = article.Author
		// todo 标签管理，需要保留自定义标签
		a.Tags = media.Tags // 文章的标签等于公众号的标签

		i64, err := strconv.ParseInt(article.PubAt, 10, 64)
		if err != nil {
			// fmt.Println(err)
			return errors.New("时间转化失败")
		}
		a.PubAt = time.Unix(i64, 0)
		a.Save()
		// panic(a.ID)
		// fmt.Println(a)
	}
	// }
	return
}

// Remove wechat.Article
func Remove(url string) (article wechat.Article, err error) {

	article, err = wechat.Find(url)
	return
}
