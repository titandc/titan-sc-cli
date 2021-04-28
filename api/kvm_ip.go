package api

import "encoding/json"

func (API *API) KVMIPGet(serverUUID string) (*APIKvmIP, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/compute/servers/"+serverUUID+"/ipkvm", nil)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}

	kvmip := &APIKvmIP{}
	if err = json.Unmarshal(apiResponseBody, kvmip); err != nil {
		return nil, err
	}

	return kvmip, nil
}
