package main

import utils_assets "exocomp-toolchain/internal/utils/assets"
import utils_root "exocomp-toolchain/internal/utils/root"
import "exocomp-toolchain/internal/utils/webapp"
import "fmt"
import "os"
import "runtime"
import "strings"

func buildPlatform(base_dir string, operating_system string, architecture string, is_building_programs bool, is_building_installer bool) {

	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, "== Build %s ==\n", operating_system)

	utils_assets.ClearAssets(base_dir)
	BuildExocomps(base_dir, operating_system, architecture)

	if is_building_programs == true {

		err := BuildPrograms(base_dir, operating_system, architecture)

		if err != nil {

			fmt.Fprintf(os.Stderr, "!! Error: %s\n", err.Error())
			os.Exit(1)

		}

	}

	if is_building_installer == true {
		BuildInstaller(base_dir, operating_system, architecture)
	}

}

func main() {

	base_dir, err0 := utils_root.GetRoot()

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

		for _, argument := range os.Args[1:] {

			flag := strings.TrimSpace(strings.ToLower(argument))

			if strings.HasPrefix(flag, "--arch=") == true {

				value := strings.TrimSpace(flag[len("--arch="):])

				if value != "" {
					architecture = value
				}

			} else {

				switch flag {

				case "--no-agent-programs":
					build_agent_programs = false

				case "--no-installer":
					build_installer = false

				case "--quick":
					build_agent_programs = false
					build_installer = false

				}

			}

		}

	}

	if err0 == nil {

		fmt.Fprint(os.Stdout, "\n")
		fmt.Fprint(os.Stdout, "== Build Web App ==\n")

		webapp.BuildWebApp(base_dir)

		if build_agent_programs == true {

			fmt.Fprint(os.Stdout, "\n")
			fmt.Fprint(os.Stdout, "== Clone programs ==\n")

			err := ClonePrograms(base_dir)

			if err != nil {

				fmt.Fprintf(os.Stderr, "!! Error: %s\n", err.Error())
				os.Exit(1)

			}

		}

		if build_linux == true {
			buildPlatform(base_dir, "linux", architecture, build_agent_programs, build_installer)
		}

		if build_darwin == true {
			buildPlatform(base_dir, "darwin", architecture, build_agent_programs, build_installer)
		}

		if build_windows == true {
			buildPlatform(base_dir, "windows", architecture, build_agent_programs, build_installer)
		}

	} else {
		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
	}

}
