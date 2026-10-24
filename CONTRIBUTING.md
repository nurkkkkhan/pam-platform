# Contributing

The [GitHub repository](https://github.com/nurkkkkhan/pam-platform) has an initial snapshot and a links update, both committed on 24 October 2026 under one Git author identity. The [YouTrack project URL](https://diploma-work.youtrack.cloud/projects/DP) is recorded, but the issue import and assignments have not been verified. The local `DIP-n` numbers in `project/youtrack-backlog.csv` are planning IDs; use each issue's actual key after checking the live project. Do not create artificial commits or assign past work to either member.

For subsequent work:

1. Take an assigned YouTrack issue with an expected outcome, acceptance criteria and deadline. Use your own Git account and configured real name and e-mail.
2. Branch from the current `main` using the real issue key: `<issue-key>-short-description`.
3. Commit one coherent change at a time with a task-linked message. The optional `scripts/contribute.sh` helper requires an existing remote `main` branch.
4. Open a pull or merge request. Describe what changed and how it was checked; link the real issue.
5. Ask the other member to review before merge. Configure branch protection on the published `main` branch.
6. Add the merged PR/MR link to the issue and mark it Done only when its acceptance criteria are met.
7. Create and push checkpoint tags (`pre-defense-1`, `pre-defense-2`, `final-pre-defense`) only after verifying the corresponding submitted state. Do not move published tags.
