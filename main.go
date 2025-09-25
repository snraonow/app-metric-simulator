package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func simulate(duration int, frequency string) {
	numCPU := runtime.NumCPU() / 2 // Use half of available cores
	if numCPU < 1 {
		numCPU = 1
	}
	fmt.Printf("Starting simulation using %d cores (50%% of available cores)...\n", numCPU)

	var wg sync.WaitGroup
	var memoryMutex sync.Mutex
	var totalMemory int64

	targetInitialMemoryGB := 1.0                                                    // Start with 1GB total
	memoryPerCore := (targetInitialMemoryGB * 1024 * 1024 * 1024) / float64(numCPU) // bytes per core

	// Run on half of the cores
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Keep allocated strings in memory
			var stringSlice []string
			baseStr := make([]byte, 1024*1024) // 1MB base string
			for i := range baseStr {
				baseStr[i] = byte(65 + (i % 26)) // Fill with repeating alphabet
			}

			end := time.Now().Add(time.Duration(duration) * time.Second)
			start := time.Now()

			// Initial memory allocation for this core
			initialSize := int(memoryPerCore)
			allocated := 0

			// Allocate initial portion of memory
			for allocated < initialSize {
				chunk := baseStr
				if allocated+len(chunk) > initialSize {
					chunk = chunk[:initialSize-allocated]
				}
				stringSlice = append(stringSlice, string(chunk))
				allocated += len(chunk)

				memoryMutex.Lock()
				totalMemory += int64(len(chunk))
				currentGB := float64(totalMemory) / (1024 * 1024 * 1024)
				memoryMutex.Unlock()

				fmt.Printf("\rCore %d: Initial Memory: %.2f GB", id, currentGB)
			}

			// Main simulation loop
			for time.Now().Before(end) {
				elapsed := time.Since(start)
				durationSecs := float64(duration)
				elapsedSecs := elapsed.Seconds()

				// Calculate intensity (60% to 90%)
				intensity := 0.6 + 0.3*(elapsedSecs/durationSecs)

				// Calculate target memory for this point in time (1GB to 10GB)
				progressRatio := elapsedSecs / durationSecs
				targetTotalMemoryGB := 1.0 + (9.0 * progressRatio) // 1GB to 10GB
				targetMemoryPerCore := (targetTotalMemoryGB * 1024 * 1024 * 1024) / float64(numCPU)

				// Allocate more memory if needed
				currentAllocated := int64(allocated)
				if float64(currentAllocated) < targetMemoryPerCore {
					// Allocate in 50MB chunks
					chunkSize := 50 * 1024 * 1024
					chunk := make([]byte, chunkSize)
					for i := range chunk {
						chunk[i] = byte(rand.Intn(256))
					}
					stringSlice = append(stringSlice, string(chunk))
					allocated += chunkSize

					memoryMutex.Lock()
					totalMemory += int64(chunkSize)
					currentGB := float64(totalMemory) / (1024 * 1024 * 1024)
					memoryMutex.Unlock()

					fmt.Printf("\rCore %d: Memory: %.2f GB / %.2f GB, CPU: %.1f%%",
						id, currentGB, targetTotalMemoryGB, intensity*100)
				}

				// CPU-intensive work
				sum := 0.0
				for j := 0; j < int(intensity*10000); j++ {
					sum += rand.Float64()
				}

				// Sleep for a short time to prevent CPU saturation
				time.Sleep(time.Duration(5*(1-intensity)) * time.Millisecond)
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Force GC to get accurate memory usage
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("\nFinal Memory Usage: %.2f GB\n", float64(m.Alloc)/(1024*1024*1024))
}

func main() {
	if len(os.Args) < 2 || os.Args[1] != "simulate" {
		fmt.Println("Usage: simulate high=memory,cpu time=20 frequency=constant|random")
		os.Exit(1)
	}

	highFlag := flag.String("high", "cpu", "Resource to stress: memory, cpu, or both (comma separated)")
	timeFlag := flag.Int("time", 10, "Duration in seconds")
	freqFlag := flag.String("frequency", "constant", "Load frequency: constant or random")

	flag.CommandLine.Parse(os.Args[2:])

	resources := strings.Split(*highFlag, ",")
	duration := *timeFlag
	frequency := *freqFlag

	rand.Seed(time.Now().UnixNano())
	fmt.Printf("Simulating high %v for %d seconds with %s frequency...\n", resources, duration, frequency)

	// Run the combined CPU and memory simulation
	simulate(duration, frequency)
	fmt.Println("\nSimulation complete.")
}
