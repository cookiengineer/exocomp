package main

import "go/ast"
import "go/token"
import "regexp"
import "strings"

var nolintPattern = regexp.MustCompile(`^//\s*nolint:\s*([a-z0-9/,._-]+)`)

func parseSuppressions(fileset *token.FileSet, file *ast.File) map[int]map[string]bool {

	result := make(map[int]map[string]bool)

	for _, group := range file.Comments {

		for _, comment := range group.List {

			matches := nolintPattern.FindStringSubmatch(comment.Text)

			if matches == nil {
				continue
			}

			line := fileset.Position(comment.Pos()).Line

			for _, rule := range strings.Split(matches[1], ",") {

				trimmed := strings.TrimSpace(rule)

				if trimmed == "" {
					continue
				}

				addSuppression(result, line, trimmed)
				addSuppression(result, line+1, trimmed)

			}

		}

	}

	return result

}

func addSuppression(target map[int]map[string]bool, line int, rule string) {

	if target[line] == nil {
		target[line] = make(map[string]bool)
	}

	target[line][rule] = true

}

func (file *File) IsSuppressed(line int, rule string) bool {

	rules := file.Suppress[line]

	if rules == nil {
		return false
	}

	if rules["all"] == true {
		return true
	}

	return rules[rule]

}
