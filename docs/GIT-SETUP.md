# Git setup and contribution evidence

See the [24 October checklist](00-checklist-pre-defense-1.md) for the remaining submission evidence. The supplied folder has no `.git` directory or remote URL, so its current files cannot establish when or by whom work was done.

The [local backlog](../project/youtrack-backlog.csv) assigns 15 planned tasks to each member. `scripts/tasks.tsv` maps task IDs to likely files; it is an ownership plan, not an authorship record. Some paths refer to future implementation and must be created before a task can be committed.

Publish the initial snapshot once, under the real Git identity of the person who imports it. Both members should then develop their actual changes on separate task branches with meaningful commits. `./scripts/contribute.sh DIP-n` can assist only after the author has made and reviewed the relevant changes; verify its file list before running it. If the live YouTrack issue ID differs from the local `DIP-n` label, use the live ID and update the mapping.

The other member reviews each PR/MR before merge. Record real review comments and link the PR or commit in the matching YouTrack issue. Tag `pre-defense-1` only after the submitted version is final and the tag can be verified on the remote. Never change commit dates or credit a team member with work they did not perform.
