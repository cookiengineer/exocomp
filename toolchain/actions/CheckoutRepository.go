package actions

import "fmt"

func CheckoutRepository(folder string, version string) error {

	target, err0 := gitOutput(folder, "rev-parse", "--verify", version+"^{commit}")

	if err0 != nil {

		err1 := runGit(folder, "fetch", "--depth", "1", "origin", "tag", version)

		if err1 == nil {
			target, err0 = gitOutput(folder, "rev-parse", "--verify", version+"^{commit}")
		} else {
			return fmt.Errorf("actions.CheckoutRepository: Cannot fetch version \"%s\".", version)
		}

	}

	if err0 == nil {

		current, err2 := gitOutput(folder, "rev-parse", "HEAD")

		if err2 == nil {

			if current != target {

				err3 := runGit(folder, "checkout", "--detach", version)

				if err3 == nil {
					current, err2 = gitOutput(folder, "rev-parse", "HEAD")
				} else {
					return fmt.Errorf("actions.CheckoutRepository: Cannot checkout version \"%s\".", version)
				}

			}

			if err2 == nil {

				if current == target {
					return nil
				} else {
					return fmt.Errorf("actions.CheckoutRepository: Version mismatch for \"%s\".", version)
				}

			} else {
				return fmt.Errorf("actions.CheckoutRepository: Cannot resolve current revision.")
			}

		} else {
			return fmt.Errorf("actions.CheckoutRepository: Cannot resolve current revision.")
		}

	} else {
		return fmt.Errorf("actions.CheckoutRepository: Cannot resolve version \"%s\".", version)
	}

}
