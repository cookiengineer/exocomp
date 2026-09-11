package fmt

import "testing"

func TestFormatUSD(t *testing.T) {

	tests := []struct {
		amount float64
		want   string
	}{
		{0.0, "$0.0000"},
		{0.0001, "$0.0001"},
		{0.0042, "$0.0042"},
		{0.01, "$0.01"},
		{1.5, "$1.50"},
		{1234.567, "$1234.57"},
	}

	for _, test := range tests {

		result := FormatUSD(test.amount)

		if result != test.want {
			t.Errorf("Expected %q, got %q", test.want, result)
		}

	}

}
