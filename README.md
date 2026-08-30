# Go API SDK for Apple Services

[![Go Reference](https://pkg.go.dev/badge/github.com/deploymenttheory/go-sdk-appleservices.svg)](https://pkg.go.dev/github.com/deploymenttheory/go-sdk-appleservices)
[![Go Report Card](https://goreportcard.com/badge/github.com/deploymenttheory/go-sdk-appleservices)](https://goreportcard.com/report/github.com/deploymenttheory/go-sdk-appleservices)
[![License](https://img.shields.io/github/license/deploymenttheory/go-sdk-appleservices)](https://github.com/deploymenttheory/go-sdk-appleservices/blob/main/LICENSE)

A collection of Go SDKs for interacting with Apple API services, device management infrastructure, and Microsoft software update feeds:

- **iTunes Search API** — search and lookup across the iTunes, App Store, iBooks Store, and Mac App Store
- **Apple Business Manager / Apple School Manager API** — device inventory and MDM server management - [Apple Business Manager API Changelog](https://developer.apple.com/documentation/apple-school-and-business-manager-api/apple-school-manager-and-apple-business-api-changelog)
- **Apple Device Management (MDM / DDM)** — typed, spec-validated construction of MDM command plists, configuration profiles and Declarative Device Management JSON, generated from [apple/device-management](https://github.com/apple/device-management)
- **Apple Update CDN** — firmware discovery and IPSW download for macOS, iOS, and iPadOS
- **Microsoft Updates** — macOS standalone app updates, Edge channels, OneDrive rings, App Store versions, and Office CVE history

## Features

- Clean, idiomatic Go API
- Fluent builder patterns for constructing API requests
- Comprehensive error handling
- Configurable logging with zap
- Automatic retries with configurable parameters
- Extensive test coverage
- Complete examples for all supported operations

---

## Supported Services

### iTunes Search API

Complete implementation of the [iTunes Search API](https://performance-partners.apple.com/search-api):

- Search for content across iTunes, App Store, iBooks Store, and Mac App Store
- Look up content by ID, UPC, EAN, ISRC, or ISBN
- Filter results by media type, entity, country, and more

---

### Apple Business Manager / Apple School Manager API

Complete implementation of the [Apple Business Manager API](https://developer.apple.com/documentation/apple-school-and-business-manager-api):

**Devices API:**
- Get organization devices with filtering and pagination
- Get detailed device information by serial number
- Support for all device types (iPhone, iPad, Mac, Apple TV, Apple Watch)

**Device Management Services API:**
- List device management services in an organization
- Get device serial numbers assigned to services
- Assign/unassign devices to/from management services
- Get device management service assignments and information
- Track device activity operations

**Key Features:**
- JWT authentication — built-in Apple JWT token generation and management
- Structured logging — comprehensive request/response logging with zap
- Type-safe structured response models
- Context-aware operations for timeouts and cancellation

**Quick Start:**

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/deploymenttheory/go-sdk-appleservices/axm"
    "github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devicemanagement"
    "github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devices"
)

func main() {
    // Method 1: Direct client creation — parse private key from PEM
    privateKey, err := axm.ParsePrivateKey([]byte(privateKeyPEM))
    if err != nil {
        log.Fatalf("Failed to parse private key: %v", err)
    }

    c, err := axm.NewClient("your-key-id", "your-issuer-id", privateKey)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Method 2: From environment variables
    // c, err := axm.NewClientFromEnv()
    // Expects: APPLE_KEY_ID, APPLE_ISSUER_ID, APPLE_PRIVATE_KEY_PATH

    // Method 3: From file
    // c, err := axm.NewClientFromFile("key-id", "issuer-id", "/path/to/key.p8")

    ctx := context.Background()

    // Get organization devices
    response, _, err := c.AXMAPI.Devices.GetV1(ctx, &devices.RequestQueryOptions{
        Fields: []string{
            devices.FieldSerialNumber,
            devices.FieldDeviceModel,
            devices.FieldStatus,
        },
        Limit: 100,
    })
    if err != nil {
        log.Fatalf("Error getting devices: %v", err)
    }

    fmt.Printf("Found %d devices\n", len(response.Data))

    // Get device management services
    mdmServers, _, err := c.AXMAPI.DeviceManagement.GetV1(ctx, &devicemanagement.RequestQueryOptions{
        Limit: 10,
    })
    if err != nil {
        log.Fatalf("Error getting MDM servers: %v", err)
    }

    fmt.Printf("Found %d MDM servers\n", len(mdmServers.Data))
}
```

📖 **[Complete Quick Start Guide →](./examples/axm/quick_start.md)**

---

### Apple Device Management (MDM / DDM)

A generated SDK built from Apple's canonical schema repo
[apple/device-management](https://github.com/apple/device-management):
typed Go values in, spec-validated Apple config out. No authentication or
network — the SDK *produces and validates* management artifacts; your MDM
server delivers them.

**What it covers:**
- **MDM commands** (`mdm/commands`) — every command as a typed payload struct; emitted as command plists
- **Configuration profiles** (`mdm/profiles`) — every payload type; emitted as `.mobileconfig` plists with the envelope assembled for you
- **DDM declarations** (`ddm/configurations`, `ddm/assets`, `ddm/activations`, `ddm/management`) — emitted as declaration JSON

**The spec is honoured:** required keys are value fields, optional keys are
pointers; allowed values are typed enums with constants; every payload's
`Validate()` enforces rangelists, numeric ranges, formats, cardinality and
`<url>`/`<hostname>`/`<email>` subtypes before anything is emitted. The
schema pipeline (pinned upstream commit → committed snapshots → offline
deterministic codegen → weekly semantic-diff PRs) tracks Apple's moving
target.

**Quick Start:**

```go
package main

import (
    "fmt"
    "log"

    "github.com/deploymenttheory/go-sdk-appleservices/device_management/ddm"
    "github.com/deploymenttheory/go-sdk-appleservices/device_management/ddm/configurations"
    "github.com/deploymenttheory/go-sdk-appleservices/device_management/mdm"
    "github.com/deploymenttheory/go-sdk-appleservices/device_management/mdm/commands"
    "github.com/deploymenttheory/go-sdk-appleservices/device_management/ptr"
)

func main() {
    // MDM command plist.
    cmd, err := mdm.NewCommand(&commands.DeviceLock{
        Message: ptr.To("Locked by IT"),
        PIN:     ptr.To("123456"),
    })
    if err != nil {
        log.Fatalf("build command: %v", err)
    }
    fmt.Println(string(cmd))

    // DDM declaration JSON — validation rejects out-of-spec values
    // (MinimumLength allows 0-16 per Apple's schema).
    decl, err := ddm.BuildDeclaration("com.example.passcode",
        &configurations.PasscodeSettings{
            RequirePasscode: ptr.To(true),
            MinimumLength:   ptr.To(int64(12)),
        })
    if err != nil {
        log.Fatalf("build declaration: %v", err)
    }
    fmt.Println(string(decl))
}
```

📖 **[Full device_management documentation →](./device_management/README.md)**

---

### Apple Update CDN

Firmware discovery and IPSW download across macOS, iOS, and iPadOS — no authentication required.

The SDK spans three external APIs:

| Service | Host | Purpose |
|---------|------|---------|
| ipsw.me API | `api.ipsw.me` | Firmware discovery with CDN download URLs, checksums, signing status |
| Apple GDMF API | `gdmf.apple.com` | Official Apple feed of currently-signed firmware versions |
| Apple CDN | `updates.cdn-apple.com` | URL parsing, HEAD metadata, streaming IPSW downloads |

**Firmware discovery (ipsw.me API):**
- `ListAllFirmwareV3` — all device platforms unfiltered
- `ListAllMacFirmwareV3` / `ListAllIOSFirmwareV3` / `ListAllIPadOSFirmwareV3` — per-platform filtered lists
- `ListUniqueMacFirmwareVersionsV3` / `ListUniqueIOSFirmwareVersionsV3` / `ListUniqueIPadOSFirmwareVersionsV3` — deduplicated versions sorted newest-first
- `GetByDeviceV4` — firmware history for a specific model identifier (e.g. `"Mac14,3"`, `"iPhone15,2"`, `"iPad14,4"`) with SHA-256 checksums

**Apple GDMF (signed version feed):**
- `GetPublicVersionsV2` — Apple's authoritative list of currently-signed versions for macOS, iOS, and visionOS including posting/expiration dates and supported device lists

**CDN utilities:**
- `ParseURL` — parse an IPSW CDN URL into its structural components (no HTTP request)
- `GetFileMetadataV1` — HEAD request returning SHA-1, SHA-256, file size, and last-modified without downloading the file
- `DownloadFileV1` — streaming IPSW download with SHA-1/SHA-256 verification and progress callback

**Quick Start:**

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn"
    "github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/tools/download_progress"
)

func main() {
    c, err := apple_update_cdn.NewDefaultClient()
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer c.Close()

    ctx := context.Background()

    // List all unique macOS firmware versions, newest first.
    versions, _, err := c.AppleUpdateCDNAPI.Firmware.ListUniqueMacFirmwareVersionsV3(ctx)
    if err != nil {
        log.Fatalf("Error listing firmware: %v", err)
    }

    fmt.Printf("Found %d unique macOS versions\n", len(versions))
    for _, fw := range versions[:3] {
        fmt.Printf("  %s (%s) signed=%v\n", fw.Version, fw.BuildID, fw.Signed)
    }

    // Download the latest signed macOS IPSW with a progress bar.
    latest := versions[0]
    destPath := filepath.Join(os.TempDir(), filepath.Base(latest.URL))

    bar := download_progress.New(os.Stderr)
    progressFn := func(written, total int64) {
        bar.Callback(filepath.Base(destPath))(written, total)
    }

    result, _, err := c.AppleUpdateCDNAPI.CDN.DownloadFileV1(ctx, latest.URL, destPath, progressFn)
    if err != nil {
        log.Fatalf("Download failed: %v", err)
    }

    fmt.Printf("\nDownloaded %.2f GB in %s — verified=%v\n",
        float64(result.BytesWritten)/1e9, result.Duration.Round(1e9), result.Verified)
}
```

**Download behaviour:**
- Issues a HEAD request first to obtain expected size and checksums
- Streams the response body directly to disk — the full file is never held in memory
- Computes SHA-1 and SHA-256 simultaneously during streaming via `io.MultiWriter`
- Removes the partial file and returns an error on checksum mismatch or GET failure
- Creates the destination directory automatically if it does not exist
- macOS IPSW files are typically 15–22 GB; ensure sufficient free space

---

### Microsoft Updates

Tracks Microsoft software releases for macOS and iOS from official Microsoft endpoints — no authentication required. Replicates the data-collection logic of the [MOFA project](https://github.com/cocopuff2u/MOFA) in pure Go.

The SDK spans multiple external APIs:

| Sub-service | Host | Purpose |
|---|---|---|
| `standalone` | `officecdnmac.microsoft.com` | Production Office CDN — plist XML per app |
| `standalone_beta` | `officecdnmac.microsoft.com` | Insider Fast (beta) channel |
| `standalone_preview` | `officecdnmac.microsoft.com` | Insider Slow (preview) channel |
| `edge` | `edgeupdates.microsoft.com` | Edge stable / beta / dev / canary |
| `onedrive` | `g.live.com` + fwlink redirects | OneDrive distribution rings |
| `appstore_macos` | `itunes.apple.com` | Microsoft apps in the macOS App Store |
| `appstore_ios` | `itunes.apple.com` | Microsoft apps in the iOS App Store |
| `update_history` | `learn.microsoft.com` | Office for Mac release table (HTML) |
| `cve_history` | `learn.microsoft.com` | Office for Mac CVE/security notes (HTML) |

**Standalone apps tracked (17):** Word, Excel, PowerPoint, Outlook, OneNote, Teams, Skype for Business, Defender (Endpoint/Consumer/Shim), Intune Company Portal, Microsoft AutoUpdate, Windows App, Microsoft 365 Copilot, Quick Assist, Remote Help, Licensing Helper Tool.

**Quick Start:**

```go
package main

import (
    "context"
    "fmt"
    "log"

    microsoft_updates "github.com/deploymenttheory/go-sdk-appleservices/microsoft_updates"
)

func main() {
    c, err := microsoft_updates.NewDefaultClient()
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer c.Close()

    ctx := context.Background()

    // Get latest production standalone app versions from the Microsoft CDN.
    resp, err := c.MicrosoftUpdatesAPI.Standalone.GetLatestV1(ctx)
    if err != nil {
        log.Fatalf("Error getting standalone apps: %v", err)
    }
    for _, pkg := range resp.Packages {
        fmt.Printf("%-40s %s\n", pkg.Title, pkg.FullVersion)
    }

    // Get Microsoft Edge across all four channels.
    edge, err := c.MicrosoftUpdatesAPI.Edge.GetAllChannelsV1(ctx)
    if err != nil {
        log.Fatalf("Error getting Edge channels: %v", err)
    }
    fmt.Printf("Edge stable: %s\n", edge.Stable.Version)
    fmt.Printf("Edge canary: %s\n", edge.Canary.Version)

    // Get OneDrive version per distribution ring.
    od, err := c.MicrosoftUpdatesAPI.OneDrive.GetAllRingsV1(ctx)
    if err != nil {
        log.Fatalf("Error getting OneDrive rings: %v", err)
    }
    for _, ring := range od.Rings {
        fmt.Printf("OneDrive %-20s %s\n", ring.Ring, ring.Version)
    }

    // Get Office CVE history.
    cves, err := c.MicrosoftUpdatesAPI.CVEHistory.GetCVEHistoryV1(ctx)
    if err != nil {
        log.Fatalf("Error getting CVE history: %v", err)
    }
    fmt.Printf("CVE history entries: %d\n", len(cves.Entries))
}
```

---

## Examples

The [examples directory](./examples) contains a runnable `main.go` for every SDK function:

```
examples/
├── axm/                         Apple Business Manager
│   ├── devices/
│   └── devicemanagement/
├── device_management/           Apple MDM / DDM config generation
│   ├── commands/                MDM command plists
│   ├── profiles/                .mobileconfig profiles
│   └── ddm/                     DDM declarations and activation sets
├── apple_update_cdn/            Apple Update CDN
│   ├── firmware/                ipsw.me firmware discovery
│   ├── gdmf/                    Apple signed-version feed
│   └── cdn/                     URL parsing, metadata, download
├── itunes_search/               iTunes Search API
└── microsoft_updates/           Microsoft Updates
    ├── standalone/              CDN app versions (production/beta/preview)
    ├── edge/                    Edge channel versions
    ├── onedrive/                OneDrive distribution rings
    ├── appstore/                macOS and iOS App Store versions
    ├── update_history/          Office for Mac update history
    └── cve_history/             Office CVE/security release notes
```

---

## Documentation

- [Go Reference Documentation](https://pkg.go.dev/github.com/deploymenttheory/go-sdk-appleservices)
- [iTunes Search API Documentation](https://performance-partners.apple.com/search-api)
- [Apple Business Manager API Documentation](https://developer.apple.com/documentation/applebusinessmanagerapi)
- [Apple Device Management Documentation](https://developer.apple.com/documentation/devicemanagement)
- [Apple device-management schema repo (source of the generated MDM/DDM SDK)](https://github.com/apple/device-management)
- [Microsoft Office CDN (MOFA reference)](https://github.com/cocopuff2u/MOFA)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the [MIT License](./LICENSE).
