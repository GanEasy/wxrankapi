package repository

import (
	"strconv"
	"strings"

	"github.com/GanEasy/wxrankapi/orm"
)

func init() {

	// var initTag = []string{
	// 	"文博",
	// 	"汽车",
	// 	"IT",
	// 	"生活",
	// 	"自媒体",
	// 	"其它",
	// }

	// for _, v := range initTag {
	// 	var tag orm.Tag
	// 	tag.GetTagByName(v)
	// 	tag.Type = "cate"
	// 	tag.Title = v
	// 	tag.IsShow = 1
	// 	tag.Save()
	// }

}

//GetTagByType ..通过属性获取标签
func GetTagByType(name string) (tags []orm.Tag, err error) {
	if name != "" {
		var tag orm.Tag
		tags = tag.GetTagsByType(name)
	}
	return
}

//GetTagsByTitle ..通过属性获取标签
func GetTagsByTitle(name string) (tags []orm.Tag, err error) {
	if name != "" {
		var tag orm.Tag
		tags = tag.GetTagsByTitle(name)
	}
	return
}

//GetTagsByIDS ..通过id获取标签
func GetTagsByIDS(idstr string) (tags []orm.Tag, err error) {
	var ids []int64
	arr := strings.Split(idstr, ",")
	for _, id := range arr {
		c, e := strconv.Atoi(id)
		if e == nil {
			ids = append(ids, int64(c))
		}
	}
	if len(ids) > 0 {
		var tag orm.Tag
		tags = tag.GetTagsByIDS(ids)
	}
	return
}
