package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/GanEasy/wxrankapi/job"
	"github.com/GanEasy/wxrankapi/repository"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	order := c.QueryParam("order")

	if limit <= 0 || limit > 100 {
		limit = 10
	}
	// limit = 10
	if offset < 0 || offset > 500 {
		offset = 0
	}

	articles, err := repository.GetArticle(limit, offset, tag, order)

	if err != nil {

	}

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

//View 阅读
func View(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	article, err := repository.View(id)

	if err != nil {

	}

	return c.JSON(http.StatusOK, article)
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

	// 用户查看文章时请求该接口
	e.GET("/view/:id", View)

	// 获取微信文章接口
	e.GET("/article", Articles)

	// 获取标签接口
	e.GET("/tags", Tags)

	e.File("/favicon.ico", "favicon.ico")

	e.Logger.Fatal(e.Start(":8005"))

	// e.Logger.Fatal(e.StartAutoTLS(":443"))

}
