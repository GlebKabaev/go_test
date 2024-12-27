package main

import (
	"fmt"
	"sync/atomic"
)

type Mutex struct {
	Count  int32
	signal chan struct{}
}

func (m *Mutex) Unlock() {
	m.signal <- struct{}{}
}

func (m *Mutex) Wait() {
	for i := int32(0); i < m.Count; i++ {
		<-m.signal
	}
}

func main() {
	var count int32
	m := Mutex{
		Count:  500,
		signal: make(chan struct{}, 500),
	}

	for i := 0; i < 500; i++ {
		go func(localI int) {
			defer m.Unlock()
			fmt.Printf("Горутина %d: Привет, мир!\n", localI)
			atomic.AddInt32(&count, 1)
		}(i)
	}

	m.Wait()

	if atomic.LoadInt32(&count) == m.Count {
		fmt.Println("Все горутины успешно завершены!")
	} else {
		fmt.Printf("Некоторые горутины не завершились: завершено %d из %d\n", count, m.Count)
	}
}

