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

	memoryPerCore := (maxMemoryGB * 1024 * 1024 * 1024) / float64(numCPU)

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var stringSlice []string
			baseStr := make([]byte, 1024*1024)
			for i := range baseStr {
				baseStr[i] = byte(65 + (i % 26))
			}

			end := time.Now().Add(time.Duration(duration) * time.Second)
			start := time.Now()

			initialSize := int(memoryPerCore)
			allocated := 0
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
				targetMemoryPerCore := (targetTotalMemoryGB * 1024 * 1024 * 1024) / float64(numCPU)

				currentAllocated := int64(allocated)
				if float64(currentAllocated) < targetMemoryPerCore {
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
					fmt.Printf("\rCore %d: Memory: %.2f GB / %.2f GB, CPU: %.1f%%", id, currentGB, targetTotalMemoryGB, intensity*100)
				}

				sum := 0.0
				for j := 0; j < int(intensity*10000); j++ {
					sum += rand.Float64()
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
