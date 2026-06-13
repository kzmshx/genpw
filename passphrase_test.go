package main

import (
	"math"
	"strings"
	"testing"
)

func TestWordlistLoaded(t *testing.T) {
	if len(effWords) != 7776 {
		t.Fatalf("wordlist size: got %d, want 7776", len(effWords))
	}
	for i, w := range effWords {
		if w == "" || strings.ContainsAny(w, " \t") {
			t.Fatalf("malformed word at %d: %q", i, w)
		}
	}
}

func TestPassphraseWordCount(t *testing.T) {
	// Use "/" as the separator: wordlist words contain only [a-z-], so "/"
	// cannot appear inside a word and split is unambiguous (unlike "-", which
	// occurs in words such as "yo-yo").
	for _, n := range []int{1, 4, 6, 10} {
		pp := Passphrase{Words: n, Separator: "/"}
		s, err := pp.Generate()
		if err != nil {
			t.Fatal(err)
		}
		if got := len(strings.Split(s, "/")); got != n {
			t.Errorf("words %d: got %d parts (%q)", n, got, s)
		}
	}
}

func TestPassphraseCapitalize(t *testing.T) {
	// "/" separator: see TestPassphraseWordCount for why "-" is unsafe here.
	pp := Passphrase{Words: 6, Separator: "/", Capitalize: true}
	s, err := pp.Generate()
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range strings.Split(s, "/") {
		if w == "" || w[0] < 'A' || w[0] > 'Z' {
			t.Fatalf("word not capitalized: %q in %q", w, s)
		}
	}
}

func TestPassphraseSeparator(t *testing.T) {
	pp := Passphrase{Words: 4, Separator: "."}
	s, err := pp.Generate()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(s, ".") != 3 {
		t.Fatalf("separator count: %q", s)
	}
}

func TestPassphraseEntropy(t *testing.T) {
	pp := Passphrase{Words: 6}
	// log2(7776) * 6 ~= 77.5 bits
	want := math.Log2(7776) * 6
	if got := pp.EntropyBits(); math.Abs(got-want) > 0.01 {
		t.Fatalf("entropy: got %.2f, want %.2f", got, want)
	}
}

func TestPassphraseValidate(t *testing.T) {
	if err := (Passphrase{Words: 0}).Validate(); err == nil {
		t.Error("expected error for 0 words")
	}
	if err := (Passphrase{Words: maxLength + 1}).Validate(); err == nil {
		t.Error("expected error for too many words")
	}
}

func TestPassphraseUnique(t *testing.T) {
	pp := Passphrase{Words: 6, Separator: "-"}
	seen := make(map[string]bool)
	for i := 0; i < 500; i++ {
		s, err := pp.Generate()
		if err != nil {
			t.Fatal(err)
		}
		if seen[s] {
			t.Fatalf("duplicate passphrase: %q", s)
		}
		seen[s] = true
	}
}
