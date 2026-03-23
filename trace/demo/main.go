package main

import (
	"os"
	"runtime/trace"
	"sync"
)

func work(id int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 50; i++ {
		mu.Lock()
		_ = id * i // имитация работы
		mu.Unlock()
	}
}

func main() {
	f, _ := os.Create("trace.prof")
	trace.Start(f)
	defer trace.Stop()

	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go work(i, &mu, &wg)
	}
	wg.Wait()
}

func workFast(id int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	local := 0
	for i := 0; i < 50; i++ {
		local += id * i // работаем локально — без блокировки
	}

	mu.Lock()
	_ = local // синхронизация один раз
	mu.Unlock()
}
