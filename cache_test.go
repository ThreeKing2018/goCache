package goCache

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"strconv"
)


func TestCache_Add(t *testing.T) {
	gc := New(DefaultExpiration)
	err := gc.Add("lang", "golang", 100 * time.Millisecond)
	assert.Nil(t, err)

	err = gc.Add("lang", "php", 50 * time.Millisecond)
	assert.NotNil(t, err)
}
func TestCache_Set(t *testing.T) {
	gc := New(DefaultExpiration)
	gc.Set("aa", 100, time.Millisecond * 50)
	v, err := gc.Get("aa")
	assert.Nil(t, err)
	assert.Equal(t, v, 100)


	<-time.After(55 * time.Millisecond)
	v, err = gc.Get("aa")
	assert.NotNil(t, err)
	assert.Nil(t, v)
}

func TestCache_Has(t *testing.T) {
	gc := New(DefaultExpiration)
	b := gc.Has("nation")
	assert.False(t, b)

	gc.Set("nation", "japan", time.Second)
	b = gc.Has("nation")
	assert.True(t, b)
}

func TestCache_Count(t *testing.T) {
	gc := NewDefault()
	count := gc.Count()
	assert.Zero(t, count)

	gc.AddDefault("person", "america")

	count1 := gc.Count()
	assert.Equal(t, 1, count1)

	for i := 0; i < 10; i ++ {
		gc.AddDefault("person" + strconv.Itoa(i), i)
	}
	count2 := gc.Count()
	assert.Equal(t,  11, count2)

	gc.Delete("person")
	count3 := gc.Count()
	assert.Equal(t,  10, count3)

	gc.Add("nation", "japan", 50 * time.Millisecond)
	gc.Add("nation2", "china", 30 * time.Millisecond)
	gc.Add("nation3", "USA", 20 * time.Millisecond)
	<- time.After(time.Millisecond * 20)

	count4 := gc.Count()
	assert.Equal(t, 12, count4)

	<- time.After(time.Millisecond * 10)
	count5 := gc.Count()
	assert.Equal(t, 11, count5)
}

func TestCache_Delete(t *testing.T) {
	gc := NewDefault()
	gc.SetDefault("lang", "golang")
	gc.SetDefault("nation", "usa")

	val, err := gc.Get("lang")
	assert.Nil(t, err)
	assert.Equal(t, "golang", val)

	gc.Delete("lang")
	val, err = gc.Get("lang")
	assert.NotNil(t, err)
	assert.NotEqual(t, "golang", val)
}

func TestCache_Flush(t *testing.T) {
	gc := NewDefault()
	gc.AddDefault("lang", "vlang")
	gc.Flush()
	count := gc.Count()
	assert.Zero(t, count)
}

func TestNew(t *testing.T) {
	gc := New(3 * time.Second)
	gc.SetDefault("lang", "php")
}
func TestNewDefault(t *testing.T) {
	gc := NewDefault()
	gc.SetDefault("lang", "python")
}
func TestCache_Info(t *testing.T) {
	gc := NewDefault()
	gc.Set("lang", "golang", 5 * time.Minute)

	val, expired, b := gc.Info("lang")
	assert.Equal(t, val, "golang")
	assert.True(t, b)
	assert.Equal(t, expired.Format("15:04"), time.Now().Add(5 *time.Minute).Format("15:04"))
}

func TestCache_Items(t *testing.T) {
	gc := NewDefault()
	gc.AddDefault("lang", "php")
	gc.AddDefault("lang2", "golang")
	gc.Add("nation", "japan", 50 * time.Millisecond)

	items := gc.Items()
	assert.Equal(t, 3, len(items))

	<- time.After(50 * time.Millisecond)
	items2 := gc.Items()
	assert.Equal(t, 2, len(items2))
}

func TestItem_IsExpired(t *testing.T) {
	gc := NewDefault()
	gc.Set("nation", "japan", time.Millisecond * 100)
	val, err := gc.Get("nation")
	assert.Nil(t, err)
	assert.Equal(t, "japan", val)

	<- time.After(time.Millisecond * 100)
	val2, err2 := gc.Get("nation")
	assert.NotNil(t, err2)
	assert.NotEqual(t, val2, "japan")
}
func BenchmarkCache_Add(b *testing.B) {
	gc := NewDefault()
	for i := 0; i < b.N; i ++ {
		gc.Add("test" + strconv.Itoa(i), i, time.Second * 100)
	}
}
func BenchmarkCache_Set(b *testing.B) {
	gc := NewDefault()
	for i := 0; i < b.N; i ++ {
		gc.Set("test" + strconv.Itoa(i), i, time.Second * 100)
	}
}


func TestCache(t *testing.T) {
	gc := NewDefault()
	gc2 := NewDefault()
	gc.SetDefault("lang", "golang")
	v, err := gc2.Get("lang")
	assert.Nil(t, err)
	assert.Equal(t, "golang", v)
}