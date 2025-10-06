package hw04lrucache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		l.Remove(nil)
		l.MoveToFront(nil)

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("one item", func(t *testing.T) {
		l := NewList()
		l.PushFront("test")

		require.Equal(t, 1, l.Len())
		require.Equal(t, l.Front(), l.Back())

		l.MoveToFront(l.Front())
		l.Remove(l.Front())

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	/// Дополнительные тесты по объекту list
	// table test
	t.Run("stringer interface", func(t *testing.T) {
		tests := []struct {
			name     string
			items    []interface{}
			expected string
		}{
			{"empty", []interface{}{}, ""},
			{"strings", []interface{}{"Test test", "12345", "\n\t"}, "Test test, 12345, \n\t"},
			{"multitype", []interface{}{"-128", 10, 11.5}, "-128, 10, 11.5"},
		}
		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				l := new(list)
				for _, v := range tc.items {
					l.PushBack(v)
				}
				require.Equal(t, tc.expected, fmt.Sprint(l))
			})
		}
	})

	t.Run("untyped list", func(t *testing.T) {
		l := new(list)

		l.PushFront(10)    // [10]
		l.PushBack("test") // [10, 20]
		l.PushBack(nil)    // [10, 20, 30]
		require.Equal(t, 3, l.Len())
		require.Equal(t, "10, test, <nil>", l.String())

		l2 := NewList()
		l2.PushBack(10.45)
		l.Back().Value = l2

		require.Equal(t, "10, test, 10.45", l.String()) // l.String(...l2.String())
	})

	t.Run("search and move", func(t *testing.T) {
		l := new(list)

		l.PushFront(10)           // [10]
		l.PushBack("test")        // [10, 20]
		l.PushBack("second test") // [10, 20, 30]
		l.PushBack(20)
		l.PushBack(30)
		require.Equal(t, 5, l.Len())
		require.Equal(t, "10, test, second test, 20, 30", l.String())

		l.MoveToFront(l.SearchFirst(30))
		l.MoveToFront(l.SearchFirst(20))
		l.MoveToFront(l.SearchFirst(10))

		require.Equal(t, "10, 20, 30, test, second test", l.String())

		l.PushBack(20)
		l.PushBack(20)

		iLst := l.SearchFirst(20)
		iLstNext, err := l.SearchNext(iLst, 20)
		require.Nil(t, err)
		require.NotEqual(t, iLst, iLstNext)

		l2 := new(list)
		l2.PushFront(10)
		l2.PushBack(20)
		l2.PushBack(20)

		// проверяем на поиск от элемента чужого списка (iList - элемент другого списка)
		iLstNext, err = l2.SearchNext(iLst, 20)
		require.Nil(t, iLstNext)
		require.EqualError(t, err, ErrEnotherListItem)

		err = l2.SafeRemove(iLst)
		require.EqualError(t, err, ErrEnotherListItem)

		err = l2.SafeMoveToFront(iLst)
		require.EqualError(t, err, ErrEnotherListItem)

		// безопасные операции по элементу своего листа
		iLst = l2.SearchFirst(20)

		err = l2.SafeMoveToFront(iLst)
		require.NoError(t, err)

		err = l2.SafeRemove(iLst)
		require.NoError(t, err)
	})
}
