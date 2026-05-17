package actions

import utils_cli "exocomp/utils/cli"
import utils_fmt "exocomp/utils/fmt"
import utils_fs "exocomp/utils/fs"
import "bufio"
import "fmt"
import "os"

var mibi_byte uint64 = 1024 * 1024
var gibi_byte uint64 = 1024 * 1024 * 1024

func Install(prefix string) {

	err := utils_cli.CheckRoot()

	if err == nil {

		available_space := utils_fs.AvailableSpace(prefix)
		install_options := make(map[string]bool, 0)

		fmt.Fprintf(os.Stdout, "Prefix:     %s\n", prefix)
		fmt.Fprintf(os.Stdout, "Disk Space: %s\n", utils_fmt.FormatBytes(available_space))

		if available_space > 1 * gibi_byte {

			install_options["exocomp"] = true
			install_options["programs"] = true

		} else if available_space > 20 * mibi_byte {

			install_options["exocomp"] = true
			install_options["programs"] = false

		} else {

			install_options["exocomp"] = false
			install_options["programs"] = false

		}

		fmt.Fprintf(os.Stdout, "\n")
		fmt.Fprintf(os.Stdout, "Select installation option:\n")
		fmt.Fprintf(os.Stdout, "\n")

		if install_options["exocomp"] == true {
			fmt.Fprintf(os.Stdout, "1) Install exocomp\n")
		} else {
			fmt.Fprintf(os.Stdout, "1) Install exocomp (requires %s more space)\n", utils_fmt.FormatBytes(20 * mibi_byte - available_space))
		}

		if install_options["programs"] == true {
			fmt.Fprintf(os.Stdout, "2) Install exocomp with agent programs\n")
		} else {
			fmt.Fprintf(os.Stdout, "2) Install exocomp with agent programs (requires %s more space)\n", utils_fmt.FormatBytes(1 * gibi_byte - available_space))
		}

		fmt.Fprintf(os.Stdout, "\n")
		fmt.Fprintf(os.Stdout, "Enter choice: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		choice := scanner.Text()
		errors := make([]error, 0)

		fmt.Fprintf(os.Stdout, "Selected choice: %s\n", choice)
		fmt.Fprintf(os.Stdout, "\n")

		if choice == "1" {

			err1 := InstallExocomp(prefix)

			if err1 != nil {
				errors = append(errors, err1)
			}

		} else if choice == "2" {

			err1 := InstallExocomp(prefix)

			if err1 != nil {
				errors = append(errors, err1)
			}

			err2 := InstallPrograms(prefix)

			if err2 != nil {
				errors = append(errors, err2)
			}

		}

		if len(errors) > 0 {

			fmt.Fprintf(os.Stdout, "\n")

			for _, err := range errors {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}

			os.Exit(1)

		} else {

			fmt.Fprintf(os.Stdout, "\n")
			fmt.Fprintf(os.Stdout, "%s\n", "Done.")

			os.Exit(0)

		}

	} else {

		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)

	}

}
