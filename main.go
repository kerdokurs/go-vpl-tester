package main

import (
	"flag"
	"log"
	"time"
)

var testTimeout = flag.Duration("timeout", 4*time.Second, "test timeout")
var maxGrade = flag.Int("max-grade", 0, "max grade")
var testerEntrypoint = flag.String("entrypoint", "./tester", "tester entrypoint")
var casesFile = flag.String("tests", "tests.cases", "test cases file")

func main() {
	flag.Parse()

	if tests, err := loadTests(); err != nil {
		log.Fatalf("Could not load tests: %v\n", err)
	} else {
		runTests(tests)
	}
}
