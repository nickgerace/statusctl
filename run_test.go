package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestRun(t *testing.T) {

	// Begin test duration timer and suppress STDOUT.
	testStart := time.Now()
	fmt.Println("Beginning test. Please wait. Avoiding STDOUT capture...")
	originalStdout := os.Stdout
	os.Stdout = nil

	// User-modified declarations. All INTs must be the same. Note that "time.Duration" wraps
	// the built-in type "int64". The "count" variable does not have an initialized value, but
	// does change its type with the "iterations" variable.
	var iterations int16 = 1000
	var count int16
	var total time.Duration = 1000
	var results [1000]time.Duration

	// Start the tests.
	for count = 0; count < iterations; count++ {
		start := time.Now()
		runAction()
		results[count] = time.Since(start)
	}

	// Calculate the average for the test results.
	var average time.Duration
	for _, result := range results {
		average += result
	}

	// Print the results. Get the average while printing the test by test results. Bring STDOUT
	// back to its original state.
	os.Stdout = originalStdout
	fmt.Printf("\nTest average over %d runs...\n%s\n", iterations, average/total)
	fmt.Printf("\nTest duration...\n%s\n\n", time.Since(testStart))
}
