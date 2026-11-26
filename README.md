# inquest

Your disk is under investigation.

`inquest` is a single-binary Go CLI that conducts a formal inquest on your
filesystem: it finds regenerable dependency/build folders, judges them against
the evidence, and executes the guilty to reclaim space.

Innocent until proven guilty — a directory is only convicted if it matches a
known artifact name **and** its parent holds the manifest that can regenerate
it. No manifest, no guilt.

## How it works

The pipeline runs in stages:

1. **Gather** — walk the tree, collect directories whose names are charged
   (`node_modules`, `venv`, `.venv`). Prune on match: only the outermost
   suspect stands accused, `.git` is never entered.
2. **Investigate** — for each suspect, check that the required manifest exists
   in the parent directory (`node_modules` → `package.json`,
   `venv`/`.venv` → `requirements.txt`). Evidence present → guilty.
3. **Weigh** — sum file sizes inside each guilty directory *(in progress)*.
4. **Execute** — `os.RemoveAll` on the guilty, after your confirmation
   *(in progress)*.

## Usage

```
inquest --path <path>    # investigate a directory
inquest -p <path>        # shorthand
inquest --version        # print version (-v)
```

Current flags:

| Flag            | Description              |
|-----------------|--------------------------|
| `--path`, `-p`  | Directory to investigate |
| `--version`, `-v` | Show version           |

Planned:

| Flag        | Description                            |
|-------------|----------------------------------------|
| `--dry-run` | Judge only, execute nothing            |
| `--yes`     | Execute all guilty without asking      |

## Install

```
go install github.com/muhammadhnf97/inquest@latest
```

Or build from source:

```
git clone https://github.com/muhammadhnf97/inquest
cd inquest
go build -o inquest .
```

## Charges (v0.1)

| Artifact       | Required evidence      |
|----------------|------------------------|
| `node_modules` | `package.json`         |
| `venv`         | `requirements.txt`     |
| `.venv`        | `requirements.txt`     |

Never accused: `dist/`, `build/`, anything under `.git/`, and any directory
lacking its manifest.

## Status

Early development. Gathering and Investigation are wired; sizing, docket
output, confirmation, and execution are next. See `SPEC.md` for the
full spec.

Inspired by [kondo](https://github.com/tbillington/kondo) and
[npkill](https://github.com/voidcosmos/npkill).

## License

MIT
