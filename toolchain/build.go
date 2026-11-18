package main

import "exocomp-toolchain/actions"
import "exocomp-toolchain/utils"
import "fmt"
import "os"
import "runtime"
import "strings"

func main() {

	base_dir, err0 := utils.GetRoot()

	architecture := runtime.GOARCH
	build_linux := true
	build_darwin := true
	build_windows := true
	build_agent_programs := true
	build_installer := true

	if len(os.Args) > 1 {

		tmp1 := strings.TrimSpace(os.Args[1])

		switch tmp1 {
		case "linux":
			build_darwin = false
			build_windows = false
		case "darwin":
			build_linux = false
			build_windows = false
		case "windows":
			build_linux = false
			build_darwin = false
		case "--quick":

			build_agent_programs = false
			build_installer = false

			if runtime.GOOS == "linux" {
				build_darwin = false
				build_windows = false
			} else if runtime.GOOS == "darwin" {
				build_linux = false
				build_windows = false
			} else if runtime.GOOS == "windows" {
				build_linux = false
				build_darwin = false
			}

		}

		if len(os.Args) >= 2 {

			tmp2 := os.Args[1:]

			for _, tmp := range tmp2 {

				flag := strings.TrimSpace(strings.ToLower(tmp))

				if strings.HasPrefix(flag, "--arch=") {

					tmp3 := strings.TrimSpace(flag[len("--arch="):])

					if tmp3 != "" {
						architecture = tmp3
					}

				} else {

					switch flag {
					case "--no-agent-programs":
						build_agent_programs = false
						break
					case "--no-installer":
						build_installer = false
						break
					case "--quick":
						build_agent_programs = false
						build_installer = false
						break
					}

				}

			}

		}

	}

	if err0 == nil {

		fmt.Fprint(os.Stdout, "\n")
		fmt.Fprint(os.Stdout, "== Build Web App ==\n")

		actions.BuildWebApp(base_dir)

		if build_agent_programs == true {

			fmt.Fprint(os.Stdout, "\n")
			fmt.Fprint(os.Stdout, "== Clone programs ==\n")

			err := actions.ClonePrograms(base_dir)

			if err != nil {
				fmt.Fprintf(os.Stderr, "!! Error: %s\n", err.Error())
				os.Exit(1)
			}

		}

		if build_linux == true {

			fmt.Fprint(os.Stdout, "\n")
			fmt.Fprint(os.Stdout, "== Build linux ==\n")

			actions.ClearAssets(base_dir)
			actions.BuildExocomps(base_dir, "linux", architecture)

			if build_agent_programs == true {

				err := actions.BuildPrograms(base_dir, "linux", architecture)

				if err != nil {
					fmt.Fprintf(os.Stderr, "!! Error: %s\n", err.Error())
					os.Exit(1)
				}

			}

			if build_installer == true {
				actions.BuildInstaller(base_dir, "linux", architecture)
			}

		}

		if build_darwin == true {

			fmt.Fprint(os.Stdout, "\n")
			fmt.Fprint(os.Stdout, "== Build darwin ==\n")

			actions.ClearAssets(base_dir)
			actions.BuildExocomps(base_dir, "darwin", architecture)

			if build_agent_programs == true {

				err := actions.BuildPrograms(base_dir, "darwin", architecture)

				if err != nil {
					fmt.Fprintf(os.Stderr, "!! Error: %s\n", err.Error())
					os.Exit(1)
				}

			}

			if build_installer == true {
				actions.BuildInstaller(base_dir, "darwin", architecture)
			}

		}

		if build_windows == true {

			fmt.Fprint(os.Stdout, "\n")
			fmt.Fprint(os.Stdout, "== Build windows ==\n")

			actions.ClearAssets(base_dir)
			actions.BuildExocomps(base_dir, "windows", architecture)

			if build_agent_programs == true {

				err := actions.BuildPrograms(base_dir, "windows", architecture)

				if err != nil {
					fmt.Fprintf(os.Stderr, "!! Error: %s\n", err.Error())
					os.Exit(1)
				}

			}

			if build_installer == true {
				actions.BuildInstaller(base_dir, "windows", architecture)
			}

		}

	} else {
		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
	}

}
