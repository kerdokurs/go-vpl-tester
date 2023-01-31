package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var testTimeout = flag.Duration("timeout", 4*time.Second, "test timeout")
var maxGrade = flag.Int("max-grade", 0, "max grade")
var testerEntrypoint = flag.String("entrypoint", "./tester", "tester entrypoint")
var casesFile = flag.String("tests", "tests.cases", "test cases file")

func main() {
	flag.Parse()

	if len(os.Args) == 2 && os.Args[1] == "help" {
		flag.Usage()
		os.Exit(0)
	}

	if tests, err := loadTests(); err != nil {
		log.Fatalf("Could not load tests: %v\n", err)
	} else {
		runTests(tests)
	}
}
