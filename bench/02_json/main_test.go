package json

import (
	"testing"
)

var testJSON = []byte(`{"id":"tx_123","amount":99.99}`)

func BenchmarkValidateV1(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ValidateV1(testJSON)
	}
}

// V1
// bad: двойной парсинг + лишние аллокации

func BenchmarkValidateV2(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ValidateV2(testJSON)
	}
}

func BenchmarkValidateV3(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ValidateV3(testJSON)
	}
}

// Параллельные версии
func BenchmarkValidateV2Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ValidateV2(testJSON)
		}
	})
}

func BenchmarkValidateV3Parallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ValidateV3(testJSON)
		}
	})
}
