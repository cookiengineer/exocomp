package lint

import "sort"

func sortDiagnostics(diagnostics []Diagnostic) {

	sort.SliceStable(diagnostics, func(a int, b int) bool {

		left := diagnostics[a]
		right := diagnostics[b]

		if left.Position.Filename != right.Position.Filename {
			return left.Position.Filename < right.Position.Filename
		} else if left.Position.Line != right.Position.Line {
			return left.Position.Line < right.Position.Line
		} else if left.Position.Column != right.Position.Column {
			return left.Position.Column < right.Position.Column
		} else {
			return left.Rule < right.Rule
		}

	})

}
