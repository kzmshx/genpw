# genpw

[![CI](https://github.com/kzmshx/genpw/actions/workflows/ci.yml/badge.svg)](https://github.com/kzmshx/genpw/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/kzmshx/genpw.svg)](https://pkg.go.dev/github.com/kzmshx/genpw)
[![Go Report Card](https://goreportcard.com/badge/github.com/kzmshx/genpw)](https://goreportcard.com/report/github.com/kzmshx/genpw)
[![Release](https://img.shields.io/github/v/release/kzmshx/genpw)](https://github.com/kzmshx/genpw/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A small, dependency-free password generator CLI.

- Randomness from `crypto/rand` only (CSPRNG).
- Unbiased selection — no modulo bias (`crypto/rand.Int` rejection sampling).
- `-min-*` constraints enforced via rejection sampling, preserving uniformity.
- Zero external dependencies (Go standard library only).

## Install

```sh
go install github.com/kzmshx/genpw@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/kzmshx/genpw/releases).

## Usage

```sh
genpw                          # 20 chars, all classes (~128 bits)
genpw -l 32                    # length 32
genpw -n 5                     # 5 candidates
genpw -no-symbols              # alphanumeric only
genpw -symbols '!@#$%'         # custom symbol set
genpw -no-ambiguous            # drop il1LoO0
genpw -min-digits 2 -min-symbols 1
genpw -copy                    # copy to clipboard, do not print
genpw -entropy                 # show strength only, no generation
```

Flags use Go's standard single-dash style (`-length`); `-h` lists them all.

## Notes

- `-copy` shells out to `pbcopy` (macOS) or `wl-copy`/`xclip` (Linux); it
  never prints the secret to stdout.
- `-entropy` reports `log2(pool) * length`; `-min-*` constraints make the
  true value marginally lower, so the figure is an upper bound.

## License

[MIT](LICENSE)
