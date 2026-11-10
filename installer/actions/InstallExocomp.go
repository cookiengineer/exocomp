package actions

import "installer/assets"
import "fmt"
import "io/fs"
import "os"
import "path/filepath"
import "runtime"

func InstallExocomp(prefix string) error {

	err0 := os.MkdirAll(filepath.Join(prefix, "bin"), 0755)

	if err0 != nil {
		return err0
	}

	extension := ""

	if runtime.GOOS == "windows" {
		extension = ".exe"
	}

	programs := []string{
		"exocomp",
		"exocomp-agent",
		"exocomp-sandbox",
		"exocomp-web",
	}

	for _, program := range programs {

		name := program + extension
		asset := fmt.Sprintf("exocomps/%s", name)

		stat, err1 := fs.Stat(assets.FS, asset)

		if err1 != nil {
			return err1
		} else if stat.IsDir() == true {
			return fmt.Errorf("InstallExocomp: %s is a directory.", asset)
		}

		buffer, err2 := fs.ReadFile(assets.FS, asset)

		if err2 != nil {
			return err2
		}

		path := filepath.Join(prefix, "bin", name)
		err3 := os.WriteFile(path, buffer, 0755)

		if err3 != nil {
			return err3
		}

		fmt.Fprintf(os.Stdout, "-> Installed %s\n", path)

	}

	return nil

}
