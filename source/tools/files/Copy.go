package files

import utils_fs "exocomp/utils/fs"
import "fmt"
import "os"

func (tool *Files) Copy(from_path string, to_path string) (string, error) {

	from_resolved, err1 := resolveFilesPath(tool.Sandbox, tool.Capabilities, from_path, false)

	if err1 == nil {

		to_resolved, err2 := resolveFilesPath(tool.Sandbox, tool.Capabilities, to_path, true)

		if err2 == nil {

			stat, err3 := os.Lstat(from_resolved)

			if err3 == nil {

				if stat.Mode()&os.ModeSymlink != 0 {

					return "", fmt.Errorf("files.Copy: Invalid path \"%s\": Symbolic links are not allowed", from_path)

				} else if stat.IsDir() {

					err4 := utils_fs.CopyAll(from_resolved, to_resolved)

					if err4 == nil {
						return fmt.Sprintf("files.Copy: Folder \"%s\" copied to \"%s\".", from_path, to_path), nil
					} else {
						return "", fmt.Errorf("files.Copy: Cannot copy \"%s\" to \"%s\".", from_path, to_path)
					}

				} else {

					err4 := utils_fs.Copy(from_resolved, to_resolved)

					if err4 == nil {
						return fmt.Sprintf("files.Copy: File \"%s\" copied to \"%s\".", from_path, to_path), nil
					} else {
						return "", fmt.Errorf("files.Copy: Cannot copy \"%s\" to \"%s\".", from_path, to_path)
					}

				}

			} else {
				return sanitizeFilesystemError("files", "Copy", "Path", from_path, err3)
			}

		} else {
			return "", fmt.Errorf("files.Copy: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("files.Copy: %s", err1.Error())
	}

}
