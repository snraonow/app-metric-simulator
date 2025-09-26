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
		fmt.Println("Usage: simulate high=memory,cpu time=20 frequency=constant|random")
		os.Exit(1)
	}

	highFlag := flag.String("high", "cpu", "Resource to stress: memory, cpu, or both (comma separated)")
	timeFlag := flag.Int("time", 10, "Duration in seconds")
	freqFlag := flag.String("frequency", "constant", "Load frequency: constant or random")
	maxCPUFlag := flag.Int("max-cpu", 60, "Maximum CPU usage percent (default 60)")
	maxMemoryFlag := flag.Float64("max-memory", 10.0, "Maximum memory usage in GB (default 10GB)")
	crashTimeFlag := flag.Int("crash-time", 0, "Time in seconds after which to simulate a crash (0 means no crash)")

	flag.CommandLine.Parse(os.Args[2:])

	resources := strings.Split(*highFlag, ",")
	duration := *timeFlag
	frequency := *freqFlag
	maxCPU := *maxCPUFlag
	maxMemoryGB := *maxMemoryFlag
	crashTime := *crashTimeFlag

	rand.Seed(time.Now().UnixNano())
	// Check if I/O simulation is requested
	simulateIO := false
	for _, r := range resources {
		if strings.ToLower(strings.TrimSpace(r)) == "io" {
			simulateIO = true
			break
		}
	}

	fmt.Printf("Simulating high %v for %d seconds with %s frequency...\n", resources, duration, frequency)

	// Run the simulation with all requested features
	simulateWithLimits(duration, frequency, maxCPU, maxMemoryGB, crashTime, simulateIO)
	fmt.Println("\nSimulation complete.")
}
