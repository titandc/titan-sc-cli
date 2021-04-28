package api

import "encoding/json"

func (API *API) PostIPPNATRuleAdd(serverUUID, ip, protocol string, transparent bool, portSrc, portDst int64) (*APIReturn, error) {
	pnatOpt := APIPNATRuleAddDel{
		IP:          ip,
		Transparent: transparent,
		Protocol:    protocol,
		PortSrc:     portSrc,
		PortDst:     portDst,
	}
	_, apiReturn, err := API.SendRequestToAPI(HTTPPost, "/compute/servers/"+serverUUID+"/pnat", pnatOpt)
	return apiReturn, err
}

func (API *API) DeleteIPPNATRule(serverUUID, ip, protocol string, transparent bool, portSrc, portDst int64) (*APIReturn, error) {
	pnatOpt := APIPNATRuleAddDel{
		IP:          ip,
		Transparent: transparent,
		Protocol:    protocol,
		PortSrc:     portSrc,
		PortDst:     portDst,
	}
	_, apiReturn, err := API.SendRequestToAPI(HTTPDelete, "/compute/servers/"+serverUUID+"/pnat", pnatOpt)
	return apiReturn, err
}

func (API *API) GetServerPNATRulesList(serverUUID string) ([]APIPNATRuleInfos, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/compute/servers/"+serverUUID+"/pnat", nil)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}
	PNATRulesList := make([]APIPNATRuleInfos, 0)
	if err = json.Unmarshal(apiResponseBody, PNATRulesList); err != nil {
		return nil, err
	}
	return PNATRulesList, nil
}
