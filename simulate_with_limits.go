package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func simulateWithLimits(duration int, frequency string, maxCPU int, maxMemoryGB float64) {
	numCPU := runtime.NumCPU() / 2
	if numCPU < 1 {
		numCPU = 1
	}
	fmt.Printf("Starting simulation using %d cores (50%% of available cores)...\n", numCPU)

	var wg sync.WaitGroup
	var memoryMutex sync.Mutex
	var totalMemory int64

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var memChunks [][]byte
			baseStr := make([]byte, 1024*1024)
			for i := range baseStr {
				baseStr[i] = byte(65 + (i % 26))
			}

			end := time.Now().Add(time.Duration(duration) * time.Second)
			start := time.Now()

			initialSize := 2678 * 1024 * 1024 // 2678 MB per core
			allocated := 0
			for allocated < initialSize {
				chunk := baseStr
				if allocated+len(chunk) > initialSize {
					chunk = chunk[:initialSize-allocated]
				}
				// Touch every page to ensure it's resident
				for page := 0; page < len(chunk); page += 4096 {
					chunk[page] = byte(page % 256)
				}
				memChunks = append(memChunks, append([]byte(nil), chunk...))
				allocated += len(chunk)
				memoryMutex.Lock()
				totalMemory += int64(len(chunk))
				currentGB := float64(totalMemory) / (1024 * 1024 * 1024)
				memoryMutex.Unlock()
				fmt.Printf("\rCore %d: Initial Memory: %.2f GB", id, currentGB)
			}

			for time.Now().Before(end) {
				elapsed := time.Since(start)
				durationSecs := float64(duration)
				elapsedSecs := elapsed.Seconds()

				// Calculate intensity (60% to 90%), but cap at maxCPU
				intensity := 0.6 + 0.3*(elapsedSecs/durationSecs)
				maxIntensity := float64(maxCPU) / 100.0
				if intensity > maxIntensity {
					intensity = maxIntensity
				}

				// Calculate target memory for this point in time, but cap at maxMemoryGB
				progressRatio := elapsedSecs / durationSecs
				targetTotalMemoryGB := 1.0 + (maxMemoryGB-1.0)*progressRatio
				if targetTotalMemoryGB > maxMemoryGB {
					targetTotalMemoryGB = maxMemoryGB
				}
				currentAllocated := float64(allocated) / (1024 * 1024 * 1024) // Convert to GB
				if currentAllocated < targetTotalMemoryGB {
					// Allocate odd MB chunk sizes: 1MB, 3MB, 5MB, ...
					// Cycle through odd values up to 19MB, then repeat
					oddMBs := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
					oddIdx := ((allocated / (1024 * 1024)) / 1) % len(oddMBs)
					chunkSize := oddMBs[oddIdx] * 1024 * 1024
					chunk := make([]byte, chunkSize)
					for page := 0; page < len(chunk); page += 4096 {
						chunk[page] = byte(page % 256)
					}
					memChunks = append(memChunks, chunk)
					allocated += chunkSize
					memoryMutex.Lock()
					totalMemory += int64(chunkSize)
					currentGB := float64(totalMemory) / (1024 * 1024 * 1024)
					memoryMutex.Unlock()
					fmt.Printf("\rCore %d: Memory: %.2f GB / %.2f GB, CPU: %.1f%%", id, currentGB, targetTotalMemoryGB, intensity*100)
				}

				// Intensive memory read/write operations combined with CPU work
				sum := 0.0
				// Number of memory operations to perform
				memOps := int(intensity * 1000)

				// Randomly read and write to allocated memory chunks while doing CPU work
				for j := 0; j < memOps; j++ {
					// CPU intensive work
					sum += rand.Float64()

					// Memory read/write operations
					if len(memChunks) > 0 {
						// Pick a random chunk
						chunkIdx := rand.Intn(len(memChunks))
						chunk := memChunks[chunkIdx]

						// Pick random pages within the chunk
						if len(chunk) > 4096 {
							pageIdx := (rand.Intn(len(chunk) / 4096)) * 4096
							// Read operation
							sum += float64(chunk[pageIdx])
							// Write operation
							chunk[pageIdx] = byte(rand.Intn(256))

							// Read and write to adjacent pages for more memory pressure
							if pageIdx+4096 < len(chunk) {
								sum += float64(chunk[pageIdx+4096])
								chunk[pageIdx+4096] = byte(rand.Intn(256))
							}
						}
					}
				}

				time.Sleep(time.Duration(5*(1-intensity)) * time.Millisecond)
			}
		}(i)
	}
	wg.Wait()
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("\nFinal Memory Usage: %.2f GB\n", float64(m.Alloc)/(1024*1024*1024))
}
