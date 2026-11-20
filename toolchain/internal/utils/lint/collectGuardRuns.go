package lint

func collectGuardRuns(file *File) []guardRun {

	result := make([]guardRun, 0)

	for _, list := range collectStatementLists(file) {

		index := 0

		if list.Depth >= 3 {
			continue
		}

		for index < len(list.Stmts) {

			if isValueGuard(list.Stmts[index]) == false || isReceiverNilGuard(list.Stmts[index], list.Receiver) == true {

				index++
				continue

			}

			end := index

			for end < len(list.Stmts) && isValueGuard(list.Stmts[end]) == true {
				end++
			}

			result = append(result, guardRun{
				Stmts:   list.Stmts,
				Start:   index,
				End:     end,
				HasTail: end < len(list.Stmts),
				Depth:   list.Depth,
			})

			index = end

		}

	}

	return result

}
