package ip_hash

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"testing"
)

func Test_ConsistentHash(t *testing.T) {
	virtualNodeList := []int{100, 150, 200}
	//测试10台服务器
	nodeNum := 10
	//测试数据量100W
	testCount := 1000000
	for _, virtualNode := range virtualNodeList {
		consistentHash := &Consistent{}
		distributeMap := make(map[string]int64)
		for i := 1; i <= nodeNum; i++ {
			serverName := "172.17.0." + strconv.Itoa(i)
			consistentHash.Add(serverName, virtualNode)
			distributeMap[serverName] = 0
		}
		//测试100W个数据分布
		for i := 0; i < testCount; i++ {
			testName := "testName"
			serverName := consistentHash.GetNode(testName + strconv.Itoa(i))
			distributeMap[serverName] = distributeMap[serverName] + 1
		}

		var keys []string
		var values []float64
		for k, v := range distributeMap {
			keys = append(keys, k)
			values = append(values, float64(v))
		}
		sort.Strings(keys)
		fmt.Printf("####测试%d个结点,一个结点有%d个虚拟结点,%d条测试数据\n", nodeNum, virtualNode, testCount)
		for _, k := range keys {
			fmt.Printf("服务器地址:%s 分布数据数:%d\n", k, distributeMap[k])
		}
		fmt.Printf("标准差:%f\n\n", getStandardDeviation(values))
	}
}

//获取标准差
func getStandardDeviation(list []float64) float64 {
	var total float64
	for _, item := range list {
		total += item
	}
	//平均值
	avg := total / float64(len(list))

	var dTotal float64
	for _, value := range list {
		dValue := value - avg
		dTotal += dValue * dValue
	}

	return math.Sqrt(dTotal / avg)
}
