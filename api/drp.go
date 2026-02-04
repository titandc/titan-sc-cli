package api

import (
	"encoding/json"
	"fmt"
)

// GetDrpStatus retrieves DRP status for a server
// GET /server/{id}/drp/status
func (API *API) GetDrpStatus(serverOID string) (*DrpStatus, error) {
	path := fmt.Sprintf("/server/%s/drp/status", serverOID)
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	var status DrpStatus
	if err := json.Unmarshal(rawData, &status); err != nil {
		return nil, err
	}
	return &status, nil
}

// DrpFailoverSoft initiates a soft (planned) failover
// POST /server/{id}/drp/failover/soft
// Precondition: VM must be stopped
func (API *API) DrpFailoverSoft(serverOID string) (*DrpOperationResult, error) {
	path := fmt.Sprintf("/server/%s/drp/failover/soft", serverOID)
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPPost, path, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	var result DrpOperationResult
	if err := json.Unmarshal(rawData, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DrpFailoverHard initiates an emergency failover to specified target site
// POST /server/{id}/drp/failover/hard
// Can be done while VM is running - use with caution!
func (API *API) DrpFailoverHard(serverOID, targetSite string) (*DrpOperationResult, error) {
	path := fmt.Sprintf("/server/%s/drp/failover/hard", serverOID)
	payload := DrpFailoverHardRequest{TargetSite: targetSite}
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPPost, path, payload)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	var result DrpOperationResult
	if err := json.Unmarshal(rawData, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DrpResync resolves a split-brain situation by choosing an authoritative site
// POST /server/{id}/drp/resync
// WARNING: Data from non-authoritative site will be LOST!
func (API *API) DrpResync(serverOID, authoritativeSite string) (*DrpOperationResult, error) {
	path := fmt.Sprintf("/server/%s/drp/resync", serverOID)
	payload := DrpResyncRequest{AuthoritativeSite: authoritativeSite}
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPPost, path, payload)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	var result DrpOperationResult
	if err := json.Unmarshal(rawData, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DrpNetworkEnable enables DRP on a private network
// POST /network/switch/{id}/drp/enable
func (API *API) DrpNetworkEnable(networkOID string) (*NetworkDetail, error) {
	path := fmt.Sprintf("/network/switch/%s/drp/enable", networkOID)
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPPost, path, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	var network NetworkDetail
	if err := json.Unmarshal(rawData, &network); err != nil {
		return nil, err
	}
	return &network, nil
}

// DrpNetworkDisable disables DRP on a private network
// POST /network/switch/{id}/drp/disable
func (API *API) DrpNetworkDisable(networkOID string) (*NetworkDetail, error) {
	path := fmt.Sprintf("/network/switch/%s/drp/disable", networkOID)
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPPost, path, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	var network NetworkDetail
	if err := json.Unmarshal(rawData, &network); err != nil {
		return nil, err
	}
	return &network, nil
}
