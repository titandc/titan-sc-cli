package api

import (
	"encoding/json"
	"errors"
)

/*
 *
 *
 **************************
 * Snapshot server function
 **************************
 *
 *
 */

const (
	SnapshotCreateErrorTooFast       = "SNAPSHOT_CREATE_FAIL_TOO_FAST"
	SnapshotCreateErrorLimitExceeded = "SNAPSHOT_CREATE_FAIL_LIMIT_EXCEEDED"
)

func (API *API) DeleteSnapshot(serverUUID, snapUUID string) (*APIReturn, error) {
	// Send request
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPDelete, "/compute/servers/"+
		serverUUID+"/snapshots/"+snapUUID, nil)

	// Communication error
	if err != nil {
		return nil, err
	}

	// Unmarshal error
	if apiReturn == nil {
		return nil, errors.New(string(rawData))
	}
	return apiReturn, nil
}

func (API *API) ListSnapshots(serverUUID string) ([]APISnapshot, *APIReturn, error) {
	// Send request
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/compute/servers/"+serverUUID+"/snapshots", nil)

	// Communication error
	if err != nil {
		return []APISnapshot{}, nil, err
	}

	// API error
	if apiReturn != nil {
		return []APISnapshot{}, apiReturn, nil
	}

	// Try to unmarshal output
	var snapshots []APISnapshot
	err = json.Unmarshal(rawData, &snapshots)
	if err != nil {
		return []APISnapshot{}, nil, err
	}
	return snapshots, nil, nil
}

func (API *API) PostCreateSnapshot(serverUUID string) (*APISnapshot, *APIReturn, error) {
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPPost, "/compute/servers/"+serverUUID+"/snapshots", nil)

	// Communication error
	if err != nil {
		return nil, nil, err
	}

	// API error
	if apiReturn != nil {
		return nil, apiReturn, nil
	}

	// Try to unmarshal output
	var snapshot APISnapshot
	err = json.Unmarshal(rawData, &snapshot)
	if err != nil {
		return nil, nil, err
	}
	return &snapshot, nil, nil
}
