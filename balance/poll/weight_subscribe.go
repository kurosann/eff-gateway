package poll

import (
	"time"
)

type TypeMuxSubscription struct {
	readC chan *TypeMuxEvent
}

type TypeMuxEvent struct {
	Time time.Time
	Data interface{}
}

func Subscribe() *TypeMuxSubscription {
	return &TypeMuxSubscription{
		readC: make(chan *TypeMuxEvent),
	}
}

func (s *TypeMuxSubscription) SendMsg(v interface{}) {
	s.readC <- &TypeMuxEvent{
		Time: time.Now(),
		Data: v,
	}
}

func (s *TypeMuxSubscription) subChan() chan *TypeMuxEvent {
	return s.readC
}

func (s *TypeMuxSubscription) WeightUpdate() {

	r := s.subChan()
	for {
		select {
		case ev := <-r:
			if ev != nil {
				data := ev.Data.(*ServerTimeFragment)
				n := *data
				newData := &n
				newData.FragmentWeight()
				data.IPList = newData.IPList
				data.Rss = newData.Rss
			}
			break
		}
	}

}
