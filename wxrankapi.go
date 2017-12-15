package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/GanEasy/wxrankapi/job"
	"github.com/GanEasy/wxrankapi/repository"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/yizenghui/reader"
)

type (
	//Stats 结构
	Stats struct {
		Uptime       time.Time      `json:"uptime"`
		RequestCount uint64         `json:"requestCount"`
		Statuses     map[string]int `json:"statuses"`
		mutex        sync.RWMutex
	}
)

//NewStats New Stats
func NewStats() *Stats {
	return &Stats{
		Uptime:   time.Now(),
		Statuses: make(map[string]int),
	}
}

// Process is the middleware function.
func (s *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.RequestCount++
		status := strconv.Itoa(c.Response().Status)
		s.Statuses[status]++
		return nil
	}
}

// Handle is the endpoint to get stats.
func (s *Stats) Handle(c echo.Context) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return c.JSON(http.StatusOK, s)
}

//Articles 文章接口
func Articles(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	tag, _ := strconv.Atoi(c.QueryParam("tag"))

	if limit <= 0 || limit > 100 {
		limit = 10
	}
	// limit = 10
	if offset < 0 || offset > 500 {
		offset = 0
	}

	articles, err := repository.GetArticle(limit, offset, tag)

	if err != nil {

	}

	return c.JSON(http.StatusOK, articles)
}

//NewArticles 最新收录文章接口
func NewArticles(c echo.Context) error {

	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	tags := repository.Str2Int64(c.QueryParam("tag"))

	id, _ := strconv.Atoi(c.QueryParam("id"))

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	articles, _ := repository.GetArticleCursorByID(id, limit, tags)

	return c.JSON(http.StatusOK, articles)
}

//HotArticles 文章接口 根据热门程序进行游标提取
func HotArticles(c echo.Context) error {

	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	tags := repository.Str2Int64(c.QueryParam("tag"))

	rank, _ := strconv.ParseFloat(c.QueryParam("rank"), 64)

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	articles, _ := repository.GetArticleCursorByRank(rank, limit, tags)

	return c.JSON(http.StatusOK, articles)
}

//Tags 标签列表接口
func Tags(c echo.Context) error {
	t := c.QueryParam("type")

	tags, err := repository.GetTagByType(t)

	if err != nil {

	}
	return c.JSON(http.StatusOK, tags)
}

//SearchTags 标签搜索
func SearchTags(c echo.Context) error {
	s := c.QueryParam("s")

	tags, err := repository.SearchTags(s)

	if err != nil {

	}
	return c.JSON(http.StatusOK, tags)
}

//Search 微信文章内容搜索入口
func Search(c echo.Context) error {
	t := c.QueryParam("s")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	rank, _ := strconv.ParseFloat(c.QueryParam("rank"), 64)
	articles, err := repository.GetArticlesByTitle(t, rank, limit)

	if err != nil {

	}
	return c.JSON(http.StatusOK, articles)
}

//Tag 标签详细
func Tag(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	tag, err := repository.Tag(id)

	if err != nil {

	}

	return c.JSON(http.StatusOK, tag)
}

//GetTagByMediaID 通过公众号ID获取标签详细
func GetTagByMediaID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	tag, err := repository.GetTagByMediaID(id)

	if err != nil {

	}

	return c.JSON(http.StatusOK, tag)
}

//View 阅读
func View(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	article, err := repository.View(id)

	if err != nil {

	}

	return c.JSON(http.StatusOK, article)
}

//Like 喜欢
func Like(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	article, err := repository.Like(id)

	if err != nil {

	}

	return c.JSON(http.StatusOK, article)
}

//Hate 讨厌
func Hate(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	article, err := repository.Hate(id)

	if err != nil {

	}

	return c.JSON(http.StatusOK, article)
}

//Media 公众号
func Media(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	media, err := repository.GetMediaByID(id)

	if err != nil {

	}

	return c.JSON(http.StatusOK, media)
}

//Fetch get 报料接口
func Fetch(c echo.Context) error {
	url := c.QueryParam("url")
	// fmt.Println(url)
	if url != "" {
		// repository.Post(url)
		// 列队任务, 防止高并发攻击
		job.JobQueue <- job.Job{
			Task: &job.TaskSpider{
				URL: url,
			},
		}
		return c.JSON(http.StatusOK, "1")
	}
	return c.JSON(http.StatusOK, "0")
}

//JsSDK 微信JS接口
func JsSDK(c echo.Context) error {
	url := c.QueryParam("url")

	js, _ := repository.GetSign(url)
	return c.JSON(http.StatusOK, js)
}

//Post 报料接口
func Post(c echo.Context) error {
	url := c.FormValue("url")
	// fmt.Println("url", url)
	if url != "" {
		err := repository.Post(url)
		if err != nil {
			return c.JSON(http.StatusOK, "0")
		}
		return c.JSON(http.StatusOK, "1")
	}
	return c.JSON(http.StatusOK, "0")
}

//GetContent 获取正文
func GetContent(c echo.Context) error {

	urlStr := c.QueryParam("url")

	info, err := reader.GetContent(urlStr)
	if err != nil {
		return c.String(http.StatusOK, "0")
	}

	input := []byte(info.Content)
	unsafe := blackfriday.MarkdownCommon(input)
	content := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	// html := fmt.Sprintf(`<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	// 						<link rel="preload" href="https://yize.gitlab.io/css/main.css" as="style" />
	// 						%v`, string(content[:]))
	// return c.HTML(http.StatusOK, html)
	info.Content = fmt.Sprintf(`%v`, string(content[:]))

	// 给图片加上 最大宽度
	info.Content = strings.Replace(info.Content, `<img src=`, `<img style="max-width:100%" src=`, -1)
	info.Content = strings.Replace(info.Content, `<section>`, `<div>`, -1)
	info.Content = strings.Replace(info.Content, `</section>`, `</div>`, -1)

	type Info struct {
		Title   string        `json:"title"`
		Content template.HTML `json:"content"`
		PubAt   string        `json:"pub_at"`
	}

	return c.JSON(http.StatusOK, Info{
		info.Title,
		template.HTML(info.Content),
		info.PubAt,
	})
}

//GetList 获取列表 临时放在这里面，做小程序测试api
func GetList(c echo.Context) error {
	urlStr := c.QueryParam("url")
	if urlStr == "" {
		return c.JSON(http.StatusOK, "0")
	}
	links, _ := reader.GetList(urlStr)
	return c.JSON(http.StatusOK, links)
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	//-------------------
	// Custom middleware
	//-------------------
	// Stats
	s := NewStats()
	e.Use(s.Process)
	e.GET("/stats", s.Handle) // Endpoint to get stats

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Welcome to api.readfollow.com !\n")
	})

	// 请求抓取
	e.GET("/fetch", Fetch)
	e.POST("/post", Post)

	// 获取公众号接口
	e.GET("/media/:id", Media)
	// 用户查看文章时请求该接口
	e.GET("/view/:id", View)
	// 赞同文章
	e.GET("/like/:id", Like)
	// 否定文章 如果否定比赞同多5票，评分为0
	e.GET("/hate/:id", Hate)

	// 获取微信文章接口
	e.GET("/article", Articles)

	// 获取微信文章接口
	e.GET("/new", NewArticles)
	e.GET("/hot", HotArticles)

	e.GET("/jssdk", JsSDK)

	// 临时做小程序api
	e.GET("/minapp/getlist", GetList)
	e.GET("/minapp/getcontent", GetContent)

	// 获取标签接口
	e.GET("/tags", Tags)
	e.GET("/searchtags", SearchTags)
	e.GET("/tag/:id", Tag)
	e.GET("/search", Search)
	e.GET("/gettagbymedia/:id", GetTagByMediaID)

	e.File("/favicon.ico", "favicon.ico")

	e.Logger.Fatal(e.Start(":8005"))

	// e.Logger.Fatal(e.StartAutoTLS(":443"))

}
