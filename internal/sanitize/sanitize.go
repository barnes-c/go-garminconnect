// Package sanitize scrubs PII from recorded VCR interactions. It runs inline
// during recording (from the test recorder's BeforeSaveHook), so a cassette is
// never written to disk unsanitized.
//
// It is a faithful Go port of the former tools/sanitize_cassettes.py and is
// idempotent: re-running it over already-sanitized content is a no-op. The
// replacements:
//
//   - Integer ID fields (name ends in Id/Pk, 6+ digit value) -> 12345678
//   - UUIDs (hyphenated and bare 32-char hex) -> one all-f constant
//   - Datetime / date-only string values -> 2026-01-01[T00:00:00]
//   - Request-URL dates after 2026-01-01 -> 2026-01-01
//   - Emails -> test@example.com
//   - Every free-text object value in a body -> "TEST" (except dates, UUIDs and
//     the "testuser"/"Test User" placeholders)
//   - Every numeric body value -> 1 (1.0 for floats), except identifier/
//     structural fields (*Id, *Pk, *Count, *Version, ...) unless the value looks
//     like an epoch-ms timestamp
//   - Volatile response headers stripped; duration -> 100ms
package sanitize

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

const (
	synthID    = "12345678"
	synthUUID  = "ffffffff-ffff-ffff-ffff-ffffffffffff"
	synthEmail = "test@example.com"
	synthDate  = "2026-01-01"
)

// Volatile response headers that vary between runs and carry no test value.
var stripHeaders = map[string]bool{
	"Cf-Ray": true, "Date": true, "Nel": true, "Report-To": true,
	"Alt-Svc": true, "Cf-Cache-Status": true, "Cache-Control": true,
	"Pragma": true, "Server": true,
}

var staticReplacements = [][2]string{
	{"garmin-connect-prod", "garmin-connect-test"},
}

var preserveText = map[string]bool{"testuser": true, "Test User": true}

var (
	datetimeRE    = regexp.MustCompile(`\d{4}-\d{2}-\d{2}([T ])\d{2}:\d{2}:\d{2}(?:\.\d+)?`)
	dateOnlyRE    = regexp.MustCompile(`"\d{4}-\d{2}-\d{2}"`)
	bareDateRE    = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	uuidRE        = regexp.MustCompile(`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	emailRE       = regexp.MustCompile(`(?i)[\w.+%-]+@[\w.-]+\.[a-z]{2,}`)
	idFieldRE     = regexp.MustCompile(`("[A-Za-z0-9_]*(?:[Ii]d|[Pp][Kk])":\s*)\d{6,}`)
	digits6RE     = regexp.MustCompile(`\d{6,}`)
	keyNumRE      = regexp.MustCompile(`"([A-Za-z_][A-Za-z0-9_]*)":\s*(-?\d+(?:\.\d+)?)`)
	arrayNumRE    = regexp.MustCompile(`[\[,]-?\d+(?:\.\d+)?`)
	preserveKeyRE = regexp.MustCompile(`(?i)(id|pk|count|index|version|number|order|sequence|priority|category|month|year|offset|zoneid|typekey)$`)
	epochMsRE     = regexp.MustCompile(`^1[5-9]\d{11}$`)
	textValueRE   = regexp.MustCompile(`:"((?:[^"\\]|\\.)*)"`)
	dateishRE     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)
	uuidishRE     = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f-]+$|^[0-9a-f]{32}$`)
	hexRunRE      = regexp.MustCompile(`(?i)[0-9a-f]{32}`)
)

// Interaction scrubs a recorded interaction in place. displayName is the real
// Garmin display name to replace (may be empty).
func Interaction(i *cassette.Interaction, displayName string) {
	stripVolatileHeaders(i.Response.Headers)
	if i.Request.Headers != nil {
		i.Request.Headers.Set("Authorization", "Bearer test")
	}
	i.Response.Duration = 100 * time.Millisecond
	i.Request.URL = URL(i.Request.URL, displayName)
	i.Request.Body = Body(i.Request.Body, displayName)
	i.Response.Body = Body(i.Response.Body, displayName)
}

func stripVolatileHeaders(h map[string][]string) {
	for k := range h {
		if stripHeaders[k] {
			delete(h, k)
		}
	}
}

// URL scrubs a request URL: display name -> testuser, UUIDs -> all-f,
// recording dates after the synthetic date -> the synthetic date, and any 6+
// digit run (an ID) -> the synthetic ID.
func URL(s, displayName string) string {
	if s == "" {
		return s
	}
	if displayName != "" {
		s = strings.ReplaceAll(s, url.PathEscape(displayName), "testuser")
		s = strings.ReplaceAll(s, displayName, "testuser")
	}
	s = scrubUUIDs(s)
	s = bareDateRE.ReplaceAllStringFunc(s, func(d string) string {
		if d > synthDate {
			return synthDate
		}
		return d
	})
	s = digits6RE.ReplaceAllString(s, synthID)
	return s
}

// Body scrubs a JSON request/response body.
func Body(s, displayName string) string {
	if s == "" {
		return s
	}
	s = zeroDatetimes(s)
	s = neutralizeMetrics(s)
	s = idFieldRE.ReplaceAllString(s, "${1}"+synthID)
	s = scrubTextValues(s)
	s = scrubUUIDs(s)
	s = applyStatic(s, displayName)
	return s
}

func zeroDatetimes(s string) string {
	s = datetimeRE.ReplaceAllString(s, synthDate+"${1}00:00:00")
	s = dateOnlyRE.ReplaceAllString(s, `"`+synthDate+`"`)
	return s
}

func placeholder(num string) string {
	if strings.Contains(num, ".") {
		return "1.0"
	}
	return "1"
}

// neutralizeMetrics replaces every numeric body value with 1, except
// identifier/structural fields (unless they hold an epoch-ms timestamp).
func neutralizeMetrics(s string) string {
	s = keyNumRE.ReplaceAllStringFunc(s, func(m string) string {
		sm := keyNumRE.FindStringSubmatch(m)
		key, num := sm[1], sm[2]
		if preserveKeyRE.MatchString(key) && !epochMsRE.MatchString(num) {
			return m
		}
		return `"` + key + `":` + placeholder(num)
	})
	return replaceArrayNumbers(s)
}

// replaceArrayNumbers replaces bare numeric array elements (a number led by '['
// or ',' and followed by ',' or ']') with the placeholder.
func replaceArrayNumbers(s string) string {
	locs := arrayNumRE.FindAllStringIndex(s, -1)
	if locs == nil {
		return s
	}
	var b strings.Builder
	prev := 0
	for _, loc := range locs {
		start, end := loc[0], loc[1]
		if end >= len(s) || (s[end] != ',' && s[end] != ']') {
			continue // not a complete array element
		}
		b.WriteString(s[prev:start])
		b.WriteByte(s[start]) // leading '[' or ','
		b.WriteString(placeholder(s[start+1 : end]))
		prev = end
	}
	b.WriteString(s[prev:])
	return b.String()
}

func scrubTextValues(s string) string {
	return textValueRE.ReplaceAllStringFunc(s, func(m string) string {
		val := textValueRE.FindStringSubmatch(m)[1]
		if val == "" || preserveText[val] || dateishRE.MatchString(val) || uuidishRE.MatchString(val) {
			return m
		}
		return `:"TEST"`
	})
}

func scrubUUIDs(s string) string {
	s = uuidRE.ReplaceAllString(s, synthUUID)
	return replaceBareHex(s)
}

// replaceBareHex replaces a 32-char hex run not adjacent to other hex/hyphen
// characters with a 32-f constant (the hyphen-stripped synthetic UUID).
func replaceBareHex(s string) string {
	repl := strings.ReplaceAll(synthUUID, "-", "")
	locs := hexRunRE.FindAllStringIndex(s, -1)
	if locs == nil {
		return s
	}
	var b strings.Builder
	prev := 0
	for _, loc := range locs {
		start, end := loc[0], loc[1]
		if start > 0 && isHexOrHyphen(s[start-1]) {
			continue
		}
		if end < len(s) && isHexOrHyphen(s[end]) {
			continue
		}
		b.WriteString(s[prev:start])
		b.WriteString(repl)
		prev = end
	}
	b.WriteString(s[prev:])
	return b.String()
}

func isHexOrHyphen(c byte) bool {
	return c == '-' || (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func applyStatic(s, displayName string) string {
	if displayName != "" {
		s = strings.ReplaceAll(s, `"`+displayName+`"`, `"Test User"`)
		s = strings.ReplaceAll(s, displayName, "Test User")
	}
	s = emailRE.ReplaceAllString(s, synthEmail)
	for _, kv := range staticReplacements {
		s = strings.ReplaceAll(s, kv[0], kv[1])
	}
	return s
}
