#!/bin/bash
# create-agents.sh - Create N isolated git worktrees for parallel agents
# Usage: ./scripts/parallel/create-agents.sh <num_agents>

set -e

NUM_AGENTS=$1
BASE_BRANCH=$(git branch --show-current)
TIMESTAMP=$(date +%s)

# Validate input
if [ -z "$NUM_AGENTS" ] || [ "$NUM_AGENTS" -lt 2 ] || [ "$NUM_AGENTS" -gt 3 ]; then
  echo "Error: Need 2-3 agents (got: $NUM_AGENTS)"
  exit 1
fi

# Check if we're in a git repo
if ! git rev-parse --git-dir > /dev/null 2>&1; then
  echo "Error: Not in a git repository"
  exit 1
fi

# Create .worktrees directory if it doesn't exist
mkdir -p .worktrees

echo "Creating $NUM_AGENTS agent worktrees from branch: $BASE_BRANCH"
echo ""

for i in $(seq 1 $NUM_AGENTS); do
  BRANCH="parallel/agent-$i-$TIMESTAMP"
  WORKTREE=".worktrees/agent-$i"

  # Remove existing worktree if it exists
  if [ -d "$WORKTREE" ]; then
    echo "⚠ Removing existing worktree: $WORKTREE"
    git worktree remove "$WORKTREE" --force 2>/dev/null || true
  fi

  # Create new worktree
  git worktree add "$WORKTREE" -b "$BRANCH" > /dev/null 2>&1

  # Verify CLAUDE.md is accessible
  if [ -f "$WORKTREE/CLAUDE.md" ]; then
    echo "✓ Agent $i: $WORKTREE (branch: $BRANCH) - CLAUDE.md accessible"
  else
    echo "⚠ Agent $i: $WORKTREE (branch: $BRANCH) - CLAUDE.md not found"
  fi
done

echo ""
echo "✓ All worktrees created successfully"
