//go:build linux

package landlock

import "fmt"

import "golang.org/x/sys/unix"

// raiseLoopback brings the loopback interface of the helper's private network
// namespace up. It must run before privileges are dropped, because changing
// interface flags needs CAP_NET_ADMIN in the user namespace that owns the
// network namespace.
func raiseLoopback() error {

	fd, err0 := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM|unix.SOCK_CLOEXEC, 0)

	if err0 != nil {
		return fmt.Errorf("Cannot open socket.")
	}

	defer unix.Close(fd)

	ifreq, err1 := unix.NewIfreq("lo")

	if err1 != nil {
		return fmt.Errorf("Cannot resolve loopback interface.")
	}

	err2 := unix.IoctlIfreq(fd, unix.SIOCGIFFLAGS, ifreq)

	if err2 != nil {
		return fmt.Errorf("Cannot read loopback flags.")
	}

	ifreq.SetUint16(ifreq.Uint16() | unix.IFF_UP)

	err3 := unix.IoctlIfreq(fd, unix.SIOCSIFFLAGS, ifreq)

	if err3 != nil {
		return fmt.Errorf("Cannot raise loopback interface.")
	}

	return nil

}
