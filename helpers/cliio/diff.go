package cliio

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/hexops/gotextdiff"
)

func PrintDiff(u gotextdiff.Unified) {
	if len(u.Hunks) == 0 {
		return
	}
	for _, hunk := range u.Hunks {
		fromCount, toCount := 0, 0
		for _, l := range hunk.Lines {
			switch l.Kind {
			case gotextdiff.Delete:
				fromCount++
			case gotextdiff.Insert:
				toCount++
			default:
				fromCount++
				toCount++
			}
		}
		fmt.Print("@@")
		if fromCount > 1 {
			fmt.Printf(" -%d,%d", hunk.FromLine, fromCount)
		} else {
			fmt.Printf(" -%d", hunk.FromLine)
		}
		if toCount > 1 {
			fmt.Printf(" +%d,%d", hunk.ToLine, toCount)
		} else {
			fmt.Printf(" +%d", hunk.ToLine)
		}
		fmt.Print(" @@\n")
		for _, l := range hunk.Lines {
			switch l.Kind {
			case gotextdiff.Delete:
				printRed("-%s", l.Content)
			case gotextdiff.Insert:
				printGreen("+%s", l.Content)
			default:
				fmt.Printf(" %s", l.Content)
			}
			if !strings.HasSuffix(l.Content, "\n") {
				fmt.Printf("\n\\ No newline at end of file\n")
			}
		}
	}
}

func printRed(format string, a ...interface{}) {
	color.New(color.FgRed).Printf(format, a...)
}

func printGreen(format string, a ...interface{}) {
	color.New(color.FgGreen).Printf(format, a...)
}
