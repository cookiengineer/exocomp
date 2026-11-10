package chat

import "bytes"
import "context"
import "io"
import "fmt"
import net_http "net/http"

func ChatRequestWithRetry(client *net_http.Client, payload []byte, url string, token string) (*net_http.Response, error) {

	ctx, cancel := context.WithTimeout(context.Background(), RetryRequestTimeout)
	defer cancel()

	for attempt := 0; attempt < RetryMaxAttempts; attempt++ {

		request, err0 := net_http.NewRequestWithContext(
			ctx,
			net_http.MethodPost,
			url,
			bytes.NewReader(payload),
		)

		if err0 == nil {

			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			if token != "" {
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}

			response, err1 := client.Do(request)

			if err1 == nil {

				if response.StatusCode == net_http.StatusOK {

					payload, err := io.ReadAll(response.Body)
					response.Body.Close()

					if err == nil {

						response.Body = io.NopCloser(bytes.NewReader(payload))
						return response, nil

					} else {
						return nil, err
					}

				} else if isRetryableStatus(response.StatusCode) == true {

					if attempt+1 < RetryMaxAttempts {

						delay := retryDelay(attempt, response)

						if response != nil && response.Body != nil {
							io.Copy(io.Discard, io.LimitReader(response.Body, 4096))
							response.Body.Close()
						}

						if err2 := sleepRetry(ctx, delay); err2 != nil {
							return nil, err2
						}

					} else {

						payload, err := io.ReadAll(response.Body)
						response.Body.Close()

						if err == nil {

							response.Body = io.NopCloser(bytes.NewReader(payload))
							return response, nil

						} else {
							return nil, err
						}

					}

				} else {

					payload, err := io.ReadAll(response.Body)
					response.Body.Close()

					if err == nil {

						response.Body = io.NopCloser(bytes.NewReader(payload))
						return response, nil

					} else {
						return nil, err
					}

				}

			} else if isRetryableError(err1) == true && attempt+1 < RetryMaxAttempts {

				delay := retryDelay(attempt, nil)

				if err2 := sleepRetry(ctx, delay); err2 != nil {
					return nil, err2
				}

			} else {
				return nil, err1
			}

		} else {
			return nil, err0
		}

	}

	return nil, fmt.Errorf("ChatRequestWithRetry: Request failed after %d attempts.", RetryMaxAttempts)

}
