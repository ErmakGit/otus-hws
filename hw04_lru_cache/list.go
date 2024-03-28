package hw04lrucache

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
	Items []*ListItem
}

func (l list) Len() int {
	return len(l.Items)
}

func (l *list) Front() *ListItem {
	if l.Len() == 0 {
		return nil
	}

	return l.Items[0]
}

func (l *list) Back() *ListItem {
	if l.Len() == 0 {
		return nil
	}

	return l.Items[l.Len()-1]
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Prev:  nil,
		Next:  l.Front(),
	}
	newList := []*ListItem{
		item,
	}
	newList = append(newList, l.Items...)

	if l.Front() != nil {
		l.Front().Prev = item
	}

	l.Items = newList

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.Back(),
	}
	if l.Back() != nil {
		l.Back().Next = item
	}

	l.Items = append(l.Items, item)

	return item
}

func (l *list) Remove(i *ListItem) {
	ind := 0
	for key, item := range l.Items {
		if item == i {
			ind = key
			break
		}
	}

	l.connectNeighbors(i)

	l.Items = append(l.Items[:ind], l.Items[ind+1:]...)
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
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
