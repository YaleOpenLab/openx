// +build all

package utils

import (
	"testing"
)

// the utils benchmarks contains stuff other than that defined in the utils package since
// we can compare different libraries, optimizations, etc withoutworrying too much about
// the run time of the benchmark tests themselves.

// go test -run=XXX -tags="all" -bench=.

func BenchmarkItoB(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ItoB(i)
	}
}

func BenchmarkTimestamp(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Timestamp()
	}
}

func BenchmarkUnix(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Unix()
	}
}

func BenchmarkI64toB(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = I64toS(int64(i))
	}
}

func BenchmarkItoS(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ItoS(i)
	}
}

func BenchmarkBToI(b *testing.B) {
	byteString := []byte("blah")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BToI(byteString)
	}
}

func BenchmarkFtoS(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FtoS(2.0)
	}
}

func BenchmarkStoF(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StoF("2.0")
	}
}

func BenchmarkStoI(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StoI("2")
	}
}

func BenchmarkSha3Hash(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SHA3hash("test")
	}
}

func BenchmarkRandomString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetRandomString(10)
	}
}

func BenchmarkGetHomeDir(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetHomeDir()
	}
}
