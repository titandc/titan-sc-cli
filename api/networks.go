package api

import (
	"encoding/json"
	"fmt"
)

func (API *API) GetNetworkList(companyOID string) (*NetworkList, error) {
	path := fmt.Sprintf("/network/switch?company_oid=%s", companyOID)
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}
	networks := &NetworkList{}
	if err := json.Unmarshal(apiResponseBody, networks); err != nil {
		return nil, err
	}
	return networks, nil
}

func (API *API) GetNetworkDetail(networkOID string) (*NetworkDetail, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/network/switch/"+networkOID, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}
	network := &NetworkDetail{}
	if err := json.Unmarshal(apiResponseBody, network); err != nil {
		return nil, err
	}
	return network, nil
}

func (API *API) CreateNetwork(reqData *NetworkCreate) (*Network, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPPost, "/network/switch", reqData)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	network := &Network{}
	err = json.Unmarshal(apiResponseBody, network)
	if err != nil {
		return nil, err
	}
	return network, nil
}

func (API *API) RemoveNetwork(networkOID string) error {
	_, apiReturn, err := API.SendRequestToAPI(HTTPDelete, "/network/switch/"+networkOID, nil)
	return handleError(apiReturn, err)
}
