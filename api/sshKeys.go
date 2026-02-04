package api

import (
	"encoding/json"
	"fmt"
)

func (API *API) GetSSHKeyList(targetOID string) ([]SSHKey, error) {
	path := fmt.Sprintf("/ssh_key?target_oid=%s", targetOID)
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	sshKeyList := make([]SSHKey, 0)
	if err = json.Unmarshal(apiResponseBody, &sshKeyList); err != nil {
		return nil, err
	}
	return sshKeyList, nil
}

func (API *API) GetSSHKey(sshKeyOID string) (*SSHKey, error) {
	path := fmt.Sprintf("/ssh_key/%s", sshKeyOID)
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	sshKey := &SSHKey{}
	if err = json.Unmarshal(apiResponseBody, sshKey); err != nil {
		return nil, err
	}
	return sshKey, nil
}

func (API *API) PostSSHKeyAdd(name, value string) (*Return, error) {
	addSSHKey := AddSSHKey{
		Value: value,
		Name:  name,
	}

	_, apiReturn, err := API.SendRequestToAPI(HTTPPost, "/ssh_key", addSSHKey)
	return apiReturn, err
}

func (API *API) DeleteSSHKey(sshKeyOID string) (*Return, error) {
	path := fmt.Sprintf("/ssh_key/%s", sshKeyOID)
	_, apiReturn, err := API.SendRequestToAPI(HTTPDelete, path, nil)
	return apiReturn, err
}
