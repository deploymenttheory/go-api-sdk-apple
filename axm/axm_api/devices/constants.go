package devices

// ProfileStatus constants for device profile status
const (
	ProfileStatusEmpty    = "empty"
	ProfileStatusAssigned = "assigned"
	ProfileStatusPushed   = "pushed"
	ProfileStatusRemoved  = "removed"
)

// DeviceModel constants for common device models
const (
	ModeliPhone     = "iPhone"
	ModeliPad       = "iPad"
	ModelMac        = "Mac"
	ModelAppleTV    = "Apple TV"
	ModelAppleWatch = "Apple Watch"
)

// DeviceColor constants for device colors
const (
	ColorSilver    = "Silver"
	ColorGold      = "Gold"
	ColorSpaceGray = "Space Gray"
	ColorRoseGold  = "Rose Gold"
	ColorBlack     = "Black"
	ColorWhite     = "White"
	ColorRed       = "Red"
	ColorBlue      = "Blue"
	ColorGreen     = "Green"
	ColorYellow    = "Yellow"
	ColorPurple    = "Purple"
)

// OrgDevice field constants for field selection
const (
	FieldSerialNumber        = "serialNumber"
	FieldAddedToOrgDateTime  = "addedToOrgDateTime"
	FieldUpdatedDateTime     = "updatedDateTime"
	FieldDeviceModel         = "deviceModel"
	FieldProductFamily       = "productFamily"
	FieldProductType         = "productType"
	FieldDeviceCapacity      = "deviceCapacity"
	FieldPartNumber          = "partNumber"
	FieldOrderNumber         = "orderNumber"
	FieldColor               = "color"
	FieldStatus              = "status"
	FieldOrderDateTime       = "orderDateTime"
	FieldIMEI                = "imei"
	FieldMEID                = "meid"
	FieldEID                 = "eid"
	FieldWiFiMACAddress      = "wifiMacAddress"
	FieldBluetoothMACAddress = "bluetoothMacAddress"
	FieldEthernetMACAddress  = "ethernetMacAddress"
	FieldPurchaseSourceId    = "purchaseSourceId"
	FieldPurchaseSourceType  = "purchaseSourceType"
	FieldAssignedServer      = "assignedServer"
	FieldAppleCareCoverage   = "appleCareCoverage"

	FieldReleasedFromOrgDateTime = "releasedFromOrgDateTime"

	FieldIsMdmMigrationCapable        = "isMdmMigrationCapable"
	FieldMdmMigrationStatus           = "mdmMigrationStatus"
	FieldMdmMigrationDeadlineDateTime = "mdmMigrationDeadlineDateTime"
)

// MdmMigrationStatus constants for the mdmMigrationStatus field. Present only when a
// device management service migration has been requested for the device.
// https://developer.apple.com/documentation/applebusinessapi/mdmmigrationstatus
const (
	MdmMigrationStatusRequested = "REQUESTED"
	MdmMigrationStatusStarted   = "STARTED"
	MdmMigrationStatusSuccess   = "SUCCESS"
	MdmMigrationStatusFailed    = "FAILED"
)

// Device status constants
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

// Product family constants
const (
	ProductFamilyiPhone = "iPhone"
	ProductFamilyiPad   = "iPad"
	ProductFamilyMac    = "Mac"
)

// AppleCare coverage field constants for field selection
const (
	FieldAppleCareStatus                 = "status"
	FieldAppleCarePaymentType            = "paymentType"
	FieldAppleCareDescription            = "description"
	FieldAppleCareAgreementNumber        = "agreementNumber"
	FieldAppleCareStartDateTime          = "startDateTime"
	FieldAppleCareEndDateTime            = "endDateTime"
	FieldAppleCareIsRenewable            = "isRenewable"
	FieldAppleCareIsCanceled             = "isCanceled"
	FieldAppleCareContractCancelDateTime = "contractCancelDateTime"
)

// AppleCare coverage status constants
const (
	AppleCareStatusActive   = "ACTIVE"
	AppleCareStatusInactive = "INACTIVE"
	AppleCareStatusExpired  = "EXPIRED"
)

// AppleCare payment type constants
const (
	PaymentTypeNone            = "NONE"
	PaymentTypeSubscription    = "SUBSCRIPTION"
	PaymentTypeABESubscription = "ABE_SUBSCRIPTION"
)
