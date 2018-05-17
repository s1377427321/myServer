package common

import (
	"os"
	"fmt"
	"errors"
	"github.com/astaxie/beego"
	"sync"
)

type CfgLog struct {
	LogDir string
	LogFileName string
	LogToConsole bool
	LogLevel int
}

var (
	CurrentLogLevel int
	logLevelLock sync.Mutex
)

func SetLogOption(logDir, namePrefix string, logLevel int, logToConsole bool) error {
	if fi, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.Mkdir(logDir, 0666)
		if err != nil {
			panic(err)
		}

		fmt.Println("The log dir:", logDir, "doesn't exist, create it!")
	} else {
		if !fi.IsDir() {
			panic(errors.New(fmt.Sprintf("The file:", logDir, "is not a directory!")))
		}
	}

	if logLevel < beego.LevelEmergency || logLevel > beego.LevelDebug {
		panic(errors.New(fmt.Sprintf("Invalid logLevel:", logLevel)))
	}

	beego.BeeLogger.SetLevel(logLevel)
	logLevelLock.Lock()
	CurrentLogLevel = logLevel
	logLevelLock.Unlock()

	str := fmt.Sprintf(`{"filename":"%s/%s.log"}`, logDir, namePrefix)
	err := beego.BeeLogger.SetLogger("file", str)

	if err != nil {
		panic(err)
	}

	if logToConsole {
		err = beego.BeeLogger.SetLogger("console", str)
		if err != nil {
			panic(err)
		}
	} else {
		beego.BeeLogger.DelLogger("console")
	}

	beego.SetLogFuncCall(true)
	return nil
}