package main

import "encoding/json"
import "fmt"
import "os"
import "strings"

func writeRecordedArguments(arguments []string) {

	path := os.Getenv("EXOCOMP_FAKE_PODMAN_ARGS")

	if path != "" {

		data, err := json.Marshal(arguments)

		if err == nil {
			os.WriteFile(path, data, 0644)
		}

	}

}

func main() {

	arguments := os.Args[1:]

	writeRecordedArguments(arguments)

	if len(arguments) == 0 {
		os.Exit(2)
	}

	command := arguments[0]
	scenario := os.Getenv("EXOCOMP_FAKE_PODMAN")

	if scenario == "inspect-missing" && command == "inspect" {

		fmt.Fprintln(os.Stderr, "Error: no such object: fake")
		os.Exit(125)

	} else if scenario == "exec-fail" && command == "exec" {

		fmt.Fprintln(os.Stdout, "command output")
		fmt.Fprintln(os.Stderr, "Error: permission denied")
		os.Exit(1)

	} else if scenario == "exec-notfound" && command == "exec" {

		fmt.Fprintln(os.Stderr, "Error: executable file not found in $PATH")
		os.Exit(127)

	} else if scenario == "run-image-missing" && command == "run" {

		fmt.Fprintln(os.Stderr, "Error: no such image: alpine:latest")
		os.Exit(125)

	} else if scenario == "run-fail" && command == "run" {

		fmt.Fprintln(os.Stderr, "Error: statfs /home/alice/secret: no such file or directory")
		os.Exit(125)

	}

	switch command {

	case "inspect":

		label := os.Getenv("EXOCOMP_FAKE_PODMAN_LABEL")

		if os.Getenv("EXOCOMP_FAKE_PODMAN_FOREIGN") == "1" {
			label = "foreign"
		}

		name := arguments[len(arguments)-1]

		data, err := json.Marshal([]map[string]any{
			{
				"Id":        "0123456789abcdef0123456789abcdef",
				"Name":      "/" + name,
				"Image":     "sha256:0123456789abcdef",
				"ImageName": "alpine:latest",
				"State": map[string]any{
					"Status":    "running",
					"Running":   true,
					"StartedAt": "2026-01-01T08:00:00Z",
				},
				"Config": map[string]any{
					"Image": "alpine:latest",
					"Labels": map[string]string{
						"io.exocomp.sandbox": label,
					},
				},
			},
		})

		if err == nil {
			fmt.Println(string(data))
		}

	case "ps":

		if os.Getenv("EXOCOMP_FAKE_PODMAN_EMPTY") == "1" {
			fmt.Println("[]")
		} else {
			fmt.Println(`[{"Id":"abc","Names":["web"],"Image":"alpine:latest","Status":"Up 2 minutes"}]`)
		}

	case "images":

		fmt.Println(`[{"Id":"abc","RepoTags":["alpine:latest","busybox:latest"]}]`)

	case "exec":

		fmt.Println(strings.Join(arguments, " "))

	case "logs":

		fmt.Println("line one")
		fmt.Println("line two")

	default:

		fmt.Println(strings.Join(arguments, " "))

	}

}
