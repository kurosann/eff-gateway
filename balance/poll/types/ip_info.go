package types

// IPTimeInfo
// 每个Ip的时间片段
type IPTimeInfo struct {
	Addr              string // ip地址
	Weight            int    // 权重
	CumulativeTime    int    // 累积计时间
	CumulativeRequest int    // 累积计请求数
	CacheTime         int    // 快照缓存累积计时间
	CacheRequest      int    // 快照缓存累积计请求数
	LastPerformance   int    // 上一次的性能
	IsActive          bool   // 该ip是否活跃
}

func NewIpInfo(Addr string, Weight int) *IPTimeInfo {
	newObj := new(IPTimeInfo)
	newObj.IsActive = true
	newObj.Addr = Addr
	newObj.Weight = Weight
	return newObj
}

// RemoveCumulative
// 移除累积
func (itf *IPTimeInfo) RemoveCumulative() bool {
	// 进行锁
	itf.CumulativeTime = 0 + itf.CacheTime
	itf.CumulativeRequest = 0 + itf.CacheRequest
	itf.CacheTime = 0
	itf.CacheRequest = 0
	return true
}

// AddCumulative
// 添加累积
func (itf *IPTimeInfo) AddCumulative(t, req int) {

	itf.CacheTime += t
	itf.CacheRequest += req
}

// WriteLastPerformance
// 添加累积
func (itf *IPTimeInfo) WriteLastPerformance(lp int) {
	itf.LastPerformance = lp
}
