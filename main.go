package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/glebmish/intervals-icu-cli/internal/api"
	"github.com/glebmish/intervals-icu-cli/internal/cliexit"
	"github.com/glebmish/intervals-icu-cli/internal/cmd"
)

// main does not print the error — cobra already does (SilenceUsage:true,
// SilenceErrors:false). Printing here on top would double every error.
func main() {
	// A panic is a bug, not a user/API/validation error: exit code 5 (internal)
	// per design-cli §7. The deferred recover converts it to a clean exit.
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "internal error: %v\n", r)
			os.Exit(5)
		}
	}()
	if err := cmd.Execute(); err != nil {
		os.Exit(exitCode(err))
	}
}

func exitCode(err error) int {
	var authErr *cliexit.AuthError
	if errors.As(err, &authErr) {
		return 2
	}
	var valErr *cliexit.ValidationError
	if errors.As(err, &valErr) {
		return 3
	}
	var discErr *cliexit.DiscoveryError
	if errors.As(err, &discErr) {
		return 4
	}
	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		if apiErr.IsAuth() {
			return 2
		}
		return 1
	}
	return 1
}
