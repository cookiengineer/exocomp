package chat

import net_http "net/http"

// isRetryableStatus reports whether a response status is worth retrying.
// Authentication, payment, permission, model and schema errors are permanent
// and return false.
func isRetryableStatus(status int) bool {

	switch status {

	case net_http.StatusRequestTimeout: // 408
		return true

	case net_http.StatusConflict: // 409
		return true

	case net_http.StatusTooEarly: // 425
		return true

	case net_http.StatusTooManyRequests: // 429
		return true

	case net_http.StatusInternalServerError: // 500
		return true

	case net_http.StatusBadGateway: // 502
		return true

	case net_http.StatusServiceUnavailable: // 503
		return true

	case net_http.StatusGatewayTimeout: // 504
		return true

	case 529: // Non-standard "site is overloaded" status
		return true

	default:
		return false

	}

}
