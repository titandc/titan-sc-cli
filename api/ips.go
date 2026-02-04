package api

import (
	"encoding/json"
	"fmt"
)

func (API *API) IPAttach(serverOID string, ipOIDs []string) (*Return, error) {
	req := IPAttachDetach{IPs: ipOIDs}
	_, apiReturn, err := API.SendRequestToAPI(HTTPPost, "/network/ip/"+serverOID, req)
	return apiReturn, err
}

func (API *API) IPDetach(serverOID string, ipOIDs []string) (*Return, error) {
	req := IPAttachDetach{IPs: ipOIDs}
	_, apiReturn, err := API.SendRequestToAPI(HTTPDelete, "/network/ip/"+serverOID, req)
	return apiReturn, err
}

func (API *API) GetCompanyIPList(companyOID string) ([]IP, error) {
	path := fmt.Sprintf("/company/ips?company_oid=%s", companyOID)
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}
	ipList := new([]IP)
	if err = json.Unmarshal(apiResponseBody, ipList); err != nil {
		return nil, err
	}
	return *ipList, nil
}

func (API *API) IPUpdateReverse(ipOID, newIPReverse string) (*Return, error) {
	req := IPUpdateRequest{
		Reverse: newIPReverse,
	}

	path := fmt.Sprintf("/ip/%s", ipOID)
	_, apiReturn, err := API.SendRequestToAPI(HTTPPut, path, &req)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	return nil, nil
}
