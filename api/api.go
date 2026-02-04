package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	// DefaultURI is the production API endpoint. Use https://staging.titandc.io/api/v2 for testing.
	DefaultURI = "https://sc.titandc.net/api/v2"
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
		Token:   token, // API uses X-API-KEY header directly, not Bearer token
		URI:     uri,
		OS:      os,
		Version: version,
	}
}

// GetLegacyV1URI derives the API v1 URI from the configured v2 URI.
// This replaces /v2 with /v1 in the URI path, allowing custom endpoints to work.
func (API *API) GetLegacyV1URI() string {
	// Replace /v2 with /v1 at the end of the URI
	if strings.HasSuffix(API.URI, "/v2") {
		return strings.TrimSuffix(API.URI, "/v2") + "/v1"
	}
	// Fallback: try replacing /api/v2 anywhere in the URI
	return strings.Replace(API.URI, "/api/v2", "/api/v1", 1)
}

// SendLegacyRequestToAPI sends a request to the API v1 endpoint.
// This is used for backward compatibility with legacy CLI commands.
func (API *API) SendLegacyRequestToAPI(method, path string, payload interface{}) ([]byte, *Return, error) {
	// Temporarily switch to v1 URI
	originalURI := API.URI
	API.URI = API.GetLegacyV1URI()
	defer func() { API.URI = originalURI }()

	return API.SendRequestToAPI(method, path, payload)
}

func (API *API) SendRequestToAPI(method, path string, payload interface{}) ([]byte, *Return, error) {
	// Transform interface to byte array
	var body []byte
	var err error

	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, nil, err
		}
	}

	// Prepare new request
	request, err := http.NewRequest(method, API.URI+path, bytes.NewBuffer(body))
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
	apiResponseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Try to unmarshal as API generic output (code, error/success)
	ret := &Return{}
	err = json.Unmarshal(apiResponseBody, ret)
	if err == nil && (ret.Error() || ret.IsSuccess()) {
		return apiResponseBody, ret, nil
	}

	// Check if response is a raw string error (e.g., "BAD_PERMISSION\n")
	var rawString string
	if json.Unmarshal(apiResponseBody, &rawString) == nil && rawString != "" {
		// Clean up the string (remove trailing newlines)
		rawString = strings.TrimSpace(rawString)
		if resp.StatusCode >= 400 || isKnownErrorString(rawString) {
			return apiResponseBody, &Return{Title: rawString}, nil
		}
	}

	// Return raw data
	return apiResponseBody, nil, nil
}

// isKnownErrorString checks if a string looks like an API error code
func isKnownErrorString(s string) bool {
	knownErrors := []string{
		"BAD_PERMISSION",
		"BAD_REQUEST",
		"UNAUTHORIZED",
		"FORBIDDEN",
		"NOT_FOUND",
		"INTERNAL_ERROR",
		"ERROR",
	}
	upper := strings.ToUpper(s)
	for _, e := range knownErrors {
		if strings.Contains(upper, e) {
			return true
		}
	}
	return false
}

func handleError(ret *Return, err error) error {
	if err != nil {
		return err
	}
	if ret != nil && ret.Error() {
		return ret.AsError()
	}
	return nil
}

// AsError converts the Return struct to a proper error with clean formatting
func (r *Return) AsError() error {
	if r == nil {
		return nil
	}

	var parts []string

	if r.Title != "" {
		parts = append(parts, r.Title)
	}

	if r.Message != "" {
		parts = append(parts, r.Message)
	}

	for _, v := range r.Data {
		parts = append(parts, fmt.Sprintf("%s: %v", v.Field, v.Value))
	}

	if len(parts) == 0 {
		return nil
	}

	return &APIError{Parts: parts}
}

// Error returns true if this Return represents an error response
// Only the "error" field (Title) indicates an actual error.
// The "message" field is used for both success and error responses.
func (r *Return) Error() bool {
	return r != nil && (r.Title != "" || len(r.Data) > 0)
}

// IsSuccess returns true if this Return represents a success response from API v1.
// API v1 returns {"code": "...", "success": "..."} for success responses.
func (r *Return) IsSuccess() bool {
	return r != nil && r.Success != ""
}

// APIError represents a structured API error
type APIError struct {
	Parts []string
}

func (e *APIError) Error() string {
	return strings.Join(e.Parts, ": ")
}

// ConcatAPIValidationError is deprecated, use Return.AsError() instead
func ConcatAPIValidationError(ret *Return) error {
	return ret.AsError()
}
