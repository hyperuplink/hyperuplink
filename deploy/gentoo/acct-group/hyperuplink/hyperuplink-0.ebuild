# Copyright 2025 Hyperuplink Authors
# Distributed under the terms of the SEGV License, Version 1.0

EAPI=8

inherit acct-group

DESCRIPTION="Group for the Hyperuplink bulletin board service"

# -1 lets Portage allocate a free GID dynamically.
ACCT_GROUP_ID=-1
