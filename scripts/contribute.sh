#!/usr/bin/env bash
# Publish a subsequent task's changes from the real author's Git account.
# Requires an existing remote repository with main and real YouTrack issues.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
TASKS=scripts/tasks.tsv

if [[ "${1:-}" == "--mine" ]]; then
  MEMBER="${2:?usage: contribute.sh --mine 'Member 1' or 'Member 2'}"
  awk -F'\t' -v m="$MEMBER" '!/^#/ && $2==m {print $1 "  " $4}' "$TASKS"
  exit 0
fi
TASK="${1:?usage: contribute.sh DIP-n}"
LINE=$(awk -F'\t' -v t="$TASK" '!/^#/ && $1==t' "$TASKS")
[[ -n "$LINE" ]] || { echo "Task $TASK not in $TASKS"; exit 1; }
FILES=$(cut -f3 <<<"$LINE"); MSG=$(cut -f4 <<<"$LINE")

[[ -n "$(git config --get user.name || true)" && -n "$(git config --get user.email || true)" ]] || { echo "Set git user.name and user.email to your real identity first."; exit 1; }
echo "Author: $(git config user.name) <$(git config user.email)>"

git remote get-url origin >/dev/null || { echo "Add the real origin remote first."; exit 1; }
git fetch -q origin main
[[ "$(git branch --show-current)" == "main" ]] || { echo "Start from main."; exit 1; }
[[ "$(git rev-parse HEAD)" == "$(git rev-parse origin/main)" ]] || { echo "Update local main before creating a task branch."; exit 1; }
git diff --cached --quiet || { echo "Unstage unrelated changes before running this helper."; exit 1; }

SLUG=$(sed -E 's/^DIP-[0-9]+: //; s/[^A-Za-z0-9]+/-/g; s/-+$//' <<<"$MSG" | tr 'A-Z' 'a-z' | cut -c1-40 | sed 's/-$//')
BRANCH="$TASK-$SLUG"
BASE=main
git switch -q -c "$BRANCH"
read -r -a PATHS <<< "$FILES"
git add -- "${PATHS[@]}"
if git diff --cached --quiet; then echo "Nothing to commit for $TASK (already committed?)"; exit 0; fi
git commit -q -m "$MSG" -m "YouTrack: $TASK"
git push -q -u origin "$BRANCH"
URL=$(git remote get-url origin | sed -E 's#git@([^:]+):#https://\1/#; s#\.git$##')
echo
echo "Pushed $BRANCH. Open a PR/MR into $BASE, ask a teammate to review, then merge:"
case "$URL" in
  *github.com*) echo "  $URL/compare/$BASE...$BRANCH?expand=1" ;;
  *) echo "  $URL/-/merge_requests/new?merge_request[source_branch]=$BRANCH" ;;
esac
echo "Then paste the PR link into YouTrack task $TASK."
