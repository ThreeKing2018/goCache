# goCache

## 文档
`https://godoc.org/github.com/ThreeKing2018/goCache`

## 获取
`go get github.com/ThreeKing2018/goCache`


## 使用
```
import "github.com/ThreeKing2018/goCache"

//默认new
c := goCache.NewDefault()
//设置缓存键, 若键存在则覆盖
c.Set("nation", "free_nation", 5 * time.Minute)
//设置缓存键,与set不同,若键存在则false,代表存储失败; 返回true则存储成功,
c.Add("lang", "golang", 1 * time.Minute)
//获取某值
v, err := c.Get("lang")
if err != nil {
    log.Fatal(err)
}
if v != "golang" {
    log.Fatalln("not equal")
}
//查看某键详情
v, t, b := c.Info("lang")
fmt.Println(v, t, b)
//查看共存储多少项
fmt.Println(c.Count())
//查看所有的缓存
fmt.Println(c.Items())
//删除键
c.Delete("lang")

//默认时间为1小时
cache := goCache.New(1 * time.Hour)
//设置键, 使用goCache.New设置的1小时过期时间
cache.SetDefault("lang", "php")
```


设置值
```
cache := NewDefault()
cache.Set("lang", "golang", 10 * time.)
```

## 标准测试
```
BenchmarkCache_Add-4     1000000              1074 ns/op             232 B/op          3 allocs/op
BenchmarkCache_Add-4     1000000              1176 ns/op             232 B/op          3 allocs/op
BenchmarkCache_Add-4     1000000              1109 ns/op             232 B/op          3 allocs/op
BenchmarkCache_Set-4     1000000              1090 ns/op             232 B/op          3 allocs/op
BenchmarkCache_Set-4     2000000              1063 ns/op             232 B/op          3 allocs/op
BenchmarkCache_Set-4     1000000              1005 ns/op             232 B/op          3 allocs/op
PASS

```