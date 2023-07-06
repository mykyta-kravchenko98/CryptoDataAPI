package cache

import "sync"

type LRUCache struct {
	Capacity int
	Cache    map[int]*cacheNode
	Head     *cacheNode
	Teil     *cacheNode
	m        *sync.RWMutex
}

type cacheNode struct {
	Key   int
	Value interface{}
	Next  *cacheNode
	Prev  *cacheNode
}

func InitLRUCache(capacity int) LRUCache {
	cache := LRUCache{
		Capacity: capacity,
		Cache:    make(map[int]*cacheNode),
		m:        &sync.RWMutex{},
	}
	return cache
}

func (this *LRUCache) Get(key int) interface{} {
	this.m.Lock()
	defer this.m.Unlock()
	if value, ok := this.Cache[key]; ok {
		this.MoveFront(value)
		return value.Value
	}

	return -1
}

func (this *LRUCache) Put(key int, value interface{}) {
	this.m.Lock()
	defer this.m.Unlock()
	if result, found := this.Cache[key]; found {
		result.Value = value
		this.MoveFront(result)
	} else {
		newNode := &cacheNode{Key: key, Value: value}
		if len(this.Cache) >= this.Capacity {
			delete(this.Cache, this.Teil.Key)
			this.RemoveTail()
		}
		this.Cache[key] = newNode
		this.AddNode(newNode)
	}
}

func (this *LRUCache) MoveFront(node *cacheNode) {
	if node == this.Head {
		return
	}

	this.RemoveNode(node)
	this.AddNode(node)
}

func (this *LRUCache) RemoveTail() {
	this.RemoveNode(this.Teil)
}

func (this *LRUCache) RemoveNode(node *cacheNode) {
	if node == this.Head {
		this.Head = node.Next
	}

	if node == this.Teil {
		this.Teil = node.Prev
	}

	if node.Prev != nil {
		node.Prev.Next = node.Next
	}

	if node.Next != nil {
		node.Next.Prev = node.Prev
	}

	node.Prev = nil
	node.Next = nil
}

func (this *LRUCache) AddNode(node *cacheNode) {
	if this.Head == nil {
		this.Head = node
		this.Teil = node
	} else {
		this.Head.Prev = node
		node.Next = this.Head
		this.Head = node
	}
}
