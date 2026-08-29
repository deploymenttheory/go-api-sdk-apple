package firmware

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/client"
	"github.com/deploymenttheory/go-sdk-appleservices/apple_update_cdn/constants"
	"resty.dev/v3"
)

// FirmwareService handles firmware discovery via the ipsw.me API.
//
// ipsw.me aggregates Apple's CDN firmware metadata and provides a stable
// JSON API for querying available IPSW restore files across Mac, iPhone, and iPad.
type FirmwareService struct {
	client client.Client
}

// NewService creates a new firmware service.
func NewService(c client.Client) *FirmwareService {
	return &FirmwareService{client: c}
}

// fetchAllFirmware retrieves the raw condensed firmware list from the ipsw.me v3
// API without any device filtering. The result includes every device class
// (Mac, iPhone, iPad, iPod, etc.) that Apple ships.
func (s *FirmwareService) fetchAllFirmware(ctx context.Context) (*AllFirmwaresV3Response, *resty.Response, error) {
	var result AllFirmwaresV3Response

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointFirmwaresAllV3)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListAllFirmwareV3 retrieves every device model and its complete firmware
// history from the ipsw.me v3 condensed endpoint, with no platform filtering.
//
// GET https://api.ipsw.me/v3/firmwares.json/condensed
func (s *FirmwareService) ListAllFirmwareV3(ctx context.Context) (*AllFirmwaresV3Response, *resty.Response, error) {
	return s.fetchAllFirmware(ctx)
}

// ListAllMacFirmwareV3 retrieves all Mac device models and their complete
// firmware history from the ipsw.me v3 condensed endpoint.
//
// The response contains all 56+ Mac model identifiers (e.g. "Mac14,3",
// "MacBookPro18,1"). Each device entry includes its full firmware history
// with Apple CDN download URLs, checksums, and signing status.
//
// Note: all current macOS versions use a single Universal IPSW URL shared
// across all Apple Silicon Mac models.
//
// GET https://api.ipsw.me/v3/firmwares.json/condensed
func (s *FirmwareService) ListAllMacFirmwareV3(ctx context.Context) (*AllFirmwaresV3Response, *resty.Response, error) {
	all, resp, err := s.fetchAllFirmware(ctx)
	if err != nil {
		return nil, resp, err
	}

	filtered := make(map[string]*MacDevice, len(all.Devices))
	for id, device := range all.Devices {
		if isMacIdentifier(id) {
			filtered[id] = device
		}
	}
	all.Devices = filtered

	return all, resp, nil
}

// ListAllIOSFirmwareV3 retrieves all iPhone device models and their complete
// firmware history from the ipsw.me v3 condensed endpoint.
//
// GET https://api.ipsw.me/v3/firmwares.json/condensed
func (s *FirmwareService) ListAllIOSFirmwareV3(ctx context.Context) (*AllFirmwaresV3Response, *resty.Response, error) {
	all, resp, err := s.fetchAllFirmware(ctx)
	if err != nil {
		return nil, resp, err
	}

	filtered := make(map[string]*MacDevice, len(all.Devices))
	for id, device := range all.Devices {
		if isIOSIdentifier(id) {
			filtered[id] = device
		}
	}
	all.Devices = filtered

	return all, resp, nil
}

// ListAllIPadOSFirmwareV3 retrieves all iPad device models and their complete
// firmware history from the ipsw.me v3 condensed endpoint.
//
// GET https://api.ipsw.me/v3/firmwares.json/condensed
func (s *FirmwareService) ListAllIPadOSFirmwareV3(ctx context.Context) (*AllFirmwaresV3Response, *resty.Response, error) {
	all, resp, err := s.fetchAllFirmware(ctx)
	if err != nil {
		return nil, resp, err
	}

	filtered := make(map[string]*MacDevice, len(all.Devices))
	for id, device := range all.Devices {
		if isIPadIdentifier(id) {
			filtered[id] = device
		}
	}
	all.Devices = filtered

	return all, resp, nil
}

// ListUniqueMacFirmwareVersionsV3 returns a deduplicated list of macOS firmware
// versions sorted from newest to oldest. Because all Mac models share the same
// Universal IPSW URL for a given version, this returns one entry per version
// rather than one per device.
//
// GET https://api.ipsw.me/v3/firmwares.json/condensed
func (s *FirmwareService) ListUniqueMacFirmwareVersionsV3(ctx context.Context) ([]*FirmwareV3, *resty.Response, error) {
	all, resp, err := s.ListAllMacFirmwareV3(ctx)
	if err != nil {
		return nil, resp, err
	}
	return uniqueVersionsSorted(all.Devices), resp, nil
}

// ListUniqueIOSFirmwareVersionsV3 returns a deduplicated list of iOS firmware
// versions sorted from newest to oldest.
//
// GET https://api.ipsw.me/v3/firmwares.json/condensed
func (s *FirmwareService) ListUniqueIOSFirmwareVersionsV3(ctx context.Context) ([]*FirmwareV3, *resty.Response, error) {
	all, resp, err := s.ListAllIOSFirmwareV3(ctx)
	if err != nil {
		return nil, resp, err
	}
	return uniqueVersionsSorted(all.Devices), resp, nil
}

// ListUniqueIPadOSFirmwareVersionsV3 returns a deduplicated list of iPadOS
// firmware versions sorted from newest to oldest.
//
// GET https://api.ipsw.me/v3/firmwares.json/condensed
func (s *FirmwareService) ListUniqueIPadOSFirmwareVersionsV3(ctx context.Context) ([]*FirmwareV3, *resty.Response, error) {
	all, resp, err := s.ListAllIPadOSFirmwareV3(ctx)
	if err != nil {
		return nil, resp, err
	}
	return uniqueVersionsSorted(all.Devices), resp, nil
}

// GetByDeviceV4 retrieves firmware for a specific device model identifier using
// the ipsw.me v4 device endpoint. Works for Mac, iPhone, and iPad identifiers.
// The v4 response is richer than v3, including SHA-256 checksums and filesize
// in addition to SHA-1/MD5.
//
// identifier is the Apple model identifier, e.g. "Mac14,3", "iPhone15,2", "iPad14,4".
//
// GET https://api.ipsw.me/v4/device/{identifier}?type=ipsw
func (s *FirmwareService) GetByDeviceV4(ctx context.Context, identifier string) (*DeviceFirmwaresV4Response, *resty.Response, error) {
	if identifier == "" {
		return nil, nil, fmt.Errorf("device identifier is required")
	}

	endpoint := fmt.Sprintf("%s/%s", constants.EndpointFirmwaresByDeviceV4, identifier)

	var result DeviceFirmwaresV4Response

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetQueryParam("type", FirmwareTypeIPSW).
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// uniqueVersionsSorted deduplicates firmware entries by BuildID across all
// devices and returns them sorted newest-first by ReleaseDate.
func uniqueVersionsSorted(devices map[string]*MacDevice) []*FirmwareV3 {
	seen := make(map[string]*FirmwareV3)
	for _, device := range devices {
		for _, fw := range device.Firmwares {
			if _, exists := seen[fw.BuildID]; !exists {
				seen[fw.BuildID] = fw
			}
		}
	}

	unique := make([]*FirmwareV3, 0, len(seen))
	for _, fw := range seen {
		unique = append(unique, fw)
	}

	sort.Slice(unique, func(i, j int) bool {
		return unique[i].ReleaseDate.After(unique[j].ReleaseDate)
	})

	return unique
}

// isMacIdentifier returns true when the identifier refers to a Mac model.
// Mac identifiers follow the pattern "MacXX,Y", "MacBookProXX,Y",
// "MacminiX,Y", "iMacXX,Y", "MacProX,Y", "MacBookAirXX,Y", or "VirtualMac".
func isMacIdentifier(identifier string) bool {
	for _, prefix := range []string{"Mac", "iMac", "VirtualMac"} {
		if strings.HasPrefix(identifier, prefix) {
			return true
		}
	}
	return false
}

// isIOSIdentifier returns true when the identifier refers to an iPhone model.
// iPhone identifiers follow the pattern "iPhoneXX,Y".
func isIOSIdentifier(identifier string) bool {
	return strings.HasPrefix(identifier, IdentifierPrefixIPhone)
}

// isIPadIdentifier returns true when the identifier refers to an iPad model.
// iPad identifiers follow the pattern "iPadXX,Y".
func isIPadIdentifier(identifier string) bool {
	return strings.HasPrefix(identifier, IdentifierPrefixIPad)
}
