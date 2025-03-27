package main

import (
	"attempt/attempt1"
	"fmt"
	"os"
	"time"
)

func main() {
	if err := runService(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func runService() error {
	startTime := time.Now()

	// Attempt here
	attempt1.Challenge()

	endTime := time.Now()

	fmt.Printf("\n\nTime Taken: %v\n", endTime.Sub(startTime))

	return nil
}
