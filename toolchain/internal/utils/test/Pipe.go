package test

type Pipe interface {
	Read([]byte) (int, error)
}
