# validate

[![Go Reference](https://pkg.go.dev/badge/github.com/uchaloop/validate.svg)](https://pkg.go.dev/github.com/uchaloop/validate) [![CI](https://github.com/uchaloop/validate/actions/workflows/ci.yml/badge.svg)](https://github.com/uchaloop/validate/actions/workflows/ci.yml) [![License: MIT](https://img.shields.io/github/license/uchaloop/validate)](LICENSE)

**Collect validation errors and report every problem in one pass.**
The zero value is ready to use. The package depends only on the standard library
and preserves error causes for `errors.Is` and `errors.As`.

[Install](#installation) · [Quick start](#quick-start) · [Contract](#how-it-works) · [Reference](#reference)

## Installation

Requires Go 1.27 or later.

```sh
go get github.com/uchaloop/validate@latest
```

## Quick start

Save this program as `main.go` in an application module after installation:

```go
package main

import (
	"fmt"
	"time"

	"github.com/uchaloop/validate"
)

type Config struct {
	Host    string
	Timeout time.Duration
}

func (c Config) Validate() error {
	var errs validate.Errors
	errs.Require(c.Timeout > 0, "timeout must be positive, got %s", c.Timeout)
	errs.Require(len(c.Host) != 0, "host is required")
	return errs.Err()
}

func main() {
	if err := (Config{Timeout: -time.Second}).Validate(); err != nil {
		fmt.Println(err)
	}
}
```

Run `go run .`. Expected output:

```text
timeout must be positive, got -1s
host is required
```

## How it works

| Method | Behavior |
|---|---|
| `Add(err)` | Record a non-nil error interface |
| `Addf(format, args...)` | Record a formatted error, including `%w` wrapping |
| `Require(isValid, format, args...)` | Record a formatted error when the condition is false |
| `Err()` | Join recorded errors in order, or return nil if empty |

`Err()` does not clear the accumulator. Later additions do not change the list
of errors in an earlier result. The underlying error objects remain shared;
they are not deep-copied. Do not copy an Errors after first use. Synchronize
access to an accumulator when any caller may add errors; there is no internal lock.

`Require` retains its name for public compatibility, as an explicit exception
to the formatting-function `f` suffix convention.

## Errors and safe checks

Arguments are evaluated before Require is called, even when the condition is
true. Guard dependent or expensive expressions with an ordinary `if`. Checking
`value != nil` inside Require does not make an argument such as `value.Field`
safe to evaluate.

`Add(nil)` is ignored. An error interface containing a typed nil pointer is
not nil and is recorded. Its Error method must support a nil receiver if the
joined result will be formatted. The package does not use reflection to discard
such values or recover panics in user error methods.

`errors.Is` and `errors.As` can inspect recorded causes, including causes wrapped
with `%w`. The accumulator does not sanitize messages or wrapped errors. Keep
passwords, tokens and other secrets out of format arguments and sanitize external
errors before exposing diagnostics.

## Reference

- [API reference](https://pkg.go.dev/github.com/uchaloop/validate)
- [Runnable examples](example_test.go)
- [CHANGELOG](CHANGELOG.md)

## License

[MIT](LICENSE)
