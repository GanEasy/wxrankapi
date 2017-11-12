package repository

import "github.com/GanEasy/wxrankapi/orm"

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
