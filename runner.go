package main

import (
	"bytes"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

func runProgram(program, input string, customFunc func(io.Writer, chan string, chan struct{}, chan string), timeout time.Duration, args ...string) string {
	var errBuf bytes.Buffer

	cmd := exec.Command(program, args...)
	cmd.Stderr = &errBuf

	inBuf, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer inBuf.Close()

	var outCh chan string
	var resultCh chan string
	var isResultChClosingCh chan struct{}

	if customFunc == nil {
		inBuf.Write([]byte(input + "\n"))
	} else {
		// custom plugin function will write the result into the resultCh
		resultCh = make(chan string)
		isResultChClosingCh = make(chan struct{})
		// output from the program will also be put into outCh for interactivity in the custom plugin function
		outCh = make(chan string, 10)
		defer func() {
			// we should notify the plugin function (which runs in a different goroutine) that it's time to wrap up
			isResultChClosingCh <- struct{}{}
			close(isResultChClosingCh)
			close(resultCh)
			close(outCh)
		}()

		// spin up the custom plugin function
		go customFunc(inBuf, outCh, isResultChClosingCh, resultCh)
	}

	outBuf := writeWatcher{
		buf: &bytes.Buffer{},
		out: outCh,
	}
	cmd.Stdout = &outBuf

	programExitedCh := make(chan error, 1)
	timer := time.NewTimer(timeout)
	defer func() {
		close(programExitedCh)
		timer.Stop()
	}()

	go func() {
		if err := cmd.Run(); err != nil && strings.Contains(err.Error(), "killed") {
			// program was killed in main Goroutine and programExitedCh is almost certainly already closed
			return
		} else {
			programExitedCh <- err
		}
	}()

	select {
	case <-timer.C:
		// if the program times out, we should kill it
		// this will cause the cmd.Run() to return an error (which states that the process was killed) and exit
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("Ma ei suutnud programmi sulgeda: %v\n", err)
		}
		return "programm jooksis liiga kaua"
	case err := <-programExitedCh:
		// this case is selected when the program has exited in the other Goroutine
		// err is nil if the program exited successfully
		// we should prioritise stderr output for error message
		if errBuf.Len() > 0 {
			return errBuf.String()
		}
		if err != nil {
			return err.Error()
		}
		break
	case result := <-resultCh:
		// note: this select case will only run when a plugin is defined
		// if this case gets selected the plugin wants to exit and should get priority
		if cmd.Process != nil {
			if err := cmd.Process.Kill(); err != nil {
				log.Printf("Ma ei suutnud programmi sulgeda: %v\n", err)
			}
		}
		return result
	}

	if customFunc == nil {
		output := outBuf.buf.String()
		return strings.Trim(output, "\n")
	} else {
		// if the program exits before the plugin returns some data on the channel
		// or get stuck in an infinite loop itself, we should not block indefinitely and error out if
		// the timer runs out
		select {
		case <-timer.C:
			return "plugin ei vastanud õigeaegselt"
		case result := <-resultCh:
			return result
		}
	}
}
