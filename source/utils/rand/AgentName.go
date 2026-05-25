package rand

import "fmt"
import "math/rand"
import "time"
import "unicode"

func AgentName(role string) string {

	randomizer     := rand.New(rand.NewSource(time.Now().UnixNano()))
	surname_stem   := createStem(role)
	surname_stem    = mutateStem(surname_stem, randomizer)
	surname_suffix := surname_suffixes[randomizer.Intn(len(surname_suffixes))]

	prename := prenames[randomizer.Intn(len(prenames))]
	surname := cleanupName(surname_stem + surname_suffix)

	return fmt.Sprintf(
		"%s %s",
		string(unicode.ToUpper(rune(prename[0]))) + prename[1:],
		string(unicode.ToUpper(rune(surname[0]))) + surname[1:],
	)

}

