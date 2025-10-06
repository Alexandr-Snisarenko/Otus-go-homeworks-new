package hw04lrucache

import (
	"errors"
	"fmt"
)

var ErrEnotherListItem = "данный элемент не принадлежит текущему списку"

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
	front  *ListItem
	back   *ListItem
	length int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) PushFront(item interface{}) *ListItem {
	newItem := &ListItem{
		Value: item,
		Next:  l.front,
	}

	if l.front == nil {
		l.front = newItem
		l.back = l.front
	} else {
		l.front.Prev = newItem
		l.front = newItem
	}

	l.length++
	return newItem
}

func (l *list) PushBack(item interface{}) *ListItem {
	newItem := &ListItem{
		Value: item,
		Prev:  l.back,
	}
	if l.back == nil {
		l.back = newItem
		l.front = l.back
	} else {
		l.back.Next = newItem
		l.back = newItem
	}

	l.length++
	return newItem
}

func (l *list) Remove(item *ListItem) {
	if item == nil {
		return
	}

	if l.length == 1 {
		l.back = nil
		l.front = nil
		l.length = 0
		return
	}

	switch item {
	case l.front:
		l.front = item.Next
		item.Next.Prev = nil
		item.Next = nil
	case l.back:
		l.back = item.Prev
		item.Prev.Next = nil
		item.Prev = nil
	default:
		item.Prev.Next = item.Next
		item.Next.Prev = item.Prev
	}

	l.length--
}

func (l *list) MoveToFront(item *ListItem) {
	if item == nil {
		return
	}
	if item == l.front {
		return
	}

	l.Remove(item) // удалили элемент length уменьшилась на 1
	l.front.Prev = item
	item.Next = l.front
	l.front = item
	l.length++ // приводим length в норму после Remove()
}

//////////////////////////////////////////////////////////////////////////////////////
// Доп ф-ии для полноценного использования листа как отдельной либы.

// Защищенный Remove.
func (l *list) SafeRemove(item *ListItem) error {
	if !l.checkItem(item) {
		return errors.New(ErrEnotherListItem)
	}

	l.Remove(item)
	return nil
}

// Защищенный MoveToFront.
func (l *list) SafeMoveToFront(item *ListItem) error {
	if !l.checkItem(item) {
		return errors.New(ErrEnotherListItem)
	}

	l.MoveToFront(item)
	return nil
}

// ищем элемент с указанным значением начиная от указанного.
func (l *list) SearchNext(startItem *ListItem, v interface{}) (*ListItem, error) {
	if !l.checkItem(startItem) {
		return nil, errors.New(ErrEnotherListItem)
	}

	for i := startItem.Next; i != nil; i = i.Next {
		if i.Value == v {
			return i, nil
		}
	}

	return nil, nil
}

// возвращаем первый найденный по содержимому элемент списка.
func (l *list) SearchFirst(v interface{}) *ListItem {
	item, _ := l.SearchNext(l.front, v)
	return item
}

// проверяем принадлежит ли указанный элемент списка текущему листу.
func (l *list) checkItem(item *ListItem) bool {
	if item == nil {
		return false
	}

	for i := l.Front(); i != nil; i = i.Next {
		if i == item {
			return true
		}
	}
	return false
}

func (l list) String() string {
	res := ""
	for i := l.Front(); i != nil; i = i.Next {
		if i == l.front {
			res = fmt.Sprintf("%v", i.Value)
		} else {
			res += fmt.Sprintf(", %v", i.Value)
		}
	}
	return res
}
