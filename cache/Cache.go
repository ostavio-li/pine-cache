package cache

import (
	"log"
	"sync"
)

type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() Stat
}

type Stat struct {
	Count     int64
	KeySize   int64
	ValueSize int64
}

func (s *Stat) add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}

func (s *Stat) del(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}

func New(typ string) Cache {
	var c Cache

	if typ == "inmemory" {
		c = newInMemoryCache()
	} else if c == nil {
		panic("unknown cache type " + typ)
	}
	log.Println(typ, "ready to serve")
	return c
}

type inMemoryCache struct {
	c     map[string][]byte
	mutex sync.RWMutex
	Stat
}

func (c *inMemoryCache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	bytes := c.c[k]
	if bytes != nil {
		c.del(k, bytes)
	}
	c.c[k] = v
	c.add(k, v)
	return nil
}

func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.c[k], nil
}

func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v := c.c[k]
	if v != nil {
		delete(c.c, k)
		c.del(k, v)
	}
	return nil
}

func (c *inMemoryCache) GetStat() Stat {
	return c.Stat
}

func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{
		make(map[string][]byte),
		sync.RWMutex{},
		Stat{},
	}
}
