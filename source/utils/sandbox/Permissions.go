package sandbox

// Permissions describes extra access granted to a sandboxed process.
type Permissions struct {
	ReadPaths  []string
	WritePaths []string
	Network    bool
}
