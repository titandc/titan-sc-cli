package api

import (
	"encoding/json"
	"fmt"
)

func (API *API) ListTemplates() ([]TemplateOSItem, *Return, error) {
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/storage/template/grouped", nil)
	// Communication error
	if err != nil {
		return nil, nil, err
	}

	// API error
	if apiReturn != nil {
		return nil, apiReturn, nil
	}

	templates := []TemplateOSItem{}
	if err = json.Unmarshal(rawData, &templates); err != nil {
		return nil, nil, err
	}
	return templates, nil, nil
}

func (API *API) GetTemplateByOID(templateOID string) (*Template, error) {
	path := fmt.Sprintf("/storage/template/%s", templateOID)
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	template := new(Template)
	if err = json.Unmarshal(rawData, &template); err != nil {
		return nil, err
	}
	return template, nil
}
