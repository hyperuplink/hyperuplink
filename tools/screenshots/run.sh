#!/usr/bin/env bash
#
# Regenerates the screenshots of the board, both the set embedded in the manual
# and the set the website's landing page shows. The board it shoots is
# disposable and lives in its own database, on its own port and its own Redis
# database index, so re-running this never touches the development board.
#
# Everything after the script's name is handed to the shooter, so `-set manual`
# or `-set site` picks one set and `-only` picks one shot out of it.

set -euo pipefail

cd "$(dirname "$0")/../.."

ROOT="$PWD"
CFG="file://$ROOT/tools/screenshots/hyperuplink.toml"
BIN="$ROOT/build/hyperuplink"
DB="hyperuplink_manual"
BOARD="http://127.0.0.1:3100"
PASSWORD='hyperhyper!'
PSQL=(psql -h localhost -p 5432 -U postgres -v ON_ERROR_STOP=1 -q)

echo "==> building"
make build >/dev/null

echo "==> recreating $DB"
"${PSQL[@]}" -d postgres -c "DROP DATABASE IF EXISTS $DB;"
"${PSQL[@]}" -d postgres -c "CREATE DATABASE $DB;"

echo "==> seeding"
tools/seed/seed.sh "$CFG" "$DB"
"${PSQL[@]}" -d "$DB" -f tools/screenshots/seed.sql

echo "==> starting the board"
LOG="$(mktemp -t hyperuplink-screenshots-XXXXXX.log)"
"$BIN" -c "$CFG" >"$LOG" 2>&1 &
SERVER=$!
echo "    board log: $LOG"
trap 'kill "$SERVER" 2>/dev/null || true' EXIT

for _ in $(seq 1 60); do
  if curl -sf -o /dev/null "$BOARD/"; then
    break
  fi
  sleep 0.25
done

echo "==> shooting"
cd tools/screenshots
go run . -board "$BOARD" -password "$PASSWORD" "$@"
