#!/usr/bin/env bash
# Run from docgen/ with default sibling layout: ../dsa, ../dsa-pavilion
# Or set DOCGEN_CODE, DOCGEN_DOCS, DOCGEN_SIDEBAR, DOCGEN_README, DOCGEN_BASE
set -euo pipefail
cd "$(dirname "$0")"
go run .
