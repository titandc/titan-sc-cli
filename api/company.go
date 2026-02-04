package api

import "encoding/json"

func (API *API) GetListOfCompanies() ([]Company, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/user/companies", nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}
	listOfCompanies := make([]Company, 0)
	if err = json.Unmarshal(apiResponseBody, &listOfCompanies); err != nil {
		return nil, err
	}
	return listOfCompanies, nil
}

func (API *API) GetCompanyDetails(companyOID string) (*Company, error) {
	apiResponseBody, apiReturn, err := API.SendRequestToAPI(HTTPGet, "/company/"+companyOID, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}
	company := &Company{}
	err = json.Unmarshal(apiResponseBody, company)
	if err != nil {
		return nil, err
	}
	return company, nil
}
