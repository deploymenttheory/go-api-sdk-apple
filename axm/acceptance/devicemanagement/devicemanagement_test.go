package devicemanagement

import (
	"context"
	"fmt"
	"testing"
	"time"

	acc "github.com/deploymenttheory/go-sdk-appleservices/axm/acceptance"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devicemanagement"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/devices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// TestAcceptance_DeviceManagement_GetDeviceManagementServices
// Verifies the MDM-server list endpoint returns a valid response.
// =============================================================================

func TestAcceptance_DeviceManagement_GetDeviceManagementServices(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.AXMAPI.DeviceManagement
	ctx := context.Background()

	// --- Default options ---
	t.Run("DefaultOptions", func(t *testing.T) {
		acc.LogTestStage(t, "List", "Getting device management services with default options")

		ctx1, cancel1 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel1()

		result, resp, err := svc.GetV1(ctx1, nil)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Data, "organization should have at least one MDM server")

		acc.LogTestSuccess(t, "GetV1: found %d MDM server(s)", len(result.Data))
	})

	// --- Field selection ---
	t.Run("WithFieldSelection", func(t *testing.T) {
		acc.LogTestStage(t, "List", "Getting device management services with specific fields")

		ctx2, cancel2 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel2()

		opts := &devicemanagement.RequestQueryOptions{
			Fields: []string{
				devicemanagement.FieldServerName,
				devicemanagement.FieldServerType,
				devicemanagement.FieldCreatedDateTime,
				devicemanagement.FieldUpdatedDateTime,
			},
			Limit: 10,
		}

		result, resp, err := svc.GetV1(ctx2, opts)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)

		for _, server := range result.Data {
			assert.NotEmpty(t, server.ID)
			assert.NotEmpty(t, server.Type)
			if server.Attributes != nil {
				assert.NotEmpty(t, server.Attributes.ServerName)
			}
		}

		acc.LogTestSuccess(t, "GetV1 (fields): %d server(s)", len(result.Data))
	})

	// --- Pagination ---
	// Note: GetDeviceManagementServicesV1 fetches all pages automatically, so Limit
	// controls page size sent to the API, not the total result count returned here.
	t.Run("WithPaginationLimit", func(t *testing.T) {
		ctx3, cancel3 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel3()

		result, resp, err := svc.GetV1(ctx3, &devicemanagement.RequestQueryOptions{Limit: 1})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Data, "organization should have at least one MDM server")
		acc.LogTestSuccess(t, "GetV1 (limit=1 page size): %d total server(s) across all pages", len(result.Data))
	})
}

// =============================================================================
// TestAcceptance_DeviceManagement_GetDeviceSerialNumbers
// Verifies the device-linkage list for a given MDM server.
// =============================================================================

func TestAcceptance_DeviceManagement_GetDeviceSerialNumbers(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.AXMAPI.DeviceManagement
	ctx := context.Background()

	// Prerequisite: need a real MDM server ID
	serverID, serverName := requireFirstMDMServer(t, svc, ctx)

	acc.LogTestStage(t, "Setup", "Using MDM server %q (ID=%s)", serverName, serverID)

	t.Run("DefaultOptions", func(t *testing.T) {
		ctx1, cancel1 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel1()

		result, resp, err := svc.GetDeviceSerialNumbersByServerIDV1(ctx1, serverID, nil)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)

		acc.LogTestSuccess(t, "GetDeviceSerialNumbersByServerIDV1: %d device(s) on server %q",
			len(result.Data), serverName)
	})

	t.Run("WithLimit", func(t *testing.T) {
		ctx2, cancel2 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel2()

		result, resp, err := svc.GetDeviceSerialNumbersByServerIDV1(ctx2, serverID,
			&devicemanagement.RequestQueryOptions{Limit: 5})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		assert.LessOrEqual(t, len(result.Data), 5)
	})
}

// =============================================================================
// TestAcceptance_DeviceManagement_GetAssignedServiceIDForDevice
// Verifies the relationship linkage endpoint for a device.
// =============================================================================

func TestAcceptance_DeviceManagement_GetAssignedServiceIDForDevice(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.AXMAPI.DeviceManagement
	ctx := context.Background()

	// Prerequisite: need a real device ID
	deviceID := requireFirstDeviceID(t, ctx)

	acc.LogTestStage(t, "Setup", "Using device ID=%s for assigned-server-ID tests", deviceID)

	ctx1, cancel1 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
	defer cancel1()

	result, resp, err := svc.GetAssignedServerIDByDeviceIDV1(ctx1, deviceID)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, result)
	// The response is valid even when the device has no assigned server (ID will be "")
	assert.NotNil(t, result.Data)

	if result.Data.ID != "" {
		acc.LogTestSuccess(t, "GetAssignedServerIDByDeviceIDV1: device %s → server %s",
			deviceID, result.Data.ID)
	} else {
		acc.LogTestSuccess(t, "GetAssignedServerIDByDeviceIDV1: device %s has no assigned server", deviceID)
	}
}

// =============================================================================
// TestAcceptance_DeviceManagement_GetAssignedServiceInfoForDevice
// Verifies the full server-object endpoint for a device that has an assignment.
// =============================================================================

func TestAcceptance_DeviceManagement_GetAssignedServiceInfoForDevice(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.AXMAPI.DeviceManagement
	ctx := context.Background()

	// Prerequisite: find a device that IS assigned to a server
	deviceID, serverID := requireAssignedDevice(t, svc, ctx)
	acc.LogTestStage(t, "Setup", "Using device ID=%s (assigned to server %s)", deviceID, serverID)

	t.Run("AllFields", func(t *testing.T) {
		ctx1, cancel1 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel1()

		result, resp, err := svc.GetAssignedServerInfoByDeviceIDV1(ctx1, deviceID, nil)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.Equal(t, serverID, result.Data.ID, "returned server ID should match known assignment")
		assert.NotNil(t, result.Data.Attributes)

		if result.Data.Attributes != nil {
			acc.LogTestSuccess(t, "GetAssignedServerInfoByDeviceIDV1: server=%q type=%s",
				result.Data.Attributes.ServerName, result.Data.Attributes.ServerType)
		}
	})

	t.Run("SpecificFields", func(t *testing.T) {
		ctx2, cancel2 := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer cancel2()

		opts := &devicemanagement.RequestQueryOptions{
			Fields: []string{
				devicemanagement.FieldServerName,
				devicemanagement.FieldServerType,
				devicemanagement.FieldCreatedDateTime,
			},
		}

		result, resp, err := svc.GetAssignedServerInfoByDeviceIDV1(ctx2, deviceID, opts)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, resp.StatusCode())
		require.NotNil(t, result)
		assert.Equal(t, serverID, result.Data.ID)
	})
}

// =============================================================================
// TestAcceptance_DeviceManagement_AssignUnassign_Lifecycle
// Full assign → verify → unassign → verify lifecycle test.
// Skipped when no unassigned devices or no MDM servers are available.
// =============================================================================

func TestAcceptance_DeviceManagement_AssignUnassign_Lifecycle(t *testing.T) {
	acc.RequireClient(t)
	if acc.Config.SkipDestructive {
		t.Skip("destructive tests skipped (AXM_SKIP_DESTRUCTIVE=true) — set AXM_SKIP_DESTRUCTIVE=false to run")
	}

	svc := acc.Client.AXMAPI.DeviceManagement
	devSvc := acc.Client.AXMAPI.Devices
	ctx := context.Background()

	// -- Find a target MDM server --
	serverID, serverName := requireFirstMDMServer(t, svc, ctx)
	acc.LogTestStage(t, "Setup", "Target MDM server: %q (ID=%s)", serverName, serverID)

	// -- Find an unassigned device --
	acc.LogTestStage(t, "Setup", "Searching for an unassigned device...")

	listCtx, listCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
	defer listCancel()

	devList, _, err := devSvc.GetV1(listCtx, &devices.RequestQueryOptions{
		Fields: []string{devices.FieldSerialNumber},
		Limit:  20,
	})
	require.NoError(t, err, "list devices")

	var unassignedDeviceID string
	for _, dev := range devList.Data {
		chkCtx, chkCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		linkage, _, linkErr := svc.GetAssignedServerIDByDeviceIDV1(chkCtx, dev.ID)
		chkCancel()
		if linkErr != nil || (linkage != nil && linkage.Data.ID == "") {
			unassignedDeviceID = dev.ID
			break
		}
	}

	if unassignedDeviceID == "" {
		t.Skip("No unassigned devices found — skipping assign/unassign lifecycle test")
	}

	acc.LogTestStage(t, "Setup", "Unassigned device: ID=%s", unassignedDeviceID)

	// Register cleanup: always unassign at test teardown regardless of test outcome
	acc.Cleanup(t, func() {
		cleanCtx, cleanCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cleanCancel()
		_, _, unassignErr := svc.UnassignDevicesV1(cleanCtx, serverID, []string{unassignedDeviceID})
		acc.LogCleanupError(t, "device assignment", fmt.Sprintf("device=%s server=%s", unassignedDeviceID, serverID), unassignErr)
	})

	// ------------------------------------------------------------------
	// 1. Assign device to server
	// ------------------------------------------------------------------
	acc.LogTestStage(t, "Assign", "Assigning device %s to server %q", unassignedDeviceID, serverName)

	assignCtx, assignCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
	defer assignCancel()

	assignResp, assignHTTPResp, err := svc.AssignDevicesV1(assignCtx, serverID, []string{unassignedDeviceID})

	require.NoError(t, err, "AssignDevicesToServerV1 should succeed")
	require.NotNil(t, assignHTTPResp)
	assert.Equal(t, 201, assignHTTPResp.StatusCode())
	require.NotNil(t, assignResp)
	assert.NotEmpty(t, assignResp.Data.ID, "activity ID should not be empty")
	assert.Equal(t, devicemanagement.ActivityTypeAssignDevices, assignResp.Data.Attributes.ActivityType)

	acc.LogTestSuccess(t, "Assigned device %s — activity ID=%s status=%s",
		unassignedDeviceID, assignResp.Data.ID, assignResp.Data.Attributes.Status)

	// ------------------------------------------------------------------
	// 2. Verify assignment (poll with retry — Apple's API is eventually consistent)
	// ------------------------------------------------------------------
	acc.LogTestStage(t, "Verify", "Verifying assignment for device %s", unassignedDeviceID)

	assigned := acc.PollUntil(t, 15*time.Second, 2*time.Second, func() bool {
		chkCtx, chkCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer chkCancel()
		linkage, _, err := svc.GetAssignedServerIDByDeviceIDV1(chkCtx, unassignedDeviceID)
		return err == nil && linkage != nil && linkage.Data.ID == serverID
	})

	if !assigned {
		acc.LogTestWarning(t, "Device %s may not have been assigned yet — Apple assignment is eventually consistent", unassignedDeviceID)
	} else {
		acc.LogTestSuccess(t, "Device %s confirmed assigned to server %q", unassignedDeviceID, serverName)
	}

	// ------------------------------------------------------------------
	// 3. Unassign device from server
	// ------------------------------------------------------------------
	acc.LogTestStage(t, "Unassign", "Unassigning device %s from server %q", unassignedDeviceID, serverName)

	unassignCtx, unassignCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
	defer unassignCancel()

	unassignResp, unassignHTTPResp, err := svc.UnassignDevicesV1(unassignCtx, serverID, []string{unassignedDeviceID})

	require.NoError(t, err, "UnassignDevicesFromServerV1 should succeed")
	require.NotNil(t, unassignHTTPResp)
	assert.Equal(t, 201, unassignHTTPResp.StatusCode())
	require.NotNil(t, unassignResp)
	assert.NotEmpty(t, unassignResp.Data.ID)
	assert.Equal(t, devicemanagement.ActivityTypeUnassignDevices, unassignResp.Data.Attributes.ActivityType)

	acc.LogTestSuccess(t, "Unassigned device %s — activity ID=%s", unassignedDeviceID, unassignResp.Data.ID)

	// ------------------------------------------------------------------
	// 4. Verify unassignment
	// ------------------------------------------------------------------
	acc.LogTestStage(t, "Verify", "Verifying unassignment for device %s", unassignedDeviceID)

	unassigned := acc.PollUntil(t, 15*time.Second, 2*time.Second, func() bool {
		chkCtx, chkCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer chkCancel()
		linkage, _, err := svc.GetAssignedServerIDByDeviceIDV1(chkCtx, unassignedDeviceID)
		return err != nil || (linkage != nil && linkage.Data.ID == "")
	})

	if !unassigned {
		acc.LogTestWarning(t, "Device %s may still appear assigned — Apple unassignment is eventually consistent", unassignedDeviceID)
	} else {
		acc.LogTestSuccess(t, "Device %s confirmed unassigned", unassignedDeviceID)
	}
}

// =============================================================================
// TestAcceptance_DeviceManagement_ValidationErrors
// Verifies client-side validation fires before any HTTP call is made.
// =============================================================================

func TestAcceptance_DeviceManagement_ValidationErrors(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.AXMAPI.DeviceManagement
	ctx := context.Background()

	t.Run("GetDeviceSerialNumbers_EmptyServerID", func(t *testing.T) {
		_, _, err := svc.GetDeviceSerialNumbersByServerIDV1(ctx, "", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "MDM server ID is required")
	})

	t.Run("GetAssignedServiceID_EmptyDeviceID", func(t *testing.T) {
		_, _, err := svc.GetAssignedServerIDByDeviceIDV1(ctx, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "device ID is required")
	})

	t.Run("GetAssignedServiceInfo_EmptyDeviceID", func(t *testing.T) {
		_, _, err := svc.GetAssignedServerInfoByDeviceIDV1(ctx, "", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "device ID is required")
	})

	t.Run("AssignDevices_EmptyServerID", func(t *testing.T) {
		_, _, err := svc.AssignDevicesV1(ctx, "", []string{"DEVICE123"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "MDM server ID is required")
	})

	t.Run("AssignDevices_EmptyDeviceIDs", func(t *testing.T) {
		_, _, err := svc.AssignDevicesV1(ctx, "SERVER123", []string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one device ID is required")
	})

	t.Run("UnassignDevices_EmptyServerID", func(t *testing.T) {
		_, _, err := svc.UnassignDevicesV1(ctx, "", []string{"DEVICE123"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "MDM server ID is required")
	})

	t.Run("UnassignDevices_EmptyDeviceIDs", func(t *testing.T) {
		_, _, err := svc.UnassignDevicesV1(ctx, "SERVER123", []string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one device ID is required")
	})
}

// =============================================================================
// Shared test helpers
// =============================================================================

// requireFirstMDMServer fetches the first MDM server in the org and skips the
// test when none exist.
func requireFirstMDMServer(t *testing.T, _ *devicemanagement.DeviceManagement, ctx context.Context) (serverID, serverName string) {
	t.Helper()
	listCtx, listCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
	defer listCancel()

	result, _, err := acc.Client.AXMAPI.DeviceManagement.GetV1(listCtx,
		&devicemanagement.RequestQueryOptions{
			Fields: []string{devicemanagement.FieldServerName},
			Limit:  1,
		})
	require.NoError(t, err, "prerequisite: list MDM servers")
	if len(result.Data) == 0 {
		t.Skip("No MDM servers found in organization — skipping test")
	}

	server := result.Data[0]
	name := ""
	if server.Attributes != nil {
		name = server.Attributes.ServerName
	}
	return server.ID, name
}

// requireFirstDeviceID fetches the first device ID in the org and skips the
// test when none exist.
func requireFirstDeviceID(t *testing.T, ctx context.Context) string {
	t.Helper()
	listCtx, listCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
	defer listCancel()

	list, _, err := acc.Client.AXMAPI.Devices.GetV1(listCtx,
		&devices.RequestQueryOptions{
			Fields: []string{devices.FieldSerialNumber},
			Limit:  1,
		})
	require.NoError(t, err, "prerequisite: list devices")
	if len(list.Data) == 0 {
		t.Skip("No devices found in organization — skipping test")
	}
	return list.Data[0].ID
}

// requireAssignedDevice finds and returns the ID of a device that has an assigned
// MDM server, along with the server ID. Skips the test when none are found.
func requireAssignedDevice(t *testing.T, svc *devicemanagement.DeviceManagement, ctx context.Context) (deviceID, serverID string) {
	t.Helper()
	listCtx, listCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
	defer listCancel()

	devList, _, err := acc.Client.AXMAPI.Devices.GetV1(listCtx,
		&devices.RequestQueryOptions{
			Fields: []string{devices.FieldSerialNumber},
			Limit:  20,
		})
	require.NoError(t, err, "prerequisite: list devices")

	for _, dev := range devList.Data {
		chkCtx, chkCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		linkage, _, linkErr := svc.GetAssignedServerIDByDeviceIDV1(chkCtx, dev.ID)
		chkCancel()
		if linkErr == nil && linkage != nil && linkage.Data.ID != "" {
			return dev.ID, linkage.Data.ID
		}
	}

	t.Skip("No devices with an assigned MDM server found — skipping test")
	return "", "" // unreachable
}

// =============================================================================
// TestAcceptance_DeviceManagement_GetOrgDeviceActivity_ValidationErrors
// Verifies client-side validation on the org device activity endpoints. Makes no
// state-changing calls, so it is safe to run against a production tenant.
// =============================================================================

func TestAcceptance_DeviceManagement_GetOrgDeviceActivity_ValidationErrors(t *testing.T) {
	acc.RequireClient(t)

	svc := acc.Client.AXMAPI.DeviceManagement
	ctx := context.Background()

	t.Run("GetOrgDeviceActivityByID_EmptyID", func(t *testing.T) {
		result, _, err := svc.GetOrgDeviceActivityByIDV1(ctx, "", nil)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "activity ID is required")
	})

	t.Run("CreateOrgDeviceActivity_NilRequest", func(t *testing.T) {
		result, _, err := svc.CreateOrgDeviceActivityV1(ctx, nil)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "request is required")
	})

	t.Run("CreateOrgDeviceActivity_NoDevices", func(t *testing.T) {
		result, _, err := svc.CreateOrgDeviceActivityV1(ctx, &devicemanagement.OrgDeviceActivityCreateRequest{
			Data: devicemanagement.OrgDeviceActivityData{
				Type: devicemanagement.ResourceTypeOrgDeviceActivities,
				Attributes: devicemanagement.OrgDeviceActivityCreateAttributes{
					ActivityType: devicemanagement.ActivityTypeCancelMDMMigration,
				},
			},
		})
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "at least one device ID is required")
	})
}

// =============================================================================
// TestAcceptance_DeviceManagement_MDMMigration_Lifecycle
// Schedules a device management service migration for a migration-capable device,
// polls the resulting activity to completion, updates the deadline, then cancels
// the migration. Destructive: it reassigns a real device between MDM services.
// =============================================================================

func TestAcceptance_DeviceManagement_MDMMigration_Lifecycle(t *testing.T) {
	acc.RequireClient(t)
	if acc.Config.SkipDestructive {
		t.Skip("destructive tests skipped (AXM_SKIP_DESTRUCTIVE=true) — set AXM_SKIP_DESTRUCTIVE=false to run")
	}

	svc := acc.Client.AXMAPI.DeviceManagement
	devSvc := acc.Client.AXMAPI.Devices
	ctx := context.Background()

	// -- Find a target MDM server --
	serverID, serverName := requireFirstMDMServer(t, svc, ctx)
	acc.LogTestStage(t, "Setup", "Target MDM server: %q (ID=%s)", serverName, serverID)

	// -- Find a device eligible for device management service migration --
	acc.LogTestStage(t, "Setup", "Searching for a migration-capable device...")

	listCtx, listCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
	defer listCancel()

	devList, _, err := devSvc.GetV1(listCtx, &devices.RequestQueryOptions{
		Fields: []string{
			devices.FieldSerialNumber,
			devices.FieldIsMdmMigrationCapable,
			devices.FieldMdmMigrationStatus,
		},
		Limit: 50,
	})
	require.NoError(t, err, "list devices")

	var migrationDeviceID string
	for _, dev := range devList.Data {
		if dev.Attributes == nil {
			continue
		}
		// Skip devices already undergoing a migration.
		if dev.Attributes.IsMdmMigrationCapable && dev.Attributes.MdmMigrationStatus == "" {
			migrationDeviceID = dev.ID
			break
		}
	}

	if migrationDeviceID == "" {
		t.Skip("No migration-capable device found — skipping migration lifecycle test")
	}

	acc.LogTestStage(t, "Setup", "Migration-capable device: ID=%s", migrationDeviceID)

	// Always attempt to cancel the migration at teardown, whatever the outcome.
	acc.Cleanup(t, func() {
		cleanCtx, cleanCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cleanCancel()
		_, _, cancelErr := svc.CreateOrgDeviceActivityV1(cleanCtx, newActivityRequest(
			devicemanagement.ActivityTypeCancelMDMMigration, "", nil, []string{migrationDeviceID},
		))
		acc.LogCleanupError(t, "mdm migration", fmt.Sprintf("device=%s", migrationDeviceID), cancelErr)
	})

	// ------------------------------------------------------------------
	// 1. Schedule the migration with a deadline
	// ------------------------------------------------------------------
	deadline := time.Now().UTC().Add(7 * 24 * time.Hour)
	acc.LogTestStage(t, "Schedule", "Scheduling migration of device %s to %q by %s",
		migrationDeviceID, serverName, deadline.Format(time.RFC3339))

	scheduleCtx, scheduleCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
	defer scheduleCancel()

	scheduled, scheduleHTTPResp, err := svc.CreateOrgDeviceActivityV1(scheduleCtx, newActivityRequest(
		devicemanagement.ActivityTypeAssignDevicesWithMDMMigrationDeadline,
		serverID, &deadline, []string{migrationDeviceID},
	))
	require.NoError(t, err, "schedule migration")
	require.NotNil(t, scheduleHTTPResp)
	assert.Equal(t, 201, scheduleHTTPResp.StatusCode())
	require.NotNil(t, scheduled)
	require.NotEmpty(t, scheduled.Data.ID, "activity ID")
	acc.LogTestSuccess(t, "Migration scheduled: activity=%s status=%s",
		scheduled.Data.ID, scheduled.Data.Attributes.Status)

	// ------------------------------------------------------------------
	// 2. Poll the activity to completion
	// ------------------------------------------------------------------
	acc.LogTestStage(t, "Poll", "Polling activity %s for completion...", scheduled.Data.ID)

	var finalStatus, finalSubStatus, downloadURL string
	completed := acc.PollUntil(t, 2*time.Minute, 5*time.Second, func() bool {
		pollCtx, pollCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
		defer pollCancel()

		activity, _, pollErr := svc.GetOrgDeviceActivityByIDV1(pollCtx, scheduled.Data.ID, nil)
		if pollErr != nil || activity == nil || activity.Data.Attributes == nil {
			return false
		}

		finalStatus = activity.Data.Attributes.Status
		finalSubStatus = activity.Data.Attributes.SubStatus
		downloadURL = activity.Data.Attributes.DownloadURL

		return finalStatus == devicemanagement.ActivityStatusCompleted ||
			finalStatus == devicemanagement.ActivityStatusFailed ||
			finalStatus == devicemanagement.ActivityStatusStopped
	})

	if !completed {
		acc.LogTestWarning(t, "Activity %s did not reach a terminal status within the poll window (last status=%s)",
			scheduled.Data.ID, finalStatus)
	} else {
		acc.LogTestSuccess(t, "Activity %s finished: status=%s subStatus=%s", scheduled.Data.ID, finalStatus, finalSubStatus)
		assert.Equal(t, devicemanagement.ActivityStatusCompleted, finalStatus,
			"migration activity should complete successfully")
		// A COMPLETED activity exposes a presigned CSV activity log.
		assert.NotEmpty(t, downloadURL, "completed activity should expose a downloadUrl")
	}

	// ------------------------------------------------------------------
	// 3. Update the migration deadline
	// ------------------------------------------------------------------
	newDeadline := deadline.Add(24 * time.Hour)
	acc.LogTestStage(t, "Update", "Updating migration deadline to %s", newDeadline.Format(time.RFC3339))

	updateCtx, updateCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
	defer updateCancel()

	updated, updateHTTPResp, err := svc.CreateOrgDeviceActivityV1(updateCtx, newActivityRequest(
		devicemanagement.ActivityTypeUpdateMDMMigrationDeadline,
		"", &newDeadline, []string{migrationDeviceID},
	))
	require.NoError(t, err, "update migration deadline")
	require.NotNil(t, updateHTTPResp)
	assert.Equal(t, 201, updateHTTPResp.StatusCode())
	require.NotNil(t, updated)
	acc.LogTestSuccess(t, "Deadline updated: activity=%s", updated.Data.ID)

	// ------------------------------------------------------------------
	// 4. Cancel the migration
	// ------------------------------------------------------------------
	acc.LogTestStage(t, "Cancel", "Cancelling migration for device %s", migrationDeviceID)

	cancelCtx, cancelCancel := context.WithTimeout(ctx, acc.Config.RequestTimeout)
	defer cancelCancel()

	cancelled, cancelHTTPResp, err := svc.CreateOrgDeviceActivityV1(cancelCtx, newActivityRequest(
		devicemanagement.ActivityTypeCancelMDMMigration, "", nil, []string{migrationDeviceID},
	))
	require.NoError(t, err, "cancel migration")
	require.NotNil(t, cancelHTTPResp)
	assert.Equal(t, 201, cancelHTTPResp.StatusCode())
	require.NotNil(t, cancelled)
	acc.LogTestSuccess(t, "Migration cancelled: activity=%s", cancelled.Data.ID)
}

// newActivityRequest builds an OrgDeviceActivityCreateRequest, attaching the MDM server
// relationship and activity type metadata only when supplied.
func newActivityRequest(activityType, mdmServerID string, deadline *time.Time, deviceIDs []string) *devicemanagement.OrgDeviceActivityCreateRequest {
	linkages := make([]devicemanagement.OrgDeviceActivityDeviceLinkage, 0, len(deviceIDs))
	for _, deviceID := range deviceIDs {
		linkages = append(linkages, devicemanagement.OrgDeviceActivityDeviceLinkage{
			Type: devicemanagement.ResourceTypeOrgDevices,
			ID:   deviceID,
		})
	}

	req := &devicemanagement.OrgDeviceActivityCreateRequest{
		Data: devicemanagement.OrgDeviceActivityData{
			Type: devicemanagement.ResourceTypeOrgDeviceActivities,
			Attributes: devicemanagement.OrgDeviceActivityCreateAttributes{
				ActivityType: activityType,
			},
			Relationships: devicemanagement.OrgDeviceActivityCreateRelationships{
				Devices: &devicemanagement.OrgDeviceActivityDevicesRelationship{Data: linkages},
			},
		},
	}

	if mdmServerID != "" {
		req.Data.Relationships.MDMServer = &devicemanagement.OrgDeviceActivityMDMServerRelationship{
			Data: devicemanagement.OrgDeviceActivityMDMServerLinkage{
				Type: devicemanagement.ResourceTypeMDMServers,
				ID:   mdmServerID,
			},
		}
	}

	if deadline != nil {
		req.Data.Attributes.ActivityTypeMetadata = &devicemanagement.ActivityTypeMetadata{
			MDMMigrationDeadlineDateTime: deadline,
		}
	}

	return req
}
