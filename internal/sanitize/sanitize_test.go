package sanitize_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"

	"github.com/barnes-c/go-garminconnect/internal/sanitize"
)

func TestBody_ScrubsPII(t *testing.T) {
	in := `{"userProfileId":127516254,"deviceId":3493638919,"displayName":"e15d1c9d-827d-4db9-b030-48792c1d1fa6",` +
		`"fullName":"Jane Doe","email":"jane@real.com","contacts":["bob@real.com"],"weight":72.5,"heartRate":61,` +
		`"calendarDate":"2025-06-17","timestamp":"2025-06-17T13:50:13.944","samples":[12.3,45.6],` +
		`"version":1700000000000,"typeId":3}`
	out := sanitize.Body(in, "e15d1c9d-827d-4db9-b030-48792c1d1fa6")

	for _, leak := range []string{"127516254", "3493638919", "Jane Doe", "jane@real.com", "bob@real.com", "72.5", "61", "2025-06-17", "1700000000000"} {
		if strings.Contains(out, leak) {
			t.Errorf("leaked %q in: %s", leak, out)
		}
	}
	for _, want := range []string{`"userProfileId":12345678`, `"deviceId":12345678`, `"fullName":"TEST"`,
		`"email":"TEST"`, `"contacts":["TEST"]`, `"weight":1.0`, `"heartRate":1`,
		`"calendarDate":"2026-01-01"`, `"samples":[1.0,1.0]`, `"version":1`, `"typeId":3`,
		"ffffffff-ffff-ffff-ffff-ffffffffffff"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in: %s", want, out)
		}
	}
}

// TestBody_LeakVectors covers the three leaks found by a full re-record that a
// naive field/value scrub misses: IDs under non-Id field names, IDs used as
// JSON object keys, and numbers in scientific notation (epoch timestamps).
func TestBody_LeakVectors(t *testing.T) {
	in := `{"userProfileNumber":127516254,"loadMap":{"3493638919":{"x":1.0}},` +
		`"metrics":[1.782637918E12,2.5,1.0E13]}`
	out := sanitize.Body(in, "")

	for _, leak := range []string{"127516254", "3493638919", "1.782637918E12", "1.0E13", "2.5"} {
		if strings.Contains(out, leak) {
			t.Errorf("leaked %q in: %s", leak, out)
		}
	}
	for _, want := range []string{`"userProfileNumber":12345678`, `"12345678":{`, `"metrics":[1.0,1.0,1.0]`} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in: %s", want, out)
		}
	}
}

// TestBody_ArrayStrings covers string elements inside JSON arrays, which the
// object-value scrub (`:"..."`) does not reach. A real re-record leaked the
// recording account's userRoles OAuth scope list through this gap.
func TestBody_ArrayStrings(t *testing.T) {
	in := `{"userRoles":["SCOPE_CONNECT_READ","SCOPE_DI_OAUTH_2_TOKEN_ADMIN"],` +
		`"dates":["2026-01-01","2026-01-01T00:00:00"],` +
		`"ids":["ffffffff-ffff-ffff-ffff-ffffffffffff"],"names":["testuser","Test User"],` +
		`"mixed":["free text",1,"2026-01-01"],"nested":[{"k":"v"},["deep","2026-01-01"]]}`
	out := sanitize.Body(in, "")

	for _, leak := range []string{"SCOPE_CONNECT_READ", "SCOPE_DI_OAUTH_2_TOKEN_ADMIN", "free text", `"deep"`, `"k":"v"`} {
		if strings.Contains(out, leak) {
			t.Errorf("leaked %q in: %s", leak, out)
		}
	}
	for _, want := range []string{`"userRoles":["TEST","TEST"]`,
		`"dates":["2026-01-01","2026-01-01T00:00:00"]`,
		`"ids":["ffffffff-ffff-ffff-ffff-ffffffffffff"]`,
		`"names":["testuser","Test User"]`,
		`"mixed":["TEST",1,"2026-01-01"]`,
		`"nested":[{"k":"TEST"},["TEST","2026-01-01"]]`} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in: %s", want, out)
		}
	}
}

func TestURL_ScrubsIDsDatesAndName(t *testing.T) {
	name := "e15d1c9d-827d-4db9-b030-48792c1d1fa6"
	// A real recording date (after testDate) must be rewritten back.
	in := "https://connectapi.garmin.com/web-gateway/solar/3493638919/2026-09-30/2026-09-30?singleDayView=true"
	out := sanitize.URL(in, name)
	if strings.Contains(out, "3493638919") {
		t.Errorf("leaked device id: %s", out)
	}
	if strings.Contains(out, "2026-09-30") {
		t.Errorf("leaked recording date: %s", out)
	}
	if !strings.Contains(out, "/solar/12345678/2026-01-01/2026-01-01") {
		t.Errorf("unexpected: %s", out)
	}

	// A profile URL embedding the display name -> testuser.
	got := sanitize.URL("https://connectapi.garmin.com/usersummary-service/usersummary/daily/"+name, name)
	if !strings.HasSuffix(got, "/daily/testuser") {
		t.Errorf("display name not scrubbed: %s", got)
	}
}

func TestURL_KeepsPastDates(t *testing.T) {
	// Synthetic test dates only ever look back from 2026-01-01 and must survive.
	in := "https://connectapi.garmin.com/x/2025-12-01/2026-01-01"
	if got := sanitize.URL(in, ""); got != in {
		t.Errorf("past dates altered: %s", got)
	}
}

// TestIdempotentOverCassettes runs the sanitizer over every committed cassette
// (all already sanitized) and asserts it makes no further change. If the Go
// logic diverges from how the cassettes were originally scrubbed, this fails.
func TestIdempotentOverCassettes(t *testing.T) {
	dir := filepath.Join("..", "..", "garminconnect", "tests", "testdata", "cassettes")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read cassettes dir: %v", err)
	}
	var checked int
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		name := strings.TrimSuffix(filepath.Join(dir, e.Name()), ".yaml")
		c, err := cassette.Load(name)
		if err != nil {
			t.Fatalf("load %s: %v", e.Name(), err)
		}
		for idx, it := range c.Interactions {
			if got := sanitize.Body(it.Response.Body, ""); got != it.Response.Body {
				t.Errorf("%s[%d] response body not a fixed point", e.Name(), idx)
			}
			if got := sanitize.Body(it.Request.Body, ""); got != it.Request.Body {
				t.Errorf("%s[%d] request body not a fixed point", e.Name(), idx)
			}
			if got := sanitize.URL(it.Request.URL, ""); got != it.Request.URL {
				t.Errorf("%s[%d] url not a fixed point: %q -> %q", e.Name(), idx, it.Request.URL, got)
			}
		}
		checked++
	}
	if checked == 0 {
		t.Fatal("no cassettes checked")
	}
	t.Logf("idempotent over %d cassettes", checked)
}
