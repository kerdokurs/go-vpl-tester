package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"io"
	"log"
	"plugin"
	"time"
)

type Test struct {
	Case      string
	Input     string
	Output    string
	Args      []string
	Program   string
	CompareTo string
	Plugin    string
	Timeout   *time.Duration
}

func loadTests() ([]Test, error) {
	tests := make([]Test, 0)
	cfg, err := ini.Load(*casesFile)
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
	programToRun := *testerEntrypoint
	if t.Program != "" {
		programToRun = t.Program
	}
	timeout := *testTimeout
	if t.Timeout != nil {
		timeout = *t.Timeout
	}

	var customFunc func(w io.Writer, outCh chan string, isResultClosingCh chan struct{}, resultCh chan string)
	if t.Plugin != "" {
		pl, err := plugin.Open(t.Plugin)
		if err != nil {
			log.Printf("pluginat %v ei leidu\n", t.Plugin)
			return false
		}
		funcSymbol, err := pl.Lookup("Run")
		if err != nil {
			log.Printf("plugin %v ei sisalda Run funktsiooni\n", t.Plugin)
			return false
		}

		customFunc = funcSymbol.(func(w io.Writer, outCh chan string, isResultClosingCh chan struct{}, resultCh chan string))
	}

	studentOutput := runProgram(programToRun, t.Input, customFunc, timeout, t.Args...)

	correctOutput := t.Output
	if correctOutput == "" && t.CompareTo != "" {
		correctOutput = runProgram(t.CompareTo, t.Input, customFunc, timeout, t.Args...)
	} else if correctOutput == "" {
		correctOutput = "OK"
	}

	if studentOutput != correctOutput {
		fmt.Println("\nComment :=>>- Ootasin väljundit:")
		preFormat(correctOutput)
		fmt.Println("Comment :=>>- Aga sain:")
		preFormat(studentOutput)
		return false
	}

	fmt.Println(" – OK.")
	return true
}

func runTests(tests []Test) {
	nTests := len(tests)
	gradePoint := float32(*maxGrade) / float32(nTests)
	var grade float32

	for i, test := range tests {
		fmt.Printf("Comment :=>>-Test %d. `%s`", i+1, test.Case)
		if ok := runTest(&test); ok {
			grade += gradePoint
		}
	}

	fmt.Printf("Grade :=>> %f\n", grade)
}
