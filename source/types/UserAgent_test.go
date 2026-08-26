package types

import "strings"
import "testing"

func TestUserAgent_Get(t *testing.T) {

	useragent := GetUserAgent("chrome-windows")

	if useragent == nil {
		t.Fatalf("Expected chrome-windows User-Agent to exist")
	}

	if useragent.Name != "chrome-windows" {
		t.Errorf("Expected name %q, got %q", "chrome-windows", useragent.Name)
	}

	if useragent.UserAgent == "" {
		t.Errorf("Expected a non-empty User-Agent string")
	}

	if useragent.Platform != "Windows" {
		t.Errorf("Expected platform %q, got %q", "Windows", useragent.Platform)
	}

}

func TestUserAgent_GetMissing(t *testing.T) {

	useragent := GetUserAgent("does-not-exist")

	if useragent != nil {
		t.Errorf("Expected %v to be nil", useragent)
	}

}

func TestUserAgent_Header(t *testing.T) {

	useragent := GetUserAgent("chrome-windows")

	if useragent == nil {
		t.Fatalf("Expected chrome-windows User-Agent to exist")
	}

	header := useragent.Header()

	if header.Get("User-Agent") != useragent.UserAgent {
		t.Errorf("Expected User-Agent header to match")
	}

	if header.Get("Sec-CH-UA") == "" {
		t.Errorf("Expected a Sec-CH-UA client hint header")
	}

	if header.Get("Sec-CH-UA-Mobile") != "?0" {
		t.Errorf("Expected Sec-CH-UA-Mobile to be %q, got %q", "?0", header.Get("Sec-CH-UA-Mobile"))
	}

}

func TestUserAgent_Presets(t *testing.T) {

	if len(UserAgents) == 0 {
		t.Errorf("Expected at least one User-Agent preset")
	}

	seen := make(map[string]bool)

	for _, useragent := range UserAgents {

		if useragent == nil {
			t.Errorf("Expected non-nil User-Agent")
			continue
		}

		if useragent.Name == "" {
			t.Errorf("Expected a non-empty User-Agent name")
		}

		if seen[useragent.Name] == true {
			t.Errorf("Expected unique User-Agent name %q", useragent.Name)
		}

		seen[useragent.Name] = true

		if useragent.UserAgent == "" {
			t.Errorf("Expected %q to have a User-Agent string", useragent.Name)
		}

		if !strings.Contains(useragent.UserAgent, useragent.Name) && useragent.Name != "safari-iphone" {
			// most presets embed their engine in the UA, but not all
		}

	}

}
