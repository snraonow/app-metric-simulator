package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

// simulateIO performs intensive I/O operations
func simulateIOOperations(tempDir string) error {
	// Create a temporary file
	f, err := os.CreateTemp(tempDir, "io-test-*.dat")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	// Create a large buffer (1MB) for I/O operations
	buffer := make([]byte, 1024*1024)
	for i := range buffer {
		buffer[i] = byte(rand.Intn(256))
	}

	// Write data in chunks
	for i := 0; i < 100; i++ {
		if _, err := f.Write(buffer); err != nil {
			return fmt.Errorf("write error: %v", err)
		}
		// Force sync to disk
		if err := f.Sync(); err != nil {
			return fmt.Errorf("sync error: %v", err)
		}
	}

	// Read back the data to generate read I/O
	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("seek error: %v", err)
	}

	readBuffer := make([]byte, 1024*1024)
	for {
		_, err := f.Read(readBuffer)
		if err != nil {
			break
		}
	}

	return nil
}

func simulateWithLimits(duration int, frequency string, maxCPU int, maxMemoryGB float64, crashTime int, simulateIO bool, batchSizeMinutes, batchWindowMinutes, simDurationMinutes int) {
	numCPU := runtime.NumCPU() / 2
	if numCPU < 1 {
		numCPU = 1
	}

	// Create temporary directory for I/O operations
	if simulateIO {
		tempDir, err := os.MkdirTemp("", "io-simulation")
		if err != nil {
			fmt.Printf("Failed to create temp directory: %v\n", err)
			return
		}
		defer os.RemoveAll(tempDir)
	}

	// Set up crash timer if crashTime is specified
	if crashTime > 0 {
		go func() {
			fmt.Printf("\nCrash simulation will trigger in %d seconds...\n", crashTime)
			time.Sleep(time.Duration(crashTime) * time.Second)
			fmt.Printf("\nTriggering crash simulation...\n")

			// Trigger an immediate panic with a stack trace
			panic("Simulated application crash triggered after " + fmt.Sprintf("%d", crashTime) + " seconds")
		}()
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

			// Set default values if not specified
			if batchSizeMinutes <= 0 {
				batchSizeMinutes = 5 // default first batch duration: 5 minutes
			}
			if batchWindowMinutes <= 0 {
				batchWindowMinutes = 20 // default batch window: 20 minutes
			}
			if simDurationMinutes <= 0 {
				simDurationMinutes = batchSizeMinutes // default simulation duration: 5 minutes
			}

			// Time windows
			firstBatchDuration := time.Duration(batchSizeMinutes) * time.Minute
			batchWindow := time.Duration(batchWindowMinutes) * time.Minute
			simDuration := time.Duration(simDurationMinutes) * time.Minute

			// Validate parameters
			if simDuration > batchWindow {
				fmt.Printf("\nWarning: Simulation duration (%d min) is larger than batch window (%d min). Adjusting simulation duration.\n",
					simDurationMinutes, batchWindowMinutes)
				simDuration = batchWindow
			}

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

			// Start I/O simulation in a separate goroutine if enabled
			if simulateIO {
				// Create temporary directory for I/O operations
				tempDir, err := os.MkdirTemp("", "simulator-io")
				if err != nil {
					fmt.Printf("\nError creating temp directory: %v\n", err)
					return
				}
				defer os.RemoveAll(tempDir)

				go func() {
					for time.Now().Before(end) {
						if err := simulateIOOperations(tempDir); err != nil {
							fmt.Printf("\nI/O simulation error: %v\n", err)
						}
						// Brief pause between I/O operations
						time.Sleep(time.Millisecond * 100)
					}
				}()
			}

			for time.Now().Before(end) {
				elapsed := time.Since(start)

				// Determine if we're in an active simulation window
				isSimulationActive := false

				if elapsed < firstBatchDuration {
					// First 5 minutes: always active
					isSimulationActive = true
				} else {
					// Calculate which 20-minute window we're in
					windowNumber := int(elapsed.Minutes() / batchWindow.Minutes())
					windowStart := start.Add(time.Duration(windowNumber) * batchWindow)

					// Generate a random 5-minute slot within this window
					// Use windowNumber as seed for consistent behavior within the same window
					r := rand.New(rand.NewSource(int64(windowNumber)))
					randomOffset := time.Duration(r.Intn(int(batchWindow.Minutes()-simDuration.Minutes()))) * time.Minute
					simStart := windowStart.Add(randomOffset)
					simEnd := simStart.Add(simDuration)

					// Check if current time falls within the simulation period
					now := time.Now()
					isSimulationActive = now.After(simStart) && now.Before(simEnd)
				}

				// Calculate base intensity
				var intensity float64
				if isSimulationActive {
					// High intensity during active simulation (60% to 90%)
					intensity = 0.6 + 0.3*(elapsed.Seconds()/float64(duration))
				} else {
					// Reduced load during inactive periods (0% to 10%)
					intensity = 0.1 * (elapsed.Seconds() / float64(duration))
				}

				// Cap at maxCPU
				maxIntensity := float64(maxCPU) / 100.0
				if intensity > maxIntensity {
					intensity = maxIntensity
				}

				// Calculate target memory based on simulation state
				progressRatio := elapsed.Seconds() / float64(duration)
				var targetTotalMemoryGB float64
				if isSimulationActive {
					// High memory usage during active simulation (50% to 100% of maxMemoryGB)
					targetTotalMemoryGB = (maxMemoryGB * 0.5) + (maxMemoryGB * 0.5 * progressRatio)
				} else {
					// Reduced memory usage during inactive periods (10% to 30% of maxMemoryGB)
					targetTotalMemoryGB = (maxMemoryGB * 0.1) + (maxMemoryGB * 0.2 * progressRatio)
				}
				if targetTotalMemoryGB > maxMemoryGB {
					targetTotalMemoryGB = maxMemoryGB
				}

				currentAllocated := float64(allocated) / (1024 * 1024 * 1024) // Convert to GB

				// Adjust memory if needed (both up and down)
				if currentAllocated > targetTotalMemoryGB*1.1 { // Allow 10% buffer before reducing
					// Release some memory by removing chunks
					numChunksToRemove := int((currentAllocated - targetTotalMemoryGB) * float64(len(memChunks)) / currentAllocated)
					if numChunksToRemove > 0 && len(memChunks) > numChunksToRemove {
						memChunks = memChunks[numChunksToRemove:]
						allocated = 0
						for _, chunk := range memChunks {
							allocated += len(chunk)
						}
					}
				} else if currentAllocated < targetTotalMemoryGB {
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
