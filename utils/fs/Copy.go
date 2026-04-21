package fs

import "io"
import "os"

func Copy(source string, target string) error {

	stat, err := os.Stat(source)

	if err == nil {

		if stat.IsDir() {

			return CopyAll(source, target)

		} else {

			src, err1 := os.Open(source)

			if err1 == nil {

				defer src.Close()

				dst, err4 := os.Create(target)

				if err4 == nil {

					defer dst.Close()

					_, err5 := io.Copy(dst, src)

					if err5 == nil {
						return nil
					} else {
						return err5
					}

				} else {
					return err4
				}

			} else {
				return err1
			}

		}

	} else {
		return err
	}

}

