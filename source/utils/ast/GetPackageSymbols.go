package ast

import "go/parser"
import "go/token"
import "io/fs"
import "path/filepath"
import "strings"

func GetPackageSymbols(folder string) map[string]map[string]string {

	result := make(map[string]map[string]string, 0)
	fileset := token.NewFileSet()

	filepath.WalkDir(folder, func(path string, entry fs.DirEntry, err error) error {

		if err == nil && entry.IsDir() == false && strings.HasSuffix(path, ".go") {

			file, parse_err := parser.ParseFile(fileset, path, nil, 0)

			if parse_err == nil {

				relative, rel_err := filepath.Rel(folder, path)

				if rel_err == nil {
					result[relative] = getSymbols(file, fileset)
				}

			}

		}

		return nil

	})

	return result

}
