#!/usr/bin/env python3
"""Sanitize VCR cassettes: dynamically detect and replace PII."""

import argparse
import hashlib
import math
import os
import re

CASSETTE_DIR = "garminconnect/tests/testdata/cassettes"

STRIP_HEADERS = {
    "Cf-Ray", "Date", "Nel", "Report-To", "Alt-Svc", "Cf-Cache-Status",
    "Cache-Control", "Pragma", "Server", "Set-Cookie", "Content-Length",
}

# Field names whose values should be replaced with a fixed synthetic ID.
# "id" catches the bare top-level Garmin user ID in social-profile responses.
_PROFILE_FIELDS = {"userProfilePk", "userProfilePK", "userId", "id"}
_DEVICE_FIELDS  = {"deviceId", "sourceDeviceId"}
_PROFILE_SYNTH  = "12345678"
_DEVICE_SYNTH   = "9876543210"

# Field names whose values get sequential synthetic IDs.
_ACTIVITY_FIELDS = {"activityId", "parentActivityId", "activitySummaryId"}
_SAMPLE_FIELDS   = {"samplePk"}
_ACTIVITY_BASE   = 10_000_001
_SAMPLE_BASE     = 1_000_000_000_001

# Detect JSON integer fields: "fieldName": 123456
_FIELD_INT_RE = re.compile(r'"([A-Za-z][A-Za-z0-9_]*)"[ \t]*:[ \t]*(\d{6,})')

_UUID_RE = re.compile(
    r'[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}', re.I
)
_UUID_BARE_RE = re.compile(r'(?<![0-9a-f-])[0-9a-f]{32}(?![0-9a-f-])', re.I)
_SYNTH_UUID_PREFIX = "aaaaaaaa-0000-0000-0000-"
_SYNTH_UUID_BARE_PREFIX = "00000000000000000000"

_EMAIL_RE    = re.compile(r'[\w.+%-]+@[\w.-]+\.[a-z]{2,}', re.I)
_SYNTH_EMAIL = "test@example.com"

# Static replacements applied after dynamic ones (longer strings first).
_STATIC = [
    ("garmin-connect-prod", "garmin-connect-test"),
]


def _synth_uuid(original: str) -> str:
    h = hashlib.sha256(original.lower().encode()).hexdigest()[:12]
    return f"{_SYNTH_UUID_PREFIX}{h}"


def _synth_uuid_bare(original: str) -> str:
    h = hashlib.sha256(original.lower().encode()).hexdigest()[:12]
    return f"{_SYNTH_UUID_BARE_PREFIX}{h}"


def _is_synthetic_uuid(s: str) -> bool:
    return s.lower().startswith(_SYNTH_UUID_PREFIX)


def _is_synthetic_uuid_bare(s: str) -> bool:
    return s.lower().startswith(_SYNTH_UUID_BARE_PREFIX)


def discover(files: list[str]) -> dict[str, str]:
    """Two-pass: collect all real PII values, then build consistent mapping."""
    profile_ids: set[str] = set()
    device_ids: set[str]  = set()
    activity_ids: set[str] = set()
    sample_ids: set[str]   = set()
    uuids: set[str]        = set()
    uuid_bares: set[str]   = set()

    for path in files:
        with open(path, encoding="utf-8") as f:
            content = f.read()

        for field, value in _FIELD_INT_RE.findall(content):
            if field in _PROFILE_FIELDS and value != _PROFILE_SYNTH and not value.startswith("1000000"):
                profile_ids.add(value)
            elif field in _DEVICE_FIELDS and value != _DEVICE_SYNTH:
                device_ids.add(value)
            elif field in _ACTIVITY_FIELDS:
                # Skip already-synthetic values (≤8 digits)
                if len(value) > 8:
                    activity_ids.add(value)
            elif field in _SAMPLE_FIELDS:
                if not value.startswith("100000000000"):
                    sample_ids.add(value)

        for m in _UUID_RE.finditer(content):
            v = m.group(0)
            if not _is_synthetic_uuid(v):
                uuids.add(v.lower())

        for m in _UUID_BARE_RE.finditer(content):
            v = m.group(0)
            if not _is_synthetic_uuid_bare(v):
                uuid_bares.add(v.lower())

    mapping: dict[str, str] = {}

    for v in profile_ids:
        mapping[v] = _PROFILE_SYNTH
    for v in device_ids:
        mapping[v] = _DEVICE_SYNTH

    for i, v in enumerate(sorted(activity_ids, reverse=True)):
        mapping[v] = str(_ACTIVITY_BASE + i)
    for i, v in enumerate(sorted(sample_ids)):
        mapping[v] = str(_SAMPLE_BASE + i)

    for v in uuids:
        mapping[v] = _synth_uuid(v)
        mapping[v.upper()] = _synth_uuid(v)
    for v in uuid_bares:
        mapping[v] = _synth_uuid_bare(v)
        mapping[v.upper()] = _synth_uuid_bare(v)

    return mapping


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


_HTTP2_PROTO_RE = re.compile(
    r"^(\s+)proto: HTTP/2\.0\n\1proto_major: 2\n\1proto_minor: 0$",
    re.MULTILINE,
)
_UNCOMPRESSED_RE = re.compile(r"^\s+uncompressed: (?:true|false)\n", re.MULTILINE)


def normalize_response_metadata(content: str) -> str:
    """Normalize HTTP/2 responses to HTTP/1.1 and drop go-vcr's uncompressed flag."""
    content = _HTTP2_PROTO_RE.sub(
        lambda m: f"{m.group(1)}proto: HTTP/1.1\n{m.group(1)}proto_major: 1\n{m.group(1)}proto_minor: 1",
        content,
    )
    content = _UNCOMPRESSED_RE.sub("", content)
    return content


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


_PRECISE_FLOAT_RE = re.compile(r'-?\d+\.\d{4,}')


def _round_2sig(v: float) -> str:
    """Round to 2 significant figures; keep 1 decimal for values < 1."""
    if v == 0.0:
        return "0.0"
    if abs(v) < 1:
        return f"{v:.1f}"
    exp = math.floor(math.log10(abs(v)))
    sig_decimals = max(0, 1 - exp)
    factor = 10 ** (exp - 1)
    # Re-round to sig_decimals to eliminate float multiplication artifacts.
    rounded = round(round(v / factor) * factor, sig_decimals)
    if sig_decimals == 0:
        return str(int(rounded))
    return f"{rounded:.{sig_decimals}f}"


def simplify_floats(content: str) -> str:
    """Replace IEEE 754 precise floats with 2-significant-figure round numbers."""
    return _PRECISE_FLOAT_RE.sub(lambda m: _round_2sig(float(m.group(0))), content)


def apply_mapping(content: str, mapping: dict[str, str]) -> str:
    # Sort longest first so substrings don't get replaced before the full value.
    for old in sorted(mapping, key=len, reverse=True):
        content = content.replace(old, mapping[old])
    return content


# Matches any field whose name ends in "fullname" (any capitalisation), plus
# locationName, activityName, serialNumber — all of which may reveal identity.
_STRING_FIELD_RE = re.compile(
    r'("[a-zA-Z]*[Ff]ull[Nn]ame"\s*:\s*)"[^"]*"'
    r'|("(?:locationName|activityName|serialNumber)"\s*:\s*)"[^"]*"'
)
_STRING_FIELD_SYNTH = {
    "locationname": "Test Location",
    "activityname": "Activity",
    "serialnumber": "TEST000000",
}


def _replace_string_fields(content: str) -> str:
    def _sub(m: re.Match) -> str:
        # group(1) matches fullname variants, group(2) matches the named fields
        prefix = m.group(1) or m.group(2)
        key = re.search(r'"(\w+)"', prefix).group(1).lower()
        synth = _STRING_FIELD_SYNTH.get(key, "Test User")
        return f'{prefix}"{synth}"'
    return _STRING_FIELD_RE.sub(_sub, content)


# Auth cassette sanitization: strip credentials and tokens from SSO / diAuth bodies.
_AUTH_JSON_RE = re.compile(
    r'"(password|access_token|refresh_token|serviceTicketId|mfaVerificationCode|username)"\s*:\s*"[^"]*"'
)
_AUTH_JSON_SYNTH = {
    "password": "test",
    "access_token": "test_access_token",
    "refresh_token": "test_refresh_token",
    "serviceTicketId": "ST-test",
    "mfaVerificationCode": "123456",
    "username": "test@example.com",
}
# URL-encoded form bodies (diAuth POSTs use application/x-www-form-urlencoded)
_AUTH_FORM_RE = re.compile(r'((?:refresh_token|service_ticket)=)[^\s&\'"]+')


def sanitize_auth_fields(content: str) -> str:
    def _sub_json(m: re.Match) -> str:
        field = re.search(r'"(\w+)"', m.group(0)).group(1)
        return f'"{field}": "{_AUTH_JSON_SYNTH[field]}"'

    content = _AUTH_JSON_RE.sub(_sub_json, content)
    content = _AUTH_FORM_RE.sub(lambda m: m.group(1) + "test", content)
    return content


def apply_static(content: str, display_name: str, email: str) -> str:
    if display_name:
        content = content.replace(f'"{display_name}"', '"Test User"')
        content = content.replace(display_name, "Test User")
    if email:
        content = content.replace(email, _SYNTH_EMAIL)
    # Replace any remaining real emails (catches addresses not passed via --email)
    content = _EMAIL_RE.sub(_SYNTH_EMAIL, content)
    content = _replace_string_fields(content)
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
    content = normalize_response_metadata(content)
    content = zero_datetimes(content)
    content = simplify_floats(content)
    content = apply_mapping(content, mapping)
    content = sanitize_auth_fields(content)
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
