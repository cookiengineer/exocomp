package main

import "exocomp-toolchain/internal/utils/build"
import "exocomp-toolchain/internal/utils/root"
import "exocomp-toolchain/internal/utils/webapp"
import "flag"
import "fmt"
import "os"
import "runtime"
import "strings"

func main() {

	architecture := flag.String("arch", runtime.GOARCH, "target architecture")
	quick := flag.Bool("quick", false, "build only for the current operating system")
	no_agent_programs := flag.Bool("no-agent-programs", false, "skip agent program binaries")
	no_installer := flag.Bool("no-installer", false, "skip the installer binary")

	flag.Parse()

	target := strings.TrimSpace(flag.Arg(0))
	operating_systems := []string{"linux", "darwin", "windows"}

	if target != "" {

		if target != "linux" && target != "darwin" && target != "windows" {

			fmt.Fprintf(os.Stderr, "!! Error: Unknown operating system \"%s\"\n", target)
			os.Exit(2)

		}

		operating_systems = []string{target}

	} else if *quick == true {

		switch runtime.GOOS {

		case "linux", "darwin", "windows":
			operating_systems = []string{runtime.GOOS}

		}

	}

	is_building_programs := *no_agent_programs == false && *quick == false
	is_building_installer := *no_installer == false && *quick == false

	base_dir, err0 := root.GetRoot()

	if err0 != nil {

		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
		os.Exit(1)

	}

	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprint(os.Stdout, "== Build Web App ==\n")

	webapp.BuildWebApp(base_dir)

	err1 := build.BuildTargets(base_dir, operating_systems, *architecture, is_building_programs, is_building_installer)

	if err1 != nil {

		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err1.Error())
		os.Exit(1)

	}

}
