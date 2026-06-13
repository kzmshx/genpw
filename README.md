# genpw

[![CI](https://github.com/kzmshx/genpw/actions/workflows/ci.yml/badge.svg)](https://github.com/kzmshx/genpw/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/kzmshx/genpw.svg)](https://pkg.go.dev/github.com/kzmshx/genpw)
[![Go Report Card](https://goreportcard.com/badge/github.com/kzmshx/genpw)](https://goreportcard.com/report/github.com/kzmshx/genpw)
[![Release](https://img.shields.io/github/v/release/kzmshx/genpw)](https://github.com/kzmshx/genpw/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A small, dependency-free password and passphrase generator CLI.

- Randomness from `crypto/rand` only (CSPRNG).
- Unbiased selection — no modulo bias (`crypto/rand.Int` rejection sampling).
- `--min-*` constraints enforced via rejection sampling, preserving uniformity.
- Diceware passphrases from the embedded EFF large wordlist (7776 words).
- Zero external dependencies (Go standard library only).

## Install

```sh
go install github.com/kzmshx/genpw@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/kzmshx/genpw/releases).

## Usage

### Passwords

```sh
genpw                          # 20 chars, all classes (~128 bits)
genpw -l 32                    # length 32
genpw -n 5                     # 5 candidates
genpw --no-symbols             # alphanumeric only
genpw --symbols '!@#$%'        # custom symbol set
genpw --no-ambiguous           # drop il1LoO0
genpw --min-digits 2 --min-symbols 1
genpw --copy                   # copy to clipboard, do not print
genpw --entropy                # show strength only, no generation
```

### Passphrases

```sh
genpw -p                       # 6-word Diceware passphrase (~77 bits)
genpw -p -w 8                  # 8 words (~103 bits)
genpw -p --capitalize          # capitalize each word
genpw -p --sep .               # custom separator
genpw -p --entropy             # show strength only
```

Run `genpw -h` for the full flag list.

## Notes

- `--copy` shells out to `pbcopy` (macOS) or `wl-copy`/`xclip` (Linux); it
  never prints the secret to stdout.
- For passwords, `--entropy` reports `log2(pool) * length`; `--min-*`
  constraints make the true value marginally lower, so the figure is an upper
  bound. For passphrases the figure is exact: `log2(7776) * words`.
- The wordlist is the [EFF large wordlist](https://www.eff.org/dice), embedded
  at build time via `go:embed`.

## License

[MIT](LICENSE)
