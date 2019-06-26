package goCache

import (
	"errors"
	"sync"
	"time"
	"fmt"
)

const (
	DefaultExpiration time.Duration = 0
	KeyNotExists = "key is not exists"
	KeyExpired = "key is expired"
	KeyExists = "key %s already exists"
)

type item struct {
	Expired int64       //过期时间
	Value   interface{} //存储的值
}

//判断是否过期
func (m item) IsExpired() bool {
	if m.Expired == 0 {
		return false
	}
	return time.Now().UnixNano() > m.Expired
}

type Cache struct {
	*goCache
}

type goCache struct {
	DefaultExpiration time.Duration
	items             map[string]item //key => Item
	lock              sync.RWMutex
}

//实例对象
func New(d time.Duration) GoCacher {
	m := make(map[string]item)
	c := &goCache{
		DefaultExpiration: d,
		items:             m,
	}
	cc := &Cache{
		goCache: c,
	}
	return cc
}
//使用默认实例对象
func NewDefault() GoCacher {
	return New(DefaultExpiration)
}
//缓存某值
//key 值
//value 任意类型数据
//ttl 缓存时间
//如果键值存在则覆盖,重新设置时间
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	var ex int64
	if ttl == DefaultExpiration {
		ttl = c.DefaultExpiration
	}
	if ttl > 0 {
		ex = time.Now().Add(ttl).UnixNano()
	}
	c.items[key] = item{
		Expired: ex,
		Value:   value,
	}
	return true
}
//缓存某值 使用默认时间
//key 值
//value 任意类型数据
//ttl 缓存时间
//如果键值存在则覆盖,重新设置时间
func (c *Cache) SetDefault(key string, value interface{}) bool {
	return c.Set(key, value, DefaultExpiration)
}
//缓存某值
//key 值
//value 任意类型数据
//ttl 缓存时间
//如果键值未过期无法写入
func (c *Cache) Add(key string, value interface{}, ttl time.Duration) error {
	if c.Has(key) {
		return fmt.Errorf(KeyExists, key)
	}
	c.Set(key, value, ttl)
	return nil
}
//默认操作
func (c *Cache) AddDefault(key string, value interface{}) error {
	return c.Add(key, value, DefaultExpiration)
}
//获取某键值
func (c *Cache) Get(key string) (reply interface{}, err error) {
	item, isExist := c.items[key]
	if !isExist {
		err = errors.New(KeyNotExists)
		return
	}
	if item.Expired > 0 {
		if item.IsExpired() {
			c.Delete(key)
			err = errors.New(KeyExpired)
			return
		}
	}
	reply = item.Value
	return
}
//获取某键详情
//值, 过期, 是否存在
func (c *Cache) Info(key string) (interface{}, time.Time, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	item, found := c.items[key]
	if !found {
		return nil, time.Time{}, false
	}
	if item.Expired > 0 {
		if item.IsExpired() {
			return nil, time.Time{}, false
		}
		return item.Value, time.Unix(0, item.Expired), true
	}
	return item.Value, time.Time{}, true
}
//获取整个缓存项
//未过期的项
func (c *Cache) Items() map[string]item {
	c.lock.Lock()
	defer c.lock.Unlock()
	items := make(map[string]item, len(c.items))
	for k, v := range c.items {
		if v.Expired > 0 {
			if v.IsExpired() {
				continue
			}
		}
		items[k] = v
	}
	return items
}
//获取缓存多少项
func (c *Cache) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	var count int
	for _, v := range c.items {
		if !v.IsExpired() {
			count ++
		}
	}
	return count
}
//刷新缓存,相当于清空
func (c *Cache) Flush() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.items = map[string]item{}
}
//删除某键值
func (c *Cache) Delete(key string) bool {
	delete(c.items, key)
	return true
}
//判断键是否存在
func (c *Cache) Has(key string) bool {
	item, found := c.items[key]
	if !found {
		return false
	}
	return !item.IsExpired()
}
