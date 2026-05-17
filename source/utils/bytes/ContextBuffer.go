package bytes

import "bytes"
import "context"
import "sync"
import "time"

type ContextBuffer struct {
	mutex       sync.Mutex
	buffer      bytes.Buffer
	last_write  time.Time
	max_bytes   int
	truncated   bool
	cancel_once sync.Once
	cancel      context.CancelFunc
}

func NewContextBuffer (max_bytes int, cancel context.CancelFunc) *ContextBuffer {

	return &ContextBuffer{
		last_write: time.Now(),
		max_bytes:  max_bytes,
		cancel:     cancel,
	}

}

func (self *ContextBuffer) Write(data []byte) (int, error) {

	self.mutex.Lock()

	remaining := self.max_bytes - self.buffer.Len()

	if remaining <= 0 {

		self.truncated = true

		self.cancel_once.Do(func() {

			if self.cancel != nil {
				self.cancel()
			}

		})

		return len(data), nil

	}

	if len(data) > remaining {

		data = data[:remaining]

		self.truncated = true

		self.cancel_once.Do(func() {

			if self.cancel != nil {
				self.cancel()
			}

		})

	}

	self.last_write = time.Now()

	result, err := self.buffer.Write(data)

	self.mutex.Unlock()

	return result, err

}

func (self *ContextBuffer) String() string {

	self.mutex.Lock()

	result := self.buffer.String()

	if self.truncated {
		result += "\n[output truncated: limit exceeded]"
	}

	self.mutex.Unlock()

	return result

}

func (self *ContextBuffer) LastWrite() time.Time {

	self.mutex.Lock()
	result := self.last_write
	self.mutex.Unlock()

	return result

}

func (self *ContextBuffer) Len() int {

	self.mutex.Lock()
	length := self.buffer.Len()
	self.mutex.Unlock()

	return length

}

func (self *ContextBuffer) IsTruncated() bool {

	self.mutex.Lock()
	truncated := self.truncated
	self.mutex.Unlock()

	return truncated

}
