package service

import (
	"context"
	"encoding/json"
	"fmt"
)

type FooMock struct {
}

func (f FooMock) SaveArticle(ctx context.Context, article *Article) (id int, err error) {
	fmt.Println("SaveArticle In", article)
	return
}

func (f FooMock) GetUserByIdStruct(ctx context.Context, in struct{ Id int }) (err error) {
	d, _ := json.Marshal(in)
	fmt.Println("GetUserByIdStruct In", string(d))
	return
}

func (f FooMock) GetUserById(ctx context.Context, id int) (user User, err error) {
	fmt.Println("GetUserById In", id)
	return
}

func (f FooMock) Ping(ctx context.Context) (out string, err error) {
	out = "ok"
	return
}

func (f FooMock) SaveUser(ctx context.Context, user User) (id int, err error) {
	id = 1
	d, _ := json.Marshal(user)
	fmt.Println("SaveUser In", string(d))
	return
}
