package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func (API *API) ServerChangeName(newServerName, serverOID string) (*Return, error) {
	updateInfos := &ServerUpdateInfos{
		Name: newServerName,
	}
	path := fmt.Sprintf("/server/%s", serverOID)
	_, apiReturn, err := API.SendRequestToAPI(HTTPPut, path, updateInfos)
	return apiReturn, err
}

func (API *API) ServerList(companyOID string) ([]ServerDetail, *Return, error) {
	// Filter to show only active servers (exclude deleted)
	// Array format uses bracket notation: states[]=value
	activeStates := "states[]=started&states[]=stopped&states[]=starting&states[]=stopping&states[]=creating&states[]=unmanaged"
	path := "/server?" + activeStates
	if companyOID != "" {
		path = "/server?" + activeStates + "&company_oid=" + companyOID
	}
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)
	// Communication error
	if err != nil {
		return nil, nil, err
	}

	// API error
	if apiReturn != nil {
		return nil, apiReturn, nil
	}

	servers := []ServerDetail{}
	if err = json.Unmarshal(rawData, &servers); err != nil {
		return nil, nil, err
	}
	return servers, nil, nil
}

func (API *API) GetServerOID(serverOID string) (*ServerDetail, *Return, error) {
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/server/"+serverOID, nil)
	// Communication error
	if err != nil {
		return nil, nil, err
	}

	// API error
	if apiReturn != nil {
		return nil, apiReturn, nil
	}

	server := &ServerDetail{}
	if err = json.Unmarshal(rawData, server); err != nil {
		return nil, nil, err
	}

	return server, nil, nil
}

func (API *API) ServerStateAction(state, serverOID string) (*Return, error) {
	act := ServerAction{
		Action: state,
	}
	path := fmt.Sprintf("/server/%s/state", serverOID)
	_, apiReturn, err := API.SendRequestToAPI(HTTPPut, path, act)
	return apiReturn, err
}

func (API *API) ServerMountISO(uriISO, serverOID string) (*Return, error) {
	if !strings.HasPrefix(uriISO, "https://") {
		return nil, errors.New("URI must be use protocol: https")
	}
	payload := &ServerMountISORequest{
		Protocol: "https",
		ISO:      uriISO,
	}
	path := fmt.Sprintf("/storage/%s/iso", serverOID)
	_, apiReturn, err := API.SendRequestToAPI(HTTPPost, path, payload)
	return apiReturn, err
}

func (API *API) ServerUmountISO(serverOID, isoOID string) (*Return, error) {
	path := fmt.Sprintf("/storage/%s/iso/%s", serverOID, isoOID)
	_, apiReturn, err := API.SendRequestToAPI(HTTPDelete, path, nil)
	return apiReturn, err
}

func (API *API) ServerScheduleTermination(serverOID, deleteReason string) (*Return, error) {
	payload := &ServerScheduleTermination{
		Reason: deleteReason,
	}
	_, apiReturn, err := API.SendRequestToAPI(HTTPDelete, "/server/"+serverOID, payload)
	return apiReturn, err
}

func (API *API) ServerReset(serverOID string, data *ResetServer) (*Return, error) {
	path := fmt.Sprintf("/server/%s/reset", serverOID)
	_, apiReturn, err := API.SendRequestToAPI(HTTPPost, path, data)
	return apiReturn, err
}
