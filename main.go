package main

import (
	"log"
)

func main() {
	loadEnv()

	if tests, err := loadTests(); err != nil {
		log.Fatalf("Could not load tests: %v\n", err)
	} else {
		runTests(tests)
	}
}
