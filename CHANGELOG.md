# Changelog

All notable changes to this module are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this module adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-09-01

Previously distributed as `github.com/uchaloop/confmaker/validate`. Update the
import path; the API is unchanged.

### Added

- `Errors`, an accumulator whose zero value is ready to use: `Add` records a
  non-nil error, `Addf` formats one, `Require` records one when a condition does
  not hold, and `Err` joins them or returns nil. `errors.Is` and `errors.As`
  reach each recorded error.
- Runnable examples, compiled by `go test`.
