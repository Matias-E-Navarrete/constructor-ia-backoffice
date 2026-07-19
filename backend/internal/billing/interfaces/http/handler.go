package http

import (
	"io"
	nethttp "net/http"

	"rimu/backend/internal/billing/application"
	"rimu/backend/internal/billing/domain"
	"rimu/backend/internal/platform/httpkit"
	"rimu/backend/internal/platform/httpmiddleware"
)

type Handler struct {
	Configured bool
	Checkout   *application.CreateCheckoutSession
	Webhook    *application.HandleWebhook
}

func (h *Handler) HandleConfig(w nethttp.ResponseWriter, r *nethttp.Request) {
	httpkit.WriteJSON(w, nethttp.StatusOK, configResponse{Configured: h.Configured})
}

func (h *Handler) HandleCheckout(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	url, err := h.Checkout.Execute(r.Context(), user.ID, user.Email)
	if err == application.ErrNotConfigured {
		httpkit.WriteError(w, nethttp.StatusServiceUnavailable, "billing_not_configured", "real checkout isn't configured on this deployment")
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadGateway, "billing_error", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, checkoutResponse{CheckoutURL: url})
}

// HandleWebhook is mounted outside RequireAuth: Stripe authenticates itself
// via the Stripe-Signature header, not our JWT.
func (h *Handler) HandleWebhook(w nethttp.ResponseWriter, r *nethttp.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "could not read request body")
		return
	}

	err = h.Webhook.Execute(r.Context(), payload, r.Header.Get("Stripe-Signature"))
	if err == domain.ErrInvalidSignature || err == domain.ErrNotConfigured {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_signature", "webhook signature verification failed")
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	w.WriteHeader(nethttp.StatusOK)
}
