#!/bin/bash
# dispatch.sh - Main entry point for parallel agent execution
# Usage: ./scripts/parallel/dispatch.sh "task description" "agent1 prompt" "agent2 prompt" ["agent3 prompt"]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TASK_DESC="$1"
shift
PROMPTS=("$@")

# Validate inputs
if [ -z "$TASK_DESC" ]; then
  echo "Error: Task description required"
  echo "Usage: $0 \"task description\" \"agent1 prompt\" \"agent2 prompt\" [\"agent3 prompt\"]"
  exit 1
fi

if [ ${#PROMPTS[@]} -lt 2 ] || [ ${#PROMPTS[@]} -gt 3 ]; then
  echo "Error: Need 2-3 agent prompts (got ${#PROMPTS[@]})"
  exit 1
fi

echo "=== Parallel Agent Dispatcher ==="
echo "Task: $TASK_DESC"
echo "Agents: ${#PROMPTS[@]}"
echo ""

# Step 1: Commit current work
echo "[1/5] Committing current work..."
if ! git diff-index --quiet HEAD --; then
  git add .
  git commit -m "wip: before parallel dispatch - $TASK_DESC"
  echo "✓ Changes committed"
else
  echo "✓ No uncommitted changes"
fi
echo ""

# Step 2: Create worktrees
echo "[2/5] Creating agent worktrees..."
"$SCRIPT_DIR/create-agents.sh" ${#PROMPTS[@]}
echo ""

# Step 3: Run agents in parallel
echo "[3/5] Dispatching agents..."
"$SCRIPT_DIR/run-parallel.sh" "${PROMPTS[@]}"
echo ""

# Step 4: Review changes
echo "[4/5] Reviewing changes..."
"$SCRIPT_DIR/review.sh"
echo ""

# Step 5: Prompt for next action
echo "[5/5] Worktrees ready"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "What would you like to do?"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  1. Let Claude Code review and merge worktrees automatically"
echo "  2. Review changes manually first (./scripts/parallel/review.sh)"
echo "  3. Merge now (./scripts/parallel/merge.sh)"
echo "  4. Exit and work on worktrees manually"
echo ""
read -p "Choose [1-4]: " -n 1 -r
echo
echo ""

case $REPLY in
  1)
    echo "🤖 Claude Code will now review and merge the worktrees..."
    echo ""
    echo "Worktrees available for review:"
    for dir in .worktrees/agent-*; do
      if [ -d "$dir" ]; then
        BRANCH=$(git -C "$dir" branch --show-current)
        echo "  • $dir (branch: $BRANCH)"
      fi
    done
    echo ""
    echo "💡 Claude Code: Please review each worktree and merge when ready"
    ;;
  2)
    "$SCRIPT_DIR/review.sh"
    echo ""
    read -p "Proceed with merge? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
      "$SCRIPT_DIR/merge.sh"
      echo ""
      read -p "Cleanup worktrees? (y/n) " -n 1 -r
      echo
      if [[ $REPLY =~ ^[Yy]$ ]]; then
        "$SCRIPT_DIR/cleanup.sh"
      fi
    fi
    ;;
  3)
    "$SCRIPT_DIR/merge.sh"
    echo ""
    read -p "Cleanup worktrees? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
      "$SCRIPT_DIR/cleanup.sh"
    fi
    ;;
  4)
    echo "Worktrees preserved. Use these commands:"
    echo "  • Monitor: ./scripts/parallel/monitor.sh"
    echo "  • Review: ./scripts/parallel/review.sh"
    echo "  • Merge: ./scripts/parallel/merge.sh"
    echo "  • Cleanup: ./scripts/parallel/cleanup.sh"
    ;;
  *)
    echo "Invalid choice. Exiting."
    ;;
esac
