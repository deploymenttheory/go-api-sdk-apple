# device_management

A generated Go SDK for **Apple MDM** and **Declarative Device Management
(DDM)**, built from Apple's canonical schema repo
[apple/device-management](https://github.com/apple/device-management).

Typed values in → spec-validated Apple config out:

```
apple/device-management @ pinned commit
        │ fetchspec (download, parse YAML, snapshot)
        ▼
metadata/specs/**.json           committed, reviewed snapshots + PROVENANCE
        │ gendm (offline, deterministic)
        ▼
mdm/commands  mdm/profiles       generated structs + Validate() + registries
ddm/{configurations,assets,
     activations,management}
        │ mdm / ddm envelope builders (the workflow engine)
        ▼
MDM command plists · .mobileconfig profiles · DDM declaration JSON
```

## Usage

```go
import (
    dm "github.com/deploymenttheory/go-sdk-appleservices/device_management"
    "github.com/deploymenttheory/go-sdk-appleservices/device_management/ddm"
    "github.com/deploymenttheory/go-sdk-appleservices/device_management/ddm/configurations"
    "github.com/deploymenttheory/go-sdk-appleservices/device_management/mdm"
    "github.com/deploymenttheory/go-sdk-appleservices/device_management/mdm/commands"
    "github.com/deploymenttheory/go-sdk-appleservices/device_management/ptr"
)

// MDM command plist
cmd, err := mdm.NewCommand(&commands.DeviceLock{PIN: ptr.To("123456")})

// Configuration profile (.mobileconfig)
prof, err := mdm.NewProfile("com.example.restrictions",
    mdm.WithDisplayName("Corporate Restrictions"),
    mdm.WithPayload(&profiles.Applicationaccess{AllowCamera: ptr.To(false)}))

// DDM declaration JSON
decl, err := ddm.BuildDeclaration("com.example.passcode",
    &configurations.PasscodeSettings{MinimumLength: ptr.To(int64(12))})
```

Every builder validates before emitting — invalid input never becomes
config. Runnable examples live in
[`examples/device_management`](../examples/device_management).

## The spec is honoured

| Apple schema construct | Generated Go |
|---|---|
| `presence: required` | value field, always serialized |
| `presence: optional` | pointer / nil-able field, `omitempty` |
| `<string> <integer> <real> <boolean> <date> <data> <any>` | `string int64 float64 bool time.Time []byte any` |
| `<array>` / `<dictionary>` + subkeys | `[]T` / named nested structs (`subkeytype` shared types dedup) |
| `rangelist` | typed constants + `Validate()` membership check |
| `range`, `format`, `repetition`, `<url> <hostname> <email>` | `Validate()` bounds / regex / cardinality / format checks |
| `requesttype` / `payloadtype` / `declarationtype` | `RequestType()` / `PayloadType()` / `DeclarationType()` + per-family `By*` registries |
| supportedOS | doc comments (`Supported: iOS 15.0+, macOS 13.0+.`) |

## The pipeline

| Stage | Command | Output |
|---|---|---|
| Acquire | `go run ./device_management/cmd/fetchspec` | `metadata/specs/**.json` + `PROVENANCE.json` |
| Generate | `go run ./device_management/cmd/gendm` | `mdm/`, `ddm/` generated packages |
| Diff | `go run ./device_management/cmd/specdiff -old <dir>` | semantic markdown change report |

Generated packages separate construct kinds per spec: `<name>.go` holds the
struct declarations, `<name>_functions.go` the wire-identifier and
`Validate()` methods, `<name>_enums.go` the allowed-value constants, and
each family's `registry.go` the identifier → factory map.

The upstream commit is **pinned** in `cmd/fetchspec/main.go`; `-discover`
resolves the latest commit on Apple's `release` branch. Snapshots are
committed, codegen is offline and byte-deterministic, and CI enforces both:

- `device-management-unit-tests.yml` — build, vet, tests, plus the
  regenerate-and-`git diff --exit-code` determinism gate.
- `device-management-spec-update.yml` — weekly: fetch latest upstream
  specs, regenerate, and open a PR whose body is the `specdiff` report —
  the moving target, reviewed change by change.

## Not covered (v1)

`declarative/status`, `declarative/protocol`, `mdm/checkin`, `mdm/errors`
and `other/` specs (read-side/protocol plumbing — same generator can add
them later), and transporting artifacts to devices (an MDM server's job).
