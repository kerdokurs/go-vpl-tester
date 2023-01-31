package main

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
	"time"
)

func runProgram(program, input string, args ...string) string {
	var inBuf bytes.Buffer
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	inBuf.Write([]byte(input))

	cmd := exec.Command(program, args...)
	cmd.Stdin = &inBuf
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	errCh := make(chan error)
	defer func() { close(errCh) }()
	timeout := time.NewTimer(testTimeout)

	go func() {
		if err := cmd.Run(); err != nil && strings.Contains(err.Error(), "killed") {
			// Program was killed in main Goroutine and errCh is almost certainly already closed
			return
		} else {
			errCh <- err
		}
	}()

	select {
	case <-timeout.C:
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("Ma ei suutnud programmi sulgeda: %v\n", err)
		}
		timeout.Stop()
		return "programm jooksis liiga kaua"
	case err := <-errCh:
		if err != nil {
			if errBuf.Len() > 0 {
				return errBuf.String()
			}
			return err.Error()
		}
		break
	}

	return strings.Trim(outBuf.String(), "\n")
}
