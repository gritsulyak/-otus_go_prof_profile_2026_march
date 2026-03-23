package main

package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	_ "net/http/pprof" // ← один импорт регистрирует /debug/pprof/*
)

//  V1: тяжёлый расчёт — строки + срез без преаллокации
func calcPortfolioRiskV1(w http.ResponseWriter, r *http.Request) {
	n := queryInt(r, "n", 5000)

	// ??!! Конкатенация строк в цикле
	log := ""
	// ??!! Срез без преаллокации
	var returns []float64

	for i := 1; i <= n; i++ {
		price := 100.0 * math.Exp(0.0002*float64(i))
		ret := math.Log(price / (price - 0.02))
		returns = append(returns, ret)                   // ?? regrowing slice
		log = log + fmt.Sprintf("step=%d ret=%.6f\n", i, ret) // ?? new string each time
	}

	mean := calcMean(returns)
	variance := calcVariance(returns, mean)
	volatility := math.Sqrt(variance) * math.Sqrt(252)

	fmt.Fprintf(w, "n=%d volatility=%.4f\nlog_len=%d\n",
		n, volatility, len(log))
}

// Вспомогательные функции расчёта
func calcMean(xs []float64) float64 {
	sum := 0.0
	for _, x := range xs { sum += x }
	return sum / float64(len(xs))
}

func calcVariance(xs []float64, mean float64) float64 {
	sum := 0.0
	for _, x := range xs {
		d := x - mean
		sum += d * d
	}
	return sum / float64(len(xs))
}

func queryInt(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if n, err := strconv.Atoi(v); err == nil { return n }
	return def
}

func main() {
	http.HandleFunc("/risk/v1", calcPortfolioRiskV1)
	http.HandleFunc("/risk/v2", calcPortfolioRiskV2)

	// pprof эндпоинты доступны автоматически:
	// GET /debug/pprof/
	// GET /debug/pprof/profile?seconds=N
	// GET /debug/pprof/heap
	// GET /debug/pprof/goroutine

	fmt.Println("Сервис запущен на :8080")
	fmt.Println("pprof доступен на :8080/debug/pprof/")
	http.ListenAndServe(":8080", nil)
}

// below оптимизированный расчёт для сравнения с тяжёлым вариантом



// V2: оптимизированный расчёт
func calcPortfolioRiskV2(w http.ResponseWriter, r *http.Request) {
	n := queryInt(r, "n", 5000)

	// strings.Builder вместо конкатенации
	var sb strings.Builder
	sb.Grow(n * 30)
	// make с известной ёмкостью
	returns := make([]float64, 0, n)

	for i := 1; i <= n; i++ {
		price := 100.0 * math.Exp(0.0002*float64(i))
		ret := math.Log(price / (price - 0.02))
		returns = append(returns, ret)
		fmt.Fprintf(&sb, "step=%d ret=%.6f\n", i, ret)
	}

	mean := calcMean(returns)
	variance := calcVariance(returns, mean)
	volatility := math.Sqrt(variance) * math.Sqrt(252)

	fmt.Fprintf(w, "n=%d volatility=%.4f\nlog_len=%d\n",
		n, volatility, sb.Len())
}



// V2: оптимизированный расчёт
func calcPortfolioRiskV2(w http.ResponseWriter, r *http.Request) {
	n := queryInt(r, "n", 5000)

	// strings.Builder вместо конкатенации
	var sb strings.Builder
	sb.Grow(n * 30)
	// make с известной ёмкостью
	returns := make([]float64, 0, n)

	for i := 1; i <= n; i++ {
		price := 100.0 * math.Exp(0.0002*float64(i))
		ret := math.Log(price / (price - 0.02))
		returns = append(returns, ret)
		fmt.Fprintf(&sb, "step=%d ret=%.6f\n", i, ret)
	}

	mean := calcMean(returns)
	variance := calcVariance(returns, mean)
	volatility := math.Sqrt(variance) * math.Sqrt(252)

	fmt.Fprintf(w, "n=%d volatility=%.4f\nlog_len=%d\n",
		n, volatility, sb.Len())
}
