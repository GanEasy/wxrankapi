package orm

import "time"

// Tag 标签属性 (不存在相同的标签)
type Tag struct {
	ID        uint `gorm:"primary_key"`
	Pid       uint
	IsShow    int
	Type      string
	Title     string
	Name      string `gorm:"type:varchar(100);unique_index"`
	Intro     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// GetTagByID 获取 Tag
func (tag *Tag) GetTagByID(id int) {
	DB().First(tag, id)
}

// GetTagByName 通过名词获取标签
func (tag *Tag) GetTagByName(name string) {
	DB().Where(Tag{Name: name}).FirstOrCreate(tag)
}

// Save 保存用户信息
func (tag *Tag) Save() {
	DB().Save(&tag)
}
