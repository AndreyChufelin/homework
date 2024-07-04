package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	val, ok := l.items[key]

	item := l.queue.PushFront(cacheItem{
		key:   key,
		value: value,
	})
	l.items[key] = item

	if ok {
		l.queue.Remove(val)
		return true
	}

	if l.queue.Len() > l.capacity {
		last := l.queue.Back()
		l.queue.Remove(last)

		delete(l.items, last.Value.(cacheItem).key)
	}

	return false
}

func (l lruCache) Get(key Key) (interface{}, bool) {
	val, ok := l.items[key]

	if ok {
		l.queue.MoveToFront(val)
		return val.Value.(cacheItem).value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
