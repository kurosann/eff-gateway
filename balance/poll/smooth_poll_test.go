package poll

import (
	"fmt"
	"testing"
)

func TestLB(t *testing.T) {
	// 测试平滑轮询
	rb := &WeightRoundRobinBalance{}
	rb.Add("127.0.0.1:2001", 1)
	rb.Add("127.0.0.1:2002", 2)
	rb.Add("127.0.0.1:2003", 3)
	rb.Add("127.0.0.1:2004", 4)

	for i := 0; i < 20; i++ {
		fmt.Println(rb.Next())
	}

}
