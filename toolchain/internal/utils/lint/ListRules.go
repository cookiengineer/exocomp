package lint

func ListRules() []RuleInfo {

	result := make([]RuleInfo, 0, len(fileRules)+len(packageRules))

	for _, rule := range fileRules {
		result = append(result, RuleInfo{ID: rule.ID, Summary: rule.Summary})
	}

	for _, rule := range packageRules {
		result = append(result, RuleInfo{ID: rule.ID, Summary: rule.Summary})
	}

	return result

}
