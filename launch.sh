#!/usr/bin/env bash
set -euo pipefail

[ -x build/ai-terminal-assistant ] || make build-daemon
[ -x ai-terminal-gui ] || go build -o ai-terminal-gui ./cmd/gui

./build/ai-terminal-assistant -port 8080 &
DAEMON_PID=$!
./ai-terminal-gui &
GUI_PID=$!

cleanup() {
  kill "$DAEMON_PID" "$GUI_PID"
}
trap cleanup EXIT

wait