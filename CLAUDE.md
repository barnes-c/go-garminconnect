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
- The sanitizer replaces every measurement value with **`1`** (`1.0` for floats) and every free-text string with **`"TEST"`**, so tests assert structure and non-zero/non-empty — never specific values. IDs are synthesized and string/date/UUID fields are preserved structurally. See `tools/sanitize_cassettes.py`.

## Adding a new API method

1. Implement the method in `garminconnect/`.
2. Write a test in `garminconnect/tests/<area>_test.go` using `newVCRClient(t)`. The record script auto-discovers it — no list to update.
3. Record the cassette (see below).
4. **Inspect the recorded cassette for sensitive data** the sanitizer doesn't already cover (new ID fields, biometric/behavioural metrics, location, schedule times). If anything personal survives sanitization, extend `tools/sanitize_cassettes.py` and re-run it before committing.
5. Add the method to the API table in `README.md` under the relevant section.

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

Some tests replay against `testDate`, which may have no data. To capture a non-empty response, set a date-override env var to a recent day that has data when recording — the sanitizer rewrites that real date back to `2026-01-01` in the cassette URL (dates after `testDate` are scrubbed) so it still replays:

- `GARMIN_SUMMARY_DATE` → `TestActivitiesForDailySummary` (a day with a logged activity)
- `GARMIN_SLEEP_DATE` → `TestDailySleepData` (a night with recorded sleep)

```bash
GARMIN_SLEEP_DATE=2026-06-17 GARMIN_SUMMARY_DATE=2026-06-18 \
  GARMIN_EMAIL=... GARMIN_PASSWORD=... bash tools/record_cassettes.sh --missing
```

## Sanitizing cassettes

`record_cassettes.sh` runs the sanitizer automatically before cassettes land on disk, so you rarely invoke it directly. It strips PII (IDs, UUIDs, emails, names, dates, epoch timestamps), rounds floats to 2 significant figures (which also coarsens GPS coordinates), and is safe to re-run (idempotent). **Never commit an unsanitized cassette.**

The replacement rules live in `tools/sanitize_cassettes.py` — change them there, not here. Run it standalone after hand-editing a cassette:

```bash
python3 tools/sanitize_cassettes.py
```

Pass `--display-name` and `--email` when recording: the script can't infer your real name and address from the data, so they must be supplied to be scrubbed.

## Tools

|             Path              |                              Purpose                              |
| ----------------------------- | ----------------------------------------------------------------- |
| `Makefile`                    | `make check` runs lint + build + test + govulncheck               |
| `tools/gettoken/`             | Logs in and prints the OAuth token; used by `record_cassettes.sh` |
| `tools/record_cassettes.sh`   | Records cassettes for all tests against a live account            |
| `tools/sanitize_cassettes.py` | Strips PII from cassettes                                         |
