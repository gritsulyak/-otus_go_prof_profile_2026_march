package main

import (
	"os"
	"runtime"
	"runtime/pprof"
)

// nok:slice
func buildPortfolioV1(prices []float64) []float64 {
	var result []float64 // cap=0, будет расти
	for _, p := range prices {
		result = append(result, p*1.13) // пересоздаёт массив
	}
	return result
}

func main() {
	prices := make([]float64, 100_000)
	for i := range prices {
		prices[i] = float64(i) * 1.05
	}

	// ── CPU профиль ──────────────────────────────────────
	cpuF, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(cpuF)

	for i := 0; i < 1000; i++ {
		buildPortfolioV2(prices)
	}

	pprof.StopCPUProfile()
	cpuF.Close()

	// ── Memory профиль ───────────────────────────────────
	memF, _ := os.Create("mem.prof")
	for i := 0; i < 1000; i++ {
		buildPortfolioV2(prices)
	}
	runtime.GC()                 // сначала вызвать GC
	pprof.WriteHeapProfile(memF) // затем записать профиль
	memF.Close()
}

// Optimized version: VVV

// ok: преаллоцированный срез
func buildPortfolioV2(prices []float64) []float64 {
	result := make([]float64, 0, len(prices)) // cap известен заранее
	for _, p := range prices {
		result = append(result, p*1.13)
	}
	return result
}
