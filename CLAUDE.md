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

`TestActivitiesForDailySummary` replays against `testDate`, which has no activities. To capture a non-empty response, set `GARMIN_SUMMARY_DATE` to a day with a logged activity when recording — the sanitizer rewrites that real date back to `2026-01-01` in the cassette URL so it still replays:

```bash
GARMIN_SUMMARY_DATE=2026-06-18 GARMIN_EMAIL=... GARMIN_PASSWORD=... bash tools/record_cassettes.sh --missing
```

## Sanitizing cassettes

`record_cassettes.sh` runs the sanitizer automatically before cassettes land on disk, so you rarely invoke it directly. It strips PII (IDs, UUIDs, emails, names, dates, epoch timestamps), rounds floats to 2 significant figures (which also coarsens GPS coordinates), and is safe to re-run (idempotent). **Never commit an unsanitized cassette.**

The replacement rules live in `tools/sanitize_cassettes.py` — change them there, not here. Run it standalone after hand-editing a cassette:

```bash
python3 tools/sanitize_cassettes.py [--display-name "Real Name"] [--email real@example.com]
```

Pass `--display-name` and `--email` when recording: the script can't infer your real name and address from the data, so they must be supplied to be scrubbed.

## Tools

|             Path              |                              Purpose                              |
| ----------------------------- | ----------------------------------------------------------------- |
| `Makefile`                    | `make check` runs lint + build + test + govulncheck               |
| `tools/gettoken/`             | Logs in and prints the OAuth token; used by `record_cassettes.sh` |
| `tools/record_cassettes.sh`   | Records cassettes for all tests against a live account            |
| `tools/sanitize_cassettes.py` | Strips PII from cassettes                                         |
