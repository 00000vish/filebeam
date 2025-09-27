package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/00000vish/filebeam/internals/config"
)

var isInit = true
var silent = false

func Create(config config.Config) bool {
	silent = config.Silent

	f, err := os.OpenFile("filebeam.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	isInit = err == nil

	if isInit {
		log.SetOutput(f)
	} else {
		fmt.Println(err)
	}

	return isInit
}

func Log(toLog string) {
	if isInit {
		log.Println(toLog)
	}

	if !silent {
		fmt.Println(toLog)
	}
}
