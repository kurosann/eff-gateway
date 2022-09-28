package balance

import (
	"fmt"
	"testing"
	"time"
)

func TestStrategy(t *testing.T) {
	sg := NewStrategy()
	sg.AddStrategy("ces")
	var myServerConf = sg.GetServer("ces")
	myServerConf.impl.Add("127.0.0.1:8001", 10)

	myServerConf.impl.SendMsg()
	myServerConf.impl.Add("127.0.0.1:8002", 14)
	myServerConf.impl.Add("127.0.0.1:8003", 13)
	myServerConf.impl.Add("127.0.0.1:8004", 14)
	myServerConf.impl.Add("127.0.0.1:8005", 13)
	myServerConf.impl.AddReqs("127.0.0.1:8004", 100, 200)
	myServerConf.impl.AddReqs("127.0.0.1:8005", 100, 200)
	myServerConf.impl.SendMsg()

	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	fmt.Println(myServerConf.impl.GetNode("ces"))
	<-time.After(5 * time.Second)
}
