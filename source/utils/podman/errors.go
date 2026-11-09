package podman

import "errors"

var errTimeout = errors.New("Command timed out.")
var errIdle = errors.New("Command stopped producing output.")
var errTruncated = errors.New("Output exceeded limit.")
var errNotOwned = errors.New("Container is not owned by this sandbox.")
