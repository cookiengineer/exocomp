package utils

import "io"
import "os"

func CopyFile(source string, target string) error {

	src_file, err1 := os.Open(source)

	if err1 == nil {

		defer src_file.Close()

		dst_file, err2 := os.Create(target)

		if err2 == nil {

			defer dst_file.Close()

			_, err3 := io.Copy(dst_file, src_file)

			if err3 == nil {

				return dst_file.Sync()

			} else {
				return err3
			}

		} else {
			return err2
		}

	} else {
		return err1
	}

}
