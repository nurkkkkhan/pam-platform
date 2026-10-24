#!/usr/bin/env python3
"""Create the PAM Platform backlog in YouTrack (assigned + prioritized tasks).

Usage:
  1. In YouTrack create a project with short name DIP (template: Scrum or Kanban).
     Add every team member and the supervisor to the project team.
  2. Profile -> Account Security -> Tokens -> New token (scope: YouTrack).
  3. Fill project/team.conf with real YouTrack logins.
  4. Run:
       export YT_URL=https://<your-instance>.youtrack.cloud
       export YT_TOKEN=perm:xxxxxxxx
       python3 project/youtrack_import.py --dry-run   # preview
       python3 project/youtrack_import.py             # create issues

Only the Python standard library is used. The token is read from the
environment and never written to disk.
"""
import argparse, csv, json, os, sys, urllib.error, urllib.request
from datetime import datetime, timezone
from pathlib import Path

HERE = Path(__file__).resolve().parent
STATE_MAP = {  # backlog state -> candidate YouTrack state names (first that exists wins)
    "Done": ["Done", "Fixed", "Verified"],
    "In Progress": ["In Progress"],
    "Open": ["Open", "To do", "Submitted"],
}

def load_team():
    team = {}
    for line in (HERE / "team.conf").read_text(encoding="utf-8").splitlines():
        if not line.strip() or line.lstrip().startswith("#"):
            continue
        ph, name, yt, gh = [p.strip() for p in line.split("|")]
        team[ph] = {"name": name, "yt": yt, "gh": gh}
    return team

def load_milestones():
    rows = csv.DictReader(open(HERE / "milestones.csv", encoding="utf-8"))
    return {r["Milestone"]: (r["Name"], r["Date"]) for r in rows}

def api(method, path, body=None):
    url = os.environ["YT_URL"].rstrip("/") + path
    req = urllib.request.Request(url, method=method,
        data=json.dumps(body).encode() if body is not None else None,
        headers={"Authorization": "Bearer " + os.environ["YT_TOKEN"],
                 "Accept": "application/json", "Content-Type": "application/json"})
    with urllib.request.urlopen(req, timeout=30) as r:
        return json.loads(r.read() or b"null")

def project_info(short):
    projects = api("GET", "/api/admin/projects?fields=id,shortName,name&$top=500")
    p = next((p for p in projects if p["shortName"] == short), None)
    if not p:
        sys.exit(f"Project {short} not found. Create it in YouTrack first.")
    fields = api("GET", f"/api/admin/projects/{p['id']}/customFields?fields=field(name),bundle(values(name))&$top=100")
    values = {f["field"]["name"]: [v["name"] for v in (f.get("bundle") or {}).get("values", [])] for f in fields}
    return p, values

def build_fields(row, team, values, due):
    cf, skipped = [], []
    if "Due Date" in values:
        ms = int(datetime.strptime(due, "%Y-%m-%d").replace(hour=12, tzinfo=timezone.utc).timestamp() * 1000)
        cf.append({"name": "Due Date", "$type": "DateIssueCustomField", "value": ms})
    else:
        skipped.append("no Due Date field (milestone date kept in description)")
    pr = row["Priority"]
    if pr in values.get("Priority", []):
        cf.append({"name": "Priority", "$type": "SingleEnumIssueCustomField", "value": {"name": pr}})
    else:
        skipped.append(f"Priority={pr}")
    member = team.get(row["Assignee"])
    if member and "Assignee" in values:
        cf.append({"name": "Assignee", "$type": "SingleUserIssueCustomField", "value": {"login": member["yt"]}})
    elif "Assignee" not in values:
        skipped.append("Assignee field missing")
    if "State" in values:
        for cand in STATE_MAP.get(row["State"], [row["State"]]):
            if cand in values["State"]:
                cf.append({"name": "State", "$type": "StateIssueCustomField", "value": {"name": cand}})
                break
        else:
            skipped.append(f"State={row['State']}")
    if row.get("Subsystem") and row["Subsystem"] in values.get("Subsystem", []):
        cf.append({"name": "Subsystem", "$type": "SingleOwnedIssueCustomField", "value": {"name": row["Subsystem"]}})
    return cf, skipped

def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--project", default="DIP")
    ap.add_argument("--dry-run", action="store_true")
    a = ap.parse_args()
    team = load_team()
    rows = list(csv.DictReader(open(HERE / "youtrack-backlog.csv", encoding="utf-8")))
    miles = load_milestones()
    bad = [r["Assignee"] for r in rows if r["Assignee"] not in team or team[r["Assignee"]]["yt"].startswith("yt-login")]
    if bad and not a.dry_run:
        sys.exit("Fill real YouTrack logins in project/team.conf first.")
    if a.dry_run:
        for r in rows:
            m = team.get(r["Assignee"], {})
            print(f"{r['Issue Id']:7} {r['Priority']:9} {r['State']:12} {m.get('yt','?'):14} {r['Milestone']} {miles[r['Milestone']][1]}  {r['Summary']}")
        print(f"\n{len(rows)} issues would be created in project {a.project}.")
        return
    for k in ("YT_URL", "YT_TOKEN"):
        if k not in os.environ:
            sys.exit(f"Set environment variable {k}.")
    p, values = project_info(a.project)
    print(f"Project {p['shortName']} ({p['name']}); fields: {', '.join(values)}")
    for r in rows:
        mname, due = miles[r["Milestone"]]
        desc = (f"**Expected outcome / acceptance criteria:** {r['Acceptance criteria']}\n\n"
                f"**Milestone:** {r['Milestone']} {mname} · **Deadline:** {due}\n\n"
                f"**Responsible:** {team[r['Assignee']]['name']}\n\n"
                f"Branch: `{r['Issue Id']}-...` · Commits/PR must reference `{r['Issue Id']}`.")
        cf, skipped = build_fields(r, team, values, due)
        body = {"project": {"id": p["id"]}, "summary": r["Summary"], "description": desc, "customFields": cf}
        try:
            res = api("POST", "/api/issues?fields=idReadable", body)
            note = f"  (skipped: {'; '.join(skipped)})" if skipped else ""
            print(f"created {res['idReadable']} (planned {r['Issue Id']}): {r['Summary']}{note}")
        except urllib.error.HTTPError as e:
            print(f"FAILED {r['Issue Id']}: {e.code} {e.read().decode()[:300]}")
    print(f"\nBacklog link for README: {os.environ['YT_URL'].rstrip('/')}/projects/{p['shortName']}")
    print("If YouTrack numbers differ from the planned IDs, update the CSV/docs, or use a fresh project.")

if __name__ == "__main__":
    main()
