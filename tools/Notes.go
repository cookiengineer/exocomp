package tools

import "exocomp/utils"
import "fmt"
import "slices"
import "strconv"
import "strings"

type Notes struct {
	Method    string
	Arguments []string
	Sandbox   string
	notes     map[uint]string
}

func NewNotes(agent string, sandbox string, tools []string, programs []string) *Notes {

	notes := &Notes{
		Method:    "",
		Arguments: make([]string, 0),
		Sandbox:   sandbox,
		notes:     make(map[uint]string),
	}

	readNotes(notes)

	return notes

}

func (tool *Notes) Call() (string, error) {

	if tool.Method == "Help" {
		return tool.Help(tool.Arguments)
	} else if tool.Method == "Create" {
		return tool.Create(tool.Arguments)
	} else if tool.Method == "List" {
		return tool.List(tool.Arguments)
	} else if tool.Method == "Remove" {
		return tool.Remove(tool.Arguments)
	} else if tool.Method == "Search" {
		return tool.Search(tool.Arguments)
	} else {
		return "", fmt.Errorf("#!tool:notes.%s: Invalid method.", tool.Method)
	}

}

func (tool *Notes) Help(arguments []string) (string, error) {

	return strings.Join([]string{
		"#!tool:notes.Create \"Note description with whitespaces\"",
		"",
		"#!tool:notes.List",
		"",
		"#!tool:notes.Search keyword",
		"",
		"#!tool:notes.Remove note-id",
	}, "\n"), nil

}

func (tool *Notes) Create(arguments []string) (string, error) {

	if len(arguments) == 1 {

		description := strings.TrimSpace(arguments[0])

		description = strings.ReplaceAll(description, "\r", "")
		description = strings.ReplaceAll(description, "\n", " ")
		description = strings.ReplaceAll(description, "\t", " ")

		if description != "" {

			found_id := uint(0)
			last_id  := uint(0)

			for n, note := range tool.notes {

				if note == description {
					found_id = n
				}

				if n > last_id {
					last_id = n
				}

			}

			if found_id == 0 {
				found_id = uint(last_id + 1)
				tool.notes[found_id] = description
			}

			err := writeNotes(tool)

			if err == nil {

				result := strings.Join([]string{
					fmt.Sprintf("#!tool:notes.Create: Note with id %d and %d bytes written.", found_id, len(description)),
				}, "\n")

				return result, nil

			} else {
				return "", fmt.Errorf("#!tool:notes.Create: %s", err.Error())
			}

		} else {
			return "", fmt.Errorf("#!tool:notes.Create: Invalid note description.")
		}

	} else {
		return "", fmt.Errorf("#!tool:notes.Create: Invalid arguments, only one argument allowed.")
	}

}

func (tool *Notes) List(arguments []string) (string, error) {

	if len(arguments) == 0 {

		ids := make([]uint, 0)

		for n, _ := range tool.notes {
			ids = append(ids, n)
		}

		slices.Sort(ids)

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("#!tool:notes.List:"))

		for _, id := range ids {

			note := tool.notes[id]
			result = append(result, fmt.Sprintf("Id: %d, Note: %s", id, note))

		}

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("#!tool:notes.List: Invalid arguments, only zero arguments allowed.")
	}

}

func (tool *Notes) Parse(text string) (Tool, [2]int, error) {

	// #!tool:notes.Create "note description"
	// #!tool:notes.List
	// #!tool:notes.Search keyword
	// #!tool:notes.Remove note-id

	lines := strings.Split(text, "\n")

	if len(lines) > 0 && strings.HasPrefix(lines[0], "#!tool:notes.") {

		fields := utils.SplitArguments(strings.TrimSpace(lines[0][len("#!tool:notes."):]))
		method := strings.ToUpper(fields[0][0:1]) + strings.ToLower(fields[0][1:])
		parsed := [2]int{0, 1}

		if method == "Help" ||
			method == "Create" ||
			method == "Create" ||
			method == "Create" ||
			method == "Create" {

			tool.Method    = method
			tool.Arguments = fields[1:]

			return Tool(tool), parsed, nil

		} else {
			return nil, [2]int{0, len(lines)}, fmt.Errorf("#!tool:notes.%s: Invalid method.", method)
		}

	} else {
		return nil, [2]int{0, len(lines)}, fmt.Errorf("Invalid Tool Call line.")
	}

}

func (tool *Notes) Search(arguments []string) (string, error) {

	if len(arguments) == 1 {

		keyword := strings.ToLower(strings.TrimSpace(arguments[0]))

		if strings.Contains(keyword, " ") == false {

			ids := make([]uint, 0)

			for n, note := range tool.notes {

				if strings.Contains(strings.ToLower(note), keyword) == true {
					ids = append(ids, n)
				}

			}

			slices.Sort(ids)

			result := make([]string, 0)
			result = append(result, fmt.Sprintf("#!tool:notes.Search:"))

			for _, id := range ids {

				note := tool.notes[id]
				result = append(result, fmt.Sprintf("Id: %d, Note: %s", id, note))

			}

			return strings.Join(result, "\n"), nil

		} else {
			return "", fmt.Errorf("#!tool:notes.Search: Invalid keyword, no whitespaces allowed.")
		}

	} else {
		return "", fmt.Errorf("#!tool:notes.Search: Invalid arguments, only one argument allowed.")
	}

}

func (tool *Notes) Remove(arguments []string) (string, error) {

	if len(arguments) == 1 {

		id, err := strconv.ParseUint(arguments[0], 10, 64)

		if err == nil {

			_, ok := tool.notes[uint(id)]

			if ok == true {
				delete(tool.notes, uint(id))
			}

			result := strings.Join([]string{
				fmt.Sprintf("#!tool:notes.Remove: Note with id %d removed.", strconv.FormatUint(id, 10)),
			}, "\n")

			return result, nil

		} else {
			return "", fmt.Errorf("#!tool:notes.Remove: Invalid argument, not a Note identifier.")
		}

	} else {
		return "", fmt.Errorf("#!tool:notes.Remove: Invalid arguments, only one argument allowed.")
	}

}

