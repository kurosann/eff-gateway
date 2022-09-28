package poll

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTypeMuxSubscription(t *testing.T) {
	stf := NewIPServer("测试服务")
	go stf.Subscription.WeightUpdate()
	var param int64
	var s = sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		s.Add(1)
		go func() {
			defer s.Done()
			atomic.AddInt64(&param, 1)
			ipaddr := fmt.Sprintf("127.0.0.1:800%d", param)
			stf.Add(ipaddr, 10)
			ifo := stf.GetAddr(ipaddr)
			if ifo != nil {
				ifo.AddCumulative(10+int(param), 1000+int(param))
			}
			stf.Subscription.SendMsg(stf)
		}()
	}
	s.Wait()
	// 新建一个定时任务对象
	// 根据cron表达式进行时间调度，cron可以精确到秒，大部分表达式格式也是从秒开始。
	//crontab := cron.New()  默认从分开始进行时间调度
	crontab := cron.New(cron.WithSeconds()) //精确到秒
	//定义定时器调用的任务函数

	task := func() {
		for _, info := range stf.IPList {
			atomic.AddInt64(&param, 10)
			info.AddCumulative(int(10+param), int(100+param))
		}

		for i := 0; i < len(stf.Rss.rss); i++ {
			fmt.Println("start")
			st := time.Now().UnixNano()
			stf.Rss.Next()
			fmt.Println("end", time.Now().UnixNano()-st)
		}
	}

	//定时任务
	spec := "*/5 * * * * ?" //cron表达式，每五秒一次
	// 添加定时任务,
	crontab.AddFunc(spec, task)
	// 启动定时器
	crontab.Start()
	// 定时任务是另起协程执行的,这里使用 select 简答阻塞.实际开发中需要
	// 根据实际情况进行控制
	select {} //阻塞主线程停止

}
