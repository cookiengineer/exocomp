package markdown

import net_url "net/url"
import "strings"

func resolveURL(base_url *net_url.URL, ref string) string {

	var result string

	ref_url, err := net_url.Parse(ref)

	if err == nil {

		resolved := base_url.ResolveReference(ref_url).String()

		if strings.HasPrefix(resolved, base_url.Scheme+"://"+base_url.Host) {
			resolved = strings.TrimPrefix(resolved, base_url.Scheme+"://"+base_url.Host)
		}

		return resolved

	}

	return result

}
