package components

import "exocomp/engine"
import "exocomp/tools"
import "exocomp/types"
import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"
import "fmt"
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
	sidebar.forceVisible = !sidebar.forceVisible
}

func (sidebar *Sidebar) IsVisible(width int) bool {

	if width <= sidebarBreakpointColumns {
		return false
	}

	return sidebar.forceVisible

}

func (sidebar *Sidebar) Width(available int) int {

	if available <= sidebarBreakpointColumns {
		return 0
	}

	return min(40, max(28, available/4))

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

func (sidebar *Sidebar) Refresh(now time.Time, force bool) {

	if force == false && now.Sub(sidebar.lastRefresh) < time.Second {
		return
	}

	sidebar.lastRefresh = now

	changelog_entries := make([]types.ChangelogEntry, 0)
	bugs := make([]types.Bug, 0)

	tool := sidebar.Session.GetTool("changelog")

	if tool != nil {

		changelog_tool, ok := tool.(*tools.Changelog)

		if ok == true {

			entries, err := changelog_tool.Snapshot()

			if err == nil {
				changelog_entries = entries
			}

		}

	}

	tool = sidebar.Session.GetTool("bugs")

	if tool != nil {

		bugs_tool, ok := tool.(*tools.Bugs)

		if ok == true {

			open_bugs, err := bugs_tool.Snapshot()

			if err == nil {
				bugs = open_bugs
			}

		}

	}

	lines := make([]sidebarLine, 0)
	width := max(20, sidebar.Width(120)-2)

	if len(changelog_entries) > 0 {

		lines = append(lines, sidebarLine{Text: "CHANGELOG", Style: styles.SidebarHeader})

		for index := len(changelog_entries) - 1; index >= 0; index-- {

			entry := changelog_entries[index]
			title := fmt.Sprintf("%s %s %s", entry.Date.Format("2006-01-02"), entry.Type, entry.File)

			if entry.Symbol != "" {
				title += ":" + entry.Symbol
			}

			for _, line := range utils_tty.Wrap(title, width) {
				lines = append(lines, sidebarLine{Text: line, Style: styles.Text})
			}

			for _, line := range utils_tty.Wrap(entry.Description, width) {
				lines = append(lines, sidebarLine{Text: "  " + line, Style: styles.TextMuted})
			}

			lines = append(lines, sidebarLine{Text: "", Style: styles.Text})

		}

	}

	if len(bugs) > 0 {

		lines = append(lines, sidebarLine{Text: "OPEN BUGS", Style: styles.SidebarHeader})

		for _, bug := range bugs {

			title := fmt.Sprintf("%s:%s", bug.File, bug.Symbol)

			for _, line := range utils_tty.Wrap(title, width) {
				lines = append(lines, sidebarLine{Text: line, Style: styles.SidebarError})
			}

			for _, line := range utils_tty.Wrap(bug.Description, width) {
				lines = append(lines, sidebarLine{Text: "  " + line, Style: styles.TextMuted})
			}

			lines = append(lines, sidebarLine{Text: "", Style: styles.Text})

		}

	}

	sidebar.lines = lines

}

func (sidebar *Sidebar) Render(buffer *utils_tty.Buffer, rect utils_tty.Rect, elapsed time.Duration) {

	if rect.Width <= 0 || rect.Height <= 0 {
		return
	}

	background := styles.Background

	buffer.Fill(rect, ' ', background)
	buffer.Fill(utils_tty.Rect{X: rect.X, Y: rect.Y, Width: 1, Height: rect.Height}, ' ', background)

	divider_style := styles.BorderRaised

	for row := rect.Y; row < rect.Y+rect.Height; row++ {
		buffer.Set(rect.X, row, "┃", divider_style)
	}

	if len(sidebar.lines) == 0 {

		message := "No changelog or open bugs yet."
		buffer.WriteString(rect.X+2, rect.Y, utils_tty.Truncate(message, rect.Width-2), styles.TextMuted)
		return

	}

	if sidebar.Scroll > len(sidebar.lines)-rect.Height {
		sidebar.Scroll = max(0, len(sidebar.lines)-rect.Height)
	}

	start := sidebar.Scroll

	for row := 0; row < rect.Height && start+row < len(sidebar.lines); row++ {

		line := sidebar.lines[start+row]
		text := utils_tty.PadRight(utils_tty.Truncate(line.Text, rect.Width-2), rect.Width-2)
		buffer.WriteString(rect.X+2, rect.Y+row, text, line.Style)

	}

}
