#!/usr/bin/env python3
"""Sanitize VCR cassettes: detect and replace PII before cassettes are committed.

Idempotent — safe to re-run on already-sanitized cassettes. Replacements:

- Integer ID fields by name: profile/user/owner IDs -> 12345678, device IDs ->
  9876543210, activity IDs -> sequential 10000001+, sample PKs -> 1000000000001+
- UUIDs (hyphenated and bare 32-char hex) -> one all-f constant; nothing is
  derived from the real value
- Datetime / date-only strings -> 2026-01-01[T00:00:00]
- Request-URL dates after 2026-01-01 -> 2026-01-01 (real recording dates;
  synthetic test dates only ever look back from testDate, so they're left alone)
- Emails -> test@example.com
- Every free-text string value in a response body -> "TEST", except dates,
  UUID-shaped values, and the "testuser"/"Test User" display-name placeholders.
  Catches names, gear/workout/device labels, descriptions, etc. generically
- Every numeric value in a response body -> 1, EXCEPT identifier/structural
  fields (*Id, *Pk, *Count, *Index, *Version, *Number, ...). This covers every
  real measurement (heart rate, sleep, SpO2, GPS, distance, calories, epoch
  timestamps, ...) generically. The constant 1 (not 0) keeps NotZero assertions
  meaningful as field-decoding checks. Number type is preserved (float -> 1.0,
  int -> 1) so Go decoding is unaffected
- Volatile response headers stripped; durations -> 100ms
"""

import argparse
import os
import re

CASSETTE_DIR = "garminconnect/tests/testdata/cassettes"

STRIP_HEADERS = {"Cf-Ray", "Date", "Nel", "Report-To", "Alt-Svc", "Cf-Cache-Status", "Cache-Control", "Pragma", "Server"}

# Identifier fields — name ends in Id/Pk with a 6+ digit value — hold real Garmin
# IDs. Every distinct value collapses to one synthetic constant. Small ids (type
# enums like typeId:1) are left alone. Idempotent: _SYNTH_ID is skipped on re-runs.
_ID_FIELD_RE = re.compile(r'"[A-Za-z0-9_]*(?:[Ii]d|[Pp][Kk])":\s*(\d{6,})')
_SYNTH_ID = "12345678"

_UUID_RE = re.compile(
    r'[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}', re.I
)
_UUID_BARE_RE = re.compile(r'(?<![0-9a-f-])[0-9a-f]{32}(?![0-9a-f-])', re.I)

_EMAIL_RE    = re.compile(r'[\w.+%-]+@[\w.-]+\.[a-z]{2,}', re.I)
_SYNTH_EMAIL = "test@example.com"

# Static replacements applied after dynamic ones (longer strings first).
_STATIC = [
    ("garmin-connect-prod", "garmin-connect-test"),
]


# Every UUID collapses to one obviously-bogus all-f value (hyphen-stripped for
# bare 32-char hex); nothing in the tests distinguishes one UUID from another,
# and nothing is derived from the real value. Idempotent: re-running rewrites
# already-synthetic UUIDs to the same value.
_SYNTH_UUID = "ffffffff-ffff-ffff-ffff-ffffffffffff"


def scrub_uuids(content: str) -> str:
    content = _UUID_RE.sub(_SYNTH_UUID, content)
    content = _UUID_BARE_RE.sub(_SYNTH_UUID.replace("-", ""), content)
    return content


def discover(files: list[str]) -> dict[str, str]:
    """Collect real Garmin IDs (any *Id/*Pk field) and map each to one synthetic
    constant. Done across all files so a value is replaced consistently, including
    where it appears in a request URL."""
    ids: set[str] = set()
    for path in files:
        with open(path, encoding="utf-8") as f:
            content = f.read()
        for value in _ID_FIELD_RE.findall(content):
            if value != _SYNTH_ID:
                ids.add(value)
    return {v: _SYNTH_ID for v in ids}


def strip_response_headers(content: str) -> str:
    lines = content.split("\n")
    result = []
    i = 0
    while i < len(lines):
        line = lines[i]
        m = re.match(r"^(\s+)(\S[^:]+):\s*$", line)
        if m and m.group(2) in STRIP_HEADERS:
            i += 1
            while i < len(lines) and re.match(r"^\s+- ", lines[i]):
                i += 1
            continue
        result.append(line)
        i += 1
    return "\n".join(result)


def normalize_duration(content: str) -> str:
    return re.sub(r"duration: [\d.]+\w*s", "duration: 100ms", content)


# Normalize datetime strings: fix date to test date, zero the time.
# Matches "2025-12-31 13:50:13", "2025-12-31T13:50:13.944", etc.
_DATETIME_RE = re.compile(r'\d{4}-\d{2}-\d{2}([T ])\d{2}:\d{2}:\d{2}(?:\.\d+)?')
# Only replace date-only strings that are JSON string values (between double quotes).
# This avoids touching dates in request URLs (which are YAML values, not JSON strings).
_DATE_ONLY_RE = re.compile(r'(?<=")\d{4}-\d{2}-\d{2}(?=")')
_SYNTH_DATE = "2026-01-01"


def zero_datetimes(content: str) -> str:
    content = _DATETIME_RE.sub(lambda m: f'{_SYNTH_DATE}{m.group(1)}00:00:00', content)
    content = _DATE_ONLY_RE.sub(_SYNTH_DATE, content)
    return content


# Scrub real dates from request URLs. Synthetic test dates are anchored at
# _SYNTH_DATE and only ever look backward (testDate, testDate.AddDate(0,-1,0),
# ...), so any URL date later than _SYNTH_DATE is a real recording date (e.g.
# the day an activity was logged) and must not be committed. ISO dates compare
# lexicographically, so a string ">" is a chronological "after".
_URL_LINE_RE = re.compile(r'^\s*url:\s*https?://')
_BARE_DATE_RE = re.compile(r'\d{4}-\d{2}-\d{2}')


def scrub_url_dates(content: str) -> str:
    def fix(line: str) -> str:
        if not _URL_LINE_RE.match(line):
            return line
        return _BARE_DATE_RE.sub(
            lambda m: _SYNTH_DATE if m.group(0) > _SYNTH_DATE else m.group(0), line
        )
    return "\n".join(fix(line) for line in content.split("\n"))


# Generic privacy rule: in a response body, every numeric value is a real
# measurement (heart rate, sleep, SpO2, GPS, distance, calories, ...) unless its
# field is an identifier or structural metadata. Replace them all with 1 so no
# real personal data is committed (1, not 0, so NotZero assertions still verify
# decoding). Being generic avoids maintaining a per-metric list and covers future
# cassettes for free. Identifier/structural fields are preserved so synthetic IDs
# still flow into chained request URLs and counts/types/versions stay coherent.
# Exception: a preserved field is still neutralized if its value looks like an
# epoch-ms timestamp — Garmin reuses structural fields like "version"/"sequence"
# to hold record timestamps (e.g. weight/latest), which would leak a real time.
_PRESERVE_KEY_RE = re.compile(
    r'(?i)(id|pk|count|index|version|number|order|sequence|priority|'
    r'category|month|year|offset|zoneid|typekey)$'
)
_EPOCH_MS_RE = re.compile(r'^1[5-9]\d{11}$')  # 13-digit, ~2017-2033
_KEY_NUM_RE = re.compile(r'"([A-Za-z_][A-Za-z0-9_]*)":\s*(-?\d+(?:\.\d+)?)')
_ARRAY_NUM_RE = re.compile(r'([\[,])(-?\d+(?:\.\d+)?)(?=[,\]])')
_BODY_LINE_RE = re.compile(r'^\s*body:\s')


def _placeholder(num: str) -> str:
    # Replace real measurements with a constant 1 (1.0 for floats), not 0, so the
    # many NotZero assertions still confirm a field decoded correctly while no
    # real value is committed. Type is preserved: Go won't decode 1.0 into an int.
    return "1.0" if "." in num else "1"


def neutralize_metrics(content: str) -> str:
    def key_sub(m: re.Match) -> str:
        key, num = m.group(1), m.group(2)
        if _PRESERVE_KEY_RE.search(key) and not _EPOCH_MS_RE.match(num):
            return m.group(0)
        return f'"{key}":{_placeholder(num)}'

    def array_sub(m: re.Match) -> str:
        return f"{m.group(1)}{_placeholder(m.group(2))}"

    def fix(line: str) -> str:
        if not _BODY_LINE_RE.match(line):
            return line
        line = _KEY_NUM_RE.sub(key_sub, line)
        line = _ARRAY_NUM_RE.sub(array_sub, line)
        return line

    return "\n".join(fix(line) for line in content.split("\n"))


def apply_mapping(content: str, mapping: dict[str, str]) -> str:
    # Sort longest first so substrings don't get replaced before the full value.
    for old in sorted(mapping, key=len, reverse=True):
        content = content.replace(old, mapping[old])
    return content


# Generic string-value scrub: any free-text value in a response body could carry
# identity (names, gear/workout/device labels, descriptions, ...), so replace
# them all with "TEST". Preserve values other steps depend on: dates (already
# normalized), UUID-shaped values (scrub_uuids collapses them), and the display
# name placeholders that login_profile.yaml and its test rely on.
_TEXT_VALUE_RE = re.compile(r':"((?:[^"\\]|\\.)*)"')
_PRESERVE_TEXT = {"testuser", "Test User"}
_DATEISH_RE = re.compile(r"^\d{4}-\d{2}-\d{2}")
_UUIDISH_RE = re.compile(r"^[0-9a-f]{8}-[0-9a-f-]+$|^[0-9a-f]{32}$", re.I)


def scrub_text_values(content: str) -> str:
    def sub(m: re.Match) -> str:
        val = m.group(1)
        if not val or val in _PRESERVE_TEXT or _DATEISH_RE.match(val) or _UUIDISH_RE.match(val):
            return m.group(0)
        return ':"TEST"'

    def fix(line: str) -> str:
        if not _BODY_LINE_RE.match(line):
            return line
        return _TEXT_VALUE_RE.sub(sub, line)

    return "\n".join(fix(line) for line in content.split("\n"))


def apply_static(content: str, display_name: str, email: str) -> str:
    if display_name:
        content = content.replace(f'"{display_name}"', '"Test User"')
        content = content.replace(display_name, "Test User")
    if email:
        content = content.replace(email, _SYNTH_EMAIL)
    # Replace any remaining real emails (catches addresses not passed via --email)
    content = _EMAIL_RE.sub(_SYNTH_EMAIL, content)
    for old, new in _STATIC:
        content = content.replace(old, new)
    return content


def sanitize_file(
    path: str, mapping: dict[str, str], display_name: str, email: str
) -> None:
    with open(path, encoding="utf-8") as f:
        content = f.read()

    content = strip_response_headers(content)
    content = normalize_duration(content)
    content = zero_datetimes(content)
    content = scrub_url_dates(content)
    content = neutralize_metrics(content)
    content = apply_mapping(content, mapping)
    content = scrub_text_values(content)
    content = scrub_uuids(content)
    content = apply_static(content, display_name, email)

    with open(path, "w", encoding="utf-8") as f:
        f.write(content)


def main() -> None:
    parser = argparse.ArgumentParser(description="Sanitize VCR cassettes.")
    parser.add_argument("--display-name", default="", help="Real display name to replace")
    parser.add_argument("--email", default="", help="Real email address to replace")
    args = parser.parse_args()

    files = sorted(
        os.path.join(CASSETTE_DIR, f)
        for f in os.listdir(CASSETTE_DIR)
        if f.endswith(".yaml")
    )

    print("Discovering PII...")
    mapping = discover(files)
    print(f"  Found {len(mapping)} values to replace.")

    for path in files:
        sanitize_file(path, mapping, args.display_name, args.email)
        print(f"  {os.path.basename(path)}")

    print(f"\nSanitized {len(files)} cassettes.")


if __name__ == "__main__":
    main()
