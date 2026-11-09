//go:build linux

package landlock

import "errors"
import "unsafe"

import "golang.org/x/sys/unix"

// NOTE: This file intentionally groups the low-level landlock syscall helpers
// as one cohesive implementation. It is a documented exception to the
// one-function-per-file rule in CODESTYLE.md.

const (
	accessExecute    = uint64(unix.LANDLOCK_ACCESS_FS_EXECUTE)
	accessWriteFile  = uint64(unix.LANDLOCK_ACCESS_FS_WRITE_FILE)
	accessReadFile   = uint64(unix.LANDLOCK_ACCESS_FS_READ_FILE)
	accessReadDir    = uint64(unix.LANDLOCK_ACCESS_FS_READ_DIR)
	accessRemoveDir  = uint64(unix.LANDLOCK_ACCESS_FS_REMOVE_DIR)
	accessRemoveFile = uint64(unix.LANDLOCK_ACCESS_FS_REMOVE_FILE)
	accessMakeChar   = uint64(unix.LANDLOCK_ACCESS_FS_MAKE_CHAR)
	accessMakeDir    = uint64(unix.LANDLOCK_ACCESS_FS_MAKE_DIR)
	accessMakeReg    = uint64(unix.LANDLOCK_ACCESS_FS_MAKE_REG)
	accessMakeSock   = uint64(unix.LANDLOCK_ACCESS_FS_MAKE_SOCK)
	accessMakeFifo   = uint64(unix.LANDLOCK_ACCESS_FS_MAKE_FIFO)
	accessMakeBlock  = uint64(unix.LANDLOCK_ACCESS_FS_MAKE_BLOCK)
	accessMakeSym    = uint64(unix.LANDLOCK_ACCESS_FS_MAKE_SYM)
	accessRefer      = uint64(unix.LANDLOCK_ACCESS_FS_REFER)
	accessTruncate   = uint64(unix.LANDLOCK_ACCESS_FS_TRUNCATE)
	accessIoctlDev   = uint64(unix.LANDLOCK_ACCESS_FS_IOCTL_DEV)
)

var errUnsupported = errors.New("landlock unsupported")

func landlockABI() (int, error) {

	result, _, errno := unix.Syscall(
		unix.SYS_LANDLOCK_CREATE_RULESET,
		0,
		0,
		uintptr(unix.LANDLOCK_CREATE_RULESET_VERSION),
	)

	if errno == 0 {
		return int(result), nil
	} else {
		return 0, errno
	}

}

func allAccessFS(abi int) uint64 {

	access := accessExecute |
		accessWriteFile |
		accessReadFile |
		accessReadDir |
		accessRemoveDir |
		accessRemoveFile |
		accessMakeChar |
		accessMakeDir |
		accessMakeReg |
		accessMakeSock |
		accessMakeFifo |
		accessMakeBlock |
		accessMakeSym

	if abi >= 2 {
		access = access | accessRefer
	}

	if abi >= 3 {
		access = access | accessTruncate
	}

	if abi >= 5 {
		access = access | accessIoctlDev
	}

	return access

}

func readAccessFS() uint64 {
	return accessExecute | accessReadFile | accessReadDir
}

func deviceAccessFS() uint64 {
	return accessReadFile | accessWriteFile | accessIoctlDev
}

func createRuleset(abi int, access uint64) (int, error) {

	attr := unix.LandlockRulesetAttr{
		Access_fs: access,
	}

	size := uintptr(8)

	if abi >= 6 {
		size = unsafe.Sizeof(attr)
	} else if abi >= 4 {
		size = 16
	}

	result, _, errno := unix.Syscall(
		unix.SYS_LANDLOCK_CREATE_RULESET,
		uintptr(unsafe.Pointer(&attr)),
		size,
		0,
	)

	if errno == 0 {
		return int(result), nil
	} else {
		return -1, errno
	}

}

func addPathRule(ruleset_fd int, path string, access uint64) error {

	path_fd, err0 := unix.Open(path, unix.O_PATH|unix.O_CLOEXEC, 0)

	if err0 == nil {

		defer unix.Close(path_fd)

		attr := unix.LandlockPathBeneathAttr{
			Allowed_access: access,
			Parent_fd:      int32(path_fd),
		}

		_, _, errno := unix.Syscall6(
			unix.SYS_LANDLOCK_ADD_RULE,
			uintptr(ruleset_fd),
			uintptr(unix.LANDLOCK_RULE_PATH_BENEATH),
			uintptr(unsafe.Pointer(&attr)),
			0,
			0,
			0,
		)

		if errno == 0 {
			return nil
		} else {
			return errno
		}

	} else {
		return err0
	}

}

func addPathRules(ruleset_fd int, paths []string, access uint64, optional bool) error {

	for _, path := range paths {

		err := addPathRule(ruleset_fd, path, access)

		if err == nil {
			continue
		} else if optional == true && isOptionalPathError(err) == true {
			// NOTE: Default allowlist entries may be absent on this system.
			continue
		} else {
			return err
		}

	}

	return nil

}

func isOptionalPathError(err error) bool {

	if errors.Is(err, unix.ENOENT) ||
		errors.Is(err, unix.ENOTDIR) ||
		errors.Is(err, unix.EACCES) ||
		errors.Is(err, unix.EINVAL) ||
		errors.Is(err, unix.ENOTSUP) {
		return true
	} else {
		return false
	}

}

func restrictPaths(policy Policy) error {

	abi, err0 := landlockABI()

	if err0 != nil {
		return errUnsupported
	}

	access := allAccessFS(abi)

	if access == 0 {
		return errUnsupported
	}

	ruleset_fd, err1 := createRuleset(abi, access)

	if err1 != nil {
		return err1
	}

	defer unix.Close(ruleset_fd)

	err2 := addPathRules(ruleset_fd, policy.ReadWrite, access, false)

	if err2 != nil {
		return err2
	}

	err3 := addPathRules(ruleset_fd, policy.ReadOnly, readAccessFS()&access, true)

	if err3 != nil {
		return err3
	}

	err4 := addPathRules(ruleset_fd, policy.Devices, deviceAccessFS()&access, true)

	if err4 != nil {
		return err4
	}

	err5 := unix.Prctl(unix.PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0)

	if err5 != nil {
		return err5
	}

	return restrictSelf(ruleset_fd, abi)

}

func restrictSelf(ruleset_fd int, abi int) error {

	flags := uintptr(0)

	if abi >= 6 {
		flags = uintptr(unix.LANDLOCK_RESTRICT_SELF_TSYNC)
	}

	_, _, errno := unix.Syscall(
		unix.SYS_LANDLOCK_RESTRICT_SELF,
		uintptr(ruleset_fd),
		flags,
		0,
	)

	if errno == 0 {
		return nil
	} else {
		return errno
	}

}
