package main

import utils_build "exocomp-toolchain/internal/utils/build"
import "fmt"
import "os"
import "path/filepath"

func BuildPrograms(base_dir string, operating_system string, architecture string) error {

	programs_dir := filepath.Join(base_dir, "installer", "assets", "programs")
	vendor_dir := filepath.Join(base_dir, "vendor")
	cache_dir := filepath.Join(base_dir, "build", "cache", operating_system, architecture)

	err01 := os.MkdirAll(programs_dir, 0755)
	err02 := os.MkdirAll(vendor_dir, 0755)
	err03 := os.MkdirAll(cache_dir, 0755)

	if err01 == nil && err02 == nil && err03 == nil {

		programs, err1 := ReadPrograms(base_dir)

		if err1 == nil {

			for _, program := range programs {

				binary_name := program.Name

				if operating_system == "windows" {
					binary_name = binary_name + ".exe"
				}

				folder := filepath.Join(vendor_dir, program.Folder)
				cache_output := filepath.Join(cache_dir, binary_name)
				asset_output := filepath.Join(programs_dir, binary_name)
				stamp_path := cache_output + ".version"
				cache_path, _ := filepath.Rel(base_dir, cache_output)
				path, _ := filepath.Rel(base_dir, asset_output)

				if isProgramCached(cache_output, program.Version) == true {

					err2 := copyProgram(cache_output, asset_output)

					if err2 == nil {

						fmt.Fprintf(os.Stdout, "-> Cached %s\n", cache_path)
						fmt.Fprintf(os.Stdout, "-> Copied %s\n", path)

					} else {
						return err2
					}

				} else {

					err2 := utils_build.Build(folder, program.Source, cache_output, []string{}, operating_system, architecture)

					if err2 == nil {

						err3 := os.WriteFile(stamp_path, []byte(program.Version), 0644)

						if err3 == nil {

							err4 := copyProgram(cache_output, asset_output)

							if err4 == nil {

								fmt.Fprintf(os.Stdout, "-> Built %s\n", cache_path)
								fmt.Fprintf(os.Stdout, "-> Copied %s\n", path)

							} else {
								return err4
							}

						} else {
							return fmt.Errorf("build.BuildPrograms: Cannot write cache stamp for %s.", path)
						}

					} else {
						return fmt.Errorf("build.BuildPrograms: Cannot build %s.", path)
					}

				}

			}

			return nil

		} else {
			return err1
		}

	} else {

		if err01 != nil {
			fmt.Fprintf(os.Stderr, "!! Error: %s\n", err01.Error())
		}

		if err02 != nil {
			fmt.Fprintf(os.Stderr, "!! Error: %s\n", err02.Error())
		}

		if err03 != nil {
			fmt.Fprintf(os.Stderr, "!! Error: %s\n", err03.Error())
		}

		return fmt.Errorf("build.BuildPrograms: Cannot create build directories.")

	}

}
