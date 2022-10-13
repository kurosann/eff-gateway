package balance

import (
	"go-gateway/balance/ip_hash"
	"go-gateway/balance/poll"
	"sync"
)

var GlobalStrategy = &StrategyRegister{}

type StrategyRegister struct {
	StrategyFunc string
	ServerMap    sync.Map //存储服务策略
}

type PollStrategy interface {
	Add(addr string, weight int) error
	GetNode(serverName string) string
	AddReqs(addr string, v ...int) error
	Update() error
	GetCycles() int
}

type Strategy struct {
	Impl PollStrategy
}

func (r *Strategy) SetStrategy(ps PollStrategy) {
	r.Impl = ps
}

func init() {
	GlobalStrategy = NewStrategy()
}

func NewStrategy() *StrategyRegister {
	return &StrategyRegister{
		StrategyFunc: "smooth_poll",
		ServerMap:    sync.Map{},
	}
}

func (sr *StrategyRegister) AddStrategy(sName string) {

	switch sr.StrategyFunc {
	case "ip_hash":
		iph := &ip_hash.Consistent{}
		traveler := Strategy{}
		traveler.SetStrategy(iph)
		sr.ServerMap.Store(sName, traveler)
		break
	case "smooth_poll":
		sr.pollServer(sName)
		break
	default:
		sr.pollServer(sName)
		break

	}
}

func (sr *StrategyRegister) pollServer(sName string) {
	stf := poll.NewIPServer(sName)
	go stf.Subscription.WeightUpdate()
	traveler := &Strategy{}
	traveler.SetStrategy(stf)
	sr.ServerMap.Store(sName, traveler)
}

func (sr *StrategyRegister) GetServer(sName string) *Strategy {

	var v, k = sr.ServerMap.Load(sName)
	if k {
		return v.(*Strategy)
	}
	return nil
}
