package main

import (
	_ "embed"
	"fmt"
	"math"
	"strings"
	"unicode"
)

// effWordlistRaw is the EFF "large" wordlist (7776 words, Diceware-compatible).
// Source: https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt
//
//go:embed eff_large_wordlist.txt
var effWordlistRaw string

// effWords holds the parsed words, indexed for uniform selection.
var effWords = parseWordlist(effWordlistRaw)

// parseWordlist extracts the word from each "<dice>\t<word>" line.
func parseWordlist(raw string) []string {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	words := make([]string, 0, len(lines))
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		words = append(words, fields[len(fields)-1])
	}
	return words
}

// Passphrase describes a Diceware-style passphrase request.
type Passphrase struct {
	Words      int
	Separator  string
	Capitalize bool
}

func (pp Passphrase) Validate() error {
	if pp.Words < 1 {
		return fmt.Errorf("words must be >= 1")
	}
	if pp.Words > maxLength {
		return fmt.Errorf("words must be <= %d", maxLength)
	}
	if len(effWords) == 0 {
		return fmt.Errorf("wordlist is empty")
	}
	return nil
}

// Generate returns one passphrase. Each word is chosen independently and
// uniformly from the wordlist using crypto/rand (no modulo bias).
func (pp Passphrase) Generate() (string, error) {
	if err := pp.Validate(); err != nil {
		return "", err
	}
	parts := make([]string, pp.Words)
	for i := range parts {
		idx, err := secureIntn(len(effWords))
		if err != nil {
			return "", err
		}
		w := effWords[idx]
		if pp.Capitalize {
			w = capitalizeFirst(w)
		}
		parts[i] = w
	}
	return strings.Join(parts, pp.Separator), nil
}

// EntropyBits is exact for passphrases: log2(wordlist) * words.
func (pp Passphrase) EntropyBits() float64 {
	if len(effWords) == 0 {
		return 0
	}
	return math.Log2(float64(len(effWords))) * float64(pp.Words)
}

func (pp Passphrase) describe() string {
	return fmt.Sprintf("words=%d, wordlist=%d", pp.Words, len(effWords))
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
