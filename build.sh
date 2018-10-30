#! /bin/bash
set -e

OUT_FILE="unidoor.new"

echo '≫ Buidling "'${OUT_FILE}'" executable...'
GOOS=linux GOARCH=arm GOARM=6 go build -o "$OUT_FILE"
