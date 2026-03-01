#!/bin/bash
# cleanup.sh - Remove all agent worktrees and branches
# Usage: ./scripts/parallel/cleanup.sh [--force]

set -e

FORCE=false
if [ "$1" == "--force" ]; then
  FORCE=true
fi

if [ ! -d ".worktrees" ]; then
  echo "No worktrees directory found. Nothing to cleanup."
  exit 0
fi

# Collect all worktrees
WORKTREES=()
BRANCHES=()

for dir in .worktrees/agent-*; do
  if [ -d "$dir" ]; then
    WORKTREES+=("$dir")
    BRANCH=$(git -C "$dir" branch --show-current 2>/dev/null || echo "")
    if [ -n "$BRANCH" ]; then
      BRANCHES+=("$BRANCH")
    fi
  fi
done

if [ ${#WORKTREES[@]} -eq 0 ]; then
  echo "No agent worktrees found"
  exit 0
fi

echo "=== Cleanup Agent Worktrees ==="
echo ""
echo "Worktrees to remove:"
for dir in "${WORKTREES[@]}"; do
  echo "  • $dir"
done
echo ""
echo "Branches to delete:"
for branch in "${BRANCHES[@]}"; do
  echo "  • $branch"
done
echo ""

# Confirm unless --force
if ! $FORCE; then
  read -p "Continue with cleanup? (y/n) " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cleanup cancelled"
    exit 0
  fi
fi

# Remove worktrees
echo "Removing worktrees..."
for dir in "${WORKTREES[@]}"; do
  AGENT=$(basename "$dir")
  if git worktree remove "$dir" --force 2>/dev/null; then
    echo "  ✓ Removed $AGENT"
  else
    echo "  ⚠ Failed to remove $AGENT"
  fi
done
echo ""

# Delete branches
echo "Deleting branches..."
for branch in "${BRANCHES[@]}"; do
  if git branch -D "$branch" 2>/dev/null; then
    echo "  ✓ Deleted $branch"
  else
    echo "  ⚠ Failed to delete $branch (may not exist or already deleted)"
  fi
done
echo ""

# Remove .worktrees directory if empty
if [ -d ".worktrees" ] && [ -z "$(ls -A .worktrees)" ]; then
  rmdir .worktrees
  echo "✓ Removed empty .worktrees directory"
fi

echo ""
echo "✓ Cleanup complete"
