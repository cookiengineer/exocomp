# Toolchain Programs

The build toolchain clones, pins and builds the external CLI programs that are
shipped with Exocomp. The list is defined once in
[toolchain/programs.json](../toolchain/programs.json).

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
| `version` | Pinned tag. The build fails if the checked-out commit differs. |
| `source` | Go package path passed to `go build`, relative to the repository. |

## Build

```bash
make linux
make darwin
make windows

# Cross-build for another architecture
make linux ARCH=arm64
```

`ClonePrograms` checks out every pinned tag below `vendor/`. `BuildPrograms`
builds each entry into the per-target cache:

```text
build/cache/<operating-system>/<architecture>/<name>
```

A cached binary is reused when its sidecar `.version` file matches the pinned
manifest version. A missing or stale cache entry is built lazily with
`GOOS=<operating-system>`, `GOARCH=<architecture>` and `CGO_ENABLED=0`.

Every cached binary is then copied into `installer/assets/programs/<name>` so
the installer embeds the correct platform and architecture. Exocomp's own
binaries are still rebuilt for each make invocation and written under
`build/<operating-system>/`.

`make clean` removes generated assets and `build/<operating-system>/` but
leaves `build/cache/` intact. Remove `build/cache/<operating-system>/<architecture>/`
manually to force a full program rebuild.

Adding a program requires:

1. adding a pinned manifest entry,
2. adding the binary to the relevant `allowed-programs` list in
   `source/agents/*.yaml`,
3. rebuilding the installer.

The toolchain builds with `CGO_ENABLED=0`. `GOARCH` defaults to the host
architecture and can be overridden with `ARCH=<goarch>` on the make command or
`--arch=<goarch>` on the toolchain command.

## Gobuster

`gobuster` (`v3.8.2`) is the chosen content/parameter fuzzer. It is pinned in
`toolchain/programs.json` and allowed for the `webscanner` role. Agents run it
directly with `-o results.txt`; `assets.Import` parses the text lines into
`endpoint` assets. No wrapper program is required.
