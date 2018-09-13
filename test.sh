#!/bin/bash

if ! pgrep -f 'go-state-channel' > /dev/null; then
    echo "Please run go-state-channel"
    exit 2
fi

resp="responses-$$"

for _ in {1..100}; do
    curl --silent http://127.0.0.1:8080 >> "$resp" &
done

pid=$(pgrep -f 'go-state-channel')
kill -s HUP "$pid"

wait

grep 'non-success' "$resp" > /dev/null
fail=$?

sort "$resp" | uniq -c
rm "$resp"

if [[ "$fail" -eq 0 ]]; then
    exit 1
fi
