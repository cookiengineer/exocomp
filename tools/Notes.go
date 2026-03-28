package tools

import "fmt"
import "slices"
import "strconv"
import "strings"

type Notes struct {
	Sandbox string
	notes   map[uint]string
}

func NewNotes(agent string, sandbox string) *Notes {

	notes := &Notes{
		Sandbox: sandbox,
		notes:   make(map[uint]string),
	}

	readNotes(notes)

	return notes

}

func (tool *Notes) Call(method string, arguments map[string]interface{}) (string, error) {

	if method == "Create" {

		note, ok := arguments["note"].(string)

		if ok == true {
			return tool.Create(note)
		} else {
			return "", fmt.Errorf("notes.Create: Invalid parameters")
		}

	} else if method == "List" {

		return tool.List()

	} else if method == "Search" {

		keyword, ok := arguments["keyword"].(string)

		if ok == true {
			return tool.Search(keyword)
		} else {
			return "", fmt.Errorf("notes.Search: Invalid parameters")
		}

	} else if method == "Remove" {

		raw_value, ok := arguments["id"]

		if ok == true {

			id := uint(0)

			switch tmp := raw_value.(type) {
			case int:
				id = uint(tmp)
			case int32:
				id = uint(tmp)
			case int64:
				id = uint(tmp)
			case float64:
				id = uint(tmp)
			case uint:
				id = uint(tmp)
			case uint32:
				id = uint(tmp)
			case uint64:
				id = uint(tmp)
			default:
				id = uint(0)
			}

			return tool.Remove(id)

		} else {
			return "", fmt.Errorf("notes.Remove: Invalid parameters")
		}

	} else {
		return "", fmt.Errorf("notes.%s: Invalid method.", method)
	}

}

func (tool *Notes) Create(note string) (string, error) {

	description := strings.TrimSpace(note)

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
				fmt.Sprintf("notes.Create: Note with id %d and %d bytes written.", found_id, len(description)),
			}, "\n")

			return result, nil

		} else {
			return "", fmt.Errorf("notes.Create: %s", err.Error())
		}

	} else {
		return "", fmt.Errorf("notes.Create: Invalid note description.")
	}

}

func (tool *Notes) List() (string, error) {

	ids := make([]uint, 0)

	for n, _ := range tool.notes {
		ids = append(ids, n)
	}

	slices.Sort(ids)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("notes.List:"))

	for _, id := range ids {

		note := tool.notes[id]
		result = append(result, fmt.Sprintf("Id: %d, Note: %s", id, note))

	}

	return strings.Join(result, "\n"), nil

}

func (tool *Notes) Search(keyword string) (string, error) {

	compare := strings.ToLower(strings.TrimSpace(keyword))

	if strings.Contains(compare, " ") == false {

		ids := make([]uint, 0)

		for n, note := range tool.notes {

			if strings.Contains(strings.ToLower(note), compare) == true {
				ids = append(ids, n)
			}

		}

		slices.Sort(ids)

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("notes.Search:"))

		for _, id := range ids {

			note := tool.notes[id]
			result = append(result, fmt.Sprintf("Id: %d, Note: %s", id, note))

		}

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("notes.Search: Invalid keyword, no whitespaces allowed.")
	}

}

func (tool *Notes) Remove(id uint) (string, error) {

	_, ok := tool.notes[uint(id)]

	if ok == true {
		delete(tool.notes, uint(id))
	}

	result := strings.Join([]string{
		fmt.Sprintf("notes.Remove: Note with id %d removed.", strconv.FormatUint(uint64(id), 10)),
	}, "\n")

	return result, nil

}

