package common

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/astaxie/beego"
	"github.com/chanyipiaomiao/hltool"
)

var (

	// Logger 日志对象
	Logger *hltool.HLogger
)

// InitLog 初始化日志
func init() {
	var logpath = beego.AppConfig.String("log::logPath")
	var currentDir string
	if !strings.HasPrefix(logpath, "/") {
		var err error
		currentDir, err = os.Getwd()
		if err != nil {
			log.Fatalf("error: %s", err)
		}
	}
	logpath = path.Join(currentDir, logpath)

	hlog, err := hltool.NewHLog(logpath)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	Logger, err = hlog.GetLogger()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

// GetLogger 返回Logger
func GetLogger() *hltool.HLogger {
	return Logger
}
