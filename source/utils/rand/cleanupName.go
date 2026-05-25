package rand

import "strings"

func cleanupName(name string) string {

	replacements := map[string]string{
		"aaa": "aa",
		"eee": "ee",
		"iii": "ii",
		"ooo": "oo",
		"uuu": "uu",
		"tt":  "t",
		"dd":  "d",
		"kk":  "k",
		"vv":  "v",
		"xx":  "x",
		"yy":  "y",
		"zz":  "z",
	}

	for from, to := range replacements {
		name = strings.ReplaceAll(name, from, to)
	}

	return name

}
