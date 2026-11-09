package landlock

// Policy describes the filesystem access rules for a sandboxed process.
type Policy struct {
	ReadWrite []string
	ReadOnly  []string
	Devices   []string
}
