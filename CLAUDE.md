# CLAUDE.md

## Running tests

```bash
make check        # lint + build + test + govulncheck
go test ./garminconnect/tests/...  # tests only
make lint         # lint only
```

Tests use go-vcr cassettes (`garminconnect/tests/testdata/cassettes/*.yaml`) — they replay recorded HTTP interactions and never hit the live API. No credentials needed.

## Test conventions

- `testDate` is fixed at `time.Date(2026, 1, 1, ...)` in `activities_test.go`. Cassette URLs are keyed to this date; changing it breaks URL matching.
- `newVCRClient(t, "cassette_name")` wires a client to a cassette for replay.
- `skipAPIError(t, err)` skips a test when the cassette captured a 4xx — the endpoint isn't available on the recorded account. Use it before `require.NoError` for optional endpoints.
- Float assertions use **2-significant-figure values** (`1500.0`, `4100.0`), not raw API precision (`1484.8990478515625`). The sanitizer rounds all floats with 4+ decimal places to 2 sig figs before cassettes are committed.

## Adding a new API method

1. Implement the method in `garminconnect/`.
2. Write a test in `garminconnect/tests/<area>_test.go` using `newVCRClient`.
3. Add the test function name to the `TESTS` array in `tools/record_cassettes.sh`.
4. Record the cassette (see below).

## Recording cassettes

Re-record all cassettes against a live account:

```bash
GARMIN_EMAIL=you@example.com GARMIN_PASSWORD=secret bash tools/record_cassettes.sh
```

Record only missing cassettes (leaves existing ones untouched):

```bash
bash tools/record_cassettes.sh --missing
```

The script logs in once, records each cassette, then runs `tools/sanitize_cassettes.py` automatically to strip PII before commit.

## Sanitizing cassettes

The sanitizer runs automatically after recording. To run it standalone:

```bash
python3 tools/sanitize_cassettes.py [--display-name "Real Name"] [--email real@example.com]
```

It is safe to re-run on already-sanitized cassettes (idempotent). What it replaces:

- Integer fields by name: profile/user IDs → `12345678`, device IDs → `9876543210`, activity IDs → sequential `10000001+`, sample PKs → sequential `1000000000001+`
- UUIDs (hyphenated and bare 32-char hex) → SHA-256-derived synthetic values
- Email addresses → `test@example.com`
- Display name (via `--display-name`) and all `*FullName` fields → `"Test User"`
- `locationName` → `"Test Location"`, `activityName` → `"Activity"`, `serialNumber` → `"TEST000000"`
- Datetime strings (`2025-12-31T13:50:13`, `2025-12-31 13:50:13.944`, etc.) → `2026-01-01T00:00:00` / `2026-01-01 00:00:00`
- Date-only JSON string values → `"2026-01-01"`
- Floats with 4+ decimal places → 2 significant figures
- Response headers stripped: `Cf-Ray`, `Date`, `Nel`, `Report-To`, `Alt-Svc`, `Cf-Cache-Status`, `Cache-Control`, `Pragma`, `Server`
- Response durations → `100ms`

## Tools

| Path | Purpose |
|---|---|
| `Makefile` | `make check` runs lint + build + test + govulncheck |
| `tools/gettoken/` | Logs in and prints the OAuth token; used by `record_cassettes.sh` |
| `tools/record_cassettes.sh` | Records cassettes for all tests against a live account |
| `tools/sanitize_cassettes.py` | Strips PII from cassettes |
