package agents

import "exocomp/types"

func NewAgent(config *types.Config) *types.Agent {

	switch config.Agent {
	case "architect":
		return NewArchitect(config)
	case "coder":
		return NewCoder(config)
	case "planner":
		return NewPlanner(config)
	case "researcher":
		// TODO: Implement researcher
		return nil
	case "summarizer":
		return NewSummarizer(config)
	case "tester":
		return NewTester(config)
	default:
		return NewPlanner(config)
	}

}

