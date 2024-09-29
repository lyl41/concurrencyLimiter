package concurrencyLimiter

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_concurrencyLimiter_Get(t *testing.T) {
	l := NewConcurrencyLimiter(2)
	wg := new(sync.WaitGroup)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			l.Get()
			defer l.Release()
			fmt.Println("#", id, "get it", time.Now().Second())
			time.Sleep(time.Second * 5)
			fmt.Println("#", id, "done", time.Now().Second())
		}(i + 1)
	}

	go func() {
		time.Sleep(time.Second * 2)
		l.Reset(1)
		fmt.Println("reset to 1")
		time.Sleep(time.Second * 4)
		l.Reset(3)
		fmt.Println("reset to 3")
	}()
	wg.Wait()
}
