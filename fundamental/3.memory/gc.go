package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

// gcObservation: đọc số liệu GC từ runtime.MemStats, ép GC chạy bằng runtime.GC().
func gcObservation() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("HeapAlloc:", m.HeapAlloc/1024, "KB")
	fmt.Println("NumGC:", m.NumGC)

	runtime.GC() // forced GC — dùng khi benchmark / cần đồng bộ
	runtime.ReadMemStats(&m)
	fmt.Println("NumGC sau runtime.GC():", m.NumGC)
}

// gcMemoryLimit: soft memory limit (Go 1.19+) — hữu ích trong container.
// SetMemoryLimit(-1) = đọc limit hiện tại mà KHÔNG đổi.
func gcMemoryLimit() {
	debug.SetMemoryLimit(64 * 1024 * 1024) // 64 MiB
	limit := debug.SetMemoryLimit(-1)
	fmt.Println("MemoryLimit:", limit/1024/1024, "MB")
}

// finalizerDemo: runtime.SetFinalizer — KHÔNG guarantee thời điểm callback chạy.
func finalizerDemo() {
	obj := new(int)
	done := make(chan struct{})

	runtime.SetFinalizer(obj, func(*int) {
		close(done)
	})

	obj = nil    // bỏ reference cuối cùng tới obj
	runtime.GC() // chu kỳ 1: phát hiện obj là rác, đưa vào finalizer queue
	runtime.GC() // chu kỳ 2: finalizer goroutine thường chạy giữa 2 chu kỳ

	select {
	case <-done:
		fmt.Println("finalizer: đã chạy sau GC")
	case <-time.After(100 * time.Millisecond):
		fmt.Println("finalizer: CHƯA chạy — thời điểm không được guarantee")
	}
}

func GarbageCollector() {
	gcObservation()
	gcMemoryLimit()
	finalizerDemo()
}
