package http

type configResponse struct {
	Configured bool `json:"configured"`
}

type checkoutResponse struct {
	CheckoutURL string `json:"checkout_url"`
}
