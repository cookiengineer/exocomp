package actions

import "installer/assets"
import "fmt"
import "io/fs"
import "os"
import "path/filepath"

func InstallExocomp(prefix string) error {

	err0 := os.MkdirAll(filepath.Join(prefix, "bin"), 0755)

	if err0 == nil {

		stat1, err10 := fs.Stat(assets.FS, "exocomps/exocomp")

		if err10 == nil && stat1.IsDir() == false {

			buffer1, err11 := fs.ReadFile(assets.FS, "exocomps/exocomp")

			if err11 == nil {

				path := filepath.Join(prefix, "bin", "exocomp")
				err12 := os.WriteFile(path, buffer1, 0755)

				if err12 == nil {
					fmt.Fprintf(os.Stdout, "-> Installed %s\n", path)
				} else {
					return err12
				}

			} else {
				return err11
			}

		} else {
			return err10
		}

		stat2, err20 := fs.Stat(assets.FS, "exocomps/exocomp-agent")

		if err20 == nil && stat2.IsDir() == false {

			buffer2, err21 := fs.ReadFile(assets.FS, "exocomps/exocomp-agent")

			if err21 == nil {

				path := filepath.Join(prefix, "bin", "exocomp-agent")
				err22 := os.WriteFile(path, buffer2, 0755)

				if err22 == nil {
					fmt.Fprintf(os.Stdout, "-> Installed %s\n", path)
				} else {
					return err22
				}

			} else {
				return err21
			}

		} else {
			return err20
		}

		stat3, err30 := fs.Stat(assets.FS, "exocomps/exocomp-sandbox")

		if err30 == nil && stat3.IsDir() == false {

			buffer3, err31 := fs.ReadFile(assets.FS, "exocomps/exocomp-sandbox")

			if err31 == nil {

				path := filepath.Join(prefix, "bin", "exocomp-sandbox")
				err32 := os.WriteFile(path, buffer3, 0755)

				if err32 == nil {
					fmt.Fprintf(os.Stdout, "-> Installed %s\n", path)
				} else {
					return err32
				}

			} else {
				return err31
			}

		} else {
			return err30
		}

	} else {
		return err0
	}

	return nil

}
