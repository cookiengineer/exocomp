package chat

import net_http "net/http"
import "testing"

func TestIsRetryableStatus(t *testing.T) {

	retryable := []int{
		net_http.StatusRequestTimeout,
		net_http.StatusConflict,
		net_http.StatusTooEarly,
		net_http.StatusTooManyRequests,
		net_http.StatusInternalServerError,
		net_http.StatusBadGateway,
		net_http.StatusServiceUnavailable,
		net_http.StatusGatewayTimeout,
		529,
	}

	for _, status := range retryable {

		if isRetryableStatus(status) != true {
			t.Errorf("Expected status %d to be retryable", status)
		}

	}

	permanent := []int{
		net_http.StatusOK,
		net_http.StatusCreated,
		net_http.StatusMovedPermanently,
		net_http.StatusBadRequest,
		net_http.StatusUnauthorized,
		net_http.StatusPaymentRequired,
		net_http.StatusForbidden,
		net_http.StatusNotFound,
		net_http.StatusUnprocessableEntity,
	}

	for _, status := range permanent {

		if isRetryableStatus(status) != false {
			t.Errorf("Expected status %d to not be retryable", status)
		}

	}

}
