package goCache

import (
	"errors"
	"sync"
	"time"
	"fmt"
	"runtime"
)

const (
	DefaultExpiration time.Duration = 0
	KeyNotExists = "key is not exists"
	KeyExpired = "key is expired"
	KeyExists = "key %s already exists"
	clearInterval time.Duration = 3 //定义清除过期元素
)

type item struct {
	Expired int64       //过期时间
	Value   interface{} //存储的值
}

//判断是否过期
func (m item) isExpired() bool {
	if m.Expired == 0 {
		return false
	}
	return time.Now().UnixNano() > m.Expired
}

type Cache struct {
	*goCache
	stop chan bool
}

type goCache struct {
	DefaultExpiration time.Duration
	items             map[string]item //key => Item
	lock              sync.RWMutex
}

//使用默认实例对象
func NewDefault() GoCacher {
	return New(DefaultExpiration)
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
	go clockClear(cc)
	runtime.SetFinalizer(cc, stopClock) //在对象被GC进程选中并从内存中移除前，SetFinalizer都不会执行，即使程序正常结束或者发生错误
	return cc
}
//发送停止信号
func stopClock(c *Cache) {
	c.stop <- true
}
//定时清除过期元素
func clockClear(c *Cache) {
	ticker := time.NewTicker(time.Second * clearInterval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-c.stop:
			ticker.Stop()
		}
	}
}

//缓存某值
//key 值
//value 任意类型数据
//ttl 缓存时间
//如果键值存在则覆盖,重新设置时间
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
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
}
//缓存某值 使用默认时间
//key 值
//value 任意类型数据
//ttl 缓存时间
//如果键值存在则覆盖,重新设置时间
func (c *Cache) SetDefault(key string, value interface{}) {
	c.Set(key, value, DefaultExpiration)
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
		if item.isExpired() {
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
		if item.isExpired() {
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
			if v.isExpired() {
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
		if !v.isExpired() {
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
//删除过期的值
func (c *Cache) DeleteExpired() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for k, v := range c.items {
		if v.Expired > 0 && v.isExpired() {
			c.Delete(k)
		}
	}
}
//判断键是否存在
func (c *Cache) Has(key string) bool {
	item, found := c.items[key]
	if !found {
		return false
	}
	return !item.isExpired()
}
