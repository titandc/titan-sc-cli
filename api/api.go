package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	DefaultURI = "https://sc.titandc.net/api/v1"
	HTTPGet    = http.MethodGet
	HTTPPut    = http.MethodPut
	HTTPPost   = http.MethodPost
	HTTPDelete = http.MethodDelete
)

type API struct {
	Token   string
	URI     string
	OS      string
	Version string
}

func NewAPI(token, uri, os, version string) *API {
	if uri == "" {
		uri = DefaultURI
	}
	return &API{
		Token:   token,
		URI:     uri,
		OS:      os,
		Version: version,
	}
}

func (API *API) SendRequestToAPI(method, path string, httpData interface{}) ([]byte, *APIReturn, error) {

	// Transform interface to byte array
	var reqData []byte
	var err error
	if httpData != nil {
		reqData, err = json.Marshal(httpData)
		if err != nil {
			return nil, nil, err
		}
	}

	// Prepare new request
	request, err := http.NewRequest(method, API.URI+path, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, nil, err
	}
	request.Header.Add("X-API-KEY", API.Token)
	request.Header.Add("Titan-Cli-Os", API.OS)
	request.Header.Add("Titan-Cli-Version", API.Version)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Read API output
	apiResponseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Try to unmarshal as API generic output (code, error/success)
	apiReturn := &APIReturn{}
	err = json.Unmarshal(apiResponseBody, apiReturn)
	if err == nil && (apiReturn.Error != "" || apiReturn.Success != "") {
		return apiResponseBody, apiReturn, nil
	}

	// Return raw data
	return apiResponseBody, nil, nil
}

func handlePotentialDoubleError(apiReturn *APIReturn, err error) error {
	if err != nil {
		return err
	}
	if apiReturn != nil && apiReturn.Error != "" {
		return errors.New(apiReturn.Error)
	}
	return nil
}

