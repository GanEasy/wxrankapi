package orm

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

// Article 文章列表
type Article struct {
	ID        uint      `gorm:"primary_key"`
	Title     string    // 标题
	Author    string    // 作者
	Cover     string    // 封面
	Intro     string    // 介绍
	PubAt     time.Time // 微信文章发布时间
	MediaID   uint
	View      int64         `gorm:"default:0"`                      // 点击次数，通过它进行计算排名
	URL       string        `gorm:"type:varchar(255);unique_index"` // 微信文章地址
	Rank      float64       `sql:"index"`                           // 排行
	Tags      pq.Int64Array `gorm:"type:int[]"`                     // 标签
	State     int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// GetArticleByID 获取Article
func (article *Article) GetArticleByID(id int) {
	DB().First(article, id)
}

// GetArticleByURL 通过url获取Article 如果没有的话进行初始化 (注：此url由文章详细页获得)
func (article *Article) GetArticleByURL(url string) {
	DB().Where(Article{URL: url}).FirstOrCreate(article)
}

// Save 保存用户信息
func (article *Article) Save() {
	DB().Save(&article)
}

// Hot 热门
func (article *Article) Hot(limit, offset int) (articles []Article) {
	DB().Offset(offset).Limit(limit).Order("rank DESC").Find(&articles)
	return
}

// New 最新
func (article *Article) New(limit, offset int) (articles []Article) {
	DB().Offset(offset).Limit(limit).Order("id DESC").Find(&articles)
	return
}

// GetArticle 获取文章
func (article *Article) GetArticle(limit, offset, tag int, order string) (articles []Article) {
	// var selectTag pq.Int64Array
	if tag != 0 {
		// selectTag = append(selectTag, int64(tag))
		//Article{Tags: selectTag}   "tags && {?}", selectTag
		DB().Where("tags @> ?", fmt.Sprintf("{%d}", tag)).Offset(offset).Limit(limit).Order(order).Find(&articles)
	} else {
		DB().Offset(offset).Limit(limit).Order(order).Find(&articles)
	}
	return
}

// GetHotArticleByTag 获取文章
func (article *Article) GetHotArticleByTag(limit, offset, tag int) (articles []Article) {
	var selectTag pq.Int64Array
	selectTag = append(selectTag, int64(tag))
	DB().Where(Article{Tags: selectTag}).Offset(offset).Limit(limit).Order("rank DESC").Find(&articles)
	return
}

// GetNewArticleByTag 获取文章
func (article *Article) GetNewArticleByTag(limit, offset, tag int) (articles []Article) {
	var selectTag pq.Int64Array
	selectTag = append(selectTag, int64(tag))
	DB().Where(Article{Tags: selectTag}).Offset(offset).Limit(limit).Order("id DESC").Find(&articles)
	return
}
