package api

import (
	"encoding/json"
)

func (API *API) HistoryByCompany(number, companyUUID string) ([]APIHistoryEvent, *APIReturn, error) {
	// Execute query
	path := "/compute/servers/events?nb=" + number + "&company_uuid=" + companyUUID
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)

	// Communication error
	if err != nil {
		return []APIHistoryEvent{}, nil, err
	}

	// API error
	if apiReturn != nil {
		return []APIHistoryEvent{}, apiReturn, nil
	}

	// Try to unmarshal output
	var events []APIHistoryEvent
	err = json.Unmarshal(rawData, &events)
	if err != nil {
		return []APIHistoryEvent{}, nil, err
	}
	return events, nil, nil
}

func (API *API) HistoryByServer(number, serverUUID string) ([]APIHistoryEvent, *APIReturn, error) {
	// Execute query
	path := "/compute/servers/" + serverUUID + "/events?nb=" + number
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)

	// Communication error
	if err != nil {
		return []APIHistoryEvent{}, nil, err
	}

	// API error
	if apiReturn != nil {
		return []APIHistoryEvent{}, apiReturn, nil
	}

	// Try to unmarshal output
	var events []APIHistoryEvent
	err = json.Unmarshal(rawData, &events)
	if err != nil {
		return []APIHistoryEvent{}, nil, err
	}
	return events, nil, nil
}
