#!/usr/bin/env bash
# Sobe o servidor do my_dbms e a interface web. Ctrl+C derruba os dois.
#
#   ./dev.sh
#   DB_DIR=../outro_my_dbms DB_ADDR=localhost:5000 WEB_ADDR=localhost:9000 ./dev.sh
set -euo pipefail

cd "$(dirname "$0")"
DB_DIR="${DB_DIR:-../my_dbms}"
DB_ADDR="${DB_ADDR:-localhost:4321}"
WEB_ADDR="${WEB_ADDR:-localhost:8080}"

bin="$(mktemp -d)"
trap 'kill $(jobs -p) 2>/dev/null; rm -rf "$bin"' EXIT

# Compila antes, para não depender de quanto tempo o go run demora e para
# sempre rodar o código atual.
(cd "$DB_DIR" && go build -o "$bin/my_dbms" ./cmd/my_dbms)
go build -o "$bin/web" .

# O banco roda na pasta dele, para usar o my_dbms.conf e o data/ de lá.
(cd "$DB_DIR" && exec "$bin/my_dbms" -listen "$DB_ADDR") &
sleep 0.5
if ! kill -0 $! 2>/dev/null; then
	echo "o my_dbms não subiu: veja o erro acima" >&2
	exit 1
fi

"$bin/web" -addr "$WEB_ADDR" -db "$DB_ADDR"
