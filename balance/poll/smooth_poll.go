package poll

import (
	"errors"
)

type WeightRoundRobinBalance struct {
	curIndex int           // 下标
	rss      []*WeightNode // 权重节点
	rsw      []int         // 权重

	conf LoadBalanceConf // 观察主体
}

type WeightNode struct {
	addr            string // 服务器地址
	weight          int    // 权重值
	currentWeight   int    // 节点当前权重
	effectiveWeight int    // 有效权重
}

// LoadBalanceConf
// 配置主体
type LoadBalanceConf interface {
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
}

func (r *WeightRoundRobinBalance) Add(addr string, weight int) error {
	if addr == "" {
		return errors.New("addr is ''")
	}

	if weight == 0 {
		return errors.New("addr is 0")
	}

	node := &WeightNode{
		addr:            addr,
		weight:          weight,
		effectiveWeight: weight,
	}
	r.rsw = append(r.rsw, weight)
	r.rss = append(r.rss, node)
	return nil
}

func (r *WeightRoundRobinBalance) Next() string {
	total := 0
	var best *WeightNode

	// currentWeight effecitveWeight  	选中的节点   		请求后current
	// (4,3,2)  		(4,3,2)		 	 	A 			  (-1,6,4)
	// (-1,6,4)  		(4,3,2)		 		b 			  (3,0,6)
	// (3,0,6)  		(4,3,2)		 		c 			  (7,3,-1)

	// 计算：(4+4-9,3+3,2+2) = (-1,6,4)
	// 计算：(-1+4,6+3-9,4+2) = (3,0,6)
	// 计算：(3+4,0+3,6+2-9) = (3,0,6)
	for i := 0; i < len(r.rss); i++ {
		w := r.rss[i]
		// 统计所有有效权重之和
		total += w.effectiveWeight
		// 变更节点临时权重为的节点临时权重+节点有效权重
		w.currentWeight += w.effectiveWeight
		// 有效权重默认与权重相同，通讯异常时-1,通讯成功+1，直到weight大小恢复
		if w.effectiveWeight < w.weight {
			w.effectiveWeight++
		}
		// 最大临时权重节点
		if best == nil || w.currentWeight > best.currentWeight {
			best = w
		}
	}

	if best == nil {
		return ""
	}
	best.currentWeight -= total
	//fmt.Println(best.addr)
	return best.addr
}

func (r *WeightRoundRobinBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}
