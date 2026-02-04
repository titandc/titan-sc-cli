package api

import "encoding/json"

func (API *API) GetUserInfos() (*User, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/user/me", nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	user := new(User)
	if err = json.Unmarshal(apiResponseBody, user); err != nil {
		return nil, err
	}
	return user, nil
}
