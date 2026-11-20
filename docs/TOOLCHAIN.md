# Toolchain Programs

The `toolchain/` module builds the external CLI programs and provides the
build, test and lint entry points. The program list is defined once in
[`toolchain/programs.json`](../toolchain/programs.json).

Each command is a `main` package under its own folder:
`toolchain/build/`, `toolchain/build-webapp/`, `toolchain/clean/`,
`toolchain/install/`, `toolchain/test/` and `toolchain/lint/`.
Flag parsing for a command lives in its `main.go`. Shared helpers live under
`toolchain/internal/utils/<domain>/`; the codestyle linter implementation lives
under `toolchain/internal/utils/lint/`, and the web app build actions live under
`toolchain/internal/utils/webapp/`.

## Manifest

Every entry describes one installed binary:

```json
{
	"name": "dnsx",
	"folder": "dnsx",
	"repository": "https://github.com/projectdiscovery/dnsx.git",
	"version": "v1.3.1",
	"source": "./cmd/dnsx"
}
```

| Field | Meaning |
|:------|:--------|
| `name` | Installed binary name and `allowed-programs` entry. |
| `folder` | Folder name below `vendor/`. Entries may share a folder. |
| `repository` | Git repository URL. |
| `version` | Pinned tag. |
| `source` | Go package path passed to `go build`, relative to the repository. |

## Make Targets

Run these from the repository root. The definitions live in the
[`Makefile`](../Makefile).

| Target | Action |
|:-------|:-------|
| `make` | Build the web app, clone and build programs, build all exocomp binaries and the installer. |
| `make quick` | Build the web app and host-platform binaries only. Skips programs and the installer. |
| `make linux` | Build for Linux. |
| `make darwin` | Build for macOS. |
| `make windows` | Build for Windows. |
| `make webapp` | Build only the WASM web app. |
| `make lint` | Run the codestyle linter for `source/`. |
| `make lint-toolchain` | Run the codestyle linter for `toolchain/`. |
| `make test` | Run the `source/` unit tests. |
| `make test-toolchain` | Run the `toolchain/` module tests. |
| `make test-integration` | Check for and start a local llama-server. The test steps are still incomplete. |
| `make install` | Run the installer. |
| `make clean` | Remove generated assets and `build/<operating-system>/`. |

Cross-build another architecture with `ARCH`:

```bash
make linux ARCH=arm64;
```

## Build and Cache

`ClonePrograms` checks out every pinned tag below `vendor/`. `BuildPrograms`
builds each entry into:

```text
build/cache/<operating-system>/<architecture>/<name>
```

A cached binary is reused when its sidecar `.version` file matches the pinned
manifest version. A missing or stale entry is built lazily with
`GOOS=<operating-system>`, `GOARCH=<architecture>` and `CGO_ENABLED=0`.

Every cached binary is copied into `installer/assets/programs/<name>` so the
installer embeds the correct platform and architecture. Exocomp's own binaries
are rebuilt on every make invocation and written under `build/<operating-system>/`.

`make clean` removes generated assets and `build/<operating-system>/` but leaves
`build/cache/` intact. Remove `build/cache/<operating-system>/<architecture>/`
manually to force a full program rebuild.

## Adding a Program

1. Add a pinned manifest entry.
2. Add the binary to the relevant `allowed-programs` list in
   `source/agents/*.yaml`.
3. Rebuild with `make`.

## Gobuster

Gobuster `v3.8.2` is pinned in `toolchain/programs.json` and allowed for the
`webscanner` role. Agents run it directly with `-o results.txt`. `assets.Import`
parses the text lines into `endpoint` assets and keeps the status and size
properties from each line.
