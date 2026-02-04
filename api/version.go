package api

import "encoding/json"

func (API *API) GetVersion() (*APIVersion, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/version/current", nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}
	version := new(APIVersion)
	if err = json.Unmarshal(apiResponseBody, version); err != nil {
		return nil, err
	}
	return version, nil
}
