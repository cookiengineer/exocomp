package lint

// PackageRule checks a whole directory.
type PackageRule struct {
	ID      string
	Summary string
	Check   func(*Context, *Package) []Diagnostic
}
