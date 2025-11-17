package model

import (
	"time"

	"gorm.io/gorm"
)

type ModelEvent struct {
	Type      string // "user_created", "user_updated", "user_deleted", etc.
	ModelType string // "user", "book"
	ModelID   uint
	Payload   interface{}
}

// 全局事件通道
var EventChannel = make(chan ModelEvent, 100)

// Book 对应表: books
type Book struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Title              string    `gorm:"column:title;type:text;not null" json:"title"`
	Author             string    `gorm:"column:author;type:text;not null" json:"author"`
	Cover              string    `gorm:"column:cover;type:text" json:"cover"`
	Category           string    `gorm:"column:category;type:text;not null" json:"category"`
	Channel            string    `gorm:"column:channel;type:text;not null;default:male" json:"channel"` // male/female
	Description        string    `gorm:"column:description;type:text" json:"description"`
	Status             string    `gorm:"column:status;type:text;not null;default:serializing" json:"status"` // serializing/completed
	WordCount          int       `gorm:"column:word_count;type:integer;default:0" json:"wordCount"`
	ClickCount         int       `gorm:"column:click_count;type:integer;default:0" json:"clickCount"`
	CommentCount       int       `gorm:"column:comment_count;type:integer;default:0" json:"commentCount"`
	LatestChapterID    *uint     `gorm:"column:latest_chapter_id;type:integer" json:"latestChapterId"` // 可为 NULL
	LatestChapterTitle string    `gorm:"column:latest_chapter_title;type:text" json:"latestChapterTitle"`
	CreateTime         time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime         time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Chapters           []Chapter `gorm:"foreignKey:BookID;references:ID" json:"chapters,omitempty"`
}

// Chapter 对应表: chapters
type Chapter struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"` // 数据库主键（自增）
	BookID     uint      `gorm:"column:book_id;not null;uniqueIndex:idx_book_chapter" json:"bookId"`
	ChapterID  uint      `gorm:"column:chapter_id;not null;uniqueIndex:idx_book_chapter" json:"chapterId"` // 章节序号：1,2,3...
	Title      string    `gorm:"column:title;type:text;not null" json:"title"`
	Content    string    `gorm:"-" json:"content"`
	IsVip      bool      `gorm:"column:is_vip;not null;default:false" json:"isVip"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

// User 对应表: users
type User struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Username   string    `gorm:"column:username;type:text;not null;unique" json:"username"`
	Password   string    `gorm:"column:password;type:text;not null" json:"password"`
	Nickname   string    `gorm:"column:nickname;type:text" json:"nickname"`
	Avatar     string    `gorm:"column:avatar;type:text" json:"avatar"`
	Type       string    `gorm:"column:type;type:text;not null;default:normal" json:"type"` // 作家/读者/管理员
	Desc       string    `gorm:"column:desc;type:text" json:"desc"`
	Level      int       `gorm:"column:level;type:integer;default:1" json:"level"`
	Sex        int       `gorm:"column:sex;type:integer;default:0" json:"sex"`
	IsVip      bool      `gorm:"column:is_vip;type:integer;default:false" json:"isVip"`
	Location   string    `gorm:"column:location;type:text" json:"location"`
	Status     int       `gorm:"column:status;type:integer;default:1" json:"status"`
	Email      string    `gorm:"column:email;type:text" json:"email"`
	Phone      string    `gorm:"column:phone;type:text" json:"phone"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

// Shelf 对应表: shelf
type Shelf struct {
	ID       uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserID   uint      `gorm:"column:user_id;type:integer;not null" json:"userId"`
	BookID   uint      `gorm:"column:book_id;type:integer;not null" json:"bookId"`
	AddTime  time.Time `gorm:"column:add_time;autoCreateTime" json:"addTime"`
	LastRead time.Time `gorm:"column:last_read;autoUpdateTime" json:"lastRead"`
	Book     Book      `gorm:"foreignKey:BookID;references:ID" json:"book,omitempty"`
}

// History 对应表: histories （注意：SQL 表名是 histories）
type History struct {
	ID              uint      `gorm:"primaryKey;autoIncrement;column:id;table:histories" json:"id"`
	UserID          uint      `gorm:"column:user_id;type:integer;not null" json:"userId"`
	BookID          uint      `gorm:"column:book_id;type:integer;not null" json:"bookId"`
	ChapterID       uint      `gorm:"column:chapter_id;type:integer;not null" json:"chapterId"`
	ReadingProgress int       `gorm:"column:reading_progress;type:integer;default:0" json:"readingProgress"`
	UpdateTime      time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Book            Book      `gorm:"foreignKey:BookID;references:ID" json:"book,omitempty"`
	Chapter         Chapter   `gorm:"foreignKey:ChapterID;references:ID" json:"chapter,omitempty"`
}

// 推荐表
type Rcmd struct {
	ID       uint   `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	RcmdType string `gorm:"column:rcmd_type;type:text;not null" json:"rcmdType"` // hot/top/curated/featured
	BookID   uint   `gorm:"column:book_id;type:integer;not null" json:"bookId"`
	Order    int    `gorm:"column:order;type:integer;default:0" json:"order"`
	Book     Book   `gorm:"foreignKey:BookID;references:ID" json:"book,omitempty"`
}

// Rank 对应表: ranks
type Rank struct {
	ID       uint   `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	RankType string `gorm:"column:rank_type;type:text;not null" json:"rankType"` // clicks/new/update/comment
	BookID   uint   `gorm:"column:book_id;type:integer;not null" json:"bookId"`
	Order    int    `gorm:"column:order;type:integer;default:0" json:"order"`
	Book     Book   `gorm:"foreignKey:BookID;references:ID" json:"book,omitempty"`
}

// 评论表
type Comment struct {
	ID           uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	BookID       uint      `gorm:"column:book_id;type:integer;not null" json:"bookId"`
	UserID       uint      `gorm:"column:user_id;type:integer;not null" json:"userId"`
	Content      string    `gorm:"column:content;type:text;not null" json:"content"`
	ReplyCount   int       `gorm:"column:reply_count;type:integer;default:0" json:"replyCount"`
	LikeCount    int       `gorm:"column:like_count;type:integer;default:0" json:"likeCount"`
	DislikeCount int       `gorm:"column:dislike_count;type:integer;default:0" json:"dislikeCount"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	User         User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Book         Book      `gorm:"foreignKey:BookID;references:ID" json:"book,omitempty"`
}

// CommentIndex 用于评论索引的数据结构
type CommentIndex struct {
	ID           uint      `json:"id"`
	BookTitle    string    `json:"bookTitle"`
	UserNickname string    `json:"userNickname"`
	Content      string    `json:"content"`
	CreateTime   time.Time `json:"createTime"`
}

// InfoJSON 非数据库结构
type InfoJSON struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Category    string `json:"category"`
	Channel     string `json:"channel"`
	Description string `json:"description"`
	Status      string `json:"status"`
	WordCount   int    `json:"wordCount"`
	Cover       string `json:"cover"`
}

func (u *User) AfterUpdate(tx *gorm.DB) error {
	EventChannel <- ModelEvent{
		Type:      "updated",
		ModelType: "user",
		ModelID:   u.ID,
	}
	return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
	EventChannel <- ModelEvent{
		Type:      "created",
		ModelType: "user",
		ModelID:   u.ID,
	}
	return nil
}

func (u *User) AfterDelete(tx *gorm.DB) error {
	EventChannel <- ModelEvent{
		Type:      "deleted",
		ModelType: "user",
		ModelID:   u.ID,
	}
	return nil
}

func (b *Book) AfterUpdate(tx *gorm.DB) error {
	EventChannel <- ModelEvent{
		Type:      "updated",
		ModelType: "book",
		ModelID:   b.ID,
	}
	return nil
}

func (b *Book) AfterCreate(tx *gorm.DB) error {
	EventChannel <- ModelEvent{
		Type:      "created",
		ModelType: "book",
		ModelID:   b.ID,
	}
	return nil
}

func (b *Book) AfterDelete(tx *gorm.DB) error {
	EventChannel <- ModelEvent{
		Type:      "deleted",
		ModelType: "book",
		ModelID:   b.ID,
	}
	return nil
}

// AfterUpdate Comment模型更新后的钩子函数
func (c *Comment) AfterUpdate(tx *gorm.DB) error {
	EventChannel <- ModelEvent{
		Type:      "updated",
		ModelType: "comment",
		ModelID:   c.ID,
	}
	return nil
}

// AfterCreate Comment模型创建后的钩子函数
func (c *Comment) AfterCreate(tx *gorm.DB) error {
	EventChannel <- ModelEvent{
		Type:      "created",
		ModelType: "comment",
		ModelID:   c.ID,
	}
	return nil
}

// AfterDelete Comment模型删除后的钩子函数
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	EventChannel <- ModelEvent{
		Type:      "deleted",
		ModelType: "comment",
		ModelID:   c.ID,
	}
	return nil
}
