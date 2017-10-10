package main

import (
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
)

func benchmark(host string, b *testing.B) {
	mc := memcache.New(host)
	mc.Set(&memcache.Item{Key: "hoge", Value: []byte("hoge!")})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Get("hoge")
	}
}

func Benchmark_memcached(b *testing.B) {
	benchmark("localhost:11222", b)
}

func Benchmark_romdb(b *testing.B) {
	benchmark("localhost:11223", b)
}
