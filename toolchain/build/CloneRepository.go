package main

import "fmt"

func CloneRepository(cwd string, url string, folder string, version string) error {

	err := runGit(
		cwd,
		"clone",
		"--single-branch",
		"--branch", version,
		"--depth", "1",
		url,
		folder,
	)

	if err == nil {

		return CheckoutRepository(folder, version)

	} else {
		return fmt.Errorf("build.CloneRepository: Cannot clone \"%s\".", url)
	}

}
