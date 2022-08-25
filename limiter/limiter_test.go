package limiter

import (
	"testing"
)

func TestNewCache(t *testing.T) {
	_, err := NewCache()
	if err != nil{
		t.Errorf("NewCache is error: %v", err)
	}
}

func TestLimiter_CountLimit(t *testing.T) {
	l, _ := NewCache()
	var y,n int
	for i:=0; i< 5; i++{
		if l.CountLimit("TestCountLimit", 2, 5) {
			y++
			continue
		}
		n++
	}
	if y != 2 || n != 3{
		t.Errorf("TestLimiter_CountLimit Request 5, Allow:%v, Limit:%v", y, n)
	}
}

func BenchmarkLimiter_CountLimit(b *testing.B) {
	l, _ := NewCache()
	for i := 0; i < b.N; i++ {
		l.CountLimit("TestBCountLimit", 2, 5)
	}
}

func TestLimiter_SyncCountLimit(t *testing.T) {
	l, _ := NewCache()
	var y,n int
	for i:=0; i< 5; i++{
		if l.SyncCountLimit("TestSyncCountLimit", 2, 5) {
			y++
			continue
		}
		n++
	}
	if y != 2 || n != 3{
		t.Errorf("TestLimiter_SyncCountLimit Request 5, Allow:%v, Limit:%v", y, n)
	}
}

func BenchmarkLimiter_SyncCountLimit(b *testing.B) {
	l, _ := NewCache()
	for i := 0; i < b.N; i++ {
		l.CountLimit("TestBSyncCountLimit", 2, 5)
	}
}

func TestLimiter_WindowLimit(t *testing.T) {
	l, _ := NewCache()
	var y,n int
	for i:=0; i< 5; i++{
		if l.WindowLimit("TestWindowLimit", 2, 5) {
			y++
			continue
		}
		n++
	}
	if y != 2 || n != 3{
		t.Errorf("TestLimiter_WindowLimit Request 5, Allow:%v, Limit:%v", y, n)
	}
}

func BenchmarkLimiter_WindowLimit(b *testing.B) {
	l, _ := NewCache()
	for i := 0; i < b.N; i++ {
		l.CountLimit("TestBWindowLimit", 2, 5)
	}
}

func TestLimiter_SyncWindowLimit(t *testing.T) {
	l, _ := NewCache()
	var y,n int
	for i:=0; i< 5; i++{
		if l.SyncWindowLimit("TestSyncWindowLimit", 2, 5) {
			y++
			continue
		}
		n++
	}

	if y != 2 || n != 3{
		t.Errorf("TestLimiter_SyncWindowLimit Request 5, Allow:%v, Limit:%v", y, n)
	}
}

func BenchmarkLimiter_SyncWindowLimit(b *testing.B) {
	l, _ := NewCache()
	for i := 0; i < b.N; i++ {
		l.CountLimit("TestBSyncWindowLimit", 2, 5)
	}
}