package components

import "testing"

func TestSidebarVisibleByDefaultOnWideTerminals(t *testing.T) {

	sidebar := NewSidebar(nil)

	if sidebar.IsVisible(200) != true {
		t.Errorf("Expected sidebar to be visible by default on wide terminals")
	}

	if sidebar.Width(200) <= 0 {
		t.Errorf("Expected sidebar to have a width on wide terminals")
	}

}

func TestSidebarHiddenOnNarrowTerminals(t *testing.T) {

	sidebar := NewSidebar(nil)

	if sidebar.IsVisible(sidebarBreakpointColumns) != false {
		t.Errorf("Expected sidebar to be hidden at the breakpoint")
	}

	if sidebar.Width(sidebarBreakpointColumns) != 0 {
		t.Errorf("Expected sidebar width to be zero at the breakpoint")
	}

}

func TestSidebarToggle(t *testing.T) {

	sidebar := NewSidebar(nil)

	sidebar.Toggle()

	if sidebar.IsVisible(200) != false {
		t.Errorf("Expected sidebar to be hidden after toggling")
	}

	sidebar.Toggle()

	if sidebar.IsVisible(200) != true {
		t.Errorf("Expected sidebar to be visible after toggling again")
	}

}
