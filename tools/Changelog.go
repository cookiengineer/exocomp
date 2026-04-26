package tools

import utils_fmt "exocomp/utils/fmt"
import "fmt"
import "sort"
import "strings"
import "time"

type changelog_entry struct {
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	File        string    `json:"file"`
	Symbol      string    `json:"symbol"`
	Description string    `json:"description"`
}

type Changelog struct {
	Sandbox    string
	Playground string
	contents   map[string]map[string][]changelog_entry // map[path][symbol]
}

func NewChangelog(playground string, sandbox string) *Changelog {

	changelog := &Changelog{
		Playground: playground,
		Sandbox:    sandbox,
		contents:   make(map[string]map[string][]changelog_entry, 0),
	}

	readChangelog(changelog)

	return changelog

}

func (tool *Changelog) Call(method string, arguments map[string]interface{}) (string, error) {

	if method == "Add" {

		path,        ok1 := arguments["path"].(string)
		symbol,      ok2 := arguments["symbol"].(string)
		description, ok3 := arguments["description"].(string)

		if ok1 == true && ok2 == true && ok3 == true {
			return tool.Add(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatSingleLine(description))
		} else if ok1 == true && ok2 == true && ok3 == false {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"description\" is not a string.")
		} else if ok1 == true && ok2 == false && ok3 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
		} else if ok1 == false && ok2 == true && ok3 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Change" {

		path,        ok1 := arguments["path"].(string)
		symbol,      ok2 := arguments["symbol"].(string)
		description, ok3 := arguments["description"].(string)

		if ok1 == true && ok2 == true && ok3 == true {
			return tool.Change(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatSingleLine(description))
		} else if ok1 == true && ok2 == true && ok3 == false {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"description\" is not a string.")
		} else if ok1 == true && ok2 == false && ok3 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
		} else if ok1 == false && ok2 == true && ok3 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Deprecate" {

		path,        ok1 := arguments["path"].(string)
		symbol,      ok2 := arguments["symbol"].(string)
		description, ok3 := arguments["description"].(string)

		if ok1 == true && ok2 == true && ok3 == true {
			return tool.Deprecate(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatSingleLine(description))
		} else if ok1 == true && ok2 == true && ok3 == false {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"description\" is not a string.")
		} else if ok1 == true && ok2 == false && ok3 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
		} else if ok1 == false && ok2 == true && ok3 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Fix" {

		path,        ok1 := arguments["path"].(string)
		symbol,      ok2 := arguments["symbol"].(string)
		description, ok3 := arguments["description"].(string)

		if ok1 == true && ok2 == true && ok3 == true {
			return tool.Fix(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatSingleLine(description))
		} else if ok1 == true && ok2 == true && ok3 == false {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"description\" is not a string.")
		} else if ok1 == true && ok2 == false && ok3 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
		} else if ok1 == false && ok2 == true && ok3 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "List" {

		return tool.List()

	} else if method == "Remove" {

		path,        ok1 := arguments["path"].(string)
		symbol,      ok2 := arguments["symbol"].(string)
		description, ok3 := arguments["description"].(string)

		if ok1 == true && ok2 == true && ok3 == true {
			return tool.Remove(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatSingleLine(description))
		} else if ok1 == true && ok2 == true && ok3 == false {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"description\" is not a string.")
		} else if ok1 == true && ok2 == false && ok3 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
		} else if ok1 == false && ok2 == true && ok3 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Search" {

		path,   ok1 := arguments["path"].(string)
		symbol, ok2 := arguments["symbol"].(string)

		if ok1 == true && ok2 == true {
			return tool.Search(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol))
		} else if ok1 == true && ok2 == false {
			return tool.Search(utils_fmt.FormatFilePath(path), "")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
		}

	} else {
		return "", fmt.Errorf("changelog.%s: Invalid method.", method)
	}

}

func (tool *Changelog) Add(path string, symbol string, description string) (string, error) {
	return tool.createEntry("Add", path, symbol, description)
}

func (tool *Changelog) Change(path string, symbol string, description string) (string, error) {
	return tool.createEntry("Change", path, symbol, description)
}

func (tool *Changelog) Deprecate(path string, symbol string, description string) (string, error) {
	return tool.createEntry("Deprecate", path, symbol, description)
}

func (tool *Changelog) Fix(path string, symbol string, description string) (string, error) {
	return tool.createEntry("Fix", path, symbol, description)
}

func (tool *Changelog) List() (string, error) {

	found := make(map[time.Time][]string, 0)

	for _, symbols := range tool.contents {

		for _, entries := range symbols {

			for _, entry := range entries {

				resolved_path, err1 := resolveSandboxPath(tool.Playground, entry.File)

				if err1 == nil {

					sandbox_path, err2 := sanitizeSandboxPath(tool.Sandbox, resolved_path)

					if err2 == nil {
						found[entry.Date] = append(found[entry.Date], fmt.Sprintf("- Date: %s, Type: %s, File: %s, Symbol: %s, Description: %s", entry.Date.Format("2006-01-02"), entry.Type, sandbox_path, entry.Symbol, entry.Description))
					}

				}

			}

		}

	}

	dates := make([]time.Time, 0)
	lines := make([]string, 0)

	for date, _ := range found {
		dates = append(dates, date)
	}

	sort.Slice(dates, func(a int, b int) bool {
		return dates[a].Before(dates[b])
	})

	for _, date := range dates {

		for _, line := range found[date] {
			lines = append(lines, line)
		}

	}

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("changelog.List: %d changelog entries.", len(lines)))

	for l := 0; l < len(lines); l++ {
		result = append(result, lines[l])
	}

	return strings.Join(result, "\n"), nil

}

func (tool *Changelog) Remove(path string, symbol string, description string) (string, error) {
	return tool.createEntry("Remove", path, symbol, description)
}

func (tool *Changelog) Search(path string, symbol string) (string, error) {

	tmp1, err1 := resolveSandboxPath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := sanitizeSandboxPath(tool.Playground, tmp1)

		if err2 == nil {

			found := make(map[time.Time][]string, 0)

			if symbol != "" {

				_, ok1 := tool.contents[internal_path]

				if ok1 == true {

					entries, ok2 := tool.contents[internal_path][symbol]

					if ok2 == true {

						for _, entry := range entries {

							sandbox_path, err3 := sanitizeSandboxPath(tool.Sandbox, entry.File)

							if err3 == nil {
								found[entry.Date] = append(found[entry.Date], fmt.Sprintf("- Date: %s, Type: %s, File: %s, Symbol: %s, Description: %s", entry.Date.Format("2006-01-02"), entry.Type, sandbox_path, entry.Symbol, entry.Description))
							}

						}

					}

				}

				dates := make([]time.Time, 0)
				lines := make([]string, 0)

				for date, _ := range found {
					dates = append(dates, date)
				}

				sort.Slice(dates, func(a int, b int) bool {
					return dates[a].Before(dates[b])
				})

				for _, date := range dates {

					for _, line := range found[date] {
						lines = append(lines, line)
					}

				}

				result := make([]string, 0)
				result = append(result, fmt.Sprintf("changelog.Search: %s#%s contains %d changelog entries.", path, symbol, len(lines)))

				for l := 0; l < len(lines); l++ {
					result = append(result, lines[l])
				}

				return strings.Join(result, "\n"), nil

			} else {

				symbols, ok1 := tool.contents[internal_path]

				if ok1 == true {

					sorted_symbols := make([]string, 0)

					for symbol, _ := range symbols {
						sorted_symbols = append(sorted_symbols, symbol)
					}

					sort.Strings(sorted_symbols)

					for _, symbol := range sorted_symbols {

						entries := tool.contents[internal_path][symbol]

						for _, entry := range entries {

							sandbox_path, err3 := sanitizeSandboxPath(tool.Sandbox, entry.File)

							if err3 == nil {
								found[entry.Date] = append(found[entry.Date], fmt.Sprintf("- Date: %s, Type: %s, File: %s, Symbol: %s, Description: %s", entry.Date.Format("2006-01-02"), entry.Type, sandbox_path, entry.Symbol, entry.Description))
							}

						}

					}

				}

				dates := make([]time.Time, 0)
				lines := make([]string, 0)

				for date, _ := range found {
					dates = append(dates, date)
				}

				sort.Slice(dates, func(a int, b int) bool {
					return dates[a].Before(dates[b])
				})

				for _, date := range dates {

					for _, line := range found[date] {
						lines = append(lines, line)
					}

				}

				result := make([]string, 0)
				result = append(result, fmt.Sprintf("changelog.Search: %s contains %d changelog entries.", path, len(lines)))

				for l := 0; l < len(lines); l++ {
					result = append(result, lines[l])
				}

				return strings.Join(result, "\n"), nil

			}

		} else {
			return "", fmt.Errorf("changelog.Search: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("changelog.Search: %s", err1.Error())
	}

}

func (tool *Changelog) createEntry(method string, path string, symbol string, description string) (string, error) {

	tmp1, err1 := resolveSandboxPath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := sanitizeSandboxPath(tool.Playground, tmp1)

		if err2 == nil {

			_, ok1 := tool.contents[internal_path]

			if ok1 == false {
				tool.contents[internal_path] = make(map[string][]changelog_entry, 0)
			}

			_, ok2 := tool.contents[internal_path][symbol]

			if ok2 == false {
				tool.contents[internal_path][symbol] = make([]changelog_entry, 0)
			}

			var found *changelog_entry = nil

			for _, entry := range tool.contents[internal_path][symbol] {

				if entry.Type == method && entry.File == internal_path && entry.Symbol == symbol && entry.Description == description {
					found = &entry
					break
				}

			}

			if found == nil {

				now   := time.Now()
				today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

				tool.contents[internal_path][symbol] = append(tool.contents[internal_path][symbol], changelog_entry{
					Date:        today,
					Type:        method,
					File:        internal_path,
					Symbol:      symbol,
					Description: description,
				})

				err3 := writeChangelog(tool)

				if err3 == nil {
					return fmt.Sprintf("changelog.%s: Log entry created for %s#%s at %s.", method, path, symbol, today.Format("2006-01-02")), nil
				} else {
					return "", fmt.Errorf("changelog.%s: %s", method, err3.Error())
				}

			} else {
				return fmt.Sprintf("changelog.%s: Log entry already exists for %s#%s at %s.", method, path, found.Symbol, found.Date.Format("2006-01-02")), nil
			}

		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, err2.Error())
		}

	} else {
		return "", fmt.Errorf("changelog.%s: %s", method, err1.Error())
	}

}
