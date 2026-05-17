package actions

import "fmt"
import "os"
import "path/filepath"

func ClonePrograms(base_dir string) {

	vendor_dir := filepath.Join(base_dir, "vendor")

	err0 := os.MkdirAll(vendor_dir, 0755)

	if err0 == nil {

		installs := []struct {
			name   string
			folder string
			url    string
		}{
			{
				name:   "amass",
				folder: filepath.Join(vendor_dir, "amass"),
				url:    "https://github.com/owasp-amass/amass.git",
			},
			{
				name:   "asnmap",
				folder: filepath.Join(vendor_dir, "asnmap"),
				url:    "https://github.com/projectdiscovery/asnmap.git",
			},
			{
				name:   "httpx",
				folder: filepath.Join(vendor_dir, "httpx"),
				url:    "https://github.com/projectdiscovery/httpx.git",
			},
			{
				name:   "katana",
				folder: filepath.Join(vendor_dir, "katana"),
				url:    "https://github.com/projectdiscovery/katana.git",
			},
			{
				name:   "naabu",
				folder: filepath.Join(vendor_dir, "naabu"),
				url:    "https://github.com/projectdiscovery/naabu.git",
			},
			{
				name:   "nuclei",
				folder: filepath.Join(vendor_dir, "nuclei"),
				url:    "https://github.com/projectdiscovery/nuclei.git",
			},
			{
				name:   "subfinder",
				folder: filepath.Join(vendor_dir, "subfinder"),
				url:    "https://github.com/projectdiscovery/subfinder.git",
			},
		}

		for _, install := range installs {

			err1    := CloneRepository(vendor_dir, install.url, install.folder)
			path, _ := filepath.Rel(base_dir, install.folder)

			if err1 == nil {
				fmt.Fprintf(os.Stdout, "-> Cloned %s\n", path)
			} else {
				fmt.Fprintf(os.Stderr, "!! Error cloning %s: %s\n", path, err1.Error())
			}

		}

	} else {
		fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
	}

}
