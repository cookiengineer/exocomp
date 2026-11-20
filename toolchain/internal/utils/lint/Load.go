package lint

import "go/ast"
import "go/token"
import "io/fs"
import "path/filepath"
import "sort"
import "strings"

var skippedFolders = map[string]bool{
	".git":         true,
	"node_modules": true,
	"testdata":     true,
	"vendor":       true,
}

// Load parses every Go file below root.
func Load(root string, is_including_tests bool) (*Context, error) {

	absolute, err0 := filepath.Abs(root)

	if err0 == nil {

		context := &Context{
			Root:   absolute,
			Fset:   token.NewFileSet(),
			Module: readModuleName(absolute),
		}

		paths := make([]string, 0)

		err1 := filepath.WalkDir(absolute, func(path string, entry fs.DirEntry, err error) error {

			if err == nil {

				if entry.IsDir() == true {

					if path != absolute && skippedFolders[entry.Name()] == true {
						return fs.SkipDir
					} else {
						return nil
					}

				} else if strings.HasSuffix(path, ".go") == false {
					return nil
				} else if is_including_tests == false && strings.HasSuffix(path, "_test.go") == true {
					return nil
				} else {

					paths = append(paths, path)

					return nil

				}

			} else {
				return err
			}

		})

		if err1 == nil {

			sort.Strings(paths)

			by_dir := make(map[string]*Package)

			for _, path := range paths {

				file, err2 := loadFile(context, absolute, path)

				if err2 != nil {

					context.ParseErrors = append(context.ParseErrors, err2.Error())
					continue

				}

				if file == nil {
					continue
				}

				directory := filepath.Dir(path)

				if by_dir[directory] == nil {

					relative, _ := filepath.Rel(absolute, directory)
					by_dir[directory] = &Package{
						Dir:   directory,
						Rel:   filepath.ToSlash(relative),
						Name:  file.AST.Name.Name,
						Types: make(map[string]*File),
					}

				}

				file.Pkg = by_dir[directory]
				by_dir[directory].Files = append(by_dir[directory].Files, file)
				context.Files = append(context.Files, file)

			}

			directories := make([]string, 0, len(by_dir))

			for directory := range by_dir {
				directories = append(directories, directory)
			}

			sort.Strings(directories)

			for _, directory := range directories {

				pkg := by_dir[directory]

				for _, file := range pkg.Files {

					for _, decl := range file.AST.Decls {

						gen_decl, ok := decl.(*ast.GenDecl)

						if ok == true && gen_decl.Tok == token.TYPE {

							for _, spec := range gen_decl.Specs {

								type_spec, ok := spec.(*ast.TypeSpec)

								if ok == true {
									pkg.Types[type_spec.Name.Name] = file
								}

							}

						}

					}

				}

				context.Packages = append(context.Packages, pkg)

			}

			return context, nil

		} else {
			return nil, err1
		}

	} else {
		return nil, err0
	}

}
