package types

import utils_fmt "exocomp/utils/fmt"
import "encoding/json"
import "fmt"
import "os"
import "path/filepath"
import "strings"

type Recovery struct {
	Playground string `json:"playground"`
}

func NewRecovery(playground string) *Recovery {

	return &Recovery{
		Playground: playground,
	}

}

func (recovery *Recovery) BackupAgent(agent *Agent) error {

	sanitized_name := utils_fmt.FormatAgentName(agent.Name)

	if sanitized_name == agent.Name {

		bytes, err1 := json.MarshalIndent(agent, "", "\t")

		if err1 == nil {

			path := filepath.Join(recovery.Playground, ".exocomp", "agents", fmt.Sprintf("%s.json", agent.Name))
			err2 := os.MkdirAll(filepath.Dir(path), 0755)

			if err2 == nil {

				err3 := os.WriteFile(path, bytes, 0666)

				if err3 == nil {
					return nil
				} else {
					return err3
				}

			} else {
				return err2
			}

		} else {
			return err1
		}

	} else {
		return fmt.Errorf("Invalid agent name \"%s\"", agent.Name)
	}

}

func (recovery *Recovery) BackupSession(session *Session) error {

	bytes, err1 := json.MarshalIndent(session, "", "\t")

	if err1 == nil {

		path := filepath.Join(recovery.Playground, ".exocomp", "session.json")
		err2 := os.MkdirAll(filepath.Dir(path), 0755)

		if err2 == nil {

			err3 := os.WriteFile(path, bytes, 0666)

			if err3 == nil {
				return recovery.BackupAgent(session.Agent)
			} else {
				return err3
			}

		} else {
			return err2
		}

	} else {
		return err1
	}

}

func (recovery *Recovery) HasBackup() bool {

	path      := filepath.Join(recovery.Playground, ".exocomp", "session.json")
	stat, err1 := os.Stat(path)

	if err1 == nil && stat.IsDir() == false {
		return true
	}

	return false

}

func (recovery *Recovery) RestoreAgents() []*Agent {

	result := make([]*Agent, 0)

	path          := filepath.Join(recovery.Playground, ".exocomp", "agents")
	entries, err1 := os.ReadDir(path)

	if err1 == nil {

		for _, entry := range entries {

			filename := entry.Name()
			ext      := filepath.Ext(filename)

			if ext == ".json" {

				basename := strings.TrimSpace(strings.TrimSuffix(filename, ext))

				if basename != "" {

					agent := recovery.RestoreAgent(basename)

					if agent != nil {
						result = append(result, agent)
					}

				}

			}

		}

	}

	return result

}

func (recovery *Recovery) RestoreAgent(name string) *Agent {

	sanitized_name := utils_fmt.FormatAgentName(name)

	if sanitized_name != "" {

		path        := filepath.Join(recovery.Playground, ".exocomp", "agents", fmt.Sprintf("%s.json", sanitized_name))
		bytes, err1 := os.ReadFile(path)

		if err1 == nil {

			agent := Agent{}
			err2  := json.Unmarshal(bytes, &agent)

			if err2 == nil {
				return &agent
			} else {
				return nil
			}

		} else {
			return nil
		}

	} else {
		return nil
	}

}

func (recovery *Recovery) RestoreSession() *Session {

	path        := filepath.Join(recovery.Playground, ".exocomp", "session.json")
	bytes, err1 := os.ReadFile(path)

	if err1 == nil {

		tmp  := Session{}
		err2 := json.Unmarshal(bytes, &tmp)

		if err2 == nil {
			return RestoreSession(recovery.Playground, tmp)
		} else {
			return nil
		}

	}

	return nil

}

func (recovery *Recovery) Snapshot(name string, raw any) error {

	if strings.Contains(name, ".") == false {

		bytes, err1 := json.MarshalIndent(raw, "", "\t")

		if err1 == nil {

			path := filepath.Join(recovery.Playground, ".exocomp", "debug", fmt.Sprintf("%s.json", name))
			err2 := os.MkdirAll(filepath.Dir(path), 0755)

			if err2 == nil {

				err3 := os.WriteFile(path, bytes, 0666)

				if err3 == nil {
					return nil
				} else {
					return err3
				}

			} else {
				return err2
			}

		} else {
			return err1
		}

	} else {
		return fmt.Errorf("invalid name \"%s\"")
	}

}

func (recovery *Recovery) SnapshotBytes(name string, raw []byte) error {

	if strings.Contains(name, ".") == false {

		var tmp interface{}

		err0 := json.Unmarshal(raw, &tmp)

		if err0 == nil {

			bytes, err1 := json.MarshalIndent(tmp, "", "\t")

			if err1 == nil {

				path := filepath.Join(recovery.Playground, ".exocomp", "debug", fmt.Sprintf("%s.json", name))
				err2 := os.MkdirAll(filepath.Dir(path), 0755)

				if err2 == nil {

					err3 := os.WriteFile(path, bytes, 0666)

					if err3 == nil {
						return nil
					} else {
						return err3
					}

				} else {
					return err2
				}

			} else {
				return err1
			}

		} else {
			return err0
		}

	} else {
		return fmt.Errorf("invalid name \"%s\"")
	}

}

