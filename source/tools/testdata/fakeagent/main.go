package main

import "encoding/json"
import "fmt"
import "os"
import "strings"
import "time"

func emitMessage(role string, content string) {

	data, err := json.Marshal(struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{Role: role, Content: content})

	if err == nil {
		fmt.Printf("schemas.Message:%s\n", string(data))
		os.Stdout.Sync()
	}

}

func emitQuit(report string) {
	emitMessage("tool", "agents.Quit: Work Report\n"+report)
}

func main() {

	scenario := os.Getenv("EXOCOMP_FAKE_SCENARIO")

	switch scenario {

	case "quit":
		emitMessage("system", "hello")
		emitQuit("my work is done: implemented the feature")
		os.Exit(0)

	case "large":
		emitMessage("tool", "files.Read: big.go\n"+strings.Repeat("x", 200*1024))
		emitQuit("my work is done: implemented the feature")
		os.Exit(0)

	case "noquit":
		emitMessage("system", "hello")
		os.Exit(0)

	case "fail":
		emitMessage("system", "hello")
		os.Exit(1)

	case "hang":
		time.Sleep(24 * time.Hour)

	default:
		emitMessage("system", "hello")
		emitQuit("my work is done: implemented the feature")
		os.Exit(0)

	}

}
