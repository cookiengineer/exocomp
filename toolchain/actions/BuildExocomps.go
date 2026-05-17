package actions

import "fmt"
import "os"
import "path/filepath"

func BuildExocomps(base_dir string, operating_system string) {

	build_dir    := filepath.Join(base_dir, "build")
	source_dir   := filepath.Join(base_dir, "source")
	exocomps_dir := filepath.Join(base_dir, "installer", "assets", "exocomps")

	err01 := os.MkdirAll(filepath.Join(build_dir, operating_system), 0755)
	err02 := os.MkdirAll(exocomps_dir, 0755)

	if err01 == nil && err02 == nil {

		builds := []struct {
			name   string
			output string
			source string
		}{
			{
				name:   "exocomp",
				output: filepath.Join(exocomps_dir, "exocomp"),
				source: "./cmds/exocomp/main.go",
			},
			{
				name:   "exocomp-agent",
				output: filepath.Join(exocomps_dir, "exocomp-agent"),
				source: "./cmds/exocomp-agent/main.go",
			},
			{
				name:   "agimus",
				output: filepath.Join(exocomps_dir, "agimus"),
				source: "./cmds/agimus/main.go",
			},
			{
				name:   "exocomp",
				output: filepath.Join(build_dir, operating_system, "exocomp"),
				source: "./cmds/exocomp/main.go",
			},
			{
				name:   "exocomp-agent",
				output: filepath.Join(build_dir, operating_system, "exocomp-agent"),
				source: "./cmds/exocomp-agent/main.go",
			},
			{
				name:   "agimus",
				output: filepath.Join(build_dir, operating_system, "agimus"),
				source: "./cmds/agimus/main.go",
			},
		}

		for _, build := range builds {

			err     := BuildBinary(source_dir, build.source, build.output, []string{}, operating_system)
			path, _ := filepath.Rel(base_dir, build.output)

			if err == nil {
				fmt.Fprintf(os.Stdout, "-> Built %s\n", path)
			} else {
				fmt.Fprintf(os.Stderr, "!! Error building %s: %s\n", path, err.Error())
			}

		}

	} else {

		if err01 != nil {
			fmt.Fprintf(os.Stderr, "!! Error: %s\n", err01.Error())
		}

		if err02 != nil {
			fmt.Fprintf(os.Stderr, "!! Error: %s\n", err02.Error())
		}

	}

}
