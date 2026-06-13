package ctxutil

import "context"

type contextKey string

const NonceKey contextKey = "csp_nonce"

// GetNonce extracts the nonce from the context for use in Templ
func GetNonce(ctx context.Context) string {
	if nonce, ok := ctx.Value(NonceKey).(string); ok {
		return nonce
	}
	return "" // Fallback if missing
}
