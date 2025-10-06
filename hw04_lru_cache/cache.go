package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

// Структура quValue для хранения в списке пар ключ, значение
// для обеспечения сложности О(1) при операциях с кешем.
type quValue struct {
	key   Key
	value interface{}
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
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

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// если запись по ключу есть - обновляем значение и переносим её в начало очереди
	if i, ok := c.items[key]; ok {
		i.Value = quValue{key, value}
		c.queue.MoveToFront(i)
		return true
	}

	// если записи по ключу нет добавляем новую в начало очереди
	// проверяем объем кеша. если количество записей равно capacity - удаляем крйнюю (с хвоста)
	if c.queue.Len() == c.capacity {
		delete(c.items, c.queue.Back().Value.(quValue).key)
		c.queue.Remove(c.queue.Back())
	}

	// Добавляем новую запись в начало кеша
	c.items[key] = c.queue.PushFront(quValue{key, value})

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// если искомая запись есть в кеше - переносим её в начало
	if i, ok := c.items[key]; ok {
		c.queue.MoveToFront(i)
		return i.Value.(quValue).value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
