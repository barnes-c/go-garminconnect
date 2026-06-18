#!/usr/bin/env bash
# Record VCR cassettes for all tests that hit the real Garmin Connect API.
#
# Usage:
#   GARMIN_EMAIL=you@example.com GARMIN_PASSWORD=secret ./record_cassettes.sh
#
# Logs in once upfront, then passes the token to every test so Garmin's SSO
# endpoint is only hit a single time regardless of how many cassettes are recorded.
#
# Cassettes already recorded with real data are preserved.
# The synthetic activities_empty cassette (used for the ErrNoData test) is also
# preserved because it cannot be recorded from a real account.
set -euo pipefail

# -m / --missing  skip deletion; only record cassettes that don't exist yet
MISSING_ONLY=false
for arg in "$@"; do
    [[ "$arg" == "-m" || "$arg" == "--missing" ]] && MISSING_ONLY=true
done

: "${GARMIN_EMAIL:?GARMIN_EMAIL must be set}"
: "${GARMIN_PASSWORD:?GARMIN_PASSWORD must be set}"

CASSETTE_DIR="garminconnect/tests/testdata/cassettes"
DELAY=5   # seconds between tests to avoid Connect API rate-limiting

# Cassettes to leave untouched.
KEEP=(
    "activities_empty"
    "login_profile"
)

keep() {
    local name="$1"
    for k in "${KEEP[@]}"; do
        [[ "$name" == "$k" ]] && return 0
    done
    return 1
}

TOKEN_FILE=".garmin_token.json"

echo "==> Logging in (once)..."
token_line=$(go run ./tools/gettoken -token-file "$TOKEN_FILE")
export GARMIN_TOKEN
export GARMIN_DISPLAY_NAME
GARMIN_TOKEN=$(printf '%s' "$token_line" | sed -n '1p')
GARMIN_DISPLAY_NAME=$(printf '%s' "$token_line" | sed -n '2p')
echo "    display_name=${GARMIN_DISPLAY_NAME}"

echo ""
if $MISSING_ONLY; then
    echo "==> --missing mode: keeping all existing cassettes, only recording new ones."
else
    echo "==> Removing placeholder cassettes..."
    for f in "$CASSETTE_DIR"/*.yaml; do
        [[ -f "$f" ]] || continue
        name=$(basename "$f" .yaml)
        if keep "$name"; then
            echo "    keeping  $name"
        else
            echo "    removing $name"
            rm "$f"
        fi
    done
fi

# One test per unique cassette.
# Tests that share a cassette with another test are omitted (they replay).
TESTS=(
    TestUserSummary
    TestAllDayStress
    TestActivities
    TestActivityDetail
    TestActivityDetails
    TestActivityTypes
    TestActivitiesForDailySummary
    TestActiveTrainingPlan
    TestCompletedTrainingPlans
    TestTrainingPlanTypes
    TestActivityCount
    TestActivitiesByDate
    TestPersonalRecords
    TestIntensityMinutes
    TestBodyBattery
    TestFloors
    TestHydration
    TestRespiration
    TestSpO2
    TestSteps
    TestRestingHeartRate
    TestDailySteps
    TestWeeklyStress
    TestWeeklyIntensityMinutes
    TestBloodPressure
    TestWeighIns
    TestDailyWeighIns
    TestBodyComposition
    TestGear
    TestGearStats
    TestGoals
    TestEarnedBadges
    TestAvailableBadges
    TestHeartRates
    TestHRVData
    TestSleepData
    TestDevices
    TestLastUsedDevice
    TestPrimaryTrainingDevice
    TestTrainingReadiness
    TestTrainingStatus
    TestMaxMetrics
    TestEnduranceScore
    TestRacePredictions
    TestHillScore
    TestLactateThreshold
    TestFitnessAge
    TestRunningTolerance
    TestCyclingFTP
    TestUserProfile
    TestUserProfileSettings
    TestWorkouts
    TestWorkout
    TestScheduledWorkouts
)

PASS=()
FAIL=()
SKIP=()
total=${#TESTS[@]}

echo ""
echo "==> Recording $total cassettes (${DELAY}s between each)..."
echo ""

for i in "${!TESTS[@]}"; do
    test="${TESTS[$i]}"
    n=$((i + 1))
    echo "--- [$n/$total] $test"

    if output=$(go test ./garminconnect/tests/... -run "^${test}$" -count=1 -v 2>&1); then
        exit_code=0
    else
        exit_code=$?
    fi
    printf '%s\n' "$output" | grep -E "^(=== RUN|--- PASS|--- FAIL|--- SKIP|FAIL\t|    .*(Error|garmin))" || true
    if printf '%s\n' "$output" | grep -q "^--- SKIP"; then
        SKIP+=("$test")
    elif [[ $exit_code -eq 0 ]]; then
        PASS+=("$test")
    else
        FAIL+=("$test")
        echo "    ^^^ FAILED"
    fi

    # Sleep between tests to avoid Connect API rate-limiting (not SSO — we only login once).
    if [[ $n -lt $total ]]; then
        sleep "$DELAY"
    fi
done

echo ""
echo "==> Sanitizing cassettes..."
sanitize_args=()
[[ -n "${GARMIN_DISPLAY_NAME:-}" ]] && sanitize_args+=(--display-name "$GARMIN_DISPLAY_NAME")
[[ -n "${GARMIN_EMAIL:-}" ]]        && sanitize_args+=(--email "$GARMIN_EMAIL")
python3 tools/sanitize_cassettes.py "${sanitize_args[@]}"

echo ""
echo "=== SUMMARY ==="
echo "PASS (${#PASS[@]}): ${PASS[*]:-none}"
echo "FAIL (${#FAIL[@]}): ${FAIL[*]:-none}"
echo "SKIP (${#SKIP[@]}): ${SKIP[*]:-none}"

[[ ${#FAIL[@]} -eq 0 ]]
