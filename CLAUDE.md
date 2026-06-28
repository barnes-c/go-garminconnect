# CLAUDE.md

## Running tests

```bash
make              # check all
make check        # lint + build + test + govulncheck
make test         # tests only
make lint         # lint only
```

Tests use go-vcr cassettes (`garminconnect/tests/testdata/cassettes/*.yaml`) — they replay recorded HTTP interactions and never hit the live API. No credentials needed.

## Test conventions

- `testDate` is fixed at `time.Date(2026, 1, 1, ...)` in `activities_test.go`. Cassette URLs are keyed to this date; changing it breaks URL matching.
- `newVCRClient(t)` wires a client to a cassette for replay. The cassette is named after the test (`t.Name()`), e.g. `TestFloors` → `testdata/cassettes/TestFloors.yaml`. One cassette per test — no shared cassettes.
- `skipAPIError(t, err)` skips a test when the cassette captured a 4xx — the endpoint isn't available on the recorded account. Use it before `require.NoError` for optional endpoints.
- The sanitizer replaces every measurement value with **`1`** (`1.0` for floats) and every free-text string with **`"TEST"`**, so tests assert structure and non-zero/non-empty — never specific values. IDs are synthesized and string/date/UUID fields are preserved structurally. It runs **inline** while recording (the recorder's `BeforeSaveHook` calls `internal/sanitize`), so a cassette is never written to disk unsanitized. See `internal/sanitize`.

## Adding a new API method

1. Implement the method in `garminconnect/`.
2. Write a test in `garminconnect/tests/<area>_test.go` using `newVCRClient(t)`. The record script auto-discovers it — no list to update.
3. Record the cassette (see below).
4. **Inspect the recorded cassette for sensitive data** the sanitizer doesn't already cover (new ID fields, biometric/behavioural metrics, location, schedule times). If anything personal survives sanitization, extend `internal/sanitize` (add a test there) and re-record before committing.
5. Add the method to the API table in `README.md` under the relevant section.

## Recording cassettes

Re-record all cassettes against a live account:

```bash
GARMIN_EMAIL=you@example.com GARMIN_PASSWORD=secret go run ./internal/record
```

Record only missing cassettes (leaves existing ones untouched):

```bash
go run ./internal/record --missing
```

Credentials may be omitted if a valid cached token exists in `.garmin_token.json`. The tool logs in once, discovers every test that calls `newVCRClient`, and records each. Cassettes are **sanitized inline** as they're written (the recorder's `BeforeSaveHook` calls `internal/sanitize`) — there is no separate scrubbing step.

Some tests replay against `testDate`, which may have no data. To capture a non-empty response, set a date-override env var to a recent day that has data when recording — the sanitizer rewrites that real date back to `2026-01-01` in the cassette URL (dates after `testDate` are scrubbed) so it still replays:

- `GARMIN_SUMMARY_DATE` → `TestActivitiesForDailySummary` (a day with a logged activity)
- `GARMIN_SLEEP_DATE` → `TestDailySleepData` (a night with recorded sleep)

```bash
GARMIN_SLEEP_DATE=2026-06-17 GARMIN_SUMMARY_DATE=2026-06-18 \
  GARMIN_EMAIL=... GARMIN_PASSWORD=... go run ./internal/record --missing
```

## Sanitizing cassettes

Sanitization runs inline during recording, so you never invoke it separately. The `internal/sanitize` package strips PII (IDs, UUIDs, emails, names, dates, epoch timestamps), replaces every measurement with `1`/`1.0` and every free-text value with `"TEST"`, and is idempotent. Its tests (`internal/sanitize/sanitize_test.go`) assert it scrubs crafted PII and is a fixed point over every committed cassette. **Never commit an unsanitized cassette** — but because scrubbing happens in the `BeforeSaveHook`, a raw cassette is never written to disk in the first place.

## Tools

|          Path           |                       Purpose                       |
| ----------------------- | --------------------------------------------------- |
| `Makefile`              | `make check` runs lint + build + test + govulncheck |
| `internal/gettoken/`       | Logs in and prints the OAuth token                  |
| `internal/record/`         | Records cassettes for all tests against a live account (sanitizes inline) |
| `internal/sanitize/`    | Inline cassette PII scrubber (used by the test recorder) |
