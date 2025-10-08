package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "simulate" {
		fmt.Println("Usage: simulate [options]")
		fmt.Println("\nResource Options:")
		fmt.Println("  --high string        Resource to stress: memory,cpu,io (comma separated)")
		fmt.Println("  --time int           Duration in seconds")
		fmt.Println("  --frequency string   Load frequency: constant or random")
		fmt.Println("\nLimit Options:")
		fmt.Println("  --max-cpu int        Maximum CPU usage percent (default 60)")
		fmt.Println("  --max-memory float   Maximum memory usage in GB (default 10GB)")
		fmt.Println("\nBatch Configuration:")
		fmt.Println("  --batch-size int     Size of the first batch in minutes (default 5)")
		fmt.Println("  --batch-window int   Size of the batch window in minutes (default 20)")
		fmt.Println("  --sim-duration int   Duration of simulation within each batch (default 5)")
		fmt.Println("\nOther Options:")
		fmt.Println("  --crash-time int     Time in seconds after which to simulate a crash (0 = no crash)")
		fmt.Println("\nExample:")
		fmt.Println("  simulate --high memory,cpu --time 3600 --batch-size 10 --batch-window 30 --sim-duration 5")
		os.Exit(1)
	}

	highFlag := flag.String("high", "cpu", "Resource to stress: memory, cpu, or both (comma separated)")
	timeFlag := flag.Int("time", 10, "Duration in seconds")
	freqFlag := flag.String("frequency", "constant", "Load frequency: constant or random")
	maxCPUFlag := flag.Int("max-cpu", 60, "Maximum CPU usage percent (default 60)")
	maxMemoryFlag := flag.Float64("max-memory", 10.0, "Maximum memory usage in GB (default 10GB)")
	crashTimeFlag := flag.Int("crash-time", 0, "Time in seconds after which to simulate a crash (0 means no crash)")
	batchSizeFlag := flag.Int("batch-size", 5, "Size of the first batch in minutes (default 5)")
	batchWindowFlag := flag.Int("batch-window", 20, "Size of the batch window in minutes (default 20)")
	simDurationFlag := flag.Int("sim-duration", 5, "Duration of simulation within each batch in minutes (default 5)")

	flag.CommandLine.Parse(os.Args[2:])

	resources := strings.Split(*highFlag, ",")
	duration := *timeFlag
	frequency := *freqFlag
	maxCPU := *maxCPUFlag
	maxMemoryGB := *maxMemoryFlag
	crashTime := *crashTimeFlag

	rand.Seed(time.Now().UnixNano())
	// Validate batch parameters
	if *batchWindowFlag <= 0 {
		fmt.Println("Error: Batch window must be greater than 0")
		os.Exit(1)
	}
	if *simDurationFlag <= 0 {
		fmt.Println("Error: Simulation duration must be greater than 0")
		os.Exit(1)
	}
	if *simDurationFlag > *batchWindowFlag {
		fmt.Printf("Warning: Simulation duration (%d) is larger than batch window (%d). Adjusting simulation duration.\n",
			*simDurationFlag, *batchWindowFlag)
		*simDurationFlag = *batchWindowFlag
	}

	// Check if I/O simulation is requested
	simulateIO := false
	for _, r := range resources {
		if strings.ToLower(strings.TrimSpace(r)) == "io" {
			simulateIO = true
			break
		}
	}

	// Print simulation configuration
	fmt.Println("\nSimulation Configuration:")
	fmt.Printf("- Resources: %v\n", resources)
	fmt.Printf("- Duration: %d seconds\n", duration)
	fmt.Printf("- Frequency: %s\n", frequency)
	fmt.Printf("- Max CPU: %d%%\n", maxCPU)
	fmt.Printf("- Max Memory: %.2f GB\n", maxMemoryGB)
	fmt.Printf("- I/O Simulation: %v\n", simulateIO)
	fmt.Printf("\nBatch Configuration:")
	fmt.Printf("- First Batch Duration: %d minutes\n", *batchSizeFlag)
	fmt.Printf("- Batch Window: %d minutes\n", *batchWindowFlag)
	fmt.Printf("- Simulation Duration: %d minutes\n", *simDurationFlag)
	if crashTime > 0 {
		fmt.Printf("\nCrash scheduled after: %d seconds\n", crashTime)
	}
	fmt.Println("\nStarting simulation...")

	// Run the simulation with all requested features
	simulateWithLimits(duration, frequency, maxCPU, maxMemoryGB, crashTime, simulateIO,
		*batchSizeFlag, *batchWindowFlag, *simDurationFlag)
	fmt.Println("\nSimulation complete.")
}
