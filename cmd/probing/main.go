package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"slices"
	"time"
)

const (
	filePath = "measurements.txt"
)

func main() {
	if err := runService(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

type dataBlob struct {
	station string
	value   float64
}

var (
	// Values we wanna probe
	maxValue = 0.0
	minValue = 100.0

	maxDecimalPlaces = 0

	sum     = 0
	records = 0
	scans   = 0

	longestName  = 0
	shortestName = 100

	maxScanSize = 0
	minScanSize = 10000000

	firstStation = ""
	lastStation  = ""
)

func runService() error {
	// Read over and collect information about the data
	fmt.Println("starting now")

	challengeFile, err := os.Open(filePath)
	if err != nil {
		return err
	}

	fileStats, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	// 10MB buffer
	const maxBuffer = 10000000
	readBuffer := make([]byte, maxBuffer)

	scanner := bufio.NewScanner(challengeFile)
	scanner.Split(shittyScannerFunc)
	scanner.Buffer(readBuffer, maxBuffer)

	scansNeeded := fileStats.Size() / maxBuffer

	startTime := time.Now()

	// Scan takes 18 seconds
	scanStart := time.Now()
	for scanner.Scan() {

		scannedBytes := scanner.Bytes()

		if len(scannedBytes) > maxScanSize {
			maxScanSize = len(scannedBytes)
		}
		if len(scannedBytes) < minScanSize {
			minScanSize = len(scannedBytes)
		}

		shittyDataParser(scannedBytes)

		scanEnd := time.Now()
		scans += 1

		fmt.Printf("Scans: %d. Per Scan: %v. Time Remaining: %v\r", scans, scanEnd.Sub(scanStart), time.Duration(int(scansNeeded)-scans)*scanEnd.Sub(scanStart))
		scanStart = time.Now()
	}

	if scanner.Err() != nil {
		return err
	}

	endTime := time.Now()

	fmt.Println("\n\n###################### RESULTS ###########################\n\n")

	fmt.Println("=== Station ===")
	fmt.Printf("First station: %v\n", firstStation)
	fmt.Printf("Last station: %v\n", lastStation)
	fmt.Printf("Longest Station Name in Bytes: %v\n", longestName)
	fmt.Printf("Shortest Station Name in Bytes: %v\n", shortestName)

	fmt.Println()

	fmt.Println("=== Value ===")
	fmt.Printf("Max Value: %v\n", maxValue)
	fmt.Printf("Min Value: %v\n", minValue)
	fmt.Printf("Max Decimal Places: %v\n", maxDecimalPlaces)
	fmt.Printf("Sum of Values: %v\n", sum)
	fmt.Printf("Avg of Values: %v\n", sum/records)

	fmt.Println()

	fmt.Println("=== Other stats ===")
	fmt.Printf("Scans Performed: %v\n", scans)
	fmt.Printf("Records Read: %v\n", records)
	fmt.Printf("Records Per Scan: %v\n", records/scans)
	fmt.Printf("Time Taken: %v\n", endTime.Sub(startTime))
	fmt.Printf("Max scan size: %v\n", maxScanSize)
	fmt.Printf("Min scan size: %v\n", minScanSize)

	return nil
}

var digitMap = map[byte]float64{
	48: 0,
	49: 1,
	50: 2,
	51: 3,
	52: 4,
	53: 5,
	54: 6,
	55: 7,
	56: 8,
	57: 9,
}

func shittyBtoF(input []byte) float64 {
	indexLocation := 0
	mutliplier := 1.0
	finalNumber := 0.0

	if input[0] == 45 { // -
		mutliplier = -1.0
		input = input[1:]
	}

	for index := range input {
		if input[index] == 46 { // .
			indexLocation = index
		}
	}

	if len(input[indexLocation+1:]) > maxDecimalPlaces {
		maxDecimalPlaces = len(input[indexLocation+1:])
	}

	input = slices.Delete(input, indexLocation, indexLocation+1)
	indexLocation -= 1

	for index := range input {
		finalNumber += digitMap[input[index]] * float64(math.Pow10(indexLocation-index))
	}

	return finalNumber * mutliplier
}

func shittyScannerFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		if len(data) == 0 {
			data = nil
		}
		return len(data), data, nil
	}

	for i := len(data) - 1; i > 0; i-- {
		if data[i] == '\n' {
			return i + 1, data[:i], nil
		}
	}

	return 0, nil, nil
}

var (
	first = true
)

// Updates stuff directly
func shittyDataParser(data []byte) {
	for {
		index := bytes.IndexByte(data, '\n')

		if index < 0 {
			index = len(data) - 1
		}

		currentLine := data[:index]
		sepIndex := bytes.IndexByte(currentLine, ';')

		// check station
		if len(currentLine[:sepIndex]) > longestName {
			longestName = len(currentLine[:sepIndex])
		}
		if len(currentLine[:sepIndex]) < shortestName {
			shortestName = len(currentLine[:sepIndex])
		}

		if first {
			firstStation = string(currentLine[:sepIndex])
			first = false
		}
		lastStation = string(currentLine[:sepIndex])

		// Check Value
		floatValue := shittyBtoF(currentLine[sepIndex+1:])
		if floatValue > maxValue {
			maxValue = floatValue
		}
		if floatValue < minValue {
			minValue = floatValue
		}

		// Floating point mathmatics... yay
		// This is probs accurate enough, Its just probing the data it will be fine
		sum += int(floatValue)
		records++

		if index == len(data)-1 {
			break
		}

		data = data[index+1:]
	}
}
