package main

import (
	"io"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
)

var r = regexp.MustCompile("kaart\\: (?P<current>\\d+)\\. Hetkeseis\\: (?P<total>\\d+)")

func Run(w io.Writer, outCh chan string, isResultClosingCh chan struct{}, resultCh chan string) {
	var isResultChClosed atomic.Int32
	go func() {
		for range isResultClosingCh {
			isResultChClosed.Add(1)
		}
	}()

	total := 0
	lastOut := ""

	for out := range outCh {
		foundMap := make(map[string]string)
		match := r.FindStringSubmatch(out)

		if len(match) > 0 {
			for i, name := range r.SubexpNames() {
				if i > 0 && i <= len(match) {
					foundMap[name] = match[i]
				}
			}

			gotCurrent, _ := strconv.Atoi(foundMap["current"])
			gotTotal, _ := strconv.Atoi(foundMap["total"])

			total += gotCurrent
			if gotTotal != total {
				if isResultChClosed.Load() == 0 {
					resultCh <- "FAIL"
				}
				return
			}

			if total < 21 {
				w.Write([]byte("j\n\n"))
			} else {
				lastOut = out
				break
			}
		}
	}

	// lõpp
	if isResultChClosed.Load() != 0 {
		return
	}

	if total == 21 && !strings.Contains(strings.ToLower(lastOut), "võit") {
		resultCh <- "FAIL"
		return
	} else if total > 21 && !strings.Contains(strings.ToLower(lastOut), "kaot") {
		resultCh <- "FAIL"
		return
	}

	resultCh <- "OK"
}
