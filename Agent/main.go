package main

import (
	"fmt"
	"sync"
	"time"

	"andreabarchietto.it/oss/go/telemetry_agent/net"
)

func mainOld() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Get initial counters across all interfaces
	prevCounters, err := net.IOCounters(false)
	if err != nil || len(prevCounters) == 0 {
		panic(err)
	}
	prevTime := time.Now()

	fmt.Println("Monitoring network speed (kbps)... Press Ctrl+C to stop.")

	for range ticker.C {
		currCounters, err := net.IOCounters(false)
		if err != nil || len(currCounters) == 0 {
			continue
		}
		currTime := time.Now()

		// Calculate elapsed time in seconds
		duration := currTime.Sub(prevTime).Seconds()

		// Bytes transmitted during the interval
		bytesSentDelta := currCounters[0].BytesSent - prevCounters[0].BytesSent
		bytesRecvDelta := currCounters[0].BytesRecv - prevCounters[0].BytesRecv

		// Convert Bytes -> Bits (x 8), then Bits -> Kilobits (/ 1000)
		// kbps = (bytes * 8) / 1000 / seconds
		txKbps := (float64(bytesSentDelta) * 8.0 / 1000.0) / duration
		rxKbps := (float64(bytesRecvDelta) * 8.0 / 1000.0) / duration

		fmt.Printf("[%s] Down: %.2f kbps | Up: %.2f kbps\n",
			currTime.Format("15:04:05"),
			rxKbps,
			txKbps,
		)

		// Update previous state for next iteration
		prevCounters = currCounters
		prevTime = currTime
	}
}

func dataExporter(takeOff chan<- string) {
	tNet := time.NewTicker(1 * time.Second)
	tCpu := time.NewTicker(1 * time.Second)
	tRam := time.NewTicker(1 * time.Second)
	tDisk := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-tNet.C:
			takeOff <- "net"
		case <-tCpu.C:
			takeOff <- "cpu"
		case <-tRam.C:
			takeOff <- "ram"
		case <-tDisk.C:
			takeOff <- "disk"
		}
	}
}

func dataSender(takeOff <-chan string) {
	for {
		data := <-takeOff
		fmt.Println("Data:", data)
	}
}

func main() {
	takeOff := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go dataExporter(takeOff)
	go dataSender(takeOff)

	wg.Wait()
	close(takeOff)
}
