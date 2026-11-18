package sandbox

// Permissions describes extra filesystem access granted to a sandboxed process.
type Permissions struct {
	ReadPaths  []string
	WritePaths []string
}
