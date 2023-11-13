package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mattzi/dataprocessor-go/rplidar"
)

func main() {
	// Initialize the RPLidar
	lidar := rplidar.NewRPLidar("/dev/rplidar", 115200, time.Second*3)

	// Connect to the RPLidar
	err := lidar.Connect()
	if err != nil {
		log.Fatalf("Error connecting to RPLidar: %v", err)
	}
	defer lidar.Disconnect()

	// Retrieve and print device information
	info, err := lidar.GetInfo()
	if err != nil {
		log.Fatalf("Error getting info: %v", err)
	}
	fmt.Printf("RPLidar Info: %+v\n", info)

	time.Sleep(time.Second * 1)

	// Start iterating over measurements
	measurements, err := lidar.IterMeasurements()
	if err != nil {
		log.Fatalf("Error starting measurements: %v", err)
	}

	SIGTERM := make(chan os.Signal, 1)
	signal.Notify(SIGTERM, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case measurement, ok := <-measurements:
			if !ok {
				return // channel closed, end the loop
			}
			fmt.Printf("Measurement: %+v\n", measurement)
		case <-SIGTERM:
			fmt.Println("SIGTERM/SIGINT Received")
			return // end after timeout
		}
	}
}
