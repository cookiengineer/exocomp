package rand

import "math/rand"
import "strings"

const vowels = "aeiou"

func mutateStem(stem string, randomizer *rand.Rand) string {

	switch randomizer.Intn(6) {

	case 0:

		if len(stem) > 0 {

			last := rune(stem[len(stem)-1])

			if strings.ContainsRune(vowels, last) == false {
				stem += "a"
			}

		}

	case 1:

		if len(stem) > 0 {

			last := rune(stem[len(stem)-1])

			if strings.ContainsRune(vowels, last) == false {
				stem += "e"
			}

		}

	case 2:

		if len(stem) > 5 {
			stem = stem[:len(stem)-1]
		}

	case 3:

		stem = strings.ReplaceAll(stem, "x", "ex")

	case 4:

		if strings.HasSuffix(stem, "r") {
			stem += "i"
		}

	}

	return stem

}
