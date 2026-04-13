package tools

import "os"
import "path/filepath"
import "testing"

func TestFiles_List(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-*")
	sandbox       := filepath.Join(playground, "sub", "package")
	tool          := NewChangelog("tester", sandbox, playground)

	if tool != nil {
	} else {
		t.Errorf("Expected tool to be not nil")
	}

	// TODO

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestFiles_Read(t *testing.T) {
	// TODO
	t.Errorf("TODO: Implement Files.Read unit test")
}

func TestFiles_Stat(t *testing.T) {
	// TODO
	t.Errorf("TODO: Implement Files.Stat unit test")
}

func TestFiles_Write(t *testing.T) {
	// TODO
	t.Errorf("TODO: Implement Files.Write unit test")
}
