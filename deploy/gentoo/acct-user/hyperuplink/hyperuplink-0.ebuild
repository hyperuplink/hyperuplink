# Copyright 2025 Hyperuplink Authors
# Distributed under the terms of the SEGV License, Version 1.1

EAPI=8

inherit acct-user

DESCRIPTION="User for the Hyperuplink bulletin board service"

# -1 lets Portage allocate a free UID dynamically.
ACCT_USER_ID=-1
ACCT_USER_GROUPS=( hyperuplink )
ACCT_USER_HOME=/var/lib/hyperuplink
ACCT_USER_SHELL=/sbin/nologin

acct-user_add_deps
