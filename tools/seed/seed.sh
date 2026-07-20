#!/usr/bin/env bash
#
# Seeds a board with a cast and something for them to have said, so that a
# database dropped after a migration was edited in place can be put back with
# one command.
#
# Use: tools/seed/seed.sh <config-url> <database>
#
# The order matters. Creating the first user is what runs the migrations and
# seeds the settings rows, the settings then have to be in place before the
# rest of the users are created, since the address type is what decides whether
# a JID is stored as a JID, and the content comes last because every post
# points at a user by username.

set -euo pipefail

cd "$(dirname "$0")/../.."

CFG="${1:-file://$PWD/hyperuplink.toml}"
DB="${2:-hyperuplink_dev}"
BIN="$PWD/build/hyperuplink"
PASSWORD='hyperhyper!'
PSQL=(psql -h localhost -p 5432 -U postgres -v ON_ERROR_STOP=1 -q)

if [ ! -x "$BIN" ]; then
  echo "$BIN is missing, run make build first" >&2
  exit 1
fi

create_user() {
  "$BIN" -c "$CFG" -create-user \
    "{\"username\":\"$1\",\"email\":\"$2\",\"password\":\"$PASSWORD\",\"email_is_jid\":${3:-false}}" \
    >/dev/null
}

# An address in PromoteAdmin lands as an admin on its own.
create_user sysop dummy1@example.com

"${PSQL[@]}" -d "$DB" -f tools/seed/settings.sql

create_user vera vera@example.org
create_user juno juno@example.org
create_user bitrot bitrot@jabber.example.org true

"${PSQL[@]}" -d "$DB" -f tools/seed/content.sql

echo "Seeded $DB: sysop (admin), vera, juno, bitrot (JID), all with $PASSWORD"
