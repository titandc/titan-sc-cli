package api

import (
	"encoding/json"
	"fmt"
)

// APIToken represents an API token
type APIToken struct {
	OID      string `json:"oid"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Expire   *int64 `json:"expire"` // Unix timestamp, nil = never expires
	OwnerOID string `json:"owner_oid"`
}

// APITokenCreate represents the request body for creating a token
type APITokenCreate struct {
	Name   string `json:"name"`
	Expire *int64 `json:"expire,omitempty"`
}

// APITokenUpdate represents the request body for updating a token
type APITokenUpdate struct {
	Name   string `json:"name,omitempty"`
	Expire *int64 `json:"expire,omitempty"`
}

// ListAPITokens retrieves all API tokens for the authenticated user
func (API *API) ListAPITokens() ([]APIToken, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/api_token", nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	var tokens []APIToken
	if err := json.Unmarshal(apiResponseBody, &tokens); err != nil {
		return nil, fmt.Errorf("failed to parse API tokens: %w", err)
	}
	return tokens, nil
}

// GetAPIToken retrieves a specific API token by OID
func (API *API) GetAPIToken(tokenOID string) (*APIToken, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/api_token/"+tokenOID, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	var token APIToken
	if err := json.Unmarshal(apiResponseBody, &token); err != nil {
		return nil, fmt.Errorf("failed to parse API token: %w", err)
	}
	return &token, nil
}

// CreateAPIToken creates a new API token
func (API *API) CreateAPIToken(create *APITokenCreate) (*APIToken, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPPost, "/api_token", create)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	var token APIToken
	if err := json.Unmarshal(apiResponseBody, &token); err != nil {
		return nil, fmt.Errorf("failed to parse API token: %w", err)
	}
	return &token, nil
}

// UpdateAPIToken updates an existing API token
func (API *API) UpdateAPIToken(tokenOID string, update *APITokenUpdate) (*APIToken, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPPut, "/api_token/"+tokenOID, update)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	var token APIToken
	if err := json.Unmarshal(apiResponseBody, &token); err != nil {
		return nil, fmt.Errorf("failed to parse API token: %w", err)
	}
	return &token, nil
}

// DeleteAPIToken deletes an API token by OID
func (API *API) DeleteAPIToken(tokenOID string) error {
	_, apiReturn, err := API.SendRequestToAPI(HTTPDelete, "/api_token/"+tokenOID, nil)
	return handleError(apiReturn, err)
}
