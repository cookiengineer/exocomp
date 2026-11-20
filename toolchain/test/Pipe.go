package main

type Pipe interface {
	Read([]byte) (int, error)
}
