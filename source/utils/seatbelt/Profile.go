package seatbelt

// Profile represents a complete SBPL sandbox profile.
type Profile struct {
	source string
}

// String returns the SBPL source.
func (profile *Profile) String() string {
	return profile.source
}
