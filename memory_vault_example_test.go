package goblin_test

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aphistic/goblin"
)

func ExampleMemoryVault() {
	// Create a new memory vault
	mVault := goblin.NewMemoryVault()

	// Create a reader to write the file from with
	// some file contents.
	fileData := bytes.NewReader([]byte("I'm a file!"))

	// Write the file contents to your/file/here.txt with
	// a file modified time of April 8, 2020 at Midnight UTC.
	_ = mVault.WriteFile(
		"your/file/here.txt", fileData,
		goblin.FileModTime(time.Date(2020, 4, 8, 0, 0, 0, 0, time.UTC)),
	)

	// Read all the files in the root of the memory vault. In
	// our case it's only the "your" directory.
	infos, _ := mVault.ReadDir(".")
	fmt.Printf("Root Files:\n")
	for _, info := range infos {
		fmt.Printf("  - %s\n", info.Name())
	}

	// Read the contents of the file we previously wrote.
	readData, _ := mVault.ReadFile("your/file/here.txt")
	fmt.Printf("File Data: %s\n", readData)

	// Output:
	// Root Files:
	//   - your
	// File Data: I'm a file!
}
