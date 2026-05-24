package agents

import "exocomp/types"
import "strings"

func NewAgent(config *types.Config) *types.Agent {

	role := config.Role

	if role == "" {
		role = "planner"
	}

	template, ok := Roles[role]

	if ok == true {

		name := strings.TrimSpace(config.Name)

		if name == "" {
			name = template.Name
		}

		model := strings.TrimSpace(config.Model)

		if model == "" {
			model = template.Model
		}

		temperature := config.Temperature

		if temperature == 0.0 {
			temperature = template.Temperature
		}

		prompt := render_prompt(name, role, template.Prompt)

		return &types.Agent{
			Name:            name,
			Role:            template.Role,
			Model:           model,
			Prompt:          prompt,
			Temperature:     temperature,
			Messages:        render_messages(prompt),
			AllowedPrograms: template.AllowedPrograms,
			AllowedTools:    template.AllowedTools,
			Sandbox:         config.Sandbox,
		}

	} else {

		return nil

	}

}

