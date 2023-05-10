package consistenthash

import (
	"hash/crc32"
	"sort"
)

type HashFunc func([]byte) uint32

type NodeMap struct {
	hashFunc HashFunc
	nodeHashes []int
	nodeHashMap map[int]string
}

func NewNodeMap(hashFunc HashFunc) *NodeMap{
	m := &NodeMap{
		hashFunc: hashFunc,
		nodeHashMap: make(map[int]string),
	}

	if hashFunc == nil{
		m.hashFunc = crc32.ChecksumIEEE
	}

	return m
}

func (m *NodeMap) IsEmpty() bool{
	return len(m.nodeHashes) == 0
}

//增加新节点
func (m *NodeMap) AddNode(keys ...string){
	for _, key := range keys{
		if key == ""{
			continue
		}
		hash := int(m.hashFunc([]byte(key)))
		m.nodeHashes = append(m.nodeHashes, hash)
		m.nodeHashMap[hash] = key
	}

	sort.Ints(m.nodeHashes)
}

//为传进来的key选节点
func (m *NodeMap) PickNode(key string) string{
	if m.IsEmpty(){
		return ""
	}

	hash := int(m.hashFunc([]byte(key)))
	idx := sort.Search(len(m.nodeHashes), func(i int) bool {
		return m.nodeHashes[i] >= hash
	})
	if idx == len(m.nodeHashes){
		idx = 0
	}

	return m.nodeHashMap[m.nodeHashes[idx]]
}