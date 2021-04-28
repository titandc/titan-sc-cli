package api

import (
	"encoding/json"
)

func (API *API) GetNetworkList(companyUUID string) (*APINetworkList, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet,
		"/compute/networks?company_uuid="+companyUUID, nil)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}
	networks := &APINetworkList{}
	if err := json.Unmarshal(apiResponseBody, networks); err != nil {
		return nil, err
	}
	return networks, nil
}

func (API *API) GetNetworkDetail(networkUUID string) (*APINetwork, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/compute/networks/"+networkUUID, nil)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}
	network := &APINetwork{}
	if err := json.Unmarshal(apiResponseBody, network); err != nil {
		return nil, err
	}
	return network, nil
}

func (API *API) CreateNetwork(companyUUID string, reqData APINetworkCreate) (*APINetwork, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPPost, "/compute/networks/?company_uuid="+
		companyUUID, reqData)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}
	network := &APINetwork{}
	err = json.Unmarshal(apiResponseBody, network)
	if err != nil {
		return nil, err
	}
	return network, nil
}

func (API *API) RemoveNetwork(networkUUID string) error {
	_, apiReturn, err := API.SendRequestToAPI(HTTPDelete, "/compute/networks/"+networkUUID, nil)
	return handlePotentialDoubleError(apiReturn, err)
}
