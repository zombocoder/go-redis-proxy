package stats

import (
	"log"
	"runtime"
)

// Log memory stats periodically
func LogMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	log.Println("---------------- /Memory Stats ----------------")
	// General stats
	log.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	log.Printf("TotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	log.Printf("Sys = %v MiB", m.Sys/1024/1024)
	log.Printf("NumGC = %v", m.NumGC)

	// Additional metrics
	log.Printf("HeapAlloc = %v MiB", m.HeapAlloc/1024/1024)
	log.Printf("HeapIdle = %v MiB", m.HeapIdle/1024/1024)
	log.Printf("HeapInuse = %v MiB", m.HeapInuse/1024/1024)
	log.Printf("HeapReleased = %v MiB", m.HeapReleased/1024/1024)
	log.Printf("StackInuse = %v MiB", m.StackInuse/1024/1024)
	log.Printf("Mallocs = %v", m.Mallocs)
	log.Printf("Frees = %v", m.Frees)

	// Number of goroutines currently running
	log.Printf("Goroutines = %v", runtime.NumGoroutine())
	log.Println("---------------- Memory Stats/ ----------------")
}
