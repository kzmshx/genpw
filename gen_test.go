package main

import (
	"strings"
	"testing"
)

func defaultPolicy() Policy {
	return Policy{Length: 20, Lower: true, Upper: true, Digits: true, Symbols: true}
}

func TestGenerateLength(t *testing.T) {
	p := defaultPolicy()
	for _, n := range []int{1, 8, 20, 64, 128} {
		p.Length = n
		s, err := p.Generate()
		if err != nil {
			t.Fatalf("len %d: %v", n, err)
		}
		if len(s) != n {
			t.Errorf("len %d: got %d", n, len(s))
		}
	}
}

func TestGenerateUnique(t *testing.T) {
	p := defaultPolicy()
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		s, err := p.Generate()
		if err != nil {
			t.Fatal(err)
		}
		if seen[s] {
			t.Fatalf("duplicate password generated: %q", s)
		}
		seen[s] = true
	}
}

func TestClassExclusion(t *testing.T) {
	p := defaultPolicy()
	p.Symbols = false
	p.Upper = false
	for i := 0; i < 200; i++ {
		s, err := p.Generate()
		if err != nil {
			t.Fatal(err)
		}
		if strings.ContainsAny(s, defaultSyms) {
			t.Fatalf("symbol leaked: %q", s)
		}
		if strings.ContainsAny(s, upperChars) {
			t.Fatalf("uppercase leaked: %q", s)
		}
	}
}

func TestExcludeAndAmbiguous(t *testing.T) {
	p := defaultPolicy()
	p.Exclude = "abc"
	p.NoAmbiguous = true
	for i := 0; i < 200; i++ {
		s, err := p.Generate()
		if err != nil {
			t.Fatal(err)
		}
		if strings.ContainsAny(s, "abc") {
			t.Fatalf("excluded char leaked: %q", s)
		}
		if strings.ContainsAny(s, ambiguousSet) {
			t.Fatalf("ambiguous char leaked: %q", s)
		}
	}
}

func TestMinimums(t *testing.T) {
	p := defaultPolicy()
	p.Length = 12
	p.MinDigits = 3
	p.MinSymbols = 2
	for i := 0; i < 500; i++ {
		s, err := p.Generate()
		if err != nil {
			t.Fatal(err)
		}
		if len(s) != 12 {
			t.Fatalf("length: %q", s)
		}
		if c := countAny(s, digitChars); c < 3 {
			t.Fatalf("min digits: got %d in %q", c, s)
		}
		if c := countAny(s, defaultSyms); c < 2 {
			t.Fatalf("min symbols: got %d in %q", c, s)
		}
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name string
		p    Policy
		ok   bool
	}{
		{"empty pool", Policy{Length: 10}, false},
		{"zero length", Policy{Length: 0, Lower: true}, false},
		{"over max length", Policy{Length: maxLength + 1, Lower: true}, false},
		{"min digits without digits", Policy{Length: 10, Lower: true, MinDigits: 1}, false},
		{"min digits but digits excluded", Policy{Length: 10, Digits: true, Exclude: digitChars, MinDigits: 1}, false},
		{"mins exceed length", Policy{Length: 2, Digits: true, Symbols: true, MinDigits: 2, MinSymbols: 2}, false},
		{"valid", defaultPolicy(), true},
		{"max length ok", Policy{Length: maxLength, Lower: true}, true},
	}
	for _, c := range cases {
		err := c.p.Validate()
		if (err == nil) != c.ok {
			t.Errorf("%s: ok=%v err=%v", c.name, c.ok, err)
		}
	}
}

func TestNoModuloBiasDistribution(t *testing.T) {
	// Single-class pool of 10 digits; over many draws every digit should
	// appear with roughly equal frequency. A modulo-biased generator would
	// skew low digits. Tolerance is generous to avoid flakiness.
	p := Policy{Length: 1, Digits: true}
	counts := make(map[rune]int)
	const N = 20000
	for i := 0; i < N; i++ {
		s, err := p.Generate()
		if err != nil {
			t.Fatal(err)
		}
		counts[rune(s[0])]++
	}
	expected := N / 10
	for d := '0'; d <= '9'; d++ {
		got := counts[d]
		if got < expected*8/10 || got > expected*12/10 {
			t.Errorf("digit %c: got %d, expected ~%d (possible bias)", d, got, expected)
		}
	}
}

func countAny(s, set string) int {
	n := 0
	for _, r := range s {
		if strings.ContainsRune(set, r) {
			n++
		}
	}
	return n
}
