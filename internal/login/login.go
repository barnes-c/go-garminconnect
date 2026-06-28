// Package login authenticates a Garmin client for the internal command-line
// tools. It reuses a cached token when available and prompts for an MFA code on
// stdin when a fresh SSO login is required.
package login

import (
	"bufio"
	"context"
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
		return nil, err
	}
	return c, nil
}

func promptMFA() (string, error) {
	fmt.Fprint(os.Stderr, "MFA code: ")
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
