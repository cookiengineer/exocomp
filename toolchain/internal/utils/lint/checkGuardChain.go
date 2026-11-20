package lint

func checkGuardChain(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		for _, run := range collectGuardRuns(file) {

			if run.HasTail == false || run.End-run.Start < 2 {
				continue
			}

			result = append(
				result,
				file.Diagnostic(
					"branch/guard-chain",
					run.Stmts[run.Start].Pos(),
					"%d consecutive guards; use one } else if { chain",
					run.End-run.Start,
				),
			)

		}

		return result

	}

}
