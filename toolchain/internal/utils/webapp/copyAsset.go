package webapp

import utils_files "exocomp-toolchain/internal/utils/files"

func copyAsset(source string, target string, is_directory bool) error {

	if is_directory == true {
		return utils_files.CopyDir(source, target)
	} else {
		return utils_files.CopyFile(source, target)
	}

}
