#!/bin/bash
# monitor.sh - Monitor progress of all agent worktrees
# Usage: ./scripts/parallel/monitor.sh

set -e

if [ ! -d ".worktrees" ]; then
  echo "No worktrees found. Run create-agents.sh first."
  exit 1
fi

echo "=== Agent Worktree Status ==="
echo ""

for dir in .worktrees/agent-*; do
  if [ ! -d "$dir" ]; then
    continue
  fi

  AGENT=$(basename "$dir")
  BRANCH=$(git -C "$dir" branch --show-current 2>/dev/null || echo "unknown")

  echo "[$AGENT] Branch: $BRANCH"

  # Check for uncommitted changes
  if git -C "$dir" diff-index --quiet HEAD -- 2>/dev/null; then
    echo "  Status: No changes yet"
  else
    echo "  Status: Has uncommitted changes"
    git -C "$dir" status --short | head -10 | sed 's/^/    /'
  fi

  # Check for commits ahead of base
  BASE_BRANCH=$(git branch --show-current)
  COMMITS_AHEAD=$(git -C "$dir" rev-list --count HEAD ^"$BASE_BRANCH" 2>/dev/null || echo "0")
  if [ "$COMMITS_AHEAD" -gt 0 ]; then
    echo "  Commits: $COMMITS_AHEAD ahead of $BASE_BRANCH"
  fi

  # Show log file if exists
  if [ -f "$dir/agent.log" ]; then
    echo "  Log: $dir/agent.log"
  fi

  echo ""
done

echo "Commands:"
echo "  Review changes: ./scripts/parallel/review.sh"
echo "  Merge results: ./scripts/parallel/merge.sh"
echo "  Cleanup: ./scripts/parallel/cleanup.sh"
