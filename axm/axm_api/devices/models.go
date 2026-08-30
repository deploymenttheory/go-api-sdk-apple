package devices

import "time"

// Shared types for pagination and links
type Meta struct {
	Paging *Paging `json:"paging,omitempty"`
}

type Paging struct {
	Total      int    `json:"total,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`
}

type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Last  string `json:"last,omitempty"`
}

// OrgDevice represents a device in the Apple Business Manager system based on the API specification
type OrgDevice struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Attributes *OrgDeviceAttributes `json:"attributes,omitempty"`
}

// OrgDeviceAttributes contains the device attributes
type OrgDeviceAttributes struct {
	SerialNumber        string     `json:"serialNumber,omitempty"`
	AddedToOrgDateTime  *time.Time `json:"addedToOrgDateTime,omitempty"`
	UpdatedDateTime     *time.Time `json:"updatedDateTime,omitempty"`
	DeviceModel         string     `json:"deviceModel,omitempty"`
	ProductFamily       string     `json:"productFamily,omitempty"`
	ProductType         string     `json:"productType,omitempty"`
	DeviceCapacity      string     `json:"deviceCapacity,omitempty"`
	PartNumber          string     `json:"partNumber,omitempty"`
	OrderNumber         string     `json:"orderNumber,omitempty"`
	Color               string     `json:"color,omitempty"`
	Status              string     `json:"status,omitempty"`
	OrderDateTime       *time.Time `json:"orderDateTime,omitempty"`
	IMEI                []string   `json:"imei,omitempty"`
	MEID                []string   `json:"meid,omitempty"`
	EID                 string     `json:"eid,omitempty"`
	WiFiMACAddress      string     `json:"wifiMacAddress,omitempty"`
	BluetoothMACAddress string     `json:"bluetoothMacAddress,omitempty"`
	EthernetMACAddress  []string   `json:"ethernetMacAddress,omitempty"`
	PurchaseSourceId    string     `json:"purchaseSourceId,omitempty"`
	PurchaseSourceType  string     `json:"purchaseSourceType,omitempty"`
	AssignedServer      string     `json:"assignedServer,omitempty"`

	// ReleasedFromOrgDateTime is set once the device has been released from the
	// organization. Only populated when querying a single device.
	ReleasedFromOrgDateTime *time.Time `json:"releasedFromOrgDateTime,omitempty"`

	// Device management service migration attributes (read-only).
	IsMdmMigrationCapable        bool       `json:"isMdmMigrationCapable,omitempty"`
	MdmMigrationStatus           string     `json:"mdmMigrationStatus,omitempty"`
	MdmMigrationDeadlineDateTime *time.Time `json:"mdmMigrationDeadlineDateTime,omitempty"`
}

// OrgDeviceResponse represents the response for a single device
type OrgDeviceResponse struct {
	Data OrgDevice `json:"data"`
}

// OrgDeviceFilter represents filter options for organization devices
type OrgDeviceFilter struct {
	SerialNumber  string `json:"serialNumber,omitempty"`
	DeviceModel   string `json:"deviceModel,omitempty"`
	ProductFamily string `json:"productFamily,omitempty"`
	Color         string `json:"color,omitempty"`
	Status        string `json:"status,omitempty"`
}

// OrgDevicesResponse represents the response for getting organization devices
type OrgDevicesResponse struct {
	Data  []OrgDevice `json:"data"`
	Meta  *Meta       `json:"meta,omitempty"`
	Links *Links      `json:"links,omitempty"`
}

// RequestQueryOptions represents the query parameters for getting organization devices
type RequestQueryOptions struct {
	// Field selection - fields to return for orgDevices
	// Possible values: serialNumber, addedToOrgDateTime, releasedFromOrgDateTime,
	// updatedDateTime, deviceModel, productFamily, productType, deviceCapacity, partNumber,
	// orderNumber, color, status, orderDateTime, imei, meid, eid, wifiMacAddress,
	// bluetoothMacAddress, ethernetMacAddress, purchaseSourceId, purchaseSourceType,
	// assignedServer, appleCareCoverage, isMdmMigrationCapable, mdmMigrationStatus,
	// mdmMigrationDeadlineDateTime
	Fields []string `json:"fields,omitempty"`

	// Limit the number of included related resources to return (max 1000)
	Limit int `json:"limit,omitempty"`
}

// AppleCareCoverage represents AppleCare coverage for a device
type AppleCareCoverage struct {
	ID         string                       `json:"id"`
	Type       string                       `json:"type"`
	Attributes *AppleCareCoverageAttributes `json:"attributes,omitempty"`
}

// AppleCareCoverageAttributes contains the AppleCare coverage attributes
type AppleCareCoverageAttributes struct {
	Status                 string     `json:"status,omitempty"`
	PaymentType            string     `json:"paymentType,omitempty"`
	Description            string     `json:"description,omitempty"`
	AgreementNumber        *string    `json:"agreementNumber,omitempty"`
	StartDateTime          *time.Time `json:"startDateTime,omitempty"`
	EndDateTime            *time.Time `json:"endDateTime,omitempty"`
	IsRenewable            bool       `json:"isRenewable,omitempty"`
	IsCanceled             bool       `json:"isCanceled,omitempty"`
	ContractCancelDateTime *time.Time `json:"contractCancelDateTime,omitempty"`
}

// AppleCareCoverageResponse represents the response for getting AppleCare coverage
type AppleCareCoverageResponse struct {
	Data  []AppleCareCoverage `json:"data"`
	Links *Links              `json:"links,omitempty"`
	Meta  *Meta               `json:"meta,omitempty"`
}
