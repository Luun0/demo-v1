#!/bin/sh
set -eu

/api serve &
api_pid=$!
/docker-entrypoint.sh nginx -g 'daemon off;' &
nginx_pid=$!

stop_children() {
  kill -TERM "$api_pid" "$nginx_pid" 2>/dev/null || true
  wait "$api_pid" 2>/dev/null || true
  wait "$nginx_pid" 2>/dev/null || true
}
trap 'stop_children; exit 0' INT TERM

while kill -0 "$api_pid" 2>/dev/null && kill -0 "$nginx_pid" 2>/dev/null; do
  sleep 1
done
stop_children
exit 1
