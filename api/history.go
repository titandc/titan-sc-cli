package api

import (
	"encoding/json"
	"fmt"
)

const (
	EventTypeServer = iota
	EventTypeCompany
)

func (API *API) GetEvents(number, offset, oid string, eventType int) ([]Event, *Return, error) {
	path := fmt.Sprintf("/server/%s/events", oid)
	if eventType == EventTypeCompany {
		path = fmt.Sprintf("/company/%s/events", oid)
	}
	path += fmt.Sprintf("?limit=%s&offset=%s", number, offset)

	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)

	// Communication error
	if err != nil || apiReturn != nil {
		return []Event{}, apiReturn, err
	}

	// Try to unmarshal output
	var events []Event
	err = json.Unmarshal(rawData, &events)
	if err != nil {
		return []Event{}, nil, err
	}
	return events, nil, nil
}
