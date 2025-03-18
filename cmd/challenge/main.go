package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if err := runService(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func getFlags() int {
	attemptNumber := flag.Int("-a", 1, "Attempt Number")

	flag.Parse()

	return *attemptNumber
}

func runService() error {
	// get attempt number
	attempt := getFlags()

	if attempt == 1 {

	}
	return nil
}
