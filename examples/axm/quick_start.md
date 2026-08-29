# AXM SDK Quick Start Guide

## Prerequisites

- An Apple Business Manager or Apple School Manager account
- An API key from the ABM/ASM portal (Key ID, Issuer ID, and a `.p8` private key file)

---

## Client Setup

### Method 1: From a PEM-encoded private key

```go
package main

import (
    "log"

    "github.com/deploymenttheory/go-sdk-appleservices/axm"
)

func main() {
    privateKeyPEM := `-----BEGIN EC PRIVATE KEY-----
<your key contents>
-----END EC PRIVATE KEY-----`

    privateKey, err := axm.ParsePrivateKey([]byte(privateKeyPEM))
    if err != nil {
        log.Fatalf("Failed to parse private key: %v", err)
    }

    c, err := axm.NewClient("your-key-id", "your-issuer-id", privateKey)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    _ = c
}
```

### Method 2: From a `.p8` key file

```go
c, err := axm.NewClientFromFile("your-key-id", "your-issuer-id", "/path/to/key.p8")
if err != nil {
    log.Fatalf("Failed to create client: %v", err)
}
```

### Method 3: From environment variables

Set the following environment variables before running:

```sh
export APPLE_KEY_ID="your-key-id"
export APPLE_ISSUER_ID="your-issuer-id"
export APPLE_PRIVATE_KEY_PATH="/path/to/key.p8"
```

```go
c, err := axm.NewClientFromEnv()
if err != nil {
    log.Fatalf("Failed to create client: %v", err)
}
```

---

## Configuration Options

All three constructors accept optional `ClientOption` values as trailing arguments:

```go
import (
    "crypto/tls"
    "time"

    "github.com/deploymenttheory/go-sdk-appleservices/axm"
    "go.uber.org/zap"
)

logger, _ := zap.NewProduction()

c, err := axm.NewClient(keyID, issuerID, privateKey,
    axm.WithLogger(logger),              // structured zap logger
    axm.WithTimeout(60*time.Second),     // per-request timeout (default: 30s)
    axm.WithRetryCount(5),               // max retries on failure (default: 3)
    axm.WithRetryWaitTime(2*time.Second),
    axm.WithRetryMaxWaitTime(30*time.Second),
    axm.WithDebug(),                     // log full request/response bodies
    axm.WithProxy("http://proxy:8080"),
    axm.WithCustomAgent("my-tool/1.0"),
    axm.WithScope("school.api"),         // default: "business.api"
    axm.WithMinTLSVersion(tls.VersionTLS13),
)
```

---

## Devices API

### List organization devices

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/deploymenttheory/go-sdk-appleservices/axm"
    "github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devices"
)

func main() {
    c, err := axm.NewClientFromEnv()
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    ctx := context.Background()

    opts := &devices.RequestQueryOptions{
        Fields: []string{
            devices.FieldSerialNumber,
            devices.FieldDeviceModel,
            devices.FieldProductFamily,
            devices.FieldStatus,
            devices.FieldAddedToOrgDateTime,
            devices.FieldUpdatedDateTime,
        },
        Limit: 100,
    }

    response, _, err := c.AXMAPI.Devices.GetV1(ctx, opts)
    if err != nil {
        log.Fatalf("Error getting organization devices: %v", err)
    }

    fmt.Printf("Found %d devices\n", len(response.Data))
    for _, device := range response.Data {
        fmt.Printf("  %s — %s (%s)\n",
            device.Attributes.SerialNumber,
            device.Attributes.DeviceModel,
            device.Attributes.Status,
        )
    }

    // Pagination: follow Links.Next for subsequent pages
    if response.Links != nil && response.Links.Next != "" {
        fmt.Printf("Next page cursor available\n")
    }
}
```

See full example: [GetOrganizationDevices/main.go](./devices/GetOrganizationDevices/main.go)

### Get a single device by ID

```go
response, _, err := c.AXMAPI.Devices.GetByDeviceIDV1(ctx, "XABC123X0ABC123X0", opts)
if err != nil {
    log.Fatalf("Error: %v", err)
}
fmt.Printf("Serial: %s\n", response.Data.Attributes.SerialNumber)
```

See full example: [GetDeviceInformationByDeviceID/main.go](./devices/GetDeviceInformationByDeviceID/main.go)

### Get AppleCare coverage for a device

```go
opts := &devices.RequestQueryOptions{
    Fields: []string{
        devices.FieldAppleCareStatus,
        devices.FieldAppleCareDescription,
        devices.FieldAppleCareStartDateTime,
        devices.FieldAppleCareEndDateTime,
    },
}

response, _, err := c.AXMAPI.Devices.GetAppleCareByDeviceIDV1(ctx, "XABC123X0ABC123X0", opts)
if err != nil {
    log.Fatalf("Error: %v", err)
}
fmt.Printf("Found %d coverage plan(s)\n", len(response.Data))
```

See full example: [GetAppleCareInformationByDeviceID/main.go](./devices/GetAppleCareInformationByDeviceID/main.go)

---

## Device Management API

### List MDM servers

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/deploymenttheory/go-sdk-appleservices/axm"
    "github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devicemanagement"
)

func main() {
    c, err := axm.NewClientFromEnv()
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    ctx := context.Background()

    opts := &devicemanagement.RequestQueryOptions{
        Fields: []string{
            devicemanagement.FieldServerName,
            devicemanagement.FieldServerType,
            devicemanagement.FieldCreatedDateTime,
        },
        Limit: 100,
    }

    response, _, err := c.AXMAPI.DeviceManagement.GetV1(ctx, opts)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Printf("Found %d MDM server(s)\n", len(response.Data))
    for _, server := range response.Data {
        fmt.Printf("  %s — %s\n", server.ID, server.Attributes.ServerName)
    }
}
```

See full example: [GetDeviceManagementServices/main.go](./devicemanagement/GetDeviceManagementServices/main.go)

### Get devices assigned to an MDM server

```go
opts := &devicemanagement.RequestQueryOptions{Limit: 100}

response, _, err := c.AXMAPI.DeviceManagement.GetDeviceSerialNumbersByServerIDV1(ctx, mdmServerID, opts)
if err != nil {
    log.Fatalf("Error: %v", err)
}
for _, linkage := range response.Data {
    fmt.Printf("  Device ID: %s\n", linkage.ID)
}
```

See full example: [GetDeviceSerialNumbersForDeviceManagementService/main.go](./devicemanagement/GetDeviceSerialNumbersForDeviceManagementService/main.go)

### Find which MDM server a device is assigned to

```go
// Returns just the server ID/type linkage
response, _, err := c.AXMAPI.DeviceManagement.GetAssignedServerIDByDeviceIDV1(ctx, deviceID)

// Returns full server attributes
opts := &devicemanagement.RequestQueryOptions{
    Fields: []string{devicemanagement.FieldServerName, devicemanagement.FieldServerType},
}
response, _, err := c.AXMAPI.DeviceManagement.GetAssignedServerInfoByDeviceIDV1(ctx, deviceID, opts)
```

See full examples:
- [GetAssignedDeviceManagementServiceID/main.go](./devicemanagement/GetAssignedDeviceManagementServiceID/main.go)
- [GetAssignedDeviceManagementServiceInfo/main.go](./devicemanagement/GetAssignedDeviceManagementServiceInfo/main.go)

### Assign devices to an MDM server

```go
mdmServerID := "1F97349736CF4614A94F624E705841AD"
deviceIDs := []string{"XABC123X0ABC123X0", "YDEF456Y1DEF456Y1"}

response, _, err := c.AXMAPI.DeviceManagement.AssignDevicesV1(ctx, mdmServerID, deviceIDs)
if err != nil {
    log.Fatalf("Error: %v", err)
}
fmt.Printf("Activity ID: %s — Status: %s\n",
    response.Data.ID,
    response.Data.Attributes.Status,
)
```

See full example: [AssignDevicesToServer/main.go](./devicemanagement/AssignDevicesToServer/main.go)

### Unassign devices from an MDM server

```go
response, _, err := c.AXMAPI.DeviceManagement.UnassignDevicesV1(ctx, mdmServerID, deviceIDs)
if err != nil {
    log.Fatalf("Error: %v", err)
}
fmt.Printf("Activity ID: %s — Status: %s\n",
    response.Data.ID,
    response.Data.Attributes.Status,
)
```

See full example: [UnassignDevicesFromServer/main.go](./devicemanagement/UnassignDevicesFromServer/main.go)

---

## Error Handling

```go
import "github.com/deploymenttheory/go-sdk-appleservices/axm"

response, _, err := c.AXMAPI.Devices.GetByDeviceIDV1(ctx, deviceID, opts)
if err != nil {
    if axm.IsNotFound(err) {
        fmt.Printf("Device %s not found\n", deviceID)
        return
    }
    log.Fatalf("Unexpected error: %v", err)
}
```
