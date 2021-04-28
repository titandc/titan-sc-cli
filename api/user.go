package api

import "encoding/json"

func (API *API) GetUserInfos() (*APIUserInfos, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/auth/user", nil)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}

	userInfos := &APIUserInfos{}
	if err = json.Unmarshal(apiResponseBody, userInfos); err != nil {
		return nil, err
	}
	return userInfos, nil
}
