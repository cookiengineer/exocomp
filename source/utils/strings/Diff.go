package strings

import "strings"

// Diff returns a line-based diff of a and b, prefixing lines only in a with
// "- ", lines only in b with "+ ", and shared lines with "  ". It returns an
// empty string when a and b are equal.
func Diff(a []string, b []string) string {

	n := len(a)
	m := len(b)

	table := make([][]int, n+1)

	for i := 0; i <= n; i++ {
		table[i] = make([]int, m+1)
	}

	for i := n - 1; i >= 0; i-- {

		for j := m - 1; j >= 0; j-- {

			if a[i] == b[j] {
				table[i][j] = table[i+1][j+1] + 1
			} else if table[i+1][j] >= table[i][j+1] {
				table[i][j] = table[i+1][j]
			} else {
				table[i][j] = table[i][j+1]
			}

		}

	}

	result := make([]string, 0)
	i := 0
	j := 0

	for i < n && j < m {

		if a[i] == b[j] {
			result = append(result, "  "+a[i])
			i++
			j++
		} else if table[i+1][j] >= table[i][j+1] {
			result = append(result, "- "+a[i])
			i++
		} else {
			result = append(result, "+ "+b[j])
			j++
		}

	}

	for i < n {
		result = append(result, "- "+a[i])
		i++
	}

	for j < m {
		result = append(result, "+ "+b[j])
		j++
	}

	return strings.Join(result, "\n")

}
