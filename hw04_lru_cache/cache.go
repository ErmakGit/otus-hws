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
	defer lc.mu.Unlock()

	return lc.get(key)
}

func (lc *lruCache) get(key Key) (interface{}, bool) {
	item, ok := lc.items[key]
	if !ok {
		return nil, false
	}

	lc.queue.MoveToFront(item)

	return lc.queue.Front().Value.(Item).Value, ok
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	_, exist := lc.get(key)

	if exist {
		firstItem := lc.queue.Front()

		updatedItem := Item{
			Key:   key,
			Value: value,
		}
		firstItem.Value = updatedItem

		lc.items[key] = firstItem
	} else {
		if len(lc.items) == lc.capacity {
			lc.removeLast()
		}

		newItem := Item{
			Key:   key,
			Value: value,
		}

		lc.queue.PushFront(newItem)
		lc.items[key] = lc.queue.Front()
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
