package tty

import "testing"

func TestFromHex(t *testing.T) {

	color, err := FromHex("#fab283")

	if err == nil {
		if color.Hex() != "#fab283" {
			t.Errorf("Expected #fab283, got %s", color.Hex())
		}
	} else {
		t.Errorf("Expected nil, got %v", err)
	}

}

func TestStringWidth(t *testing.T) {

	if StringWidth("abc") != 3 {
		t.Errorf("Expected 3, got %d", StringWidth("abc"))
	}

	if StringWidth("界") != 2 {
		t.Errorf("Expected 2, got %d", StringWidth("界"))
	}

}

func TestBufferWriteString(t *testing.T) {

	buffer := NewBuffer(8, 1)
	style := NewStyle(FromInts(255, 255, 255, 255), FromInts(0, 0, 0, 255))

	buffer.WriteString(0, 0, "hello", style)

	if buffer.Get(0, 0).Text != "h" {
		t.Errorf("Expected h, got %s", buffer.Get(0, 0).Text)
	}

	if buffer.Get(4, 0).Text != "o" {
		t.Errorf("Expected o, got %s", buffer.Get(4, 0).Text)
	}

}

func TestActivityFrames(t *testing.T) {

	activity := NewActivity(FromInts(250, 178, 131, 255), FromInts(10, 10, 10, 255))
	frames := activity.Frames(4)

	if len(frames) == 0 {
		t.Errorf("Expected activity frames")
	}

	if StringWidth(frames[0]) != 4 {
		t.Errorf("Expected 4 cells, got %d", StringWidth(frames[0]))
	}

}
