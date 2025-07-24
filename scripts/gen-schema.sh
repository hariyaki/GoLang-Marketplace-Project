#!/usr/bin/env bash
schema=ProjectSchema.txt
ignore='node_modules|dist|uploads|docs/swagger.*|\.git|\.idea|\.vscode|\.DS_Store'
tree -a -F --noreport -I "$ignore" > "$schema"
echo "✓ ProjectSchema.txt updated (Unix)"