package components

import "exocomp/engine"
import "exocomp/tools/bugs"
import "exocomp/tools/changelog"
import "exocomp/types"
import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"
import "fmt"
import "strings"
import "time"

const sidebarBreakpointColumns = 124

type sidebarLine struct {
	Text  string
	Style utils_tty.Style
}

type Sidebar struct {
	Session      *engine.Session
	Scroll       int
	forceVisible bool
	lastRefresh  time.Time
	lines        []sidebarLine
}

func NewSidebar(session *engine.Session) *Sidebar {

	return &Sidebar{
		Session:      session,
		forceVisible: true,
		lines:        make([]sidebarLine, 0),
		lastRefresh:  time.Time{},
	}

}

func (sidebar *Sidebar) Toggle() {
	sidebar.forceVisible = (sidebar.forceVisible == false)
}

func (sidebar *Sidebar) IsVisible(width int) bool {

	if width > sidebarBreakpointColumns {
		return sidebar.forceVisible
	} else {
		return false
	}
}

func (sidebar *Sidebar) Width(available int) int {

	if available <= sidebarBreakpointColumns {
		return 0
	} else {
		return min(40, max(28, available/4))
	}
}

func (sidebar *Sidebar) Height(width int) int {
	return 0
}

func (sidebar *Sidebar) ScrollBy(delta int) {

	sidebar.Scroll += delta

	if sidebar.Scroll < 0 {
		sidebar.Scroll = 0
	}

}

func (sidebar *Sidebar) Refresh(now time.Time, is_force bool) {

	if is_force == false && now.Sub(sidebar.lastRefresh) < time.Second {
		return
	}

	sidebar.lastRefresh = now

	changelog_entries := make([]types.ChangelogEntry, 0)
	bug_entries := make([]types.Bug, 0)

	tool := sidebar.Session.GetTool("changelog")

	if tool != nil {

		changelog_tool, ok := tool.(*changelog.Changelog)

		if ok == true {

			entries, err := changelog_tool.Snapshot()

			if err == nil {
				changelog_entries = entries
			}

		}

	}

	tool = sidebar.Session.GetTool("bugs")

	if tool != nil {

		bugs_tool, ok := tool.(*bugs.Bugs)

		if ok == true {

			open_bugs, err := bugs_tool.Snapshot()

			if err == nil {
				bug_entries = open_bugs
			}

		}

	}

	lines := make([]sidebarLine, 0)
	width := max(20, sidebar.Width(120)-2)

	if len(changelog_entries) > 0 {

		lines = append(lines, sidebarLine{Text: "CHANGELOG", Style: styles.SidebarSection})

		for index := len(changelog_entries) - 1; index >= 0; index-- {

			entry := changelog_entries[index]
			title := fmt.Sprintf("%s %s %s", entry.Date.Format("2006-01-02"), entry.Type, entry.File)

			if entry.Symbol != "" {
				title += ":" + entry.Symbol
			}

			for _, line := range utils_tty.Wrap(title, width) {
				lines = append(lines, sidebarLine{Text: line, Style: styles.SidebarText})
			}

			for _, line := range utils_tty.Wrap(entry.Description, width) {
				lines = append(lines, sidebarLine{Text: "  " + line, Style: styles.SidebarMuted})
			}

			lines = append(lines, sidebarLine{Text: "", Style: styles.SidebarPanel})

		}

	}

	if len(bug_entries) > 0 {

		lines = append(lines, sidebarLine{Text: "OPEN BUGS", Style: styles.SidebarSection})

		for _, bug := range bug_entries {

			title := fmt.Sprintf("%s:%s", bug.File, bug.Symbol)

			for _, line := range utils_tty.Wrap(title, width) {
				lines = append(lines, sidebarLine{Text: line, Style: styles.SidebarPanelError})
			}

			for _, line := range utils_tty.Wrap(bug.Description, width) {
				lines = append(lines, sidebarLine{Text: "  " + line, Style: styles.SidebarMuted})
			}

			lines = append(lines, sidebarLine{Text: "", Style: styles.SidebarPanel})

		}

	}

	sidebar.lines = lines

}

func (sidebar *Sidebar) Render(buffer *utils_tty.Buffer, rect utils_tty.Rect, elapsed time.Duration) {

	if rect.Width <= 0 || rect.Height <= 0 {
		return
	}

	buffer.Fill(rect, ' ', styles.SidebarPanel)

	content_x := rect.X + 2
	content_width := max(1, rect.Width-5)
	content_y := rect.Y + 1

	title := "Session"

	if sidebar.Session != nil && sidebar.Session.Agent != nil {

		if strings.TrimSpace(sidebar.Session.Agent.Name) != "" {
			title = sidebar.Session.Agent.Name
		}

	}

	if rect.Height >= 5 {

		buffer.WriteString(content_x, rect.Y+1, utils_tty.Truncate(title, max(0, rect.Width-4)), styles.SidebarTitle)
		content_y = rect.Y + 3

	}

	content_height := max(0, rect.Y+rect.Height-1-content_y)

	if len(sidebar.lines) == 0 {

		if content_height > 0 {

			message := "No changelog or open bugs yet."
			buffer.WriteString(content_x, content_y, utils_tty.Truncate(message, content_width), styles.SidebarMuted)

		}

		return

	}

	if sidebar.Scroll > len(sidebar.lines)-content_height {
		sidebar.Scroll = max(0, len(sidebar.lines)-content_height)
	}

	start := sidebar.Scroll

	for row := 0; row < content_height && start+row < len(sidebar.lines); row++ {

		line := sidebar.lines[start+row]
		text := utils_tty.PadRight(utils_tty.Truncate(line.Text, content_width), content_width)
		buffer.WriteString(content_x, content_y+row, text, line.Style)

	}

	drawSidebarScrollbar(
		buffer,
		rect.X+rect.Width-1,
		content_y,
		content_height,
		len(sidebar.lines),
		sidebar.Scroll,
	)

}
