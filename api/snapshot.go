package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	SnapshotCreateErrorTooFast       = "SNAPSHOT_CREATE_FAIL_TOO_FAST"
	SnapshotCreateErrorLimitExceeded = "SNAPSHOT_CREATE_FAIL_LIMIT_EXCEEDED"
)

func (API *API) DeleteSnapshot(snapOID string) (*Return, error) {
	// Send request
	path := fmt.Sprintf("/storage/snapshot/%s", snapOID)
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPDelete, path, nil)

	// Communication error
	if err != nil {
		return nil, err
	}
	if apiReturn != nil {
		return apiReturn, nil
	}

	// Unmarshal error
	if string(rawData) == "\"SUCCESS\"" {
		return nil, nil
	}
	return nil, errors.New("unknown error")
}

// DeleteSnapshotLegacy deletes a snapshot using API v1 format (requires both server and snapshot UUID).
// API v1 path: DELETE /compute/servers/{server_uuid}/snapshots/{snapshot_uuid}
// This is for backward compatibility with v3.x CLI and will be removed in a future version.
func (API *API) DeleteSnapshotLegacy(serverUUID, snapUUID string) (*Return, error) {
	path := fmt.Sprintf("/compute/servers/%s/snapshots/%s", serverUUID, snapUUID)
	rawData, apiReturn, err := API.SendLegacyRequestToAPI(HTTPDelete, path, nil)

	// Communication error
	if err != nil {
		return nil, err
	}
	if apiReturn != nil {
		return apiReturn, nil
	}

	// Check for success
	if string(rawData) == "\"SUCCESS\"" {
		return nil, nil
	}
	return nil, errors.New("unknown error")
}

// ListSnapshotsLegacy retrieves all snapshots for a server using API v1.
// API v1 path: GET /compute/servers/{server_uuid}/snapshots
// This is for backward compatibility with v3.x CLI and will be removed in a future version.
func (API *API) ListSnapshotsLegacy(serverUUID string) ([]Snapshot, *Return, error) {
	path := fmt.Sprintf("/compute/servers/%s/snapshots", serverUUID)
	rawData, apiReturn, err := API.SendLegacyRequestToAPI(HTTPGet, path, nil)

	// Communication error
	if err != nil {
		return []Snapshot{}, nil, err
	}

	// API error
	if apiReturn != nil {
		return []Snapshot{}, apiReturn, nil
	}

	// Try to unmarshal output
	var snapshots []Snapshot
	err = json.Unmarshal(rawData, &snapshots)
	if err != nil {
		return []Snapshot{}, nil, err
	}

	return snapshots, nil, nil
}

// CreateSnapshotLegacy creates a new snapshot for a server using API v1.
// API v1 path: POST /compute/servers/{server_uuid}/snapshots
// This is for backward compatibility with v3.x CLI and will be removed in a future version.
func (API *API) CreateSnapshotLegacy(serverUUID string) (*SnapshotDetail, *Return, error) {
	path := fmt.Sprintf("/compute/servers/%s/snapshots", serverUUID)
	rawData, apiReturn, err := API.SendLegacyRequestToAPI(HTTPPost, path, nil)

	// Communication error
	if err != nil {
		return nil, nil, err
	}

	// API error
	if apiReturn != nil {
		return nil, apiReturn, nil
	}

	// Try to unmarshal output
	var snapshot SnapshotDetail
	err = json.Unmarshal(rawData, &snapshot)
	if err != nil {
		return nil, nil, err
	}
	return &snapshot, nil, nil
}

// RestoreSnapshotLegacy restores a server snapshot using API v1.
// API v1 path: PUT /compute/servers/{server_uuid}/snapshots/{snapshot_uuid}/restore
// This is for backward compatibility with v3.x CLI and will be removed in a future version.
func (API *API) RestoreSnapshotLegacy(serverUUID, snapUUID string) (*Return, error) {
	path := fmt.Sprintf("/compute/servers/%s/snapshots/%s/restore", serverUUID, snapUUID)
	rawData, apiReturn, err := API.SendLegacyRequestToAPI(HTTPPut, path, nil)

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

func (API *API) ListSnapshots(serverOID string) ([]Snapshot, *Return, error) {
	path := fmt.Sprintf("/storage/%s/snapshot", serverOID)
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)

	// Communication error
	if err != nil {
		return []Snapshot{}, nil, err
	}

	// API error
	if apiReturn != nil {
		return []Snapshot{}, apiReturn, nil
	}

	// Try to unmarshal output
	var snapshots []Snapshot
	err = json.Unmarshal(rawData, &snapshots)
	if err != nil {
		return []Snapshot{}, nil, err
	}

	return snapshots, nil, nil
}

func (API *API) CreateSnapshot(serverOID string) (*SnapshotDetail, *Return, error) {
	path := fmt.Sprintf("/storage/%s/snapshot", serverOID)
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPPost, path, nil)

	// Communication error
	if err != nil {
		return nil, nil, err
	}

	// API error
	if apiReturn != nil {
		return nil, apiReturn, nil
	}

	// Try to unmarshal output
	var snapshot SnapshotDetail
	err = json.Unmarshal(rawData, &snapshot)
	if err != nil {
		return nil, nil, err
	}
	return &snapshot, nil, nil
}

func (API *API) RestoreSnapshot(snapshotOID string) (*Return, error) {
	// Send request
	path := fmt.Sprintf("/storage/snapshot/%s/restore", snapshotOID)
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPPut, path, nil)

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
