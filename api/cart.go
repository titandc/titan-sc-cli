package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (API *API) CreateServerCart(cart *AddServerCart) (string, error) {
	path := fmt.Sprintf("/cart/server")
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPPost, path, cart)
	if err = handleError(apiReturn, err); err != nil {
		return "", err
	}

	var reply []CreateServerCart
	if err = json.Unmarshal(rawData, &reply); err != nil {
		return "", err
	}

	if len(reply) > 0 && reply[0].CartOID != "" {
		return reply[0].CartOID, nil
	}
	return "", errors.New("cartOID not found")
}

// GetCartPrice retrieves the price preview for a cart
func (API *API) GetCartPrice(cartOID string) (*CartPrice, error) {
	path := fmt.Sprintf("/cart/getPrice?cart_oid=%s", cartOID)
	rawData, apiReturn, err := API.SendRequestToAPI(HTTPGet, path, nil)
	if err = handleError(apiReturn, err); err != nil {
		return nil, err
	}

	price := &CartPrice{}
	if err = json.Unmarshal(rawData, price); err != nil {
		return nil, err
	}
	return price, nil
}

// BuyCart purchases the cart using the specified payment method
// If subscriptionOID is provided, the server is added to the existing subscription
// Otherwise, a new subscription is created automatically
func (API *API) BuyCart(cartOID, paymentMethodOID, subscriptionOID string) error {
	req := BuyCartRequest{
		CartOID:          cartOID,
		PaymentMethodOID: paymentMethodOID,
		SubscriptionOID:  subscriptionOID,
	}

	path := "/cart/buy"
	_, apiReturn, err := API.SendRequestToAPI(HTTPPost, path, req)
	if err = handleError(apiReturn, err); err != nil {
		return err
	}
	return nil
}
