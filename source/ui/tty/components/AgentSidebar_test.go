package components

import "exocomp/engine"
import "exocomp/types"
import "exocomp/ui/tty/styles"
import "strings"
import "testing"
import "time"

func newAgentSidebarTestSession(t *testing.T) *engine.Session {

	t.Helper()

	agent := &types.Agent{
		Name:         "Planner",
		Role:         "planner",
		Model:        "deepseek-v4-pro",
		AllowedTools: []string{"agents.List"},
	}

	config := types.NewConfig(
		"Planner",
		"planner",
		"deepseek-v4-pro",
		"",
		0.0,
		t.TempDir(),
		t.TempDir(),
		nil,
		false,
	)

	session := engine.NewSession(agent, config)
	session.InstallTools()

	session.InstallAgents([]*types.Agent{
		{
			Name:   "Coder",
			Role:   "coder",
			Status: "working",
		},
		{
			Name:   "Tester",
			Role:   "tester",
			Status: "finished",
		},
	})

	return session

}

func agentSidebarLineText(sidebar *AgentSidebar) string {

	builder := strings.Builder{}

	for _, line := range sidebar.lines {
		builder.WriteString(line.Text)
		builder.WriteString("\n")
	}

	return builder.String()

}

func TestAgentSidebarVisibleOnWideTerminals(t *testing.T) {

	sidebar := NewAgentSidebar(nil)

	if sidebar.IsVisible(200) != true {
		t.Errorf("Expected the agent sidebar to be visible on wide terminals")
	}

	if sidebar.Width(200) <= 0 {
		t.Errorf("Expected the agent sidebar to have a width on wide terminals")
	}

}

func TestAgentSidebarHiddenOnNarrowTerminals(t *testing.T) {

	sidebar := NewAgentSidebar(nil)

	if sidebar.IsVisible(agentSidebarBreakpointColumns) != false {
		t.Errorf("Expected the agent sidebar to be hidden at the breakpoint")
	}

	if sidebar.Width(agentSidebarBreakpointColumns) != 0 {
		t.Errorf("Expected the agent sidebar width to be zero at the breakpoint")
	}

}

func TestAgentSidebarListsPlannerAndHiredAgents(t *testing.T) {

	session := newAgentSidebarTestSession(t)
	sidebar := NewAgentSidebar(session)
	sidebar.Selected = "Coder"
	sidebar.Refresh(time.Now(), true)

	text := agentSidebarLineText(sidebar)

	if strings.Contains(text, "AGENTS") == false {
		t.Errorf("Expected an agents section, got %q", text)
	}

	planner_index := strings.Index(text, "Planner")
	coder_index := strings.Index(text, "Coder")
	tester_index := strings.Index(text, "Tester")

	if planner_index == -1 || coder_index == -1 || tester_index == -1 {
		t.Errorf("Expected all agents to be listed, got %q", text)
	}

	if planner_index > coder_index || coder_index > tester_index {
		t.Errorf("Expected planner first and hired agents sorted, got %q", text)
	}

	if strings.Contains(text, "planner") == false || strings.Contains(text, "working") == false || strings.Contains(text, "finished") == false {
		t.Errorf("Expected agent roles and statuses to be listed, got %q", text)
	}

	selected_found := false

	for _, line := range sidebar.lines {

		if strings.Contains(line.Text, "Coder") == true {

			if line.Style != styles.SidebarSelected {
				t.Errorf("Expected the selected agent to use the selected style")
			}

			selected_found = true

		}

	}

	if selected_found == false {
		t.Errorf("Expected the selected agent line")
	}

}

func TestAgentSidebarRenderEmptyState(t *testing.T) {

	sidebar := NewAgentSidebar(nil)
	sidebar.Refresh(time.Now(), true)

	if strings.Contains(agentSidebarLineText(sidebar), "No agents yet.") == false {
		t.Errorf("Expected an empty agent sidebar state")
	}

}
