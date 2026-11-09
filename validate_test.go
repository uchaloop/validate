package validate

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
)

func TestErrorsAggregate(t *testing.T) {
	var e Errors

	e.Add(nil) // ignored
	e.Add(errors.New("first"))
	e.Addf("second %d", 2)
	e.Require(false, "third %s", "x")
	e.Require(true, "not added")

	err := e.Err()
	if err == nil {
		t.Fatal("Err() = nil, want aggregated error")
	}

	msg := err.Error()
	for _, want := range []string{"first", "second 2", "third x"} {
		if !strings.Contains(msg, want) {
			t.Errorf("aggregated error %q does not contain %q", msg, want)
		}
	}
	if strings.Contains(msg, "not added") {
		t.Errorf("satisfied Require must not add an error: %q", msg)
	}
}

func TestEmptyIsNil(t *testing.T) {
	var e Errors
	if err := e.Err(); err != nil {
		t.Fatalf("empty Err() = %v, want nil", err)
	}
}

func TestErrorsIsCompatible(t *testing.T) {
	sentinel := errors.New("sentinel")

	var e Errors
	e.Add(sentinel)
	e.Addf("other")

	if !errors.Is(e.Err(), sentinel) {
		t.Fatal("errors.Is should find the recorded sentinel")
	}
}

func TestErrorsPreserveCauses(t *testing.T) {
	cause := &strconv.NumError{
		Func: "ParseInt",
		Num:  "invalid",
		Err:  strconv.ErrSyntax,
	}

	var errs Errors
	errs.Addf("port: %w", cause)
	errs.Require(false, "nested: %w", io.EOF)

	result := errs.Err()
	var got *strconv.NumError
	if !errors.As(result, &got) || got != cause {
		t.Fatal("typed cause was lost")
	}

	if !errors.Is(result, strconv.ErrSyntax) || !errors.Is(result, io.EOF) {
		t.Fatal("wrapped causes were lost")
	}
}

func TestErrorsSnapshots(t *testing.T) {
	var errs Errors
	errs.Addf("first")
	errs.Addf("second")

	first := errs.Err()
	if first.Error() != "first\nsecond" || errs.Err().Error() != first.Error() {
		t.Fatal("order changed or Err consumed errors")
	}

	errs.Addf("third")
	if first.Error() != "first\nsecond" {
		t.Fatal("previous result changed after Add")
	}

	if errs.Err().Error() != "first\nsecond\nthird" {
		t.Fatal("new error missing")
	}
}

func TestAddNilLeavesEmpty(t *testing.T) {
	var errs Errors
	errs.Add(nil)
	if err := errs.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}
}

type nilReceiverError struct{}

func (*nilReceiverError) Error() string { return "typed nil error" }

func TestAddPreservesTypedNil(t *testing.T) {
	var cause *nilReceiverError
	var errs Errors
	errs.Add(cause)
	result := errs.Err()
	if result == nil {
		t.Fatal("typed nil error was discarded")
	}
	var got *nilReceiverError
	if !errors.As(result, &got) || got != nil {
		t.Fatal("typed nil cause was not preserved")
	}
	if result.Error() != "typed nil error" {
		t.Fatalf("unexpected message: %q", result.Error())
	}
}
