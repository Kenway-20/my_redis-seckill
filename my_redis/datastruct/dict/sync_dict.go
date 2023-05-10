package dict

import "sync"

type SyncDict struct {
	m sync.Map
}

func MakeSyncDict() *SyncDict{
	return &SyncDict{}
}

func (dict *SyncDict) Get(key string) (val interface{}, exists bool) {
	val, ok := dict.m.Load(key)
	return val, ok
}

func (dict *SyncDict) Len() int{
	cnt := 0
	//range里面传方法表示对每个k-v执行一般这个方法，返回true时就继续执行，直到返回false或者遍历结束，这里用来获取sync.map的k-v数
	dict.m.Range(func(key, val interface{}) bool {
		cnt++
		return true
	})
	return cnt
}

func (dict *SyncDict) Put(key string, val interface{}) (result int) {
	_, ok := dict.m.Load(key)
	dict.m.Store(key, val)
	if ok{
		return 0
	}
	return 1
}

func (dict *SyncDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, ok := dict.m.Load(key)
	if ok{
		return 0
	}
	dict.m.Store(key, val)
	return 1
}

func (dict *SyncDict) PutIfExists(key string, val interface{}) (result int) {
	_, ok := dict.m.Load(key)
	if ok{
		dict.m.Store(key, val)
		return 0
	}
	return 1
}

func (dict *SyncDict) Remove(key string) (result int) {
	_, ok := dict.m.Load(key)
	if ok{
		dict.m.Delete(key)
		return 1
	}
	return 0
}

func (dict *SyncDict) ForEach(consumer Consumer) {
	dict.m.Range(func(key, value interface{}) bool {
		consumer(key.(string), value)
		return true
	})
}

func (dict *SyncDict) Keys() []string {
	result := make([]string, dict.Len())
	index := 0
	dict.m.Range(func(key, value interface{}) bool {
		result[index] = key.(string)
		index++
		return true
	})

	return  result
}

func (dict *SyncDict) RandomKeys(limit int) []string {
	result := make([]string, dict.Len())
	for i :=0; i < limit; i++{
		dict.m.Range(func(key, value interface{}) bool {
			result[i] = key.(string)
			return false
		})
	}

	return result
}

func (dict *SyncDict) RandomDistinctKeys(limit int) []string {
	result := make([]string, dict.Len())
	index := 0
	dict.m.Range(func(key, value interface{}) bool {
		result[index] = key.(string)
		index++
		if index >= limit{
			return false
		}
		return true
	})
	return result
}

func (dict *SyncDict) Clear() {
	//直接重新生成dict指向的对象即可，旧的让GC去自动回收
	*dict = *MakeSyncDict()
}

