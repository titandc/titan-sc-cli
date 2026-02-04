package api

import (
	"encoding/json"
	"fmt"
)

// GetSubscriptionList retrieves all subscriptions for a company
func (API *API) GetSubscriptionList(companyOID string, activeOnly bool) ([]Subscription, error) {
	path := "/subscription"
	if activeOnly {
		path = "/subscription?states[]=ongoing"
	}
	if companyOID != "" {
		if activeOnly {
			path = "/subscription?states[]=ongoing&company_oid=" + companyOID
		} else {
			path = "/subscription?company_oid=" + companyOID
		}
	}

	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	var subscriptions []Subscription
	if err = json.Unmarshal(rawData, &subscriptions); err != nil {
		return nil, err
	}
	return subscriptions, nil
}

// GetSubscription retrieves a single subscription by OID
func (API *API) GetSubscription(subscriptionOID string) (*Subscription, error) {
	path := fmt.Sprintf("/subscription/%s", subscriptionOID)
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	var subscription Subscription
	if err = json.Unmarshal(rawData, &subscription); err != nil {
		return nil, err
	}
	return &subscription, nil
}
