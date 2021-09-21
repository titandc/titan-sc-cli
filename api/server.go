package api

import (
	"encoding/json"
)

func (API *API) ServerChangeName(newServerName, serverUUID string) (*APIReturn, error) {
	updateInfos := &APIServerUpdateInfos{
		Name: newServerName,
	}
	path := "/compute/servers/" + serverUUID
	_, apiReturn, err := API.SendRequestToAPI(HTTPPut, path, updateInfos)
	return apiReturn, err
}

func (API *API) ServerChangeReverse(newServerReverse, serverUUID string) (*APIReturn, error) {
	updateInfos := &APIServerUpdateInfos{
		Reverse: newServerReverse,
	}
	path := "/compute/servers/" + serverUUID
	_, apiReturn, err := API.SendRequestToAPI(HTTPPut, path, updateInfos)
	return apiReturn, err
}

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

func (API *API) ServerStateAction(state, serverUUID string) (*APIReturn, error) {
	// send request
	act := APIServerAction{
		Action: state,
	}
	path := "/compute/servers/" + serverUUID + "/action"
	_, apiReturn, err := API.SendRequestToAPI(HTTPPut, path, act)
	return apiReturn, err
}

func (API *API) ServerLoadISO(uriISO, serverUUID string) (*APIReturn, error) {
	reqStruct := &APIServerLoadISORequest{
		Protocol: "https",
		ISO:      uriISO,
	}
	path := "/compute/servers/" + serverUUID + "/iso"
	_, apiReturn, err := API.SendRequestToAPI(HTTPPost, path, reqStruct)
	return apiReturn, err
}

func (API *API) ServerUnloadISO(serverUUID string) (*APIReturn, error) {
	path := "/compute/servers/" + serverUUID + "/iso"
	_, apiReturn, err := API.SendRequestToAPI(HTTPDelete, path, nil)
	return apiReturn, err
}

func (API *API) ServerListTemplates() ([]APITemplateFullInfos, *APIReturn, error) {
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/compute/templates", nil)
	// Communication error
	if err != nil {
		return nil, nil, err
	}

	// API error
	if apiReturn != nil {
		return nil, apiReturn, nil
	}

	templates := []APITemplateFullInfos{}
	if err = json.Unmarshal(rawData, &templates); err != nil {
		return nil, nil, err
	}
	return templates, nil, nil
}

func (API *API) ServerDelete(serverUUID, deleteReason string) (*APIReturn, error) {
	deleteRequest := &APIDeleteServer{
		Reason: deleteReason,
	}
	_, apiReturn, err := API.SendRequestToAPI(HTTPDelete, "/compute/servers/"+serverUUID, deleteRequest)
	return apiReturn, err
}

func (API *API) ServerReset(serverUUID string, resetData *APIResetServer) (*APIReturn, error) {
	_, apiReturn, err := API.SendRequestToAPI(HTTPPut, "/compute/servers/"+serverUUID+"/reset", resetData)
	return apiReturn, err
}

func (API *API) ServerCreate(createData *APICreateServers) (*APIReturn, error) {
	_, apiReturn, err := API.SendRequestToAPI(HTTPPost, "/compute/servers", createData)
	return apiReturn, err
}

func (API *API) GetServersOfCompany(companyUUID string) ([]APIServer, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet,
		"/compute/servers?company_uuid="+companyUUID, nil)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}

	var servers []APIServer
	if err := json.Unmarshal(apiResponseBody, &servers); err != nil {
		return nil, err
	}
	return servers, nil
}
