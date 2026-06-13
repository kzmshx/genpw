# genpw

A small, dependency-free password generator CLI.

- Randomness from `crypto/rand` only (CSPRNG).
- Unbiased index selection — no modulo bias.
- Zero external dependencies (Go standard library only).

## Install

```sh
go install github.com/kzmshx/genpw@latest
```

## Usage

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

Run `genpw -h` for the full flag list.

## Notes

- `--copy` shells out to `pbcopy` (macOS) or `wl-copy`/`xclip` (Linux); it
  never prints the password to stdout.
- `--entropy` reports `log2(pool) * length`; `--min-*` constraints make the
  true value marginally lower, so the figure is an upper bound.
