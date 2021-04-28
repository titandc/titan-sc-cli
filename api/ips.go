package api

import "encoding/json"

func (API *API) PostIPAttach(serverUUID string, ipOpt []APIIPAttachDetach) (*APIReturn, error) {
	_, apiReturn, err := API.SendRequestToAPI(HTTPPost, "/compute/servers/"+serverUUID+"/ips", ipOpt)
	return apiReturn, err
}

func (API *API) DeleteIPDetach(serverUUID string, ipOpt APIIPAttachDetach) (*APIReturn, error) {
	_, apiReturn, err := API.SendRequestToAPI(HTTPDelete, "/compute/servers/"+serverUUID+"/ips", ipOpt)
	return apiReturn, err
}

func (API *API) GetIPList() ([]APIIPAttachDetach, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/compute/ips", nil)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}
	ipList := make([]APIIPAttachDetach, 0)
	if err = json.Unmarshal(apiResponseBody, &ipList); err != nil {
		return nil, err
	}
	return ipList, nil
}

func (API *API) GetCompanyIPList(companyUUID string) ([]APIIPAttachDetach, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/companies/"+companyUUID+"/ips", nil)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}
	ipList := new([]APIIPAttachDetach)
	if err = json.Unmarshal(apiResponseBody, ipList); err != nil {
		return nil, err
	}
	return *ipList, nil
}
