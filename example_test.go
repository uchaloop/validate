package validate_test

import (
	"fmt"
	"time"

	"github.com/uchaloop/validate"
)

type Config struct {
	Host    string
	Timeout time.Duration
}

// Validate reports every problem in one pass. A config is fixed in a deployment
// and rolled out, so one report is one round trip.
func (c Config) Validate() error {
	var errs validate.Errors

	errs.Require(c.Timeout > 0, "timeout must be positive, got %s", c.Timeout)
	errs.Require(len(c.Host) != 0, "host is required")

	return errs.Err()
}

func Example() {
	err := Config{Timeout: -time.Second}.Validate()

	fmt.Println(err)
	// Output:
	// timeout must be positive, got -1s
	// host is required
}

func ExampleErrors_Err_valid() {
	err := Config{Host: "store:9000", Timeout: 30 * time.Second}.Validate()

	fmt.Println(err)
	// Output:
	// <nil>
}
