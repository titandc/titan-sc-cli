package api

import (
	"encoding/json"
	"fmt"
)

func (API *API) ServerAddon(serverOID string) (*ServerAddonInfo, error) {
	path := fmt.Sprintf("/server/%s/addon/info", serverOID)
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}
	addons := &ServerAddonInfo{}
	if err = json.Unmarshal(apiResponseBody, &addons); err != nil {
		return nil, err
	}
	return addons, err
}
