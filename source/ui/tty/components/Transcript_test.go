package components

import "exocomp/engine"
import "exocomp/schemas"
import "exocomp/types"
import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"
import net_url "net/url"
import "strings"
import "testing"

func newTranscriptTestSession(t *testing.T) *engine.Session {

	t.Helper()

	url, err := net_url.Parse("http://127.0.0.1:1/v1")

	if err != nil {
		t.Fatalf("Expected a valid URL, got %s", err.Error())
	}

	agent := &types.Agent{
		Name:  "Peanut Hamper",
		Role:  "planner",
		Model: "deepseek-v4-pro",
	}

	config := &types.Config{
		Model:      "deepseek-v4-pro",
		URL:        url,
		Playground: t.TempDir(),
	}

	return engine.NewSession(agent, config)

}

func transcriptLineText(line transcriptLine) string {

	builder := strings.Builder{}

	for _, span := range line.Spans {
		builder.WriteString(span.Text)
	}

	return builder.String()

}

func toggleTranscriptLine(t *testing.T, transcript *Transcript, buffer *utils_tty.Buffer, rect utils_tty.Rect, part string) {

	t.Helper()

	index := -1

	for line_index, line := range transcript.cached {

		if line.Toggle != nil && line.Toggle.Part == part {
			index = line_index
			break
		}

	}

	if index == -1 {
		t.Fatalf("Expected a toggle for %s", part)
	}

	y := rect.Y + index - transcript.lastStart

	if transcript.HandleMouse(0, y, rect) == false {
		t.Fatalf("Expected the toggle to handle a click")
	}

	transcript.Render(buffer, rect, 0)

}

func TestTranscriptAddsBottomGap(t *testing.T) {

	session := newTranscriptTestSession(t)
	session.Agent.Messages = append(session.Agent.Messages, &schemas.Message{
		Role:    "user",
		Content: "hello",
	})

	transcript := NewTranscript(session)
	buffer := utils_tty.NewBuffer(40, 10)
	transcript.Render(buffer, utils_tty.Rect{X: 0, Y: 0, Width: 40, Height: 10}, 0)

	if len(transcript.cached) == 0 {
		t.Fatalf("Expected transcript lines")
	}

	last := transcript.cached[len(transcript.cached)-1]

	if len(last.Spans) != 0 || last.HasFill == true {
		t.Errorf("Expected an empty bottom gap line")
	}

}

func TestTranscriptCollapsesSystemMessage(t *testing.T) {

	session := newTranscriptTestSession(t)
	session.Agent.Messages = append(session.Agent.Messages, &schemas.Message{
		Role:    "system",
		Content: "You are an Exocomp robot.\nThis is the full system prompt.",
	})

	transcript := NewTranscript(session)
	buffer := utils_tty.NewBuffer(60, 20)
	rect := utils_tty.Rect{X: 0, Y: 0, Width: 60, Height: 20}
	transcript.Render(buffer, rect, 0)

	found := false

	for _, line := range transcript.cached {

		if line.Toggle != nil && line.Toggle.Part == "system" {

			found = true

			if strings.Contains(transcriptLineText(line), "▸ System") == false {
				t.Errorf("Expected a collapsed system header, got %q", transcriptLineText(line))
			}

			if line.HasBorder == false || line.Border != styles.BorderSystem {
				t.Errorf("Expected the system border style")
			}

			if strings.Contains(transcriptLineText(line), "full system prompt") == true {
				t.Errorf("Expected only the first line in the collapsed header")
			}

		}

	}

	if found == false {
		t.Fatalf("Expected a system toggle")
	}

	toggleTranscriptLine(t, transcript, buffer, rect, "system")

	expanded := false

	for _, line := range transcript.cached {

		if line.Toggle != nil && line.Toggle.Part == "system" {

			if strings.Contains(transcriptLineText(line), "▾ System") == true {
				expanded = true
			}

		}

	}

	if expanded == false {
		t.Errorf("Expected the system message to expand")
	}

}

func TestTranscriptCollapsesReasoning(t *testing.T) {

	session := newTranscriptTestSession(t)
	session.Agent.Messages = append(session.Agent.Messages, &schemas.Message{
		Role:             "assistant",
		Content:          "Here is the answer.",
		ReasoningContent: "I should inspect the files first.",
	})

	transcript := NewTranscript(session)
	buffer := utils_tty.NewBuffer(60, 20)
	rect := utils_tty.Rect{X: 0, Y: 0, Width: 60, Height: 20}
	transcript.Render(buffer, rect, 0)

	toggleTranscriptLine(t, transcript, buffer, rect, "reasoning")

	expanded := false

	for _, line := range transcript.cached {

		if line.Toggle != nil && line.Toggle.Part == "reasoning" {

			if line.HasBorder == false || line.Border != styles.BorderThinking {
				t.Errorf("Expected the reasoning border style")
			}

			if strings.Contains(transcriptLineText(line), "▾ Thought") == true {
				expanded = true
			}

		}

	}

	if expanded == false {
		t.Errorf("Expected the reasoning message to expand")
	}

}

func TestTranscriptCollapsesToolMessage(t *testing.T) {

	session := newTranscriptTestSession(t)
	session.Agent.Messages = append(session.Agent.Messages, &schemas.Message{
		Role:     "tool",
		ToolName: "files.Read",
		Content:  "first line\nsecond line",
	})

	transcript := NewTranscript(session)
	buffer := utils_tty.NewBuffer(60, 20)
	rect := utils_tty.Rect{X: 0, Y: 0, Width: 60, Height: 20}
	transcript.Render(buffer, rect, 0)

	toggleTranscriptLine(t, transcript, buffer, rect, "tool")

	expanded := false

	for _, line := range transcript.cached {

		if line.Toggle != nil && line.Toggle.Part == "tool" {

			if line.HasBorder == false || line.Border != styles.BorderTool {
				t.Errorf("Expected the tool border style")
			}

			if strings.Contains(transcriptLineText(line), "▾ Tool · files.Read") == true {
				expanded = true
			}

		}

	}

	if expanded == false {
		t.Errorf("Expected the tool message to expand")
	}

}

func TestTranscriptRendersThinkingSpinner(t *testing.T) {

	session := newTranscriptTestSession(t)
	session.Waiting = true

	transcript := NewTranscript(session)
	buffer := utils_tty.NewBuffer(60, 10)
	transcript.Render(buffer, utils_tty.Rect{X: 0, Y: 0, Width: 60, Height: 10}, 0)

	found := false

	for _, line := range transcript.cached {

		if line.Spinner == true {

			found = true

			if line.HasBorder == false || line.Border != styles.BorderThinking {
				t.Errorf("Expected the thinking border style")
			}

		}

	}

	if found == false {
		t.Errorf("Expected a thinking spinner line")
	}

}

func TestTranscriptRendersUserBorder(t *testing.T) {

	session := newTranscriptTestSession(t)
	session.Agent.Messages = append(session.Agent.Messages, &schemas.Message{
		Role:    "user",
		Content: "hello",
	})

	transcript := NewTranscript(session)
	buffer := utils_tty.NewBuffer(40, 10)
	rect := utils_tty.Rect{X: 0, Y: 0, Width: 40, Height: 10}
	transcript.Render(buffer, rect, 0)

	index := -1

	for line_index, line := range transcript.cached {

		if strings.Contains(transcriptLineText(line), "[user]") {
			index = line_index
			break
		}

	}

	if index == -1 {
		t.Fatalf("Expected a user line")
	}

	if transcript.cached[index].HasBorder == false || transcript.cached[index].Border != styles.BorderUser {
		t.Errorf("Expected the user border style")
	}

	y := rect.Y + index - transcript.lastStart

	if buffer.Get(rect.X, y).Text != "┃" {
		t.Errorf("Expected the left border glyph, got %q", buffer.Get(rect.X, y).Text)
	}

	if buffer.Get(rect.X+2, y).Text != "[" {
		t.Errorf("Expected the content to start after the border, got %q", buffer.Get(rect.X+2, y).Text)
	}

}
