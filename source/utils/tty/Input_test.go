package tty

import "testing"

func TestDecodeSGRMousePress(t *testing.T) {

	key, err := decodeSGRMouse("<0;10;5", 'M')

	if err != nil {
		t.Fatalf("Expected a mouse key, got %s", err.Error())
	}

	if key.Type != KeyMouse {
		t.Errorf("Expected KeyMouse")
	}

	if key.MouseButton != 0 || key.MouseX != 9 || key.MouseY != 4 {
		t.Errorf("Unexpected mouse key: %+v", key)
	}

	if key.MouseRelease == true {
		t.Errorf("Expected a press")
	}

}

func TestDecodeSGRMouseRelease(t *testing.T) {

	key, err := decodeSGRMouse("<0;10;5", 'm')

	if err != nil {
		t.Fatalf("Expected a mouse key, got %s", err.Error())
	}

	if key.MouseRelease == false {
		t.Errorf("Expected a release")
	}

}

func TestDecodeSGRMouseWheel(t *testing.T) {

	key, err := decodeSGRMouse("<64;1;1", 'M')

	if err != nil {
		t.Fatalf("Expected a mouse key, got %s", err.Error())
	}

	if key.MouseWheel != -1 {
		t.Errorf("Expected wheel up, got %d", key.MouseWheel)
	}

}
