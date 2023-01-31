package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
)

type Test struct {
	Case   string
	Input  string
	Output string
	Args   []string
}

func loadTests() ([]Test, error) {
	p := "test.cases"
	if len(os.Args) > 1 {
		p = os.Args[1]
	}

	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tests := make([]Test, 0)
	cfg, err := ini.Load(file)
	if err != nil {
		panic(err)
	}

	for _, section := range cfg.Sections() {
		if section.Name() == "DEFAULT" {
			continue
		}

		var test Test
		if err = section.MapTo(&test); err != nil {
			log.Printf("Could not decode test: %v\n", err)
			continue
		}
		test.Case = section.Name()
		tests = append(tests, test)
	}

	return tests, err
}

func runTest(t *Test) bool {
	studentOutput := runProgram(testerEndpoint, t.Input, t.Args...)

	if studentOutput != t.Output {
		fmt.Println("\nComment :=>>- Ootasin väljundit:")
		preFormat(t.Output)
		fmt.Println("Comment :=>>- Aga sain:")
		preFormat(studentOutput)
		return false
	}

	fmt.Println(" – OK.")
	return true
}

func runTests(tests []Test) {
	nTests := len(tests)
	gradePoint := float32(maxGrade) / float32(nTests)
	var grade float32

	for i, test := range tests {
		fmt.Printf("Comment :=>>-Test %d. `%s`", i+1, test.Case)
		if ok := runTest(&test); ok {
			grade += gradePoint
		}
	}

	fmt.Printf("Grade :=>> %f\n", grade)
}
