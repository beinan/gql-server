package cache

import (
	"sync"

	"github.com/beinan/gql-server/concurrent/future"
)

type ID = string
type Value = interface{}

type Future = future.Future

type Cache struct {
	sync.RWMutex
	data map[ID]Future
}

func MkCache() *Cache {
	return &Cache{
		data: make(map[ID]Future),
	}
}

func (c *Cache) LoadOrElse(
	key ID,
	producer func() (Value, error),
) Future {
	if value, ok := c.Load(key); ok {
		return value
	}
	c.Lock() //lock for write
	if value, ok := c.data[key]; ok {
		return value
	}
	value := future.MakeFuture(producer)
	c.data[key] = value
	c.Unlock() //write unlock
	return value
}

func (c *Cache) Load(key ID) (value Future, ok bool) {
	c.RLock()
	result, ok := c.data[key]
	c.RUnlock()
	return result, ok
}

/*
func (c *Cache) Delete(key ID) {
	c.Lock()
	delete(c.data, key)
	c.Unlock()
}

func (c *Cache) Store(key ID, value Future) {
	c.Lock()
	c.data[key] = value
	c.Unlock()
}
*/
