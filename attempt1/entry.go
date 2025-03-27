package attempt1

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

// A very simple look at the challenge. Just reading the file, parsing the
// content and returning. We done this with probing, so its just how much we can get out
// of the Go std without having to dig through source code.

// So the first attempt will do a couple things. First, we are not using concurrency.
// That will likely be the last thing we do, since getting it as fast as possible
// is the best place to start.

const (
	// Size of the datablobs we read at a time.
	blobSize = 10000000
)

type values struct {
	count float32
	sum   float32
	max   float32
	min   float32
}

func (v *values) string() string {
	return fmt.Sprintf("%.2f/%.2f/%.2f", v.min, (v.sum / v.count), v.max)
}

func Challenge() {
	// Max size of 10000 station names. We can allocate all of that now
	memMap := make(map[string]values, 10000)

	// Read file
	file, err := os.Open("measurements.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(scannerFunc)
	scanner.Buffer(make([]byte, blobSize), blobSize)

	for scanner.Scan() {
		parseData(memMap, scanner.Bytes())
	}

	// Print
	// "Abha=-23.0/18.0/59.2", but in JSON like format

	fmt.Print("\n\n{")
	for station, values := range memMap {
		fmt.Printf("\"%v=%v\",", station, values.string())
	}
	fmt.Println("}")
}

// We know the scanner's default readlines reads up to the next newline. We wanna
// Minimise reads so we're gonna pull down X Bytes at a time.
func scannerFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	fmt.Print(".")
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

// We pass the output into the function to avoid the GC
// We need the min, max and avg of the station
func parseData(output map[string]values, data []byte) {
	// Define all the vars so they don't get reassigned in the loop, and its easier
	// to keep track

	var nextEnd int
	var currentLine []byte

	var sepIndex int
	var stationName []byte
	var value []byte
	var floatValue float32

	var ok bool
	var currentStationStats values

	for {
		nextEnd = bytes.IndexByte(data, '\n')
		if nextEnd < 0 {
			currentLine = data
		} else {
			currentLine = data[:nextEnd]
		}

		if len(data) == 0 {
			break
		}

		// Now we have the current line, we need to split that on ';'
		sepIndex = bytes.IndexByte(data, ';')
		stationName = currentLine[:sepIndex]
		value = currentLine[sepIndex+1:]
		floatValue = toFloat(value)

		// We need the station name to be a string, and the value to be a float
		currentStationStats, ok = output[string(stationName)]

		if !ok {
			currentStationStats = values{min: 100}
		}

		currentStationStats.sum += floatValue
		currentStationStats.count += 1
		if currentStationStats.max < floatValue {
			currentStationStats.max = floatValue
		}
		if currentStationStats.min > floatValue {
			currentStationStats.min = floatValue
		}

		// Reassign values back to map
		output[string(stationName)] = currentStationStats

		if nextEnd < 0 {
			break
		}
		data = data[nextEnd+1:]
	}
}

// Kinda gross, but the byte is defined by its index. e.g. byte 48 is "0"
// Should be quicker than a map since its so small.
var digitMap = []byte{
	48,
	49,
	50,
	51,
	52,
	53,
	54,
	55,
	56,
	57,
}

// Small enough to just allocate. Plus it doesn't loop much.
func toFloat(input []byte) float32 {
	// Check for negative
	var multiplier float32 = 1.0

	if input[0] == '-' {
		multiplier = -1.0
		input = input[1:]
	}

	// Check digit count. We know it will be either 4 (XX.X) or 3 (X.X)
	var startTimes float32 = 10
	if len(input) == 3 {
		startTimes = 1
	}

	var endFloat float32
	var byteIndex int

	for index := range input {
		if input[index] == '.' {
			continue
		}

		for byteIndex = range digitMap {
			if digitMap[byteIndex] == input[index] {
				break
			}
		}

		endFloat += (float32(byteIndex) * startTimes)
		startTimes /= 10
	}

	return endFloat * multiplier
}
