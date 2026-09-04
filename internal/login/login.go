// Package login authenticates a Garmin client for the internal command-line
// tools. It reuses a cached token when available and prompts for an MFA code on
// stdin when a fresh SSO login is required.
package login

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
)

// Client logs in using a cached token if available, otherwise the
// GARMIN_EMAIL / GARMIN_PASSWORD environment variables, prompting for an MFA
// code on stdin if the account requires one, and returns the client.
func Client(tokenFile string) (*gc.Client, error) {
	c := gc.NewClient(tokenFile, gc.WithMFAPrompt(promptMFA))
	if err := c.Login(context.Background(), os.Getenv("GARMIN_EMAIL"), os.Getenv("GARMIN_PASSWORD")); err != nil {
		return nil, withBody(err)
	}
	return c, nil
}

// withBody appends the response body Garmin returned, which APIError keeps out
// of its message to stay terse. For the command-line tools it is often the only
// diagnostic on offer — a rate-limited login carries a request-id and nothing
// else.
func withBody(err error) error {
	var apiErr *gc.APIError
	if errors.As(err, &apiErr) && apiErr.Body != "" {
		return fmt.Errorf("%w (body: %s)", err, apiErr.Body)
	}
	return err
}

func promptMFA() (string, error) {
	fmt.Fprint(os.Stderr, "MFA code: ")
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	code := strings.TrimSpace(line)
	// ReadString returns the data alongside io.EOF when stdin ends without a
	// newline (e.g. `echo -n 123456 |`); only fail if no code was read.
	if err != nil && code == "" {
		return "", err
	}
	return code, nil
}
