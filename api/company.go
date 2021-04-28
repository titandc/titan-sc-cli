package api

import "encoding/json"

func (API *API) GetListOfCompanies() ([]APICompany, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/companies", nil)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}
	listOfCompanies := make([]APICompany, 0)
	if err = json.Unmarshal(apiResponseBody, &listOfCompanies); err != nil {
		return nil, err
	}
	return listOfCompanies, nil
}

func (API *API) GetCompanyDetails(companyUUID string) (*APICompanyDetail, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/companies/"+companyUUID, nil)
	if err = handlePotentialDoubleError(apiReturn, err); err != nil {
		return nil, err
	}
	companyDetails := &APICompanyDetail{}
	err = json.Unmarshal(apiResponseBody, companyDetails)
	if err != nil {
		return nil, err
	}
	return companyDetails, nil
}
