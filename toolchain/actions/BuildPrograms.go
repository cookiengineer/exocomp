package actions

import "fmt"
import "os"
import "path/filepath"

func BuildPrograms(base_dir string, operating_system string) {

	programs_dir := filepath.Join(base_dir, "installer", "assets", "programs")
	vendor_dir   := filepath.Join(base_dir, "vendor")

	err01 := os.MkdirAll(programs_dir, 0755)
	err02 := os.MkdirAll(vendor_dir, 0755)

	if err01 == nil && err02 == nil {

		installs := []struct {
			name   string
			output string
			source string
			url    string
		}{
			{
				name:   "amass",
				output: filepath.Join(programs_dir, "amass"),
				source: "./cmd/amass",
				url:    "https://github.com/owasp-amass/amass.git",
			},
			{
				name:   "asnmap",
				output: filepath.Join(programs_dir, "asnmap"),
				source: "./cmd/asnmap",
				url:    "https://github.com/projectdiscovery/asnmap.git",
			},
			{
				name:   "httpx",
				output: filepath.Join(programs_dir, "httpx"),
				source: "./cmd/httpx",
				url:    "https://github.com/projectdiscovery/httpx.git",
			},
			{
				name:   "katana",
				output: filepath.Join(programs_dir, "katana"),
				source: "./cmd/katana",
				url:    "https://github.com/projectdiscovery/katana.git",
			},
			{
				name:   "naabu",
				output: filepath.Join(programs_dir, "naabu"),
				source: "./cmd/naabu",
				url:    "https://github.com/projectdiscovery/naabu.git",
			},
			{
				name:   "nuclei",
				output: filepath.Join(programs_dir, "nuclei"),
				source: "./cmd/nuclei",
				url:    "https://github.com/projectdiscovery/nuclei.git",
			},
			{
				name:   "subfinder",
				output: filepath.Join(programs_dir, "subfinder"),
				source: "./cmd/subfinder",
				url:    "https://github.com/projectdiscovery/subfinder.git",
			},
		}

		for _, install := range installs {

			repo_dir := filepath.Join(vendor_dir, install.name)
			err1     := CloneRepository(vendor_dir, install.url, repo_dir)

			if err1 == nil {

				fmt.Fprintf(os.Stdout, "-> Cloned %s\n", repo_dir)

				err2    := BuildBinary(repo_dir, install.source, install.output, []string{}, operating_system)
				path, _ := filepath.Rel(base_dir, install.output)

				if err2 == nil {
					fmt.Fprintf(os.Stdout, "-> Built %s\n", path)
				} else {
					fmt.Fprintf(os.Stderr, "!! Error building %s: %s\n", path, err2.Error())
				}

			} else {
				fmt.Fprintf(os.Stderr, "!! Error cloning %s: %s\n", repo_dir, err1.Error())
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
