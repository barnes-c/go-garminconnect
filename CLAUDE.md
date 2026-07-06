# CLAUDE.md

Go client for Garmin Connect's private API (reverse-engineered — no official spec).

## Commands

- `make` / `make check` — lint + build + test + govulncheck
- `make test`, `make lint` — individual steps

## Layout

One file per API area in `garminconnect/` (`wellness.go`, `sleep.go`, …). All HTTP
goes through the helpers in `client.go` (`get`, `post`, `getBytes`, `upload`);
endpoint methods stay thin. Errors: sentinels (`ErrUnauthorized`, `ErrRateLimit`,
`ErrNoData`, `ErrMFARequired`) plus `*APIError{StatusCode, Path}`.

## Constraints — do not "fix" these

- The Go floor is **1.24**. `golang.org/x/crypto` and `golang.org/x/sys` are pinned
  below latest because newer versions require Go 1.25. Don't bulk-upgrade deps;
  verify floor changes with `GOTOOLCHAIN=go1.24.4 go build ./... && go test ./...`.
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

## Adding an API method

1. Method in `garminconnect/<area>.go`; test in `garminconnect/tests/<area>_test.go`
   using `newVCRClient(t)` (the recorder auto-discovers it).
2. Record: `GARMIN_EMAIL=… GARMIN_PASSWORD=… go run ./internal/record --missing`
   (a cached token in `.garmin_token.json` also works, no credentials needed).
3. **Inspect the new cassette** for PII the sanitizer doesn't cover yet (new ID
   fields, biometrics, location, schedule times). If anything survives, extend
   `internal/sanitize` + its tests and re-record.
4. Add the method to the README API table.

If a test replays against `testDate` and that day has no data on the recording
account, set `GARMIN_SUMMARY_DATE` / `GARMIN_SLEEP_DATE` to a recent day with data;
the sanitizer rewrites it back to 2026-01-01 in the cassette URL.

## Tools

|         Path         |                             Purpose                              |
| -------------------- | ---------------------------------------------------------------- |
| `internal/gettoken/` | Logs in and prints the OAuth token                               |
| `internal/record/`   | Re-records cassettes against a live account (sanitizes inline)   |
| `internal/smoke/`    | Live API smoke test, run weekly by `.github/workflows/smoke.yml` |
