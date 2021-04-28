package api

func (API *API) PostManagedServices(companyUUID string) (*APIReturn, error) {
	managedServicesOpts := APIManagedServices{
		Company: companyUUID,
	}
	_, apiReturn, err := API.SendRequestToAPI(HTTPPost, "/compute/managed_services", managedServicesOpts)
	return apiReturn, err
}
