# Contributing

The supplied folder has no Git history or live YouTrack link. The local `DIP-n` numbers in `project/youtrack-backlog.csv` are a plan; use the IDs assigned by the real YouTrack project after import. Publish the initial snapshot once under the real identity of the person who performs that import. Do not create artificial commits or assign past work to either member.

For subsequent work:

1. Take an assigned YouTrack issue with an expected outcome, acceptance criteria and deadline. Use your own Git account and configured real name and e-mail.
2. Branch from the current `main` using the real issue ID: `DIP-n-short-description`.
3. Commit one coherent change at a time with a task-linked message. The optional `scripts/contribute.sh` helper requires an existing remote `main` branch.
4. Open a pull or merge request. Describe what changed and how it was checked; link the real issue.
5. Ask the other member to review before merge. Enable branch protection on `main` when the remote repository is configured.
6. Add the merged PR/MR link to the issue and mark it Done only when its acceptance criteria are met.
7. Create and push checkpoint tags (`pre-defense-1`, `pre-defense-2`, `final-pre-defense`) only after verifying the corresponding submitted state. Do not move published tags.
