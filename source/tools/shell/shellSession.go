package shell

import utils_shell "exocomp/utils/shell"
import "context"
import "fmt"
import "io"
import "os"
import "os/exec"
import "sync"
import "time"

import pty "github.com/creack/pty"

type shellSession struct {
	ID           string
	Program      string
	Arguments    []string
	Workdir      string
	Columns      int
	Rows         int
	Status       string
	ExitCode     int
	StartedAt    time.Time
	FinishedAt   time.Time
	LastActivity time.Time

	cmd    *exec.Cmd
	pty    *os.File
	cancel context.CancelFunc

	logFile      *os.File
	logPath      string
	logSize      int64
	logTruncated bool
	dropped      int64

	output []byte
	filter utils_shell.TerminalFilter

	mutex      sync.Mutex
	writeMutex sync.Mutex
	notify     chan struct{}
	done       chan struct{}
	finishOnce sync.Once
}

func (session *shellSession) run() {

	buffer := make([]byte, 4096)

	for {

		read, err := session.pty.Read(buffer)

		if read > 0 {
			session.appendOutput(buffer[0:read])
		}

		if err == nil {
			continue
		} else {
			break
		}

	}

	err := session.cmd.Wait()
	session.finish(err)

}

func (session *shellSession) appendOutput(data []byte) {

	cleaned := session.filter.Write(data)

	if len(cleaned) == 0 {

		session.touch()

	} else {

		session.mutex.Lock()

		session.LastActivity = time.Now()
		session.output = append(session.output, cleaned...)

		if len(session.output) > shellMaxBuffer {

			drop := len(session.output) - shellMaxBuffer
			session.output = session.output[drop:]
			session.dropped += int64(drop)

		}

		if session.logFile != nil && session.logTruncated == false {

			remaining := int64(shellMaxLog) - session.logSize

			if remaining > 0 {

				chunk := cleaned

				if int64(len(chunk)) > remaining {
					chunk = chunk[0:int(remaining)]
				}

				written, err := session.logFile.Write(chunk)

				if written > 0 {
					session.logSize += int64(written)
				}

				session.logTruncated = err != nil || int64(len(cleaned)) > remaining

			} else {
				session.logTruncated = true
			}

		}

		session.mutex.Unlock()
		session.notifyOutput()

	}

}

func (session *shellSession) read(wait time.Duration, limit int) ([]byte, string, int) {

	if wait > 0 {

		timer := time.NewTimer(wait)

		for {

			session.mutex.Lock()
			ready := len(session.output) > 0 || (session.Status != "starting" && session.Status != "running")
			session.mutex.Unlock()

			if ready == true {

				timer.Stop()
				break

			}

			select {

			case <-session.notify:
			case <-session.done:
			case <-timer.C:
				session.mutex.Lock()
				data, status, exit_code := session.readLocked(limit)
				session.mutex.Unlock()
				return data, status, exit_code

			}

		}

	}

	session.mutex.Lock()
	data, status, exit_code := session.readLocked(limit)
	session.mutex.Unlock()

	return data, status, exit_code

}

func (session *shellSession) readLocked(limit int) ([]byte, string, int) {

	size := limit

	if size > len(session.output) {
		size = len(session.output)
	}

	data := make([]byte, size)
	copy(data, session.output[0:size])
	session.output = session.output[size:]
	session.LastActivity = time.Now()

	return data, session.Status, session.ExitCode

}

func (session *shellSession) write(data []byte) (int, error) {

	session.mutex.Lock()
	status := session.Status
	session.mutex.Unlock()

	if status == "starting" || status == "running" {

		session.writeMutex.Lock()
		defer session.writeMutex.Unlock()

		written := 0

		for len(data) > 0 {

			size, err := session.pty.Write(data)

			if size > 0 {

				written += size
				data = data[size:]

			}

			if err == nil {
				continue
			} else {

				session.touch()
				return written, err

			}

		}

		session.touch()

		return written, nil

	} else {
		return 0, fmt.Errorf("Session is not running.")
	}

}

func (session *shellSession) resize(columns int, rows int) error {

	session.mutex.Lock()
	status := session.Status
	session.mutex.Unlock()

	if status == "starting" || status == "running" {

		err := pty.Setsize(session.pty, &pty.Winsize{
			Cols: uint16(columns),
			Rows: uint16(rows),
		})

		if err == nil {

			session.mutex.Lock()
			session.Columns = columns
			session.Rows = rows
			session.LastActivity = time.Now()
			session.mutex.Unlock()

			return nil

		} else {
			return err
		}

	} else {
		return fmt.Errorf("Session is not running.")
	}

}

func (session *shellSession) signal(name string) (string, error) {

	session.mutex.Lock()
	status := session.Status
	pid := 0

	if session.cmd != nil && session.cmd.Process != nil {
		pid = session.cmd.Process.Pid
	}

	session.mutex.Unlock()

	if status == "starting" || status == "running" {

		signal, err := utils_shell.SendSignalByName(pid, name)

		if err == nil {

			session.touch()
			return signal, nil

		} else {
			return "", fmt.Errorf("Cannot send signal \"%s\".", name)
		}

	} else {
		return "", fmt.Errorf("Session is not running.")
	}

}

func (session *shellSession) close(is_force bool) ([]byte, string, int) {

	session.mutex.Lock()
	status := session.Status
	pid := 0

	if session.cmd != nil && session.cmd.Process != nil {
		pid = session.cmd.Process.Pid
	}

	session.mutex.Unlock()

	if status == "starting" || status == "running" {

		if is_force == true {

			_ = utils_shell.KillProcess(pid)

		} else {

			_, _ = utils_shell.SendSignalByName(pid, "TERM")

			select {

			case <-session.done:
			case <-time.After(shellCloseGrace):
				_ = utils_shell.KillProcess(pid)
				select {

				case <-session.done:
				case <-time.After(shellCloseGrace):
					_ = session.pty.Close()

				}

			}

		}

	}

	select {

	case <-session.done:
	case <-time.After(shellCloseWait):
		_ = session.pty.Close()

	}

	data, status, exit_code := session.read(0, shellMaxRead)

	session.mutex.Lock()

	if session.Status == "starting" || session.Status == "running" {

		session.Status = "closed"
		status = "closed"

	}

	if session.logFile != nil {

		session.logFile.Close()
		session.logFile = nil

	}

	session.mutex.Unlock()

	return data, status, exit_code

}

func (session *shellSession) timeout() {

	session.mutex.Lock()

	if session.Status == "starting" || session.Status == "running" {
		session.Status = "timed_out"
	}

	pid := 0

	if session.cmd != nil && session.cmd.Process != nil {
		pid = session.cmd.Process.Pid
	}

	session.mutex.Unlock()

	_ = utils_shell.KillProcess(pid)

}

func (session *shellSession) finish(err error) {

	session.finishOnce.Do(func() {

		session.mutex.Lock()

		session.FinishedAt = time.Now()

		if session.Status == "starting" || session.Status == "running" {

			if err == nil {

				session.Status = "exited"
				session.ExitCode = 0

			} else if exit_error, ok := err.(*exec.ExitError); ok == true {

				session.Status = "exited"
				session.ExitCode = exit_error.ExitCode()

			} else {

				session.Status = "failed"
				session.ExitCode = -1

			}

		}

		if session.logFile != nil {

			session.logFile.Close()
			session.logFile = nil

		}

		session.mutex.Unlock()

		close(session.done)

		if session.cancel != nil {
			session.cancel()
		}

		session.notifyOutput()

	})

}

func (session *shellSession) touch() {

	session.mutex.Lock()
	session.LastActivity = time.Now()
	session.mutex.Unlock()

}

func (session *shellSession) status() (string, int) {

	session.mutex.Lock()
	defer session.mutex.Unlock()

	return session.Status, session.ExitCode

}

func (session *shellSession) outputStats() (int64, int64, bool) {

	session.mutex.Lock()
	defer session.mutex.Unlock()

	return session.dropped, session.logSize, session.logTruncated

}

func (session *shellSession) pendingOutput() int {

	session.mutex.Lock()
	defer session.mutex.Unlock()

	return len(session.output)

}

func (session *shellSession) notifyOutput() {

	select {

	case session.notify <- struct{}{}:
	default:

	}

}

func (session *shellSession) readLog(offset int64, limit int) ([]byte, int64, bool, error) {

	session.mutex.Lock()
	path := session.logPath
	size := session.logSize
	truncated := session.logTruncated
	session.mutex.Unlock()

	if offset > size {
		offset = size
	}

	length := int64(limit)

	if offset+length > size {
		length = size - offset
	}

	if length > 0 {

		file, err0 := os.Open(path)

		if err0 == nil {

			defer file.Close()

			buffer := make([]byte, length)
			read, err1 := file.ReadAt(buffer, offset)

			if err1 == nil || err1 == io.EOF {
				return buffer[0:read], size, truncated, nil
			} else {
				return nil, size, truncated, fmt.Errorf("Cannot read session log.")
			}

		} else {
			return nil, size, truncated, fmt.Errorf("Cannot open session log.")
		}

	} else {
		return []byte{}, size, truncated, nil
	}

}
