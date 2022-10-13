package balance

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"math/rand"
	"strconv"
	"testing"
)

func TestStrategy(t *testing.T) {
	sg := NewStrategy()
	sg.AddStrategy("ces")
	var myServerConf = sg.GetServer("ces")
	for i := 1; i < 6; i++ {
		myServerConf.Impl.Add("127.0.0.1:800"+strconv.Itoa(i), rand.Intn(100)*i*i)
		myServerConf.Impl.Update()
	}

	task := func() {

		for i := 1; i < 6; i++ {
			myServerConf.Impl.AddReqs("127.0.0.1:800"+strconv.Itoa(i), rand.Intn(100)*i*i, rand.Intn(1000)*i*i)
		}
		myServerConf.Impl.Update()

		for i := 0; i < myServerConf.Impl.GetCycles(); i++ {
			fmt.Println("服务器地址：", myServerConf.Impl.GetNode("ces"))
		}
	}
	crontab := cron.New(cron.WithSeconds()) //精确到秒
	//定时任务
	spec := "*/5 * * * * ?" //cron表达式，每五秒一次
	// 添加定时任务,
	crontab.AddFunc(spec, task)
	// 启动定时器
	crontab.Start()
	select {} //阻塞主线程停止
}
