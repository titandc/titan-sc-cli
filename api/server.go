package api

import (
	"encoding/json"
)

// ServerList retrieves all servers for the given company UUID.
// GET /compute/servers/detail?company_uuid=...
func (API *API) ServerList(companyUUID string) ([]APIServer, *APIReturn, error) {
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet,
		"/compute/servers/detail?company_uuid="+companyUUID, nil)
	// Communication error
	if err != nil {
		return nil, nil, err
	}

	// API error
	if apiReturn != nil {
		return nil, apiReturn, nil
	}

	servers := []APIServer{}
	if err = json.Unmarshal(rawData, &servers); err != nil {
		return nil, nil, err
	}
	return servers, nil, nil
}

// GetServerUUID retrieves details for a specific server by UUID.
// GET /compute/servers/{uuid}
func (API *API) GetServerUUID(serverUUID string) (*APIServer, *APIReturn, error) {
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/compute/servers/"+serverUUID, nil)
	// Communication error
	if err != nil {
		return nil, nil, err
	}

	// API error
	if apiReturn != nil {
		return nil, apiReturn, nil
	}

	server := &APIServer{}
	if err = json.Unmarshal(rawData, server); err != nil {
		return nil, nil, err
	}

	return server, nil, nil
}
