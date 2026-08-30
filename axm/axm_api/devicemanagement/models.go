package devicemanagement

import "time"

// ====== MDM SERVER TYPES ======

// MDMServer represents an MDM server in the Apple Business Manager system
type MDMServer struct {
	ID            string                  `json:"id"`
	Type          string                  `json:"type"`
	Attributes    *MDMServerAttributes    `json:"attributes,omitempty"`
	Relationships *MDMServerRelationships `json:"relationships,omitempty"`
}

// MDMServerAttributes contains the MDM server attributes
type MDMServerAttributes struct {
	ServerName             string     `json:"serverName,omitempty"`
	ServerType             string     `json:"serverType,omitempty"`
	EnableMdmDisownFlag    bool       `json:"enableMdmDisownFlag,omitempty"`
	DefaultProductFamilies []string   `json:"defaultProductFamilies,omitempty"`
	Status                 string     `json:"status,omitempty"`
	DeviceCount            int        `json:"deviceCount,omitempty"`
	LastConnectedDateTime  *time.Time `json:"lastConnectedDateTime,omitempty"`
	LastConnectedIp        string     `json:"lastConnectedIp,omitempty"`
	CreatedDateTime        *time.Time `json:"createdDateTime,omitempty"`
	UpdatedDateTime        *time.Time `json:"updatedDateTime,omitempty"`
	Devices                []string   `json:"devices,omitempty"`
}

// MDMServerRelationships contains the MDM server relationships
type MDMServerRelationships struct {
	Devices *MDMServerDevicesRelationship `json:"devices,omitempty"`
}

// MDMServerDevicesRelationship contains the devices relationship links
type MDMServerDevicesRelationship struct {
	Links *MDMServerDevicesLinks `json:"links,omitempty"`
}

// MDMServerDevicesLinks contains the navigation links for devices
type MDMServerDevicesLinks struct {
	Self string `json:"self,omitempty"`
}

// ResponseMDMServers represents the response for getting MDM servers
type ResponseMDMServers struct {
	Data  []MDMServer `json:"data"`
	Meta  *Meta       `json:"meta,omitempty"`
	Links *Links      `json:"links,omitempty"`
}

// MDMServerResponse represents the response for getting a single MDM server
type MDMServerResponse struct {
	Data  MDMServer `json:"data"`
	Links *Links    `json:"links,omitempty"`
}

// ====== MDM SERVER REQUEST TYPES ======

// MDMServerCertificate represents a server certificate for MDM server creation
type MDMServerCertificate struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// MDMServerCreateRequest is the request body for creating a new MDM server
type MDMServerCreateRequest struct {
	Data MDMServerCreateRequestData `json:"data"`
}

// MDMServerCreateRequestData is the data object for an MDM server create request
type MDMServerCreateRequestData struct {
	Type       string                           `json:"type"` // must be "mdmServers"
	Attributes MDMServerCreateRequestAttributes `json:"attributes"`
}

// MDMServerCreateRequestAttributes contains the attributes for creating an MDM server.
// ServerName and ServerCertificate are required.
type MDMServerCreateRequestAttributes struct {
	ServerName          string               `json:"serverName"`
	ServerCertificate   MDMServerCertificate `json:"serverCertificate"`
	EnableMdmDisownFlag bool                 `json:"enableMdmDisownFlag,omitempty"`
}

// MDMServerUpdateRequest is the request body for updating an MDM server
type MDMServerUpdateRequest struct {
	Data MDMServerUpdateRequestData `json:"data"`
}

// MDMServerUpdateRequestData is the data object for an MDM server update request
type MDMServerUpdateRequestData struct {
	Type       string                           `json:"type"` // must be "mdmServers"
	ID         string                           `json:"id"`
	Attributes MDMServerUpdateRequestAttributes `json:"attributes"`
}

// MDMServerUpdateRequestAttributes contains the attributes for updating an MDM server.
// Only provided fields are changed.
type MDMServerUpdateRequestAttributes struct {
	ServerName             string   `json:"serverName,omitempty"`
	EnableMdmDisownFlag    *bool    `json:"enableMdmDisownFlag,omitempty"`
	DefaultProductFamilies []string `json:"defaultProductFamilies,omitempty"`
}

// ====== DEVICE LINKAGE TYPES ======

// ResponseMDMServerDevicesLinkages represents the response for getting device linkages for an MDM server
type ResponseMDMServerDevicesLinkages struct {
	Data  []MDMServerDeviceLinkage `json:"data"`
	Links *Links                   `json:"links,omitempty"`
	Meta  *Meta                    `json:"meta,omitempty"`
}

// MDMServerDeviceLinkage represents a device linkage in the MDM server relationships
type MDMServerDeviceLinkage struct {
	Type string `json:"type"` // Should be "orgDevices"
	ID   string `json:"id"`   // Device ID
}

// ResponseOrgDeviceAssignedServerLinkage represents the response for getting assigned server linkage
type ResponseOrgDeviceAssignedServerLinkage struct {
	Data  OrgDeviceAssignedServerLinkage `json:"data"`
	Links *AssignedServerLinks           `json:"links,omitempty"`
}

// OrgDeviceAssignedServerLinkage represents the linkage between a device and its assigned server
type OrgDeviceAssignedServerLinkage struct {
	Type string `json:"type"` // Should be "mdmServers"
	ID   string `json:"id"`   // MDM Server ID
}

// AssignedServerLinks contains linkage navigation links
type AssignedServerLinks struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}

// ====== DEVICE ACTIVITY TYPES ======

// OrgDeviceActivity represents a device activity (assign/unassign operations)
type OrgDeviceActivity struct {
	ID         string                       `json:"id"`
	Type       string                       `json:"type"`
	Attributes *OrgDeviceActivityAttributes `json:"attributes,omitempty"`
	Links      *OrgDeviceActivityLinks      `json:"links,omitempty"`
}

// OrgDeviceActivityAttributes contains the activity attributes.
//
// CompletedDateTime and DownloadURL are populated only once Status is COMPLETED.
// DownloadURL is a presigned URL for the activity log in CSV format.
type OrgDeviceActivityAttributes struct {
	Status            string     `json:"status,omitempty"`
	SubStatus         string     `json:"subStatus,omitempty"`
	CreatedDateTime   *time.Time `json:"createdDateTime,omitempty"`
	CompletedDateTime *time.Time `json:"completedDateTime,omitempty"`
	DownloadURL       string     `json:"downloadUrl,omitempty"`
	ActivityType      string     `json:"activityType,omitempty"`
}

// OrgDeviceActivityLinks contains activity navigation links
type OrgDeviceActivityLinks struct {
	Self string `json:"self,omitempty"`
}

// ResponseOrgDeviceActivity represents the response for creating an org device activity
type ResponseOrgDeviceActivity struct {
	Data  OrgDeviceActivity `json:"data"`
	Links *Links            `json:"links,omitempty"`
}

// ====== DEVICE ACTIVITY REQUEST TYPES ======

// OrgDeviceActivityCreateRequest represents the request for creating a device activity
type OrgDeviceActivityCreateRequest struct {
	Data OrgDeviceActivityData `json:"data"`
}

// OrgDeviceActivityData contains the activity data for the request
type OrgDeviceActivityData struct {
	Type          string                               `json:"type"`
	Attributes    OrgDeviceActivityCreateAttributes    `json:"attributes"`
	Relationships OrgDeviceActivityCreateRelationships `json:"relationships"`
}

// OrgDeviceActivityCreateAttributes contains the activity creation attributes
type OrgDeviceActivityCreateAttributes struct {
	ActivityType         string                `json:"activityType"`
	ActivityTypeMetadata *ActivityTypeMetadata `json:"activityTypeMetadata,omitempty"`
}

// ActivityTypeMetadata carries the additional data required by certain activity types.
// It is required for ASSIGN_DEVICES_WITH_MDM_MIGRATION_DEADLINE and
// UPDATE_MDM_MIGRATION_DEADLINE, and unused by the other activity types.
type ActivityTypeMetadata struct {
	// MDMMigrationDeadlineDateTime is the deadline by which a device needs to
	// complete its device management service migration. The API rejects deadlines
	// more than 90 days in the future.
	MDMMigrationDeadlineDateTime *time.Time `json:"mdmMigrationDeadlineDateTime,omitempty"`
}

// OrgDeviceActivityCreateRelationships contains the relationships for activity creation
type OrgDeviceActivityCreateRelationships struct {
	MDMServer *OrgDeviceActivityMDMServerRelationship `json:"mdmServer,omitempty"`
	Devices   *OrgDeviceActivityDevicesRelationship   `json:"devices,omitempty"`
}

// OrgDeviceActivityMDMServerRelationship represents the MDM server relationship
type OrgDeviceActivityMDMServerRelationship struct {
	Data OrgDeviceActivityMDMServerLinkage `json:"data"`
}

// OrgDeviceActivityMDMServerLinkage represents the MDM server linkage
type OrgDeviceActivityMDMServerLinkage struct {
	Type string `json:"type"` // Should be "mdmServers"
	ID   string `json:"id"`   // MDM Server ID
}

// OrgDeviceActivityDevicesRelationship represents the devices relationship
type OrgDeviceActivityDevicesRelationship struct {
	Data []OrgDeviceActivityDeviceLinkage `json:"data"`
}

// OrgDeviceActivityDeviceLinkage represents a device linkage
type OrgDeviceActivityDeviceLinkage struct {
	Type string `json:"type"` // Should be "orgDevices"
	ID   string `json:"id"`   // Device ID
}

// ====== SHARED TYPES ======

// Meta represents pagination metadata
type Meta struct {
	Paging *Paging `json:"paging,omitempty"`
}

// Paging contains pagination information
type Paging struct {
	Total      int    `json:"total,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`
}

// Links contains navigation links for API responses
type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Last  string `json:"last,omitempty"`
}

// RequestQueryOptions represents the query parameters for getting MDM servers
type RequestQueryOptions struct {
	// Field selection - fields to return for mdmServers
	// Possible values: serverName, serverType, createdDateTime, updatedDateTime, devices
	Fields []string `json:"fields,omitempty"`

	// Limit the number of included related resources to return (max 1000)
	Limit int `json:"limit,omitempty"`
}

// GetOrgDeviceActivityQueryOptions represents the query parameters for
// GetOrgDeviceActivityByIDV1.
type GetOrgDeviceActivityQueryOptions struct {
	// Fields specifies which fields to return. Use FieldActivity* constants.
	Fields []string
}
