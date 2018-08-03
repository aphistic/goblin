package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"

	"github.com/aphistic/goblin"
)

var (
	flagName     string
	flagPackage  string
	flagOut      string
	flagIncludes []string
)

func main() {
	appGoblin := kingpin.New("goblin", "Goblin")
	cmdCreate := appGoblin.Command("create", "Create a vault").Default()
	cmdCreate.Flag("name", "Name of the vault to create").Short('n').
		Required().StringVar(&flagName)
	cmdCreate.Flag("package", "Name of the package for the output file").Short('p').
		StringVar(&flagPackage)
	cmdCreate.Flag("out", "Name to use for the output file").Short('o').
		StringVar(&flagOut)
	cmdCreate.Flag("include", "Files to include in the vault").Short('i').
		StringsVar(&flagIncludes)

	_, err := appGoblin.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse command line arguments: %s\n", err)
		os.Exit(1)
	}
	if flagPackage == "" {
		flagPackage = flagName
	}
	if flagOut == "" {
		flagOut = fmt.Sprintf("goblin_%s.go", flagName)
	}

	b := goblin.NewBuilder(flagName)
	err = b.Include(flagIncludes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error including files: %s\n", err)
		os.Exit(1)
	}

	f, err := os.OpenFile(flagOut, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open output file %s: %s\n", flagOut, err)
		os.Exit(1)
	}
	defer f.Close()

	err = b.Write(flagPackage, f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		os.Exit(1)
	}
}
