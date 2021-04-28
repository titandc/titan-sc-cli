package run

import (
	"encoding/json"
	"titan-sc/api"
)

func (run *RunMiddleware) IsAdmin() (bool, *api.APIReturn, error) {
	rawData, apiReturn, err := run.API.SendRequestToAPI(api.HTTPGet, "/auth/user/isadmin", nil)
	if err != nil {
		return false, apiReturn, err
	}

	buffer := api.IsAdminStruct{}
	if err := json.Unmarshal(rawData, &buffer); err != nil {
		return false, apiReturn, err
	}
	return buffer.Admin, apiReturn, nil
}
