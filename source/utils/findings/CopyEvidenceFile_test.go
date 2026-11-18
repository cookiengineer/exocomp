package findings

import "crypto/sha256"
import "encoding/hex"
import "os"
import "path/filepath"
import "testing"

func TestCopyEvidenceFile(t *testing.T) {

	playground := newFindingsTestPlayground(t)
	source := filepath.Join(playground, "proof_test.go")
	content := "package findings\n"

	err0 := os.WriteFile(source, []byte(content), 0644)

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	artifact, hash, size, err1 := CopyEvidenceFile(playground, "fnd_4a91c2e0", "ev_4a91c2e0", source, "proof_test.go")

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	sum := sha256.Sum256([]byte(content))
	expected_hash := hex.EncodeToString(sum[:])

	if hash != expected_hash {
		t.Errorf("Expected SHA-256 %s, got %s", expected_hash, hash)
	}

	if size != int64(len(content)) {
		t.Errorf("Expected size %d, got %d", len(content), size)
	}

	if artifact == "" {
		t.Errorf("Expected evidence artifact path")
	}

}

func TestCopyEvidenceFile_MissingSource(t *testing.T) {

	playground := newFindingsTestPlayground(t)

	_, _, _, err0 := CopyEvidenceFile(playground, "fnd_4a91c2e0", "ev_4a91c2e0", filepath.Join(playground, "missing"), "missing.txt")

	if err0 == nil {
		t.Errorf("Expected missing evidence artifact to fail")
	}

}
