package main

import (
	"bytes"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

type writeWatcher struct {
	buf *bytes.Buffer
	out chan string
}

func (tw *writeWatcher) Write(b []byte) (int, error) {
	n, err := tw.buf.Write(b)

	if tw.out != nil {
		tw.out <- string(b)
	}

	return n, err
}

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
	var isResultChClosing chan struct{}

	if customFunc == nil {
		inBuf.Write([]byte(input + "\n"))
	} else {
		resultCh = make(chan string)
		outCh = make(chan string, 10)
		isResultChClosing = make(chan struct{})
		defer func() {
			isResultChClosing <- struct{}{}
			close(isResultChClosing)
			close(resultCh)
			close(outCh)
		}()
		go customFunc(inBuf, outCh, isResultChClosing, resultCh)
	}

	outBuf := writeWatcher{
		buf: &bytes.Buffer{},
		out: outCh,
	}
	cmd.Stdout = &outBuf

	errCh := make(chan error, 1)
	timer := time.NewTimer(timeout)
	defer func() {
		close(errCh)
		timer.Stop()
	}()

	go func() {
		if err := cmd.Run(); err != nil && strings.Contains(err.Error(), "killed") {
			// Program was killed in main Goroutine and errCh is almost certainly already closed
			return
		} else {
			errCh <- err
		}
	}()

	select {
	case <-timer.C:
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("Ma ei suutnud programmi sulgeda: %v\n", err)
		}
		return "programm jooksis liiga kaua"
	case err := <-errCh:
		if err != nil {
			if errBuf.Len() > 0 {
				return errBuf.String()
			}
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
		select {
		case <-timer.C:
			return "plugin ei vastanud õigeaegselt"
		case result := <-resultCh:
			return result
		}
	}
}
