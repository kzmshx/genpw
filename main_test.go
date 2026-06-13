package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunDefault(t *testing.T) {
	var buf bytes.Buffer
	if err := run(nil, &buf); err != nil {
		t.Fatal(err)
	}
	got := strings.TrimRight(buf.String(), "\n")
	if len(got) != 20 {
		t.Fatalf("default length: got %d (%q)", len(got), got)
	}
}

func TestRunCount(t *testing.T) {
	var buf bytes.Buffer
	if err := run([]string{"-n", "5", "-l", "12"}, &buf); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 5 {
		t.Fatalf("count: got %d lines", len(lines))
	}
	for _, l := range lines {
		if len(l) != 12 {
			t.Fatalf("line length: %q", l)
		}
	}
}

func TestRunEntropy(t *testing.T) {
	var buf bytes.Buffer
	if err := run([]string{"--entropy"}, &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "bits") {
		t.Fatalf("entropy output: %q", buf.String())
	}
}

func TestRunEmptyPoolErrors(t *testing.T) {
	var buf bytes.Buffer
	err := run([]string{"--no-lower", "--no-upper", "--no-digits", "--no-symbols"}, &buf)
	if err == nil {
		t.Fatal("expected error for empty pool")
	}
}

func TestRunCopyRejectsMultiple(t *testing.T) {
	var buf bytes.Buffer
	err := run([]string{"--copy", "-n", "2"}, &buf)
	if err == nil {
		t.Fatal("expected error: --copy with count > 1")
	}
}
