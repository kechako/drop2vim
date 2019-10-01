package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	appName = "drop2vim"
	version = "0.1.0"
)

const versionMsg = "%s %s\n"

const helpMsg = `Usage: %s <file_names>...

  -v, --version  print the version number
  -h, --help     show this message
`

func printVersion() {
	fmt.Printf(versionMsg, appName, version)
}

func showHelp() {
	printVersion()
	fmt.Printf(helpMsg, appName)
}

type flags int

const (
	helpFlag flags = 1 << iota
	versionFlag

	noFlags  flags = 0
	allFlags       = helpFlag | versionFlag
)

func getFlags(args []string) flags {
	var f flags
	for _, arg := range args {
		switch arg {
		case "-h", "--help":
			f |= helpFlag
		case "-v", "--version":
			f |= versionFlag
		}
		if f == allFlags {
			break
		}
	}

	return f
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "No file names are specified.")
		os.Exit(2)
	}

	if f := getFlags(args); f&helpFlag > 0 {
		showHelp()
		return
	} else if f&versionFlag > 0 {
		printVersion()
		return
	}

	var cmdList []string

	for _, file := range args {
		path, err := filepath.Abs(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "The specified file path is not valid: %s\n", path)
			os.Exit(1)
		}

		if _, err := os.Stat(path); err != nil {
			fmt.Fprintf(os.Stderr, "The specified file is not found: %s\n", file)
			os.Exit(1)
		}

		cmdList = append(cmdList, "drop", path)
	}

	cmdMsg, err := json.Marshal(cmdList)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to generate a command message")
		os.Exit(1)
	}

	fmt.Printf("\x1b]51;%s\x07", cmdMsg)
}
