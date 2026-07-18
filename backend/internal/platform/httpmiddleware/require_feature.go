package httpmiddleware

import "net/http"

// FlagChecker is satisfied by featureflags/application.Registry. Defined
// here (rather than imported) to avoid a dependency from this shared
// middleware package back into the featureflags package's application layer.
type FlagChecker interface {
	IsEnabled(key string) bool
}

// RequireFeatureEnabled is the module kill switch: if the flag is off, every
// route behind it returns 503 instead of running.
func RequireFeatureEnabled(flags FlagChecker, key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !flags.IsEnabled(key) {
				writeJSONError(w, http.StatusServiceUnavailable, "feature_disabled", key+" is temporarily disabled")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
