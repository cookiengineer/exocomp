package tools

import "exocomp/utils"
import "fmt"
import "sort"
import "strings"
import "time"

type Changelog struct {
	Sandbox  string
	contents map[time.Time][]string // map[2025-12-31 10:20:30][]string{changelog_description}
}

func NewChangelog(agent string, sandbox string) *Changelog {

	changelog := &Changelog{
		Sandbox:  sandbox,
		contents: make(map[time.Time][]string, 0),
	}

	readChangelog(changelog)

	return changelog

}

func (tool *Changelog) Call(method string, arguments map[string]interface{}) (string, error) {

	if method == "Add" {

		path, ok1 := arguments["path"].(string)
		note, ok2 := arguments["note"].(string)

		if ok1 == true && ok2 == true {
			return tool.Add(utils.FormatFilePath(path), utils.FormatSingleLine(note))
		} else if ok1 == true && ok2 == false {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"note\" is not a string.")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Change" {

		path, ok1 := arguments["path"].(string)
		note, ok2 := arguments["note"].(string)

		if ok1 == true && ok2 == true {
			return tool.Change(utils.FormatFilePath(path), utils.FormatSingleLine(note))
		} else if ok1 == true && ok2 == false {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"note\" is not a string.")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Deprecate" {

		path, ok1 := arguments["path"].(string)
		note, ok2 := arguments["note"].(string)

		if ok1 == true && ok2 == true {
			return tool.Deprecate(utils.FormatFilePath(path), utils.FormatSingleLine(note))
		} else if ok1 == true && ok2 == false {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"note\" is not a string.")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Fix" {

		path, ok1 := arguments["path"].(string)
		note, ok2 := arguments["note"].(string)

		if ok1 == true && ok2 == true {
			return tool.Fix(utils.FormatFilePath(path), utils.FormatSingleLine(note))
		} else if ok1 == true && ok2 == false {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"note\" is not a string.")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Remove" {

		path, ok1 := arguments["path"].(string)
		note, ok2 := arguments["note"].(string)

		if ok1 == true && ok2 == true {
			return tool.Remove(utils.FormatFilePath(path), utils.FormatSingleLine(note))
		} else if ok1 == true && ok2 == false {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"note\" is not a string.")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
		}


	} else if method == "Search" {

		path, ok := arguments["path"].(string)

		if ok == true {
			return tool.Search(utils.FormatFilePath(path))
		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string")
		}

	} else {
		return "", fmt.Errorf("changelog.%s: Invalid method.", method)
	}

}

func (tool *Changelog) Add(path string, note string) (string, error) {

	_, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		if strings.HasPrefix(note, "Added ") {

			now   := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

			_, ok := tool.contents[today]

			if ok == false {
				tool.contents[today] = make([]string, 0)
			}

			message := "`" + path + "`: " + note
			tool.contents[today] = append(tool.contents[today], message)
			writeChangelog(tool)

			result := strings.Join([]string{
				fmt.Sprintf("changelog.Add: Note with %d B written.", len(message)),
			}, "\n")

			return result, nil

		} else {
			return "", fmt.Errorf("changelog.Add: Note must start with \"Added\".")
		}

	} else {
		return "", fmt.Errorf("changelog.Add: %s", err0.Error())
	}

}

func (tool *Changelog) Change(path string, note string) (string, error) {

	_, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		if strings.HasPrefix(note, "Changed ") {

			now   := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

			_, ok := tool.contents[today]

			if ok == false {
				tool.contents[today] = make([]string, 0)
			}

			message := "`" + path + "`: " + note
			tool.contents[today] = append(tool.contents[today], message)
			writeChangelog(tool)

			result := strings.Join([]string{
				fmt.Sprintf("changelog.Change: Note with %d B written.", len(message)),
			}, "\n")

			return result, nil

		} else {
			return "", fmt.Errorf("changelog.Change: Note must start with \"Changed\".")
		}

	} else {
		return "", fmt.Errorf("changelog.Change: %s", err0.Error())
	}

}

func (tool *Changelog) Deprecate(path string, note string) (string, error) {

	_, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		if strings.HasPrefix(note, "Deprecated ") {

			now   := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

			_, ok := tool.contents[today]

			if ok == false {
				tool.contents[today] = make([]string, 0)
			}

			message := "`" + path + "`: " + note
			tool.contents[today] = append(tool.contents[today], message)
			writeChangelog(tool)

			result := strings.Join([]string{
				fmt.Sprintf("changelog.Deprecate: Note with %d B written.", len(message)),
			}, "\n")

			return result, nil

		} else {
			return "", fmt.Errorf("changelog.Deprecate: Note must start with \"Deprecated\".")
		}

	} else {
		return "", fmt.Errorf("changelog.Deprecate: %s", err0.Error())
	}

}

func (tool *Changelog) Fix(path string, note string) (string, error) {

	_, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		if strings.HasPrefix(note, "Fixed ") {

			now   := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

			_, ok := tool.contents[today]

			if ok == false {
				tool.contents[today] = make([]string, 0)
			}

			message := "`" + path + "`: " + note
			tool.contents[today] = append(tool.contents[today], message)
			writeChangelog(tool)

			result := strings.Join([]string{
				fmt.Sprintf("changelog.Fix: Note with %d B written.", len(message)),
			}, "\n")

			return result, nil

		} else {
			return "", fmt.Errorf("changelog.Fix: Note must start with \"Fixed\".")
		}

	} else {
		return "", fmt.Errorf("changelog.Fix: %s", err0.Error())
	}

}

func (tool *Changelog) Remove(path string, note string) (string, error) {

	_, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		if strings.HasPrefix(note, "Removed ") {

			now   := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

			_, ok := tool.contents[today]

			if ok == false {
				tool.contents[today] = make([]string, 0)
			}

			message := "`" + path + "`: " + note
			tool.contents[today] = append(tool.contents[today], message)
			writeChangelog(tool)

			result := strings.Join([]string{
				fmt.Sprintf("changelog.Remove: Note with %d B written.", len(message)),
			}, "\n")

			return result, nil

		} else {
			return "", fmt.Errorf("changelog.Remove: Note must start with \"Removed\".")
		}

	} else {
		return "", fmt.Errorf("changelog.Remove: %s", err0.Error())
	}

}

func (tool *Changelog) Search(path string) (string, error) {

	_, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		results := make(map[time.Time][]string, 0)

		for date, entries := range tool.contents {

			for _, entry := range entries {

				if strings.HasPrefix(entry, fmt.Sprintf("`%s`: ", path)) {

					_, ok := results[date]

					if ok == false {
						results[date] = make([]string, 0)
					}

					results[date] = append(results[date], entry)

				}

			}

		}

		dates  := make([]time.Time, 0)
		result := make([]string, 0)

		result = append(result, fmt.Sprintf("changelog.Search: %s", path))
		result = append(result, "")

		for date, _ := range results {
			dates = append(dates, date)
		}

		sort.Slice(dates, func(a int, b int) bool {
			return dates[a].After(dates[b])
		})

		for _, date := range dates {

			entries := results[date]

			result = append(result, fmt.Sprintf("Date: %s", date.Format("2006-01-02 15:04:05")))
			result = append(result, "")

			sort.Strings(entries)

			for _, note := range entries {
				result = append(result, fmt.Sprintf("%s", note))
			}

			result = append(result, "")

		}

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("changelog.Search: Invalid file path \"%s\".", path)
	}

}

