package utils

import "fmt"

//先求两个数的最大公约数（使用辗转相除法）
func getMaxCommonDivisor(a, b int) int {
	//定义一个交换站值
	var temp = 0
	fmt.Println("a, b", a, b)
	for a%b != 0 {
		temp = a % b
		a = b
		b = temp
	}
	return b
}

//求两个数的最小公倍数（两个数相乘 = 这两个数的最大公约数和最小公倍数的 积）
func getMinCommonMultiple(a, b int) int {
	return a * b / getMaxCommonDivisor(a, b)
}

//求多个数的最小公倍数
func GetMinMultiCommonMultiple(a []int) int {
	var val = a[0]
	//实现原理：拿前两个数的最小公约数和后一个数比较，求他们的公约数以此来推。。。
	for i := 1; i < len(a); i++ {
		val = getMinCommonMultiple(val, a[i])
	}
	return val
}
