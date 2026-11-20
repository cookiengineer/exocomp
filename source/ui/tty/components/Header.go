package components

import "exocomp/engine"
import "exocomp/ui/tty/styles"
import utils_fmt "exocomp/utils/fmt"
import utils_tty "exocomp/utils/tty"
import "fmt"
import "strings"
import "time"

type Header struct {
	Session  *engine.Session
	Role     string
	Selected string
}

func NewHeader(session *engine.Session) *Header {

	return &Header{
		Session: session,
		Role:    "user",
	}

}

func (header *Header) Width(width int) int {
	return width
}

func (header *Header) Height(width int) int {
	return 4
}

func (header *Header) Render(buffer *utils_tty.Buffer, rect utils_tty.Rect, elapsed time.Duration) {

	if rect.Width >= 0 && rect.Height >= 0 {

		buffer.Fill(rect, ' ', styles.Background)

		if header.Session != nil && header.Session.Agent != nil && header.Session.Config != nil {

			agent := header.Session.GetAgentInfo(header.Selected)

			if agent != nil {

				config := header.Session.Config
				usage := agent.ContextUsage
				status := agent.Status
				is_planner := header.Selected == "" || header.Selected == header.Session.Agent.Name

				if is_planner == true {

					status = "idle"

					if header.Session.IsWaiting() == true {
						status = "thinking"
					}

				} else if strings.TrimSpace(status) == "" {
					status = "unknown"
				}

				model := agent.Model

				if strings.TrimSpace(model) == "" {
					model = config.Model
				}

				sandbox := agent.Sandbox

				if strings.TrimSpace(sandbox) == "" {
					sandbox = config.Sandbox
				}

				buffer.WriteString(
					rect.X,
					rect.Y+0,
					utils_tty.PadRight(
						utils_tty.Truncate(fmt.Sprintf(
							"Agent:   %s (as %s) is %s at %s (%s of %s)",
							agent.Name,
							agent.Role,
							status,
							strings.TrimSpace(utils_fmt.FormatPercentage(float64(usage.Tokens), float64(usage.Length))),
							utils_fmt.FormatTokens(usage.Tokens),
							utils_fmt.FormatTokens(usage.Length),
						), rect.Width),
						rect.Width,
					),
					styles.Text,
				)

				buffer.WriteString(
					rect.X,
					rect.Y+1,
					utils_tty.PadRight(
						utils_tty.Truncate(fmt.Sprintf(
							"Model:   %s with temperature %.2f",
							model,
							agent.Temperature,
						), rect.Width),
						rect.Width,
					),
					styles.Text,
				)

				buffer.WriteString(
					rect.X,
					rect.Y+2,
					utils_tty.PadRight(
						utils_tty.Truncate(fmt.Sprintf(
							"Sandbox: %s",
							sandbox,
						), rect.Width),
						rect.Width,
					),
					styles.Text,
				)

			}

			utils_tty.DrawRule(
				buffer,
				rect.X,
				rect.Y+rect.Height-1,
				rect.Width,
				styles.Border,
			)

		}

	}

}
