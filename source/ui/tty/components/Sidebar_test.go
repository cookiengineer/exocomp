package components

import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"
import "fmt"
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

func TestSidebarRenderPanelWithoutDivider(t *testing.T) {

	sidebar := NewSidebar(nil)
	sidebar.lines = make([]sidebarLine, 0)

	for index := 0; index < 20; index++ {
		sidebar.lines = append(sidebar.lines, sidebarLine{
			Text:  fmt.Sprintf("line %d", index),
			Style: styles.SidebarText,
		})
	}

	buffer := utils_tty.NewBuffer(40, 10)
	sidebar.Render(buffer, utils_tty.Rect{X: 0, Y: 0, Width: 40, Height: 10}, 0)

	if buffer.Get(0, 0).Style != styles.SidebarPanel {
		t.Errorf("Expected the raised sidebar panel background")
	}

	for y := 0; y < 10; y++ {

		if buffer.Get(0, y).Text == "┃" {
			t.Errorf("Expected no sidebar divider at row %d", y)
		}

	}

	if buffer.Get(2, 1).Text != "S" {
		t.Errorf("Expected the sidebar title, got %q", buffer.Get(2, 1).Text)
	}

	if buffer.Get(2, 1).Style.Attributes&utils_tty.AttrBold == 0 {
		t.Errorf("Expected a bold sidebar title")
	}

	found_thumb := false

	for y := 3; y < 9; y++ {

		if buffer.Get(39, y).Text == "█" {
			found_thumb = true
		}

	}

	if found_thumb == false {
		t.Errorf("Expected a visible sidebar scrollbar")
	}

}

func TestSidebarRenderEmptyState(t *testing.T) {

	sidebar := NewSidebar(nil)
	buffer := utils_tty.NewBuffer(40, 10)
	sidebar.Render(buffer, utils_tty.Rect{X: 0, Y: 0, Width: 40, Height: 10}, 0)

	line := renderTestLine(buffer, 3, 40)

	if line == "" {
		t.Errorf("Expected an empty state line")
	}

	if buffer.Get(0, 0).Style != styles.SidebarPanel {
		t.Errorf("Expected the raised sidebar panel background")
	}

}
