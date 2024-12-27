package main

import (
	"sync/atomic"
	"testing"
)

func TestMainLogic(t *testing.T) {
	var count int32
	m := Mutex{
		Count:  500,
		signal: make(chan struct{}, 500),
	}

	done := make(chan struct{}) // Канал для уведомления о завершении всех горутин

	for i := 0; i < 500; i++ {
		go func(localI int) {
			defer m.Unlock()
			atomic.AddInt32(&count, 1)
			if atomic.LoadInt32(&count) == m.Count {
				close(done) // Закрываем канал, когда все горутины завершены
			}
		}(i)
	}

	m.Wait()

	select {
	case <-done:
		if atomic.LoadInt32(&count) != m.Count {
			t.Errorf("Expected %d goroutines to complete, but got %d", m.Count, count)
		}
	default:
		t.Error("Test failed: not all goroutines completed")
	}
}
