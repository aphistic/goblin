// Package main is the utility used to create and manage goblin vaults.
package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"

	"github.com/aphistic/goblin"
	"github.com/aphistic/goblin/internal/logging"
)

var (
	flagName         string
	flagPackage      string
	flagOut          string
	flagIncludeRoot  string
	flagIncludes     []string
	flagExportLoader bool
	flagBinary       bool
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
	cmdCreate.Flag("include-root", "Root path to use when including files in the vault").Short('r').
		StringVar(&flagIncludeRoot)
	cmdCreate.Flag("include", "Files to include in the vault").Short('i').
		StringsVar(&flagIncludes)
	cmdCreate.Flag("export-loader", "Export loader in generated code").Short('e').
		BoolVar(&flagExportLoader)
	cmdCreate.Flag("binary", "Write out binary data").Short('b').BoolVar(&flagBinary)

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

	logger := logging.NewPrintfLogger()

	b := goblin.NewMemoryBuilder(
		goblin.MemoryBuilderLogger(logger),
		goblin.MemoryBuilderExportLoader(flagExportLoader),
	)
	err = b.Include(flagIncludeRoot, flagIncludes)
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

	if flagBinary {
		err = b.WriteBinary(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing binary file: %s\n", err)
			os.Exit(1)
		}

	} else {
		err = b.WriteLoader(flagPackage, flagName, f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing code file: %s\n", err)
			os.Exit(1)
		}
	}
}
