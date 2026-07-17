package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aymanbagabas/go-udiff"
	"github.com/mattn/go-isatty"
	"github.com/spf13/pflag"
)

const (
	redSeq   = "\033[31m"
	greenSeq = "\033[32m"
	resetSeq = "\033[0m"
)

var (
	contextLines int
	color        string
)

func init() {
	pflag.Usage = usage
	pflag.IntVarP(&contextLines, "context", "C", udiff.DefaultContextLines, "number of context lines")
	pflag.StringVarP(&color, "color", "", "auto", "colorize the output; can be 'always', 'never', or 'auto'")
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] file1 file2\n", os.Args[0])
	pflag.PrintDefaults()
}

func main() {
	pflag.Parse()
	args := pflag.Args()
	if len(args) != 2 {
		pflag.Usage()
		os.Exit(1)
	}

	var colorize bool
	switch strings.ToLower(color) {
	case "always":
		colorize = true
	case "auto":
		colorize = isatty.IsTerminal(os.Stdout.Fd())
	}

	f1, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't open file: %s\n", err)
		os.Exit(1)
	}

	defer f1.Close()
	f2, err := os.Open(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't open file: %s\n", err)
		os.Exit(1)
	}

	defer f2.Close()
	s1, err := io.ReadAll(f1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't read file: %s\n", err)
		os.Exit(1)
	}

	s2, err := io.ReadAll(f2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't read file: %s\n", err)
		os.Exit(1)
	}

	edits := udiff.Strings(string(s1), string(s2))
	u, err := udiff.ToUnifiedDiff(f1.Name(), f2.Name(), string(s1), edits, contextLines)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't generate diff: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(toString(u, colorize))
}

// String converts a unified diff to the standard textual form for that diff.
// The output of this function can be passed to tools like patch.
func toString(u udiff.UnifiedDiff, colorize bool) string {
	if len(u.Hunks) == 0 {
		return ""
	}
	b := new(strings.Builder)
	fmt.Fprintf(b, "--- %s\n", u.From)
	fmt.Fprintf(b, "+++ %s\n", u.To)
	for _, hunk := range u.Hunks {
		fromCount, toCount := 0, 0
		for _, l := range hunk.Lines {
			switch l.Kind {
			case udiff.Delete:
				fromCount++
			case udiff.Insert:
				toCount++
			default:
				fromCount++
				toCount++
			}
		}
		fmt.Fprint(b, "@@")
		if fromCount > 1 {
			fmt.Fprintf(b, " -%d,%d", hunk.FromLine, fromCount)
		} else if hunk.FromLine == 1 && fromCount == 0 {
			// Match odd GNU diff -u behavior adding to empty file.
			fmt.Fprintf(b, " -0,0")
		} else {
			fmt.Fprintf(b, " -%d", hunk.FromLine)
		}
		if toCount > 1 {
			fmt.Fprintf(b, " +%d,%d", hunk.ToLine, toCount)
		} else if hunk.ToLine == 1 && toCount == 0 {
			// Match odd GNU diff -u behavior adding to empty file.
			fmt.Fprintf(b, " +0,0")
		} else {
			fmt.Fprintf(b, " +%d", hunk.ToLine)
		}
		fmt.Fprint(b, " @@\n")
		for _, l := range hunk.Lines {
			switch l.Kind {
			case udiff.Delete:
				if colorize {
					fmt.Fprint(b, redSeq)
				}
				fmt.Fprintf(b, "-%s", l.Content)
				if colorize {
					fmt.Fprint(b, resetSeq)
				}
			case udiff.Insert:
				if colorize {
					fmt.Fprint(b, greenSeq)
				}
				fmt.Fprintf(b, "+%s", l.Content)
				if colorize {
					fmt.Fprint(b, resetSeq)
				}
			default:
				fmt.Fprintf(b, " %s", l.Content)
			}
			if !strings.HasSuffix(l.Content, "\n") {
				fmt.Fprintf(b, "\n\\ No newline at end of file\n")
			}
		}
	}
	return b.String()
}
