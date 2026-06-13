package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strings"
)

const (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars   = "0123456789"
	defaultSyms  = "!@#$%^&*()-_=+[]{};:,.?"
	ambiguousSet = "il1LoO0" // visually confusable characters

	maxLength   = 4096  // upper bound on requested length (self-DoS guard)
	maxAttempts = 10000 // rejection-sampling retry cap for --min-* constraints
)

// Policy describes which characters are allowed and any minimum requirements.
type Policy struct {
	Length      int
	Lower       bool
	Upper       bool
	Digits      bool
	Symbols     bool
	SymbolSet   string // overrides defaultSyms when non-empty
	Exclude     string // characters to remove from the final pool
	NoAmbiguous bool
	MinDigits   int
	MinSymbols  int
}

// classes returns the enabled character classes after applying Exclude and
// NoAmbiguous filtering. Each returned class is non-overlapping by construction.
func (p Policy) classes() (lower, upper, digits, symbols string) {
	syms := p.SymbolSet
	if syms == "" {
		syms = defaultSyms
	}
	drop := p.Exclude
	if p.NoAmbiguous {
		drop += ambiguousSet
	}
	filter := func(s string) string {
		out := make([]rune, 0, len(s))
		for _, r := range s {
			if !containsRune(drop, r) {
				out = append(out, r)
			}
		}
		return string(out)
	}
	if p.Lower {
		lower = filter(lowerChars)
	}
	if p.Upper {
		upper = filter(upperChars)
	}
	if p.Digits {
		digits = filter(digitChars)
	}
	if p.Symbols {
		symbols = filter(syms)
	}
	return
}

func (p Policy) pool() string {
	l, u, d, s := p.classes()
	return l + u + d + s
}

// Validate reports configuration errors before any generation is attempted.
func (p Policy) Validate() error {
	if p.Length < 1 {
		return fmt.Errorf("length must be >= 1")
	}
	if p.Length > maxLength {
		return fmt.Errorf("length must be <= %d", maxLength)
	}
	_, _, digits, symbols := p.classes()
	if len(p.pool()) == 0 {
		return fmt.Errorf("character pool is empty: enable at least one class")
	}
	if p.MinDigits > 0 && !p.Digits {
		return fmt.Errorf("--min-digits requires digits to be enabled")
	}
	if p.MinSymbols > 0 && !p.Symbols {
		return fmt.Errorf("--min-symbols requires symbols to be enabled")
	}
	// A class can be enabled yet filtered empty by --exclude / --no-ambiguous.
	if p.MinDigits > 0 && len(digits) == 0 {
		return fmt.Errorf("--min-digits set but all digits were excluded")
	}
	if p.MinSymbols > 0 && len(symbols) == 0 {
		return fmt.Errorf("--min-symbols set but all symbols were excluded")
	}
	if p.MinDigits+p.MinSymbols > p.Length {
		return fmt.Errorf("minimum counts (%d) exceed length (%d)", p.MinDigits+p.MinSymbols, p.Length)
	}
	return nil
}

// Generate returns one password satisfying the policy. Randomness comes from
// crypto/rand and index selection is unbiased (no modulo bias).
//
// Every character is drawn independently from the full pool, which is already
// a uniform distribution. When --min-* constraints are set, candidates that do
// not satisfy them are rejected and redrawn (rejection sampling), which keeps
// the result uniform over the space of policy-satisfying passwords. With no
// constraints the first candidate always passes (one iteration).
func (p Policy) Generate() (string, error) {
	if err := p.Validate(); err != nil {
		return "", err
	}
	_, _, digits, symbols := p.classes()
	pool := p.pool()

	buf := make([]byte, 0, p.Length)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		buf = buf[:0]
		if err := appendRandom(&buf, pool, p.Length); err != nil {
			zero(buf)
			return "", err
		}
		if countByteAny(buf, digits) >= p.MinDigits && countByteAny(buf, symbols) >= p.MinSymbols {
			out := string(buf)
			zero(buf) // best-effort: clear the working buffer (the string copy remains)
			return out, nil
		}
	}
	zero(buf)
	return "", fmt.Errorf("could not satisfy --min-* constraints in %d attempts; reduce minimums or widen classes", maxAttempts)
}

// EntropyBits is an approximation: log2(pool) * length. Minimum-count
// constraints reduce the true value slightly, so this is an upper bound.
func (p Policy) EntropyBits() float64 {
	n := len(p.pool())
	if n == 0 {
		return 0
	}
	return math.Log2(float64(n)) * float64(p.Length)
}

func appendRandom(dst *[]byte, set string, n int) error {
	for i := 0; i < n; i++ {
		idx, err := secureIntn(len(set))
		if err != nil {
			return err
		}
		*dst = append(*dst, set[idx])
	}
	return nil
}

// countByteAny counts bytes of b that appear in set.
func countByteAny(b []byte, set string) int {
	n := 0
	for _, c := range b {
		if strings.IndexByte(set, c) >= 0 {
			n++
		}
	}
	return n
}

// zero overwrites b with zeros (best-effort secret hygiene).
func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// secureIntn returns a uniform integer in [0, n) using crypto/rand. big.Int
// rejection sampling guarantees no modulo bias.
func secureIntn(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("invalid range: %d", n)
	}
	v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0, err
	}
	return int(v.Int64()), nil
}

func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}
