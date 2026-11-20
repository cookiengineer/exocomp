package main

// FileRule checks a single parsed file.
type FileRule struct {
	ID      string
	Summary string
	Check   func(*Context, *File) []Diagnostic
}
