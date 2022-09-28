package poll

import (
	"errors"
	"fmt"
	"go-gateway/balance/poll/types"
)

// ServerTimeFragment
// 服务组的时间片段
type ServerTimeFragment struct {
	ServerName    string                   // 服务名
	IPList        []*types.IPTimeInfo      // ip地址-map
	Rss           *WeightRoundRobinBalance // 权重循环平衡
	Subscription  *TypeMuxSubscription     // 通知计算
	DefaultWeight int
	TotalLP       int
}

// NewIPServer
// 初始化服务
func NewIPServer(serverName string) *ServerTimeFragment {

	return &ServerTimeFragment{
		ServerName:   serverName,
		Subscription: Subscribe(),
		Rss:          new(WeightRoundRobinBalance),
		IPList:       make([]*types.IPTimeInfo, 0),
	}
}

func (it *ServerTimeFragment) Add(addr string, weight int) error {
	if addr == "" {
		return errors.New("addr is ''")
	}
	if weight == 0 && it.DefaultWeight > 0 {
		weight = it.DefaultWeight
	}

	it.IPList = append(it.IPList, types.NewIpInfo(addr, weight))
	it.Rss.Add(addr, weight)
	return nil
}

func (it *ServerTimeFragment) AddReqs(ipAddr string, v ...int) error {
	info := it.GetAddr(ipAddr)
	if info != nil {
		info.AddCumulative(v[0], v[1])
	}
	return nil
}
func (it *ServerTimeFragment) SendMsg() error {

	it.Subscription.SendMsg(it)
	return nil
}

func (it *ServerTimeFragment) RemoveCumulative() {
	for _, v := range it.IPList {
		v.RemoveCumulative()
	}
}

// FragmentWeight
// 片段权重
func (it *ServerTimeFragment) FragmentWeight() {
	// 计算方法：
	// 不同服务器的性能对每个请求的处理时间都是不一样的
	// 每个请求需要进行的处理过程也不一样，因此需要进行多次计算取值
	// 平均每个请求的性能 = [(CumulativeRequest/CumulativeTime ) + lastPerformance]/2
	// 即：p = (t/r+lp)/2
	if it.IPList == nil {
		return
	}
	it.TotalLP = 0
	for i := 0; i < len(it.IPList); i++ {
		info := it.IPList[i]
		if info.LastPerformance == 0 || info.CumulativeTime == 0 || info.CumulativeRequest == 0 {
			// 第一次运行默认为 1
			fmt.Println("第一次运行 ", info.Addr)
			info.LastPerformance = 1
			it.TotalLP += 1
			continue
		}
		fmt.Println("CumulativeTime：", info.CumulativeTime)
		fmt.Println("CumulativeRequest：", info.CumulativeRequest)
		fmt.Println("LastPerformance：", info.LastPerformance)
		// 已经有上次的片段开始计算
		info.LastPerformance = ((info.CumulativeRequest/info.CumulativeTime)+info.LastPerformance)/2 + 1
		it.IPList[i] = info
		// 存储每个服务器的性能值
		it.TotalLP += info.LastPerformance
	}
	it.weightProportionSort()
}

// WeightProportionSort
// 性能权重比例、下一次的ip执行请求的排序
func (it *ServerTimeFragment) weightProportionSort() {
	it.Rss = new(WeightRoundRobinBalance)
	it.DefaultWeight = 0
	for _, info := range it.IPList {
		info.Weight = info.LastPerformance * 100 / it.TotalLP
		fmt.Println("LastPerformance: ", info.LastPerformance)
		fmt.Println("TotalLP: ", it.TotalLP)
		fmt.Println("ip: ", info.Addr)
		fmt.Println("weight: ", info.Weight)
		it.DefaultWeight += info.Weight
		it.Rss.Add(info.Addr, info.Weight)
	}
	it.DefaultWeight = it.DefaultWeight / len(it.IPList)
	fmt.Println("平均权重: ", it.DefaultWeight)
}

func (it *ServerTimeFragment) GetNode(serverName string) string {
	if it.ServerName != serverName {
		return ""
	}
	return it.Rss.Next()
}

func (it *ServerTimeFragment) GetAddr(ipAddr string) *types.IPTimeInfo {
	for _, info := range it.IPList {
		if info.Addr == ipAddr {

			if info == nil {
				fmt.Println("info:nil")
			}
			if info.CumulativeRequest == 0 {
				fmt.Println("CumulativeRequest:0")
			}
			if info.CumulativeTime == 0 {
				fmt.Println("CumulativeTime:0")
			}
			return info
		}
	}
	return nil
}