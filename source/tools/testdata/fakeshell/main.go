package main

import "bufio"
import "fmt"
import "os"
import "strings"
import "time"

import "golang.org/x/term"

func main() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("prompt> ")

	for {

		line, err := reader.ReadString('\n')

		if err == nil {

			command := strings.TrimSpace(line)

			if command == "exit" {
				fmt.Print("bye\n")
				return
			} else if command == "fail" {
				fmt.Print("failing\n")
				os.Exit(3)
			} else if command == "hello" {
				fmt.Print("hello world\n")
			} else if command == "ansi" {
				fmt.Print("\x1b[31mred\x1b[0m\n")
			} else if command == "big" {
				fmt.Print(strings.Repeat("A", 10000) + "\n")
			} else if command == "size" {
				columns, rows, size_err := term.GetSize(int(os.Stdout.Fd()))
				if size_err == nil {
					fmt.Printf("size: %dx%d\n", columns, rows)
				} else {
					fmt.Print("size: unknown\n")
				}
			} else if command == "wait" {
				fmt.Print("waiting\n")
				for {
					time.Sleep(time.Hour)
				}
			} else if strings.HasPrefix(command, "echo ") {
				fmt.Print(strings.TrimPrefix(command, "echo ") + "\n")
			} else if command != "" {
				fmt.Print("unknown\n")
			}

			fmt.Print("prompt> ")

		} else {
			return
		}

	}

}
