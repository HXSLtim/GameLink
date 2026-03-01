#!/bin/bash
# merge.sh - Intelligently merge all agent branches with conflict detection
# Usage: ./scripts/parallel/merge.sh [--dry-run]

set -e

DRY_RUN=false
if [ "$1" == "--dry-run" ]; then
  DRY_RUN=true
  echo "=== Dry Run Mode - No actual merges ==="
  echo ""
fi

if [ ! -d ".worktrees" ]; then
  echo "No worktrees found. Nothing to merge."
  exit 0
fi

BASE_BRANCH=$(git branch --show-current)

# Collect all branches to merge
BRANCHES=()
for dir in .worktrees/agent-*; do
  if [ -d "$dir" ]; then
    BRANCH=$(git -C "$dir" branch --show-current 2>/dev/null || echo "")
    if [ -n "$BRANCH" ]; then
      BRANCHES+=("$BRANCH")
    fi
  fi
done

if [ ${#BRANCHES[@]} -eq 0 ]; then
  echo "No branches to merge"
  exit 0
fi

echo "=== Merge Plan ==="
echo "Base: $BASE_BRANCH"
echo "Branches to merge:"
for branch in "${BRANCHES[@]}"; do
  echo "  • $branch"
done
echo ""

# Dry run: check for conflicts
if $DRY_RUN; then
  echo "Checking for potential conflicts..."
  echo ""

  for branch in "${BRANCHES[@]}"; do
    echo "Testing merge: $branch"
    if git merge --no-commit --no-ff "$branch" 2>&1 | grep -qi "conflict"; then
      echo "  ⚠ CONFLICTS DETECTED"
      git merge --abort 2>/dev/null || true
    else
      echo "  ✓ No conflicts"
      git merge --abort 2>/dev/null || true
    fi
  done

  echo ""
  echo "Dry run complete. Run without --dry-run to perform actual merge."
  exit 0
fi

# Actual merge
echo "Starting sequential merge..."
echo ""

MERGED=0
FAILED=0

for branch in "${BRANCHES[@]}"; do
  echo "Merging: $branch"

  # Attempt merge
  if git merge "$branch" --no-edit -m "feat: merge $branch" 2>&1; then
    echo "  ✓ Merged successfully"
    MERGED=$((MERGED + 1))
  else
    echo "  ✗ MERGE FAILED - Conflicts detected"
    echo ""
    echo "Conflicting files:"
    git status --short | grep "^UU" | sed 's/^/    /'
    echo ""
    echo "To resolve:"
    echo "  1. Edit conflicting files manually"
    echo "  2. git add <resolved-files>"
    echo "  3. git commit -m 'feat: merge $branch with conflict resolution'"
    echo "  4. Re-run this script to continue merging remaining branches"
    FAILED=$((FAILED + 1))
    exit 1
  fi
  echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Merge Complete"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  ✓ Successfully merged: $MERGED branches"
if [ $FAILED -gt 0 ]; then
  echo "  ✗ Failed: $FAILED branches"
fi
echo ""
echo "Next steps:"
echo "  • Run tests: make test (backend) or pnpm test (frontend)"
echo "  • Cleanup worktrees: ./scripts/parallel/cleanup.sh"
echo "  • Push changes: git push origin $BASE_BRANCH"
