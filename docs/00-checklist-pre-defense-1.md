# Minimum requirements for 24 October 2026

This checklist distinguishes files present in the supplied folder from evidence that requires real people and external services. A local draft does not establish supervisor review, a YouTrack project, or individual Git authorship.

| # | Requirement | Local material | Status as supplied |
|---|---|---|---|
| 1 | Completed project README | [`README.md`](../README.md) covers title, problem, users, outcomes, two-member responsibilities, scope, technologies and setup. | Partial: real member names, supervisor, repository URL and YouTrack URL remain to be supplied. |
| 2 | Project proposal and scope | [`01-project-proposal.md`](01-project-proposal.md) | Draft present; supervisor review not evidenced. |
| 3a | Literature review | [`03-literature-review.md`](03-literature-review.md) includes a source table. | Draft present; source quality and citation checks remain. |
| 3b | Existing system comparison | [`04-existing-systems-comparison.md`](04-existing-systems-comparison.md) compares five systems. | Draft present; unresolved `?` cells need evidence. |
| 3c | Technology comparison | [`05-technology-comparison.md`](05-technology-comparison.md) and [`adr/`](adr/) | Drafts present. |
| 4 | Initial requirements and acceptance criteria | [`06-requirements-specification.md`](06-requirements-specification.md) and [`tests/test-cases.md`](../tests/test-cases.md) | Drafts present; validate against user evidence. |
| 5 | Architectural diagrams and editable sources | [`../diagrams/src/`](../diagrams/src/) DOT files, [`../diagrams/export/`](../diagrams/export/) exports, and [`07-architecture.md`](07-architecture.md) | Context, use-case and draft component diagrams present locally. |
| 6 | Prototype materials or link | [`../prototype/index.html`](../prototype/index.html), [`08-ux-prototype.md`](08-ux-prototype.md), [`../api/openapi.yaml`](../api/openapi.yaml) | Local prototype present; no hosted link needed if the repository provides the file. |
| 7 | Verification and evaluation plan | [`09-verification-and-evaluation-plan.md`](09-verification-and-evaluation-plan.md) | Draft present; measurements are future work. |
| 8 | YouTrack link with assigned and prioritized tasks | [`../project/youtrack-backlog.csv`](../project/youtrack-backlog.csv) has 30 planned tasks and a 15/15 two-member split. | Missing live project URL and verified issue assignments. |
| 9 | Progressive contributions from each team member | [`../scripts/tasks.tsv`](../scripts/tasks.tsv) records planned ownership. | Missing Git history and peer-review evidence; this supplied folder has no `.git` directory. |

The full project rubric also asks for real user or client evidence and substantive supervisor feedback. [`02-evidence-of-relevance.md`](02-evidence-of-relevance.md) has an interview plan, and [`11-supervisor-feedback.md`](11-supervisor-feedback.md) is a template. Neither should be presented as completed primary research or review.

## Actions before submission

1. Add both real names, Git identities, supervisor details, and any confirmed client details. Keep placeholders where facts are unavailable.
2. Review the literature citations, comparison cells, requirements, diagrams, and prototype; correct gaps found in that review.
3. Publish this folder to a GitHub or GitLab repository under the identity of the person who actually performs the initial import. Give the second member and supervisor access. Do not split the pre-existing package into artificial commits credited to different people.
4. Create the YouTrack project, enter the two real logins in `project/team.conf`, dry-run and import the backlog, then verify the actual issue IDs, assignees, priorities and deadlines. Add the real URL to `README.md`.
5. Each member should perform and commit their own subsequent tasks on task branches. The other member reviews each PR/MR before merge; link the real PR or commit in the matching YouTrack issue.
6. Record real interviews or a substantive client agreement and supervisor feedback when they occur. Keep planned work clearly labelled until then.
7. After confirming the submitted repository state, create and push the `pre-defense-1` tag and fill `SUBMISSION.md` with the real URL, tag, completed items, limits and outstanding work.
