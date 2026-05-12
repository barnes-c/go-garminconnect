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

: "${GARMIN_EMAIL:?GARMIN_EMAIL must be set}"
: "${GARMIN_PASSWORD:?GARMIN_PASSWORD must be set}"

CASSETTE_DIR="garminconnect/tests/testdata/cassettes"
DELAY=5   # seconds between tests to avoid Connect API rate-limiting

# Cassettes to leave untouched.
KEEP=(
    "user_summary"
    "all_day_stress"
    "activities_empty"
)

keep() {
    local name="$1"
    for k in "${KEEP[@]}"; do
        [[ "$name" == "$k" ]] && return 0
    done
    return 1
}

echo "==> Logging in (once)..."
token_line=$(go run ./cmd/gettoken)
export GARMIN_TOKEN
export GARMIN_DISPLAY_NAME
GARMIN_TOKEN=$(printf '%s' "$token_line" | sed -n '1p')
GARMIN_DISPLAY_NAME=$(printf '%s' "$token_line" | sed -n '2p')
echo "    display_name=${GARMIN_DISPLAY_NAME}"

echo ""
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

# One test per unique cassette.
# Tests that share a cassette with another test are omitted (they replay).
TESTS=(
    TestActivities
    TestActivityDetail
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

    if go test ./garminconnect/tests/... -run "^${test}$" -count=1 -v 2>&1 \
        | grep -E "^(=== RUN|--- PASS|--- FAIL|--- SKIP|FAIL\t|    .*(Error|garmin))"; then
        PASS+=("$test")
    else
        exit_code=${PIPESTATUS[0]}
        if [[ $exit_code -eq 0 ]]; then
            SKIP+=("$test")
        else
            FAIL+=("$test")
            echo "    ^^^ FAILED"
        fi
    fi

    # Sleep between tests to avoid Connect API rate-limiting (not SSO — we only login once).
    if [[ $n -lt $total ]]; then
        sleep "$DELAY"
    fi
done

echo ""
echo "=== SUMMARY ==="
echo "PASS (${#PASS[@]}): ${PASS[*]:-none}"
echo "FAIL (${#FAIL[@]}): ${FAIL[*]:-none}"
echo "SKIP (${#SKIP[@]}): ${SKIP[*]:-none}"

[[ ${#FAIL[@]} -eq 0 ]]
