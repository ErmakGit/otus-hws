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

type Item struct {
	Key   Key
	Value interface{}
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
	item, ok := lc.items[key]
	lc.mu.Unlock()

	if !ok {
		return nil, false
	}

	lc.mu.Lock()
	lc.queue.MoveToFront(item)
	lc.mu.Unlock()

	return lc.queue.Front().Value.(Item).Value, ok
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	// lc.mu.Lock()
	_, exist := lc.Get(key)
	// lc.mu.Unlock()

	if exist {
		lc.mu.Lock()
		firstItem := lc.queue.Front()
		lc.mu.Unlock()

		updatedItem := Item{
			Key:   key,
			Value: value,
		}
		firstItem.Value = updatedItem

		lc.mu.Lock()
		lc.items[key] = firstItem
		lc.mu.Unlock()
	} else {
		lc.mu.Lock()
		if len(lc.items) == lc.capacity {
			lc.removeLast()
		}
		lc.mu.Unlock()

		newItem := Item{
			Key:   key,
			Value: value,
		}

		lc.mu.Lock()
		lc.queue.PushFront(newItem)
		lc.items[key] = lc.queue.Front()
		lc.mu.Unlock()
	}

	return exist
}

func (lc *lruCache) removeLast() {
	lastItem := lc.queue.Back()
	lc.queue.Remove(lastItem)

	key := lastItem.Value.(Item).Key
	delete(lc.items, key)
}

func (lc *lruCache) Clear() {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}
