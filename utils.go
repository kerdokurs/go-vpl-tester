package main

import (
	"fmt"
	"strings"
)

func preFormat(str string) {
	for _, line := range strings.Split(str, "\n") {
		fmt.Printf("Comment :=>> >%s\n", line)
	}
}
