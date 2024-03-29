package hw04lrucache

import "sync"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	items     map[*ListItem]*ListItem
	firstItem *ListItem
	lastItem  *ListItem
	once      sync.Once
}

func (l *list) Len() int {
	return len(l.items)
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Back() *ListItem {
	return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.once.Do(l.init)
	item := &ListItem{
		Value: v,
		Prev:  nil,
		Next:  l.Front(),
	}

	if l.Front() != nil {
		l.Front().Prev = item
	}

	l.firstItem = item
	if l.Len() == 0 {
		l.lastItem = item
	}
	l.items[item] = item

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.once.Do(l.init)
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.Back(),
	}
	if l.Back() != nil {
		l.Back().Next = item
	}

	l.lastItem = item
	if l.Len() == 0 {
		l.firstItem = item
	}
	l.items[item] = item

	return item
}

func (l *list) init() {
	if l.Len() == 0 {
		l.items = make(map[*ListItem]*ListItem, 0)
	}
}

func (l *list) Remove(i *ListItem) {
	l.connectNeighbors(i)
	delete(l.items, i)
	l.updateBorderItems(i)
}

func (l *list) updateBorderItems(i *ListItem) {
	if i == l.firstItem {
		l.firstItem = i.Next
	}

	if i == l.lastItem {
		l.lastItem = i.Prev
	}
}

func (l *list) connectNeighbors(i *ListItem) {
	prev := i.Prev
	next := i.Next
	if i.Prev != nil {
		prev.Next = next
	}

	if i.Next != nil {
		next.Prev = prev
	}
}

func (l *list) MoveToFront(i *ListItem) {
	value := i.Value
	l.Remove(i)
	l.PushFront(value)
}

func NewList() List {
	return new(list)
}
