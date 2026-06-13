// genpw is a small, dependency-free password generator.
//
// Randomness comes exclusively from crypto/rand and index selection is
// unbiased. By default all four character classes are enabled; disable any
// with the -no-* flags.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "genpw:", err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("genpw", flag.ContinueOnError)
	fs.Usage = func() { usage(fs) }

	var (
		length      int
		count       int
		noLower     bool
		noUpper     bool
		noDigits    bool
		noSymbols   bool
		symbolSet   string
		exclude     string
		noAmbiguous bool
		minDigits   int
		minSymbols  int
		showEntropy bool
		copyClip    bool
	)
	// Both short and long forms bind to the same variable.
	fs.IntVar(&length, "length", 20, "password length")
	fs.IntVar(&length, "l", 20, "password length (shorthand)")
	fs.IntVar(&count, "count", 1, "number of passwords to generate")
	fs.IntVar(&count, "n", 1, "number of passwords (shorthand)")
	fs.BoolVar(&noLower, "no-lower", false, "exclude lowercase letters")
	fs.BoolVar(&noUpper, "no-upper", false, "exclude uppercase letters")
	fs.BoolVar(&noDigits, "no-digits", false, "exclude digits")
	fs.BoolVar(&noSymbols, "no-symbols", false, "exclude symbols")
	fs.StringVar(&symbolSet, "symbols", "", "custom symbol set (overrides default)")
	fs.StringVar(&exclude, "exclude", "", "characters to exclude from the pool")
	fs.BoolVar(&noAmbiguous, "no-ambiguous", false, "exclude confusable chars (il1LoO0)")
	fs.IntVar(&minDigits, "min-digits", 0, "minimum number of digits")
	fs.IntVar(&minSymbols, "min-symbols", 0, "minimum number of symbols")
	fs.BoolVar(&showEntropy, "entropy", false, "print entropy (bits) and exit")
	fs.BoolVar(&copyClip, "copy", false, "copy to clipboard, do not print")

	if err := fs.Parse(args); err != nil {
		// -h/-help: usage was already printed by Parse; exit cleanly.
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	p := Policy{
		Length:      length,
		Lower:       !noLower,
		Upper:       !noUpper,
		Digits:      !noDigits,
		Symbols:     !noSymbols,
		SymbolSet:   symbolSet,
		Exclude:     exclude,
		NoAmbiguous: noAmbiguous,
		MinDigits:   minDigits,
		MinSymbols:  minSymbols,
	}
	if err := p.Validate(); err != nil {
		return err
	}

	if showEntropy {
		_, err := fmt.Fprintf(out, "%.1f bits (%s)\n", p.EntropyBits(), p.describe())
		return err
	}

	pwds := make([]string, 0, count)
	for i := 0; i < count; i++ {
		s, err := p.Generate()
		if err != nil {
			return err
		}
		pwds = append(pwds, s)
	}

	if copyClip {
		if count != 1 {
			return fmt.Errorf("-copy requires -count 1")
		}
		if err := clipboardCopy(pwds[0]); err != nil {
			return err
		}
		_, err := fmt.Fprintf(out, "copied to clipboard (%.1f bits)\n", p.EntropyBits())
		return err
	}

	_, err := fmt.Fprintln(out, strings.Join(pwds, "\n"))
	return err
}

// clipboardCopy shells out to the platform clipboard tool (zero Go deps).
func clipboardCopy(s string) error {
	var name string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		name = "pbcopy"
	case "linux":
		if _, err := exec.LookPath("wl-copy"); err == nil {
			name = "wl-copy"
		} else {
			name, args = "xclip", []string{"-selection", "clipboard"}
		}
	default:
		return fmt.Errorf("-copy not supported on %s", runtime.GOOS)
	}
	cmd := exec.Command(name, args...)
	cmd.Stdin = strings.NewReader(s)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clipboard via %s: %w", name, err)
	}
	return nil
}

// usage is hand-written rather than using fs.PrintDefaults(): Go's flag package
// has no real long-option concept (-x and --x are equivalent, help always
// renders single-dash), and PrintDefaults lists each shorthand on its own line
// (e.g. -l and -length separately). Keep this list in sync with run().
func usage(fs *flag.FlagSet) {
	_, _ = fmt.Fprint(fs.Output(), `genpw - secure password generator

Usage:
  genpw [flags]

Examples:
  genpw                       20 chars, all classes
  genpw -l 32                 32 chars
  genpw -n 5                  5 candidates
  genpw -no-symbols           alphanumeric only
  genpw -symbols '!@#$%'      custom symbol set
  genpw -no-ambiguous         drop il1LoO0
  genpw -min-digits 2 -min-symbols 1
  genpw -copy                 copy instead of print
  genpw -entropy              show strength only

Flags:
  -l, -length int    password length (default 20)
  -n, -count int     number of passwords (default 1)
  -no-lower          exclude lowercase letters
  -no-upper          exclude uppercase letters
  -no-digits         exclude digits
  -no-symbols        exclude symbols
  -symbols string    custom symbol set (overrides default)
  -exclude string    characters to exclude from the pool
  -no-ambiguous      exclude confusable chars (il1LoO0)
  -min-digits int    minimum number of digits
  -min-symbols int   minimum number of symbols
  -copy              copy to clipboard, do not print
  -entropy           print entropy (bits) and exit
`)
}
