package main

import (
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
)

func benchmark(mc *memcache.Client, b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Get("hoge")
	}
}

func Benchmark_memcached(b *testing.B) {
	mc := memcache.New("localhost:11223")
	// let memcached server have key "hoge"
	mc.Set(&memcache.Item{Key: "hoge", Value: []byte("hoge!")})
	benchmark(mc, b)
}

func Benchmark_romdb(b *testing.B) {
	mc := memcache.New("localhost:11224")
	benchmark(mc, b)
}
