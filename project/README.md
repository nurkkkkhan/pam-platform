# Local backlog for a two-member team

`youtrack-backlog.csv` contains 30 planned work items, split evenly: 15 assigned to Member 1 (Backend) and 15 to Member 2 (Frontend and UI). Member 2 also owns project coordination, research and documentation tasks; Member 1 owns backend implementation, the OpenAPI contract, gateways, CI and performance testing. Member 2 leads the evaluation report using Member 1's security and performance evidence. The `DIP-1`…`DIP-30` labels are local planning IDs. A real YouTrack import may assign different issue IDs; update links and branch names to the actual IDs after creation.

Individual responsibility and handoff plans: [Member 1 — Backend] and [Member 2 — Frontend and UI].

No YouTrack project URL, account names, or import result was supplied with this folder. The CSV is a planning snapshot, not evidence that issues exist in YouTrack. Fill the two real identities and YouTrack logins in `team.conf` before importing, then verify assignment, priority, state, and due dates in the live project.

`In Progress` means a local draft or initial code exists but review or acceptance is still pending. `Open` means the acceptance criteria have not been demonstrated. Mark `Done` only after those criteria are met; do not imply a supervisor review, tag, remote CI run, or repository access without evidence.

The milestone dates in `milestones.csv` are the proposed task deadlines. Only the 24 October 2026 checkpoint is established by the supplied requirements; confirm the later dates with the supervisor or academic calendar. The import tool uses milestone dates as YouTrack due dates.
