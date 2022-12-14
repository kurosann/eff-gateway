package glog

import (
	"eff-gateway/setting"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	InfoLog  *log.Logger
	WarnFile *log.Logger
	ErrorLog *log.Logger
)

func init() {
	errFile, err := os.OpenFile(fmt.Sprintf("./%s/errors.log", setting.Config.Log.Path), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	warnFile, err := os.OpenFile(fmt.Sprintf("./%s/warning.log", setting.Config.Log.Path), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	infoFile, err := os.OpenFile(fmt.Sprintf("./%s/info.log", setting.Config.Log.Path), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开日志文件失败")
	}
	InfoLog = log.New(io.MultiWriter(os.Stderr, infoFile), "[Info]:", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	WarnFile = log.New(io.MultiWriter(os.Stderr, warnFile), "[Warn]:", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	ErrorLog = log.New(io.MultiWriter(os.Stderr, errFile), "[Err]:", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
}
