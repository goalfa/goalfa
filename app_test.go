package goalfa

import (
	"context"
	"testing"
)

type UserService struct {
}

func (u UserService) Register(ctx context.Context, username string) (err error) {
	return
}

func (u UserService) Forgot(ctx context.Context, username string) (err error) {
	return
}

func Test_registerService(t *testing.T) {
	//route := Route{}
	//service := UserService{}
	//out := make([]Route, 0)
	//app := New()
	//app.prepareService(&route, service, &out)
}
