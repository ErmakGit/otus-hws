package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	mu       sync.Mutex
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	item, ok := lc.items[key]
	if !ok {
		return nil, false
	}

	lc.queue.MoveToFront(item)

	return lc.queue.Front().Value, ok
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	_, exist := lc.Get(key)

	lc.mu.Lock()
	defer lc.mu.Unlock()

	if exist {
		firstItem := lc.queue.Front()
		firstItem.Value = value
		lc.items[key] = firstItem
	} else {
		lc.removeLast()
		lc.queue.PushFront(value)
		lc.items[key] = lc.queue.Front()
	}

	return exist
}

func (lc *lruCache) removeLast() {
	if len(lc.items) == lc.capacity {
		lastItem := lc.queue.Back()
		lc.queue.Remove(lastItem)

		var key Key
		for k, item := range lc.items {
			if item == lastItem {
				key = k
			}
		}
		delete(lc.items, key)
	}
}

func (lc *lruCache) Clear() {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}
