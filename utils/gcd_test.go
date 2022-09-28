package utils

import (
	"fmt"
	"testing"
)

func TestGcd(t *testing.T) {
	// 最小公倍数计算
	var arr = []int{1, 2, 3, 4}
	var s = GetMinMultiCommonMultiple(arr)
	fmt.Println(s)
}
