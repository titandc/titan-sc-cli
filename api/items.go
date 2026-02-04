package api

import (
	"encoding/json"
	"fmt"
)

func (API *API) ListItems() ([]ItemLimited, error) {
	path := fmt.Sprintf("/item")
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	var items []ItemLimited
	if err = json.Unmarshal(apiResponseBody, &items); err != nil {
		return nil, err
	}

	return items, err
}
