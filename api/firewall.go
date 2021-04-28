package api

import "encoding/json"

func (API *API) PostFirewallAdd(networkUUID, serverUUID, protocol, port, source string) (*APIReturn, error) {
	firewallOpt := APINetworkFirewallRule{
		ServerUUID: serverUUID,
		Protocol:   protocol,
		Port:       port,
		Source:     source,
	}

	_, apiReturn, err := API.SendRequestToAPI(HTTPPost, "/compute/networks/"+networkUUID+"/firewall", firewallOpt)
	return apiReturn, err
}

func (API *API) DeleteFirewall(networkUUID, serverUUID, protocol, port, source string) (*APIReturn, error) {
	firewallOpt := APINetworkFirewallRule{
		ServerUUID: serverUUID,
		Protocol:   protocol,
		Port:       port,
		Source:     source,
	}

	_, apiReturn, err := API.SendRequestToAPI(HTTPDelete, "/compute/networks/"+networkUUID+"/firewall", firewallOpt)
	return apiReturn, err
}

func (API *API) GetFireWallFullInfos(networkUUID string) (*APINetworkFullInfosFirewall, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/compute/networks/"+networkUUID+"/firewall", nil)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}
	firewallFullInfos := new(APINetworkFullInfosFirewall)
	if err = json.Unmarshal(apiResponseBody, firewallFullInfos); err != nil {
		return nil, err
	}
	return firewallFullInfos, err
}
