package findings

import "strings"
import "testing"

func TestRedactSecrets(t *testing.T) {

	header := RedactSecrets("Authorization: Bearer supersecret")
	query := RedactSecrets("https://example.com/?token=supersecret&next=1")
	userinfo := RedactSecrets("https://user:password@example.com/")
	bearer := RedactSecrets("Bearer eyJhbGciOiJIUzI1NiJ9.payload.signature")

	if strings.Contains(header, "supersecret") == true {
		t.Errorf("Expected authorization header to be redacted, got %q", header)
	}

	if strings.Contains(query, "supersecret") == true {
		t.Errorf("Expected query token to be redacted, got %q", query)
	}

	if strings.Contains(userinfo, "password") == true {
		t.Errorf("Expected URL userinfo to be redacted, got %q", userinfo)
	}

	if strings.Contains(bearer, "eyJhbGciOiJIUzI1NiJ9") == true {
		t.Errorf("Expected bearer token to be redacted, got %q", bearer)
	}

}
