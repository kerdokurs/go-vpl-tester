package main

import (
	"os"
	"strconv"
	"time"
)

var testTimeout = 4 * time.Second
var maxGrade = 0
var testerEndpoint = "./tester"

func loadEnv() {
	maxGrade, _ = strconv.Atoi(os.Getenv("MAX_GRADE"))

	if testTimeOutInt, err := strconv.Atoi(os.Getenv("TEST_TIMEOUT")); err == nil {
		testTimeout = time.Duration(testTimeOutInt) * time.Second
	}

	if os.Getenv("TESTER_ENDPOINT") != "" {
		testerEndpoint = os.Getenv("TESTER_ENDPOINT")
	}
}
