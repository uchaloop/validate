/*
Package validate provides a small error accumulator, so a Validate method
reports every problem in one pass instead of failing on the first. It depends on
nothing outside the standard library.

# Why accumulate

A method that returns on its first problem reports one problem. Whoever has to
fix them finds the next one only by fixing this one and running again, which
turns three mistakes into three round trips. Accumulating turns them into one
report:

	func (c Config) Validate() error {
		var errs validate.Errors

		errs.Require(c.Timeout > 0, "timeout must be positive, got %s", c.Timeout)
		errs.Require(len(c.Host) != 0, "host is required")

		return errs.Err()
	}

	timeout must be positive, got -1s
	host is required

That matters most where the round trip is expensive - a value fixed somewhere
else and applied again, rather than a line edited and rerun.

# Recording

The zero value is ready to use, so an Errors is declared and written to; there
is nothing to construct.

Add records an error when it is non-nil, which is what makes the common "call
something, keep what it said" line a single statement. Addf formats a new one,
and %w wraps the cause the way fmt.Errorf does. Require records a formatted
error when a condition does not hold - it states what has to be true, rather
than what must not be, which is how the check reads in the report it produces.

Err joins what was recorded, or returns nil when nothing was, so a Validate ends
in one line whatever it found. The join keeps errors.Is and errors.As reaching
each recorded error, and renders them one per line.

The accumulator does not redact messages or wrapped causes. Callers must avoid
putting secrets into errors or format arguments and sanitize diagnostics at the
application boundary.
*/
package validate

import (
	"errors"
	"fmt"
)

// Errors collects validation errors in insertion order. The zero value is ready
// to use. Do not copy an Errors after first use. Calls that access the same
// accumulator must be synchronized when any caller may add errors.
type Errors struct {
	errs []error
}

// Add records err when the error interface is non-nil. An interface containing
// a typed nil pointer is recorded; its Error method must support that receiver
// if the joined result is formatted. No reflection-based nil filtering is done.
func (e *Errors) Add(err error) {
	if err != nil {
		e.errs = append(e.errs, err)
	}
}

// Addf records a new formatted error using fmt.Errorf, including %w wrapping.
func (e *Errors) Addf(format string, args ...any) {
	e.errs = append(e.errs, fmt.Errorf(format, args...))
}

// Require records a formatted error when isValid is false. It is the convenient
// path for "field X is required" style checks. Arguments are evaluated before
// the call even when isValid is true; guard dependent or expensive checks with if.
// Require is a compatibility exception to the formatting-function f suffix.
func (e *Errors) Require(isValid bool, format string, args ...any) {
	if !isValid {
		e.Addf(format, args...)
	}
}

// Err returns all accumulated errors joined into one, or nil when there are
// none. It does not clear the accumulator. Later additions do not change earlier
// results, but the underlying error values are not deep-copied. The result works
// with errors.Is/As over each recorded error.
func (e *Errors) Err() error {
	return errors.Join(e.errs...)
}
