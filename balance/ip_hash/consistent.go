package ip_hash

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type Consistent struct {
	//排序的hash虚拟结点
	hashSortedNodes []uint32
	//虚拟结点对应的结点信息
	circle map[uint32]string
	//已绑定的结点
	nodes map[string]bool
	//map读写锁
	sync.RWMutex
	//虚拟结点数
	virtualNodeCount int
}

// hashKey
// key值运算
func (c *Consistent) hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

// Add
// 增加节点
func (c *Consistent) Add(node string, virtualNodeCount int) error {
	if node == "" {
		return nil
	}
	c.Lock()
	defer c.Unlock()

	if c.circle == nil {
		c.circle = map[uint32]string{}
	}
	if c.nodes == nil {
		c.nodes = map[string]bool{}
	}

	if _, ok := c.nodes[node]; ok {
		return errors.New("node already existed")
	}
	c.nodes[node] = true
	//增加虚拟结点
	for i := 0; i < virtualNodeCount; i++ {
		virtualKey := c.hashKey(node + strconv.Itoa(i))
		c.circle[virtualKey] = node
		c.hashSortedNodes = append(c.hashSortedNodes, virtualKey)
	}

	//虚拟结点排序
	sort.Slice(c.hashSortedNodes, func(i, j int) bool {
		return c.hashSortedNodes[i] < c.hashSortedNodes[j]
	})

	return nil
}

// GetNode
// 获取节点
func (c *Consistent) GetNode(key string) string {
	c.RLock()
	defer c.RUnlock()

	hash := c.hashKey(key)
	i := c.getPosition(hash)

	return c.circle[c.hashSortedNodes[i]]
}

func (c *Consistent) AddReqs(ipAddr string, v ...int) error {

	return nil
}
func (c *Consistent) SendMsg() error {

	return nil
}

// getPosition
// 获取位置
func (c *Consistent) getPosition(hash uint32) int {
	i := sort.Search(len(c.hashSortedNodes), func(i int) bool { return c.hashSortedNodes[i] >= hash })

	if i < len(c.hashSortedNodes) {
		if i == len(c.hashSortedNodes)-1 {
			return 0
		} else {
			return i
		}
	} else {
		return len(c.hashSortedNodes) - 1
	}
}
