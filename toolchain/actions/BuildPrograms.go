package actions

import "exocomp-toolchain/utils"
import "fmt"
import "os"
import "path/filepath"

func BuildPrograms(base_dir string, operating_system string) {

	programs_dir := filepath.Join(base_dir, "installer", "assets", "programs")
	vendor_dir   := filepath.Join(base_dir, "vendor")

	err01 := os.MkdirAll(programs_dir, 0755)
	err02 := os.MkdirAll(vendor_dir, 0755)

	if err01 == nil && err02 == nil {

		builds := []struct {
			name   string
			folder string
			output string
			source string
		}{
			{
				name:   "amass",
				folder: filepath.Join(vendor_dir, "amass"),
				output: filepath.Join(programs_dir, "amass"),
				source: "./cmd/amass",
			},
			{
				name:   "asnmap",
				folder: filepath.Join(vendor_dir, "asnmap"),
				output: filepath.Join(programs_dir, "asnmap"),
				source: "./cmd/asnmap",
			},
			{
				name:   "httpx",
				folder: filepath.Join(vendor_dir, "httpx"),
				output: filepath.Join(programs_dir, "httpx"),
				source: "./cmd/httpx",
			},
			{
				name:   "katana",
				folder: filepath.Join(vendor_dir, "katana"),
				output: filepath.Join(programs_dir, "katana"),
				source: "./cmd/katana",
			},
			{
				name:   "naabu",
				folder: filepath.Join(vendor_dir, "naabu"),
				output: filepath.Join(programs_dir, "naabu"),
				source: "./cmd/naabu",
			},
			{
				name:   "nuclei",
				folder: filepath.Join(vendor_dir, "nuclei"),
				output: filepath.Join(programs_dir, "nuclei"),
				source: "./cmd/nuclei",
			},
			{
				name:   "subfinder",
				folder: filepath.Join(vendor_dir, "subfinder"),
				output: filepath.Join(programs_dir, "subfinder"),
				source: "./cmd/subfinder",
			},
		}

		for _, build := range builds {

			err     := utils.Build(build.folder, build.source, build.output, []string{}, operating_system)
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
