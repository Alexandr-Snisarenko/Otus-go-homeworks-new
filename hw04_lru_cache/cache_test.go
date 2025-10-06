package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)

		c.Clear()
		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)
		c.Set("1", 10)
		c.Set("2", 20)
		c.Set("3", 30)
		c.Set("4", 40)

		val, ok := c.Get("1")
		require.False(t, ok)
		require.Nil(t, val)

		for i := 4; i >= 2; i-- {
			val, ok = c.Get(Key(strconv.Itoa(i)))
			require.True(t, ok)
			require.Equal(t, i*10, val)
		}

		// 4 переместилась в конец
		// добавляем еще одну запись
		c.Set("5", 50)
		// 4 - должна уйти
		val, ok = c.Get("4")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear cache", func(t *testing.T) {
		c := NewCache(3)
		c.Set("1", 10)
		c.Set("2", 20)
		c.Set("3", 30)

		val, ok := c.Get("1")
		require.True(t, ok)
		require.Equal(t, val, 10)

		c.Clear()

		val, ok = c.Get("1")
		require.False(t, ok)
		require.Nil(t, val)

		c.Set("1", 10)

		val, ok = c.Get("1")
		require.True(t, ok)
		require.Equal(t, val, 10)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Run("async test", func(t *testing.T) {
		c := NewCache(10)
		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			for i := 0; i < 1_000_000; i++ {
				c.Set(Key(strconv.Itoa(i)), i)
			}
		}()

		go func() {
			defer wg.Done()
			for i := 0; i < 1_000_000; i++ {
				c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
			}
		}()

		wg.Wait()
	})
}
