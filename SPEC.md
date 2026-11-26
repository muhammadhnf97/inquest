# inquest — v0.1 spec

Single-binary Go CLI that conducts a formal inquest on your disk: finds
regenerable dependency/build folders, judges them, and executes the guilty
to reclaim space.

Name: **inquest** — one word, the formal investigation it performs.
Binary: `inquest`.

## Command surface

```
inquest --path <path>           # run inquest: scan path, judge, ask before executing
inquest --path <path> --dry-run # judge only, execute nothing
inquest --path <path> --yes     # execute all found guilty without asking
inquest --version / -v
```

`--path` has `-p` shorthand. Default (no --dry-run/--yes) = inquest, list the
accused, confirm before execution.

## Guilt rules (whitelist only — innocent until proven guilty)

A directory is guilty ONLY if name matches AND its parent contains the manifest
(the evidence):

| folder        | required manifest in parent dir          | restored by        | status      |
|---------------|------------------------------------------|--------------------|-------------|
| node_modules  | package.json                             | npm/yarn install   | implemented |
| .venv, venv   | requirements.txt or pyproject.toml       | pip install        | partial*    |
| __pycache__   | any *.py file                            | auto on next run   | planned     |
| target        | Cargo.toml                               | cargo build        | planned     |
| .next, .nuxt  | package.json                             | next/nuxt build    | planned     |
| vendor        | go.mod                                   | go mod vendor      | planned     |

*current `charges` map holds one manifest per name; .venv/venv convict on
requirements.txt only — pyproject.toml still to add.

Nested junk (e.g. node_modules inside node_modules) is skipped — only the
outermost dir stands accused. Never descend into .git. Never match inside an
already-accused dir.

ACQUITTED (never offered): dist/, build/, anything under .git/, folders lacking
the manifest. (dist/build are intentionally excluded — they may be committed
artifacts; no manifest, no guilt.)

## Output format

```
inquest v0.1.0 — investigating /path

  SIZE    VERDICT       TYPE          PATH
  306 MB  guilty        node_modules  netra/client/node_modules
  1.4 GB  guilty        .venv         witness/server/.venv
  ...                                 (evidence: package.json / requirements.txt)

9 folders stand accused. Sentence: 2.9 GB reclaimed.
Execute? [y/N]
```

Sorted by size descending. Sizes humanized (KB/MB/GB).

## Behavior details

- Walk via filepath.WalkDir; prune on name-match (outermost only) and .git.
- Size = sum of file info.Size() under the dir (follow no symlinks).
- Permission-denied dirs: skip silently, count and report at end ("N dirs
  withheld evidence — permission denied" is optional flavor).
- No accused: print "Inquest complete. All innocent." exit 0.
- Execution: os.RemoveAll per guilty dir; report per-dir success/failure.
- Interrupt (Ctrl-C): halt inquest, print partial docket, exit 130.

## Architecture

Pipeline — four single-job stages:

```
ValidatePath(path)          → os.Stat: exists + is a dir, else exit 2
Gathering(path)  []Suspect  → filepath.WalkDir; name-match the charges table,
                              collect Suspect{Path, Artifact, Manifest},
                              SkipDir on match (outermost only) + .git
Investigate(suspects)       → os.Stat(parent/manifest) per suspect — the
                              trial: evidence present → guilty
Weigh(guilty)               → recursive size-walk per guilty dir (the only
                              O(disk) step — goroutines + channel later)
Execute(guilty)             → os.RemoveAll per guilty Path, per-dir report
```

Types:

```
Suspect{Path, Artifact, Manifest}
  Path     — the artifact's own path (.../client/node_modules) — delete target
  Artifact — dir name that matched (node_modules, .venv, ...)
  Manifest — required evidence filename in parent (package.json, ...)
```

Design notes:

- Gather prunes on NAME alone — it can't know guilt yet. Accepts the
  (vanishingly rare) miss of a guilty dir nested inside an innocent one.
- Investigate is a pure list filter — no walking, one Stat per suspect.
- "Manifest only" never reaches Investigate — no artifact dir, no case.
- walks are for discovery, stats are for verification — Investigate must
  not re-walk the tree.

## Tests (table-driven, t.TempDir)

- node_modules + package.json → guilty
- node_modules WITHOUT package.json → innocent
- nested node_modules → outer only accused
- vendor + go.mod → guilty; vendor alone → innocent
- dist/ + package.json → innocent (never accused)
- .git contents → never reported
- --dry-run → nothing deleted
- permission-denied dir → skipped, reported

## Release checklist

1. README: what/why/install + "inspired by kondo/npkill" + flavor line
   ("your disk is under investigation")
2. LICENSE (MIT) — done
3. git tag v0.1.0 → goreleaser: darwin/linux/windows binaries + sha256
4. homebrew-tap repo + goreleaser brews formula
5. go install github.com/<user>/inquest@v0.1.0 works
6. Test on own machine first (expected: finds ~2.5GB in FUNPROJECT —
   dry run already surfaces ~11 suspects there)

## Out of scope (v0.2+)

TUI picker, staleness scoring ("the accused has lain dormant 6 months"),
config/ignore file, scheduled inquests, dist/build opt-in flag, brew core,
package managers beyond brew, duplicate-file detection (separate tool —
breaks the "everything convicted is rebuildable" invariant), signs-of-life
check (recent git activity → leniency warning).

## Versioning ladder

- **v0.1.0** — core loop proven: gather → verdict → docket → prompt →
  execute. Three charges, --dry-run/--yes, Ctrl-C, tests.
- **v0.2.x** — TUI picker, config/ignore file, completed charges table
  (all 8 rows + manifest alternates), scheduled inquests.
- **v0.3.x** — sort flags, staleness scoring, polish; possibly the
  AI-assisted preliminary hearing on uncharged dirs.
- **v1.0.0** — the contract freezes: flag surface, exit codes, output
  format, and guilt rules are all stable. Semantic promise: upgrades
  within 1.x never break scripts or muscle memory.
