package lint

func checkGuard(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		for _, run := range collectGuardRuns(file) {

			if run.HasTail == false || run.End-run.Start != 1 {
				continue
			}

			result = append(
				result,
				file.Diagnostic(
					"branch/guard",
					run.Stmts[run.Start].Pos(),
					"guard clause has %d statements after it; invert it into a positive branch",
					len(run.Stmts)-1-run.Start,
				),
			)

		}

		return result

	}

}
