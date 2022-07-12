package logger

import (
	"fmt"
	"github.com/fatih/color"
	"log"
)

func Panic(err error) {
	log.Panicln(color.RedString(fmt.Sprintf("%s", err)))
}

func Error(err error) {
	log.Println(color.RedString(fmt.Sprintf("%s", err)))
}
