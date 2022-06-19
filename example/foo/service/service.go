package service

import (
	"context"
	"time"
)

type FooService interface {
	Ping(ctx context.Context) (out string, err error)
	SaveUser(ctx context.Context, user User) (id int, err error)
	GetUserById(ctx context.Context, id int) (user User, err error)
	GetUserByIdStruct(ctx context.Context, in struct {
		Id   int
		Name string
	}) (err error)
	SaveArticle(ctx context.Context, article *Article) (id int, err error)
	//GetArticlesByTags(ctx context.Context, tags []string) (articles []Article, err error)
	//GetUserArticlesByTags(ctx context.Context, userId int, tags []string) (articles []Article, err error)
}

//type CID struct {
//	value string
//}
//
//func (p CID) String() string {
//	return p.value
//}

type User struct {
	Id        int
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Article struct {
	Id      int
	Title   string
	Meta    *Meta
	Tags    []Tag
	Content string
}

type Tag struct {
	Name  string
	Color string
}

type Meta struct {
	KeyWords []string
	Mapping  struct {
		Source string
		Target string
	}
}
