package agents

import "exocomp/types"

func NewAgent(config *types.Config) *types.Agent {

	switch config.Role {
	case "planner":
		return NewPlanner(config)

	// Development
	case "architect":
		return NewArchitect(config)
	case "coder":
		return NewCoder(config)
	case "researcher":
		// TODO: Implement this
		// return NewResearcher(config)
		return nil
	case "summarizer":
		return NewSummarizer(config)
	case "tester":
		return NewTester(config)
	
	// Pentesting
	case "exploiter":
		return NewExploiter(config)
	case "reverser":
		// TODO: Implement this
		// return NewReverser(config)
		return nil
	case "threathunter":
		// TODO: Implement this
		// return NewThreatHunter(config)
		return nil
	case "webscanner":
		return NewWebScanner(config)
	
	// Default
	default:
		return NewPlanner(config)
	}

}

