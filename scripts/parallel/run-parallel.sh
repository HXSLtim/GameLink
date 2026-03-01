#!/bin/bash
# run-parallel.sh - Execute agents in parallel with their prompts
# Usage: ./scripts/parallel/run-parallel.sh "prompt1" "prompt2" ["prompt3"]

set -e

PROMPTS=("$@")
NUM_AGENTS=${#PROMPTS[@]}

if [ $NUM_AGENTS -lt 2 ] || [ $NUM_AGENTS -gt 3 ]; then
  echo "Error: Need 2-3 prompts (got: $NUM_AGENTS)"
  exit 1
fi

echo "Starting $NUM_AGENTS agents in parallel..."
echo ""

# Store PIDs for monitoring
PIDS=()

for i in "${!PROMPTS[@]}"; do
  AGENT_NUM=$((i + 1))
  WORKTREE=".worktrees/agent-$AGENT_NUM"
  PROMPT="${PROMPTS[$i]}"
  LOG_FILE="$WORKTREE/agent.log"

  if [ ! -d "$WORKTREE" ]; then
    echo "Error: Worktree $WORKTREE does not exist"
    exit 1
  fi

  echo "Agent $AGENT_NUM: Starting in $WORKTREE"
  echo "  Prompt: $PROMPT"

  # Run agent in background
  # Note: This creates a marker file that Claude Code can detect
  (
    cd "$WORKTREE"
    echo "=== Agent $AGENT_NUM Started at $(date) ===" > "$LOG_FILE"
    echo "Prompt: $PROMPT" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"
    echo "Context: Read CLAUDE.md for project patterns and conventions" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"
    echo "Full Task: $PROMPT" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"
    echo "Status: READY_FOR_CLAUDE_CODE" >> "$LOG_FILE"
    echo "Worktree: $WORKTREE" >> "$LOG_FILE"
    echo "Branch: $(git branch --show-current)" >> "$LOG_FILE"
  ) &

  PIDS+=($!)
done

echo ""
echo "✓ All agents initialized"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Agent Worktrees Ready for Implementation"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Output worktree info in a structured format for Claude Code to parse
for i in $(seq 1 $NUM_AGENTS); do
  WORKTREE=".worktrees/agent-$i"
  BRANCH=$(git -C "$WORKTREE" branch --show-current)
  echo "Agent $i:"
  echo "  Path: $WORKTREE"
  echo "  Branch: $BRANCH"
  echo "  Task: ${PROMPTS[$((i-1))]}"
  echo "  Status: PENDING_IMPLEMENTATION"
  echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📋 Next Steps:"
echo "  1. Implement tasks in each worktree (manually or with Claude Code)"
echo "  2. Monitor progress: ./scripts/parallel/monitor.sh"
echo "  3. Review changes: ./scripts/parallel/review.sh"
echo "  4. Merge results: ./scripts/parallel/merge.sh"
echo ""
echo "💡 Tip: Claude Code can now review and merge these worktrees automatically"
