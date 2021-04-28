package api

import "encoding/json"

func (API *API) GetSSHKeyList() ([]APIUserSSHKey, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/auth/user/sshkeys", nil)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}
	sshKeyList := make([]APIUserSSHKey, 0)
	if err = json.Unmarshal(apiResponseBody, &sshKeyList); err != nil {
		return nil, err
	}
	return sshKeyList, nil
}

func (API *API) PostSSHKeyAdd(name, value string) (*APIReturn, error) {

	addSSHKey := APIAddUserSSHKey{
		Value: value,
		Name:  name,
	}
	_, apiReturn, err := API.SendRequestToAPI(HTTPPost, "/auth/user/sshkeys", addSSHKey)
	return apiReturn, err
}

func (API *API) DeleteSSHKey(name string) (*APIReturn, error) {
	delSSHKey := APIDeleteUserSSHKey{
		Name: name,
	}
	_, apiReturn, err := API.SendRequestToAPI(HTTPDelete, "/auth/user/sshkeys", delSSHKey)
	return apiReturn, err
}
