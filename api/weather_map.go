package api

import "encoding/json"

func (API *API) GetWeatherMap() (*APIWeatherMap, error) {

	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet,
		"/weather", nil)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}
	weatherMap := &APIWeatherMap{}
	if err = json.Unmarshal(apiResponseBody, weatherMap); err != nil {
		return nil, err
	}
	return weatherMap, nil
}
