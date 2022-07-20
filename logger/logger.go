package logger

import (
	"fmt"
	"github.com/fatih/color"
)

func Panic(err error) {
	panic(color.RedString(fmt.Sprintf("%s", err)))
}

func Error(err error) {
	fmt.Println(color.RedString("❌[错误]:"), err)
}

func Warn(msg string) {
	fmt.Println(color.YellowString("⚠️[警告]:"), msg)
}

func Success(msg string) {
	fmt.Println(msg)
}
