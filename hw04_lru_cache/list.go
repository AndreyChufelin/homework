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
	length int
	front  *ListItem
	back   *ListItem
}

func (l list) Len() int {
	return l.length
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.length++
	old := l.Front()
	newItem := &ListItem{Value: v, Next: old, Prev: nil}
	l.front = newItem

	if old == nil {
		l.back = newItem
	} else {
		old.Prev = newItem
	}

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.length++
	old := l.Back()
	newItem := &ListItem{Value: v, Prev: old, Next: nil}
	l.back = newItem

	if old == nil {
		l.front = newItem
	} else {
		old.Next = newItem
	}

	return newItem
}

func excludeItem(i *ListItem) {
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
}

func (l *list) Remove(i *ListItem) {
	if i == l.front {
		l.front = i.Next
	}
	if i == l.back {
		l.back = i.Prev
	}

	l.length--
	excludeItem(i)
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}
	if i == l.back {
		l.back = l.back.Prev
	}

	excludeItem(i)
	i.Next = l.front
	i.Prev = nil
	l.front.Prev = i
	l.front = i
}

func NewList() List {
	return new(list)
}
