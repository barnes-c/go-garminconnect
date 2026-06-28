// Command record refreshes the VCR cassettes that back the test suite by
// running each cassette-backed test against the live Garmin Connect API.
//
// It logs in once, discovers every test that calls newVCRClient, then runs them
// one at a time (with a short delay to avoid rate-limiting). Cassettes are
// sanitized inline by the test recorder's BeforeSaveHook, so there is no
// separate scrubbing step.
//
// Usage:
//
//	GARMIN_EMAIL=you@example.com GARMIN_PASSWORD=secret go run ./tools/record
//	go run ./tools/record --missing   # only record cassettes that don't exist yet
//
// Credentials may be omitted if a valid cached token exists in the token file.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
)

const (
	testDir     = "garminconnect/tests"
	cassetteDir = "garminconnect/tests/testdata/cassettes"
	tokenFile   = ".garmin_token.json"
	delay       = 5 * time.Second
)

// Cassettes recorded specially (not via newVCRClient), preserved on full re-record.
var keep = map[string]bool{"TestLogin_FetchesProfile": true}

func main() {
	missingOnly := flag.Bool("missing", false, "only record cassettes that don't exist yet")
	flag.BoolVar(missingOnly, "m", false, "shorthand for --missing")
	flag.Parse()

	if err := run(*missingOnly); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(missingOnly bool) error {
	token, displayName, err := login()
	if err != nil {
		return err
	}
	fmt.Printf("==> Logged in (display_name=%s)\n", displayName)

	tests, err := discoverTests()
	if err != nil {
		return err
	}
	if len(tests) == 0 {
		return fmt.Errorf("no cassette-backed tests found under %s", testDir)
	}

	if missingOnly {
		fmt.Println("==> --missing mode: keeping existing cassettes")
		tests = filterMissing(tests)
		if len(tests) == 0 {
			fmt.Println("    nothing to record")
			return nil
		}
	} else if err := removeCassettes(); err != nil {
		return err
	}

	env := append(os.Environ(),
		"GARMIN_TOKEN="+token,
		"GARMIN_DISPLAY_NAME="+displayName,
	)

	var pass, fail, skip []string
	fmt.Printf("\n==> Recording %d cassettes (%s between each)...\n\n", len(tests), delay)
	for i, test := range tests {
		fmt.Printf("--- [%d/%d] %s\n", i+1, len(tests), test)
		switch runTest(test, env) {
		case resultSkip:
			skip = append(skip, test)
		case resultPass:
			pass = append(pass, test)
		default:
			fail = append(fail, test)
			fmt.Println("    ^^^ FAILED")
		}
		if i < len(tests)-1 {
			time.Sleep(delay)
		}
	}

	fmt.Println("\n=== SUMMARY ===")
	fmt.Printf("PASS (%d): %s\n", len(pass), strings.Join(pass, " "))
	fmt.Printf("FAIL (%d): %s\n", len(fail), strings.Join(fail, " "))
	fmt.Printf("SKIP (%d): %s\n", len(skip), strings.Join(skip, " "))
	if len(fail) > 0 {
		return fmt.Errorf("%d test(s) failed", len(fail))
	}
	return nil
}

// login reuses a cached token if present, otherwise authenticates with
// GARMIN_EMAIL / GARMIN_PASSWORD, and returns the access token and display name.
func login() (token, displayName string, err error) {
	c := gc.NewClient(tokenFile)
	if err := c.Login(context.Background(), os.Getenv("GARMIN_EMAIL"), os.Getenv("GARMIN_PASSWORD")); err != nil {
		return "", "", fmt.Errorf("login: %w", err)
	}
	return c.Token(), c.DisplayName(), nil
}

var (
	funcRE    = regexp.MustCompile(`^func (Test\w+)\(`)
	usesVCRRE = regexp.MustCompile(`newVCRClient\(t\)`)
)

// discoverTests returns every test function that calls newVCRClient, i.e. every
// test that owns a cassette.
func discoverTests() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(testDir, "*_test.go"))
	if err != nil {
		return nil, err
	}
	var tests []string
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		var cur string
		added := false
		for line := range strings.SplitSeq(string(data), "\n") {
			if m := funcRE.FindStringSubmatch(line); m != nil {
				cur, added = m[1], false
			}
			if cur != "" && !added && usesVCRRE.MatchString(line) {
				tests = append(tests, cur)
				added = true
			}
		}
	}
	sort.Strings(tests)
	return tests, nil
}

func filterMissing(tests []string) []string {
	var out []string
	for _, test := range tests {
		if _, err := os.Stat(filepath.Join(cassetteDir, test+".yaml")); os.IsNotExist(err) {
			out = append(out, test)
		}
	}
	return out
}

func removeCassettes() error {
	files, err := filepath.Glob(filepath.Join(cassetteDir, "*.yaml"))
	if err != nil {
		return err
	}
	fmt.Println("==> Removing cassettes (except kept)...")
	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".yaml")
		if keep[name] {
			continue
		}
		if err := os.Remove(file); err != nil {
			return err
		}
	}
	return nil
}

type result int

const (
	resultFail result = iota
	resultPass
	resultSkip
)

func runTest(test string, env []string) result {
	cmd := exec.Command("go", "test", "./"+testDir+"/...", "-run", "^"+test+"$", "-count=1", "-v")
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	text := string(out)
	for line := range strings.SplitSeq(text, "\n") {
		if strings.HasPrefix(line, "=== RUN") || strings.HasPrefix(line, "--- ") || strings.Contains(line, "Error") {
			fmt.Println("   ", strings.TrimSpace(line))
		}
	}
	switch {
	case strings.Contains(text, "--- SKIP"):
		return resultSkip
	case err == nil:
		return resultPass
	default:
		return resultFail
	}
}
