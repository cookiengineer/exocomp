package findings

import "crypto/sha256"
import "encoding/hex"
import "os"
import "strings"
import "testing"

func TestWriteEvidenceContent(t *testing.T) {

	playground := newFindingsTestPlayground(t)
	content := "HTTP/1.1 200 OK"

	artifact, hash, size, err0 := WriteEvidenceContent(playground, "fnd_4a91c2e0", "ev_4a91c2e0", "evidence.txt", content)

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	sum := sha256.Sum256([]byte(content))
	expected_hash := hex.EncodeToString(sum[:])

	if hash != expected_hash {
		t.Errorf("Expected SHA-256 %s, got %s", expected_hash, hash)
	}

	if size != int64(len(content)) {
		t.Errorf("Expected size %d, got %d", len(content), size)
	}

	if strings.HasPrefix(artifact, "./.exocomp/evidence/") == false {
		t.Errorf("Expected relative artifact path, got %s", artifact)
	}

	resolved := strings.TrimPrefix(artifact, "./")
	buffer, err1 := os.ReadFile(playground + "/" + resolved)

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	} else if string(buffer) != content {
		t.Errorf("Expected evidence content to match")
	}

}

func TestEvidenceArtifactPath_SanitizesFilename(t *testing.T) {

	playground := newFindingsTestPlayground(t)

	absolute, display, err0 := EvidenceArtifactPath(playground, "fnd_4a91c2e0", "ev_4a91c2e0", "../../secret?.txt")

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	if strings.Contains(absolute, "..") == true {
		t.Errorf("Expected sanitized absolute path, got %s", absolute)
	}

	if strings.HasPrefix(display, "./.exocomp/evidence/") == false {
		t.Errorf("Expected relative artifact display path, got %s", display)
	}

}
