package shell

import "testing"

func TestSanitizeProgramArgumentNetworkValues(t *testing.T) {

	sandbox := t.TempDir()

	valid := []string{
		"https://example.com/",
		"http://example.com:8080/admin?x=1",
		"wss://example.com/socket",
		"10.0.0.0/8",
		"10.0.0.0/8,192.168.0.0/16",
		"::1/128",
	}

	for _, argument := range valid {

		result, err := SanitizeProgramArgument(sandbox, nil, argument)

		if err != nil {
			t.Fatalf("SanitizeProgramArgument(%q) returned error: %v", argument, err)
		}

		if result != argument {
			t.Fatalf("SanitizeProgramArgument(%q) = %q, want unchanged", argument, result)
		}

	}

}

func TestSanitizeProgramArgumentPaths(t *testing.T) {

	sandbox := t.TempDir()

	cases := []struct {
		argument string
		expected string
	}{
		{"results.jsonl", "results.jsonl"},
		{"./results.jsonl", "./results.jsonl"},
		{"scans/nuclei/results.jsonl", "./scans/nuclei/results.jsonl"},
		{"http/cves/test.yaml", "./http/cves/test.yaml"},
	}

	for _, test := range cases {

		result, err := SanitizeProgramArgument(sandbox, nil, test.argument)

		if err != nil {
			t.Fatalf("SanitizeProgramArgument(%q) returned error: %v", test.argument, err)
		}

		if result != test.expected {
			t.Fatalf("SanitizeProgramArgument(%q) = %q, want %q", test.argument, result, test.expected)
		}

	}

}

func TestSanitizeProgramArgumentRejectsEscapes(t *testing.T) {

	sandbox := t.TempDir()

	invalid := []string{
		"../escape.txt",
		"/etc/passwd",
		"../../etc/passwd",
	}

	for _, argument := range invalid {

		if _, err := SanitizeProgramArgument(sandbox, nil, argument); err == nil {
			t.Fatalf("SanitizeProgramArgument(%q) returned no error", argument)
		}

	}

}
