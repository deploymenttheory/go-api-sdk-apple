package devicemanagement

// JSON:API resource type identifiers used in orgDeviceActivities request payloads
const (
	ResourceTypeOrgDeviceActivities = "orgDeviceActivities"
	ResourceTypeMDMServers          = "mdmServers"
	ResourceTypeOrgDevices          = "orgDevices"
)

// Activity type constants for the orgDeviceActivities activityType attribute.
//
// ActivityTypeReleaseDevices is available on the Apple Business API only (2.4+);
// the Apple School Manager API supports the other five.
// https://developer.apple.com/documentation/applebusinessapi/orgdeviceactivitytype
const (
	ActivityTypeAssignDevices                         = "ASSIGN_DEVICES"
	ActivityTypeUnassignDevices                       = "UNASSIGN_DEVICES"
	ActivityTypeAssignDevicesWithMDMMigrationDeadline = "ASSIGN_DEVICES_WITH_MDM_MIGRATION_DEADLINE"
	ActivityTypeUpdateMDMMigrationDeadline            = "UPDATE_MDM_MIGRATION_DEADLINE"
	ActivityTypeCancelMDMMigration                    = "CANCEL_MDM_MIGRATION"
	ActivityTypeReleaseDevices                        = "RELEASE_DEVICES"
)

// Activity status constants for the orgDeviceActivities status attribute
const (
	ActivityStatusCompleted  = "COMPLETED"
	ActivityStatusInProgress = "IN_PROGRESS"
	ActivityStatusStopped    = "STOPPED"
	ActivityStatusFailed     = "FAILED"
)

// Activity sub-status constants for the orgDeviceActivities subStatus attribute
const (
	ActivitySubStatusSubmitted                     = "SUBMITTED"
	ActivitySubStatusPreProcessing                 = "PRE_PROCESSING"
	ActivitySubStatusPending                       = "PENDING"
	ActivitySubStatusProcessing                    = "PROCESSING"
	ActivitySubStatusPostProcessing                = "POST_PROCESSING"
	ActivitySubStatusStopping                      = "STOPPING"
	ActivitySubStatusCompletedWithSuccess          = "COMPLETED_WITH_SUCCESS"
	ActivitySubStatusCompletedWithError            = "COMPLETED_WITH_ERROR"
	ActivitySubStatusCompletedWithFailure          = "COMPLETED_WITH_FAILURE"
	ActivitySubStatusCompletedPostProcessingFailed = "COMPLETED_POST_PROCESSING_FAILED"
)

// OrgDeviceActivity field constants for the fields[orgDeviceActivities] query parameter.
// Prefixed with Activity to avoid colliding with the MDM server field constants below.
const (
	FieldActivityStatus            = "status"
	FieldActivitySubStatus         = "subStatus"
	FieldActivityCreatedDateTime   = "createdDateTime"
	FieldActivityCompletedDateTime = "completedDateTime"
	FieldActivityDownloadURL       = "downloadUrl"
)

// MDM Server field constants for field selection
const (
	FieldServerName             = "serverName"
	FieldServerType             = "serverType"
	FieldEnableMdmDisownFlag    = "enableMdmDisownFlag"
	FieldDefaultProductFamilies = "defaultProductFamilies"
	FieldStatus                 = "status"
	FieldDeviceCount            = "deviceCount"
	FieldLastConnectedDateTime  = "lastConnectedDateTime"
	FieldLastConnectedIp        = "lastConnectedIp"
	FieldCreatedDateTime        = "createdDateTime"
	FieldUpdatedDateTime        = "updatedDateTime"
	FieldDevices                = "devices"
)

// MDM server status constants
const (
	MDMServerStatusActive   = "ACTIVE"
	MDMServerStatusInactive = "INACTIVE"
)
