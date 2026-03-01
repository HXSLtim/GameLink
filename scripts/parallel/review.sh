#!/bin/bash
# review.sh - Generate comprehensive diff summary for all agent worktrees
# Usage: ./scripts/parallel/review.sh

set -e

if [ ! -d ".worktrees" ]; then
  echo "No worktrees found. Run create-agents.sh first."
  exit 1
fi

BASE_BRANCH=$(git branch --show-current)

echo "=== Agent Changes Review ==="
echo "Base branch: $BASE_BRANCH"
echo ""

# Collect all worktrees for review
WORKTREES=()
for dir in .worktrees/agent-*; do
  if [ -d "$dir" ]; then
    WORKTREES+=("$dir")
  fi
done

if [ ${#WORKTREES[@]} -eq 0 ]; then
  echo "No agent worktrees found"
  exit 0
fi

# Review each worktree
for dir in "${WORKTREES[@]}"; do
  AGENT=$(basename "$dir")
  BRANCH=$(git -C "$dir" branch --show-current 2>/dev/null || echo "unknown")

  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "[$AGENT] Branch: $BRANCH"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo ""

  # Show prompt from log
  if [ -f "$dir/agent.log" ]; then
    echo "Task:"
    grep "^Prompt:" "$dir/agent.log" | sed 's/^Prompt: /  /'
    echo ""
  fi

  # Check if there are any changes
  if git -C "$dir" diff-index --quiet HEAD -- 2>/dev/null && \
     [ "$(git -C "$dir" rev-list --count HEAD ^"$BASE_BRANCH" 2>/dev/null || echo 0)" -eq 0 ]; then
    echo "  ⚠ No changes detected"
    echo ""
    continue
  fi

  # Show diff stats
  echo "Changes:"
  git -C "$dir" diff --stat "$BASE_BRANCH" 2>/dev/null | sed 's/^/  /' || echo "  (unable to generate diff)"
  echo ""

  # Show commit log if any
  COMMITS=$(git -C "$dir" rev-list --count HEAD ^"$BASE_BRANCH" 2>/dev/null || echo "0")
  if [ "$COMMITS" -gt 0 ]; then
    echo "Commits ($COMMITS):"
    git -C "$dir" log --oneline "$BASE_BRANCH..HEAD" 2>/dev/null | sed 's/^/  /' || true
    echo ""
  fi

  # Show file list
  echo "Modified files:"
  git -C "$dir" diff --name-only "$BASE_BRANCH" 2>/dev/null | sed 's/^/  /' || echo "  (none)"
  echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Review Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Worktrees ready for merge:"
for dir in "${WORKTREES[@]}"; do
  AGENT=$(basename "$dir")
  BRANCH=$(git -C "$dir" branch --show-current 2>/dev/null || echo "unknown")
  echo "  • $AGENT → $BRANCH"
done
echo ""
echo "Next steps:"
echo "  • Merge all: ./scripts/parallel/merge.sh"
echo "  • View detailed diff: git -C .worktrees/agent-N diff $BASE_BRANCH"
echo "  • Cancel and cleanup: ./scripts/parallel/cleanup.sh"
