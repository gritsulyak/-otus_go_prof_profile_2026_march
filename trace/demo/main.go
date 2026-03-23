package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"
)

// workContended имитирует высокую борьбу за мьютекс (lock contention)
func workContended(id int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	// mu.Lock() // UNCOMMENT to see serial execution
	add := 0
	for i := id; i < 1000*1000; i++ {
		// --- Критическая секция ---
		// Делаем что-то, что занимает время, пока держим лок.
		// Чем дольше время тут, тем выше contention.
		mu.Lock()
		add += i
		if i%1000 == 0 {
			time.Sleep(time.Nanosecond)
		}
		mu.Unlock()
		// --------------------------
	}
	fmt.Println(id, add)
	// mu.Unlock() // UNCOMMENT
}

func main() {
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)

	f, err := os.Create("trace.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Начинаем трассировку
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	var mu sync.Mutex
	var wg sync.WaitGroup

	// Запускаем много горутин
	numGoroutines := 100
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go workContended(i, &mu, &wg)
	}

	wg.Wait()
}
