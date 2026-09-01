# validate

[![CI](https://github.com/uchaloop/validate/actions/workflows/ci.yml/badge.svg)](https://github.com/uchaloop/validate/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/uchaloop/validate.svg)](https://pkg.go.dev/github.com/uchaloop/validate)
[![License: MIT](https://img.shields.io/github/license/uchaloop/validate)](LICENSE)

A small error accumulator, so a `Validate` method reports every problem in one
pass instead of failing on the first.

- **One report instead of three round trips** - fixing one problem should not be
  how you discover the next.
- **The zero value is ready to use** - nothing to construct.
- **No dependencies** outside the standard library.

```bash
go get github.com/uchaloop/validate
```

## Use

```go
func (c Config) Validate() error {
	var errs validate.Errors

	errs.Require(c.Timeout > 0, "timeout must be positive, got %s", c.Timeout)
	errs.Require(len(c.Host) != 0, "host is required")

	return errs.Err()
}
```

```text
timeout must be positive, got -1s
host is required
```

`Add` records an error when it is non-nil, `Addf` formats a new one, and
`Require` records one when a condition does not hold. `Err` joins what was
recorded, or returns nil when nothing was, and `errors.Is` and `errors.As` still
reach each recorded error.

## Documentation

What each method records and the reasons behind the design are in the package
documentation:
**[pkg.go.dev/github.com/uchaloop/validate](https://pkg.go.dev/github.com/uchaloop/validate)**.

## License

[MIT](LICENSE)
