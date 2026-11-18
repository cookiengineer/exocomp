package capabilities

import "strings"

// IsValidDomain validates an exact DNS domain name.
func IsValidDomain(domain string) bool {

	domain = strings.TrimSuffix(domain, ".")

	if domain == "" || len(domain) > 253 {
		return false
	}

	for _, label := range strings.Split(domain, ".") {

		if len(label) == 0 || len(label) > 63 {
			return false
		}

		for index, value := range label {

			is_letter := value >= 'a' && value <= 'z'
			is_digit := value >= '0' && value <= '9'
			is_hyphen := value == '-'

			if is_letter == false && is_digit == false && is_hyphen == false {
				return false
			}

			if is_hyphen == true && (index == 0 || index == len(label)-1) {
				return false
			}

		}

	}

	return true

}
