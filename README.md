# goCache

## 文档

## 获取
`go get `


## 使用

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