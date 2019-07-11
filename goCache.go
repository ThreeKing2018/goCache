package goCache

import "time"

type GoCacher interface {
	Set(key string, value interface{}, ttl time.Duration) //设置键值
	SetDefault(key string, value interface{})//设置键值,使用默认值
	Add(key string, value interface{}, ttl time.Duration) error //添加键值
	AddDefault(key string, value interface{}) error //添加键值,使用默认时间
	Get(key string) (reply interface{}, err error)//获取值
	Info(key string) (interface{}, time.Time, bool)//获取某键详情
	Items() map[string]item //获取所有缓存项
	Count() int //获取缓存项目个数
	Flush()//刷新缓存
	Delete(key string) bool//删除某键
	Has(key string) bool//判断键是否存在
}