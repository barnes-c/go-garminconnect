# CLAUDE.md

Go client for Garmin Connect's private API (reverse-engineered — no official spec).

## Commands

- `make` / `make check` — lint + build + test + govulncheck
- `make test`, `make lint` — individual steps

## Layout

One file per API area in `garminconnect/` (`wellness.go`, `sleep.go`, …). All HTTP
goes through the helpers in `client.go` (`get`, `post`, `put`, `del`, `getBytes`,
`upload`); endpoint methods stay thin. Errors: sentinels (`ErrUnauthorized`,
`ErrRateLimit`, `ErrNoData`, `ErrMFARequired`, `ErrCaptchaRequired`) plus
`*APIError{StatusCode, Path, Body}`.

## Constraints — do not "fix" these

- The Go floor is **1.26** (oldest upstream-supported release). Don't bump it
  further just because a dep wants newer; verify floor changes with
  `GOTOOLCHAIN=go1.26.0 go build ./... && go test ./...`.
- `newUTLSClient` uses a uTLS Android fingerprint to pass Garmin's bot detection.
  A plain `http.Transport` breaks login.
- `testDate` (2026-01-01, in `activities_test.go`) is baked into cassette URLs.
  Changing it breaks replay matching.

## Tests

Tests replay go-vcr cassettes (`garminconnect/tests/testdata/cassettes/`) — never
live, no credentials needed.

- `newVCRClient(t)` — one cassette per test, named `<t.Name()>.yaml`. Never shared.
- `skipAPIError(t, err)` before `require.NoError` for endpoints not available on
  every account (a recorded 4xx skips the test).
- Cassettes are sanitized **inline at record time** (`BeforeSaveHook` →
  `internal/sanitize`): every measurement becomes `1`/`1.0`, every free-text string
  `"TEST"`, IDs synthesized. So tests assert structure and non-zero/non-empty —
  never specific values. A raw cassette never touches disk; never commit one.
  Text/plain and file bodies (multipart uploads, GPX/FIT/CSV/ZIP downloads) become
  `"TEST"` entirely — assert downloads are non-empty, nothing more.
- Recording uses `ModeRecordOnly` and paces live requests 1s apart; replay serves
  each recorded interaction once, in order, so repeated identical requests are fine.
- URL dates after 2026-01-01 are rewritten to 2026-01-01, so every date in a URL
  must be on or before `testDate`. Use `dataDate(t, "GARMIN_SUMMARY_DATE")` (or
  `GARMIN_SLEEP_DATE`) for data-dependent days, and query a single day for range
  endpoints — a multi-day range collapses and won't replay.
- Don't build a URL from a string read out of a response: it is `"TEST"` on replay.
  IDs are fine (they become `12345678` in bodies and URLs alike).
- `newServerClient` (httptest) tests are for client behaviour (auth, errors,
  request encoding). Don't use them to cover an endpoint: several mocks asserted a
  wrong path or body and passed for months. Record it instead.

## Adding an API method

1. **Probe the endpoint live first** (curl with a token from `internal/gettoken`).
   python-garminconnect is the best reference but its paths and bodies are often
   wrong or stale; Garmin's 400 messages usually name the missing field.
2. Method in `garminconnect/<area>.go`; test in `garminconnect/tests/<area>_test.go`
   using `newVCRClient(t)` (the recorder auto-discovers it). Writes get one
   lifecycle test that creates, uses and deletes its own data, so recording leaves
   the account unchanged.
3. Record: `go run ./internal/record --missing` with a cached token in
   `.garmin_token.json` (or `GARMIN_EMAIL`/`GARMIN_PASSWORD`). Without `--missing` it
   deletes and re-records **every** cassette. To record one test:
   `GARMIN_TOKEN=… GARMIN_DISPLAY_NAME=… go test ./garminconnect/tests/ -run '^TestX$'`.
   Record writes against a throwaway account, never a real one.
4. **Inspect the new cassette** for PII the sanitizer doesn't cover yet (new ID
   fields, biometrics, location, schedule times). If anything survives, extend
   `internal/sanitize` + its tests, then `go run ./internal/record --sanitize`
   re-applies it to every cassette offline.
5. If it covers something new, mention it in its area's row of the README API table.

Garmin rate-limits bursts with 429s. course-service also rejects deleting a course
right after creating it; tests that must wait sleep only while recording
(`os.Getenv("GARMIN_TOKEN") != ""`), so replay stays fast.

## Tools

|         Path         |                             Purpose                              |
| -------------------- | ---------------------------------------------------------------- |
| `internal/gettoken/` | Logs in and prints the OAuth token                               |
| `internal/record/`   | Re-records cassettes against a live account (sanitizes inline)   |
| `internal/smoke/`    | Live API smoke test, run weekly by `.github/workflows/smoke.yml` |
| `internal/login/`    | Shared login helper for the tools (cached token, env creds, MFA) |
| `internal/sanitize/` | Cassette sanitizer used by the recorder                          |
