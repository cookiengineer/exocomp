package tools

import "exocomp/utils"
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

func NewChangelog(agent string, sandbox string, playground string) *Changelog {

	changelog := &Changelog{
		Sandbox:    sandbox,
		Playground: playground,
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
			return tool.Add(utils.FormatFilePath(path), utils.FormatSymbol(symbol), utils.FormatSingleLine(description))
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
			return tool.Change(utils.FormatFilePath(path), utils.FormatSymbol(symbol), utils.FormatSingleLine(description))
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
			return tool.Deprecate(utils.FormatFilePath(path), utils.FormatSymbol(symbol), utils.FormatSingleLine(description))
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
			return tool.Fix(utils.FormatFilePath(path), utils.FormatSymbol(symbol), utils.FormatSingleLine(description))
		} else if ok1 == true && ok2 == true && ok3 == false {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"description\" is not a string.")
		} else if ok1 == true && ok2 == false && ok3 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
		} else if ok1 == false && ok2 == true && ok3 == true {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Remove" {

		path,        ok1 := arguments["path"].(string)
		symbol,      ok2 := arguments["symbol"].(string)
		description, ok3 := arguments["description"].(string)

		if ok1 == true && ok2 == true && ok3 == true {
			return tool.Remove(utils.FormatFilePath(path), utils.FormatSymbol(symbol), utils.FormatSingleLine(description))
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
			return tool.Search(utils.FormatFilePath(path), utils.FormatSymbol(symbol))
		} else if ok1 == true && ok2 == false {
			return tool.Search(utils.FormatFilePath(path), "")
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

func (tool *Changelog) Remove(path string, symbol string, description string) (string, error) {
	return tool.createEntry("Remove", path, symbol, description)
}

func (tool *Changelog) Search(path string, symbol string) (string, error) {
	return "", fmt.Errorf("changelog.Search: %s", "Not implemented")
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

			found := false

			for _, entry := range tool.contents[internal_path][symbol] {

				if entry.Type == method && entry.File == internal_path && entry.Symbol == symbol && entry.Description == description {
					found = true
					break
				}

			}

			if found == false {

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
					return fmt.Sprintf("changelog.%s: Entry created for %s#%s at %s.", method, path, symbol, today.Format("2006-01-02")), nil
				} else {
					return "", fmt.Errorf("changelog.%s: %s", method, err3.Error())
				}

			} else {
				return fmt.Sprintf("changelog.%s: Entry already exists for %s#%s.", method, path, symbol), nil
			}

		} else {
			return "", fmt.Errorf("changelog.%s: %s", method, err2.Error())
		}

	} else {
		return "", fmt.Errorf("changelog.%s: %s", method, err1.Error())
	}

}
