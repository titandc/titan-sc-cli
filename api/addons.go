package api

import "encoding/json"

func (API *API) GetAllAddons() ([]APIAddonsItem, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/compute/addons", nil)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}
	addons := make([]APIAddonsItem, 0)
	if err = json.Unmarshal(apiResponseBody, &addons); err != nil {
		return nil, err
	}
	return addons, err
}
