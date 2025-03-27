package attempt1

import (
	"testing"
)

func TestToFloat(t *testing.T) {
	testCases := map[string]struct {
		value float32
		bytes []byte
	}{
		"32.5": {
			value: 32.5,
			bytes: []byte{'3', '2', '.', '5'},
		},
		"-32.5": {
			value: -32.5,
			bytes: []byte{'-', '3', '2', '.', '5'},
		},
		"2.5": {
			value: 2.5,
			bytes: []byte{'2', '.', '5'},
		},
		"99.9": {
			value: 99.9,
			bytes: []byte{'9', '9', '.', '9'},
		},
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			if newValue := toFloat(tC.bytes); newValue != tC.value {
				t.Fatalf("Incorrect value. Wanted %v, got %v", tC.value, newValue)
			}
		})
	}
}

func BenchmarkTest(b *testing.B) {
	toFloat([]byte{'7', '4', '.', '2'})
}

func TestParseData(t *testing.T) {
	testCases := map[string]struct {
		stationName string
		max         float32
		min         float32
		count       float32
		sum         float32
		bytes       []byte
	}{
		"1 value": {
			stationName: "test",
			max:         32.4,
			min:         32.4,
			count:       1.0,
			sum:         32.4,
			bytes:       []byte("test;32.4\n"),
		},
		"1 value, no newline": {
			stationName: "test",
			max:         32.4,
			min:         32.4,
			count:       1.0,
			sum:         32.4,
			bytes:       []byte("test;32.4"),
		},
		"2 values": {
			stationName: "test",
			max:         50.0,
			min:         10.0,
			count:       2.0,
			sum:         60.0,
			bytes:       []byte("test;50.0\ntest;10.0\n"),
		},
		"2 values, no newline": {
			stationName: "test",
			max:         50.0,
			min:         20.0,
			count:       2.0,
			sum:         70.0,
			bytes:       []byte("test;50.0\ntest;20.0"),
		},
	}

	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			memMap := make(map[string]values, 10000)
			parseData(memMap, tC.bytes)

			newValues, ok := memMap[tC.stationName]
			if !ok {
				t.Fatalf("station not in mem")
			}

			if newValues.min != tC.min {
				t.Fatalf("Incorrect min value. Wanted %v, got %v", tC.min, newValues.min)
			}
			if newValues.max != tC.max {
				t.Fatalf("Incorrect max value. Wanted %v, got %v", tC.max, newValues.max)
			}
			if newValues.sum != tC.sum {
				t.Fatalf("Incorrect sum value. Wanted %v, got %v", tC.sum, newValues.sum)
			}
			if newValues.count != tC.count {
				t.Fatalf("Incorrect count value. Wanted %v, got %v", tC.count, newValues.count)
			}
		})
	}
}

func BenchmarkParseData(b *testing.B) {

	benchmarks := map[string][]byte{
		"1 item":   []byte("test;50.0\n"),
		"10 item":  []byte("test;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\n"),
		"100 item": []byte("test;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\ntest;50.0\n"),
	}
	for name, testBytes := range benchmarks {
		b.Run(name, func(b *testing.B) {
			parseData(make(map[string]values, 10000), testBytes)
		})
	}
}
