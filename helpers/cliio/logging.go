package cliio

import (
	"fmt"
)

var LogVerboseEnabled bool

func LogVerbose(message string) {
	if LogVerboseEnabled {
		fmt.Println(message)
	}
}

func LogVerbosef(format string, a ...interface{}) {
	LogVerbose(fmt.Sprintf(format, a...))
}
