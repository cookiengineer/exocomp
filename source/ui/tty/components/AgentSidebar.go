package components

import "exocomp/engine"
import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"
import "fmt"
import "strings"
import "time"

const agentSidebarBreakpointColumns = 100

type AgentSidebar struct {
	Session      *engine.Session
	Selected     string
	Scroll       int
	forceVisible bool
	lastRefresh  time.Time
	lines        []sidebarLine
}

func NewAgentSidebar(session *engine.Session) *AgentSidebar {

	return &AgentSidebar{
		Session:      session,
		forceVisible: true,
		lines:        make([]sidebarLine, 0),
		lastRefresh:  time.Time{},
	}

}

func (sidebar *AgentSidebar) Toggle() {
	sidebar.forceVisible = !sidebar.forceVisible
}

func (sidebar *AgentSidebar) IsVisible(width int) bool {

	if width > agentSidebarBreakpointColumns {
		return sidebar.forceVisible
	} else {
		return false
	}
}

func (sidebar *AgentSidebar) Width(available int) int {

	if available <= agentSidebarBreakpointColumns {
		return 0
	} else {
		return min(32, max(20, available/5))
	}
}

func (sidebar *AgentSidebar) Height(width int) int {
	return 0
}

func (sidebar *AgentSidebar) ScrollBy(delta int) {

	sidebar.Scroll += delta

	if sidebar.Scroll < 0 {
		sidebar.Scroll = 0
	}

}

func (sidebar *AgentSidebar) Refresh(now time.Time, force bool) {

	if force == false && now.Sub(sidebar.lastRefresh) < time.Second {
		return
	}

	sidebar.lastRefresh = now

	agents := make([]string, 0)
	roles := make([]string, 0)
	statuses := make([]string, 0)

	if sidebar.Session != nil {

		for _, agent := range sidebar.Session.GetAgents() {

			status := agent.Status

			if sidebar.Session.Agent != nil && agent.Name == sidebar.Session.Agent.Name {
				status = "idle"

				if sidebar.Session.IsWaiting() == true {
					status = "thinking"
				}
			} else if strings.TrimSpace(status) == "" {
				status = "unknown"
			}

			agents = append(agents, agent.Name)
			roles = append(roles, agent.Role)
			statuses = append(statuses, status)

		}

	}

	lines := make([]sidebarLine, 0)

	if len(agents) > 0 {

		lines = append(lines, sidebarLine{Text: "AGENTS", Style: styles.SidebarSection})

		for index, name := range agents {

			prefix := "  "

			if name == sidebar.Selected {
				prefix = "▸ "
			}

			title := fmt.Sprintf("%s (%s, %s)", name, roles[index], statuses[index])
			style := styles.SidebarText

			if name == sidebar.Selected {
				style = styles.SidebarSelected
			}

			for _, line := range utils_tty.Wrap(prefix+title, max(20, sidebar.Width(120)-2)) {
				lines = append(lines, sidebarLine{Text: line, Style: style})
			}

		}

	} else {
		lines = append(lines, sidebarLine{Text: "No agents yet.", Style: styles.SidebarMuted})
	}

	sidebar.lines = lines

}

func (sidebar *AgentSidebar) Render(buffer *utils_tty.Buffer, rect utils_tty.Rect, elapsed time.Duration) {

	if rect.Width <= 0 || rect.Height <= 0 {
		return
	}

	buffer.Fill(rect, ' ', styles.SidebarPanel)

	content_x := rect.X + 2
	content_width := max(1, rect.Width-4)
	content_y := rect.Y + 1

	if rect.Height >= 5 {

		buffer.WriteString(content_x, rect.Y+1, utils_tty.Truncate("Agents", max(0, rect.Width-4)), styles.SidebarTitle)
		content_y = rect.Y + 3

	}

	content_height := max(0, rect.Y+rect.Height-1-content_y)

	if len(sidebar.lines) == 0 {
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
