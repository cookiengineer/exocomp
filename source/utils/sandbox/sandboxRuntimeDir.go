package sandbox

import "crypto/sha256"
import "encoding/hex"
import "fmt"
import "os"
import "path/filepath"

// sandboxRuntimeDir returns the private scratch directory for a sandbox. It is
// kept outside the sandbox so the private .exocomp state folder can be masked
// from executed programs without taking away HOME, TMPDIR and build caches.
func sandboxRuntimeDir(sandbox string) (string, error) {

	absolute, err0 := filepath.Abs(sandbox)

	if err0 != nil {
		return "", fmt.Errorf("Cannot resolve sandbox path.")
	}

	sum := sha256.Sum256([]byte(absolute))
	hash := hex.EncodeToString(sum[:])[0:16]

	return filepath.Join(os.TempDir(), "exocomp-runtime", hash), nil

}
