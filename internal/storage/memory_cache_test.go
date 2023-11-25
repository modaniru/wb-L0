package storage

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutAndGet10(t *testing.T) {
	cache := NewInmemoryCache()
	for i := 0; i < 10; i++ {
		str := strconv.Itoa(i)
		cache.Put(str, []byte(str))
	}
	for i := 0; i < 10; i++ {
		str := strconv.Itoa(i)
		actual, ok := cache.Get(str)
		assert.True(t, ok)
		assert.Equal(t, str, string(actual))
	}
	for i := 10; i < 20; i++ {
		str := strconv.Itoa(i)
		actual, ok := cache.Get(str)
		assert.False(t, ok)
		assert.Nil(t, actual)
	}
	assert.Equal(t, 10, cache.Len())
}

func TestPutAndGet100Async(t *testing.T) {
	cache := NewInmemoryCache()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < 100; i++ {
			str := strconv.Itoa(i)
			cache.Put(str, []byte(str))
		}
		wg.Done()
	}()
	wg.Wait()
	for i := 0; i < 100; i++ {
		str := strconv.Itoa(i)
		actual, ok := cache.Get(str)
		assert.True(t, ok)
		assert.Equal(t, str, string(actual))
	}
	assert.Equal(t, 100, cache.Len())
}

func TestPutAndGet1000Async(t *testing.T) {
	cache := NewInmemoryCache()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < 1000; i++ {
			str := strconv.Itoa(i)
			cache.Put(str, []byte(str))
		}
		wg.Done()
	}()
	wg.Wait()
	for i := 0; i < 1000; i++ {
		str := strconv.Itoa(i)
		actual, ok := cache.Get(str)
		assert.True(t, ok)
		assert.Equal(t, str, string(actual))
	}
	assert.Equal(t, 1000, cache.Len())
}

func TestPutAndGet10000Async(t *testing.T) {
	cache := NewInmemoryCache()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < 10000; i++ {
			str := strconv.Itoa(i)
			cache.Put(str, []byte(str))
		}
		wg.Done()
	}()
	wg.Wait()
	for i := 0; i < 10000; i++ {
		str := strconv.Itoa(i)
		actual, ok := cache.Get(str)
		assert.True(t, ok)
		assert.Equal(t, str, string(actual))
	}
}
