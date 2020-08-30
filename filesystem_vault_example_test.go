package goblin_test

import (
	"fmt"

	"github.com/aphistic/goblin"
)

func ExampleFilesystemVault() {
	fsysVault := goblin.NewFilesystemVault("/")
	infos, _ := fsysVault.ReadDir(".")

	fmt.Printf("Files in /\n")
	for _, info := range infos {
		fmt.Printf("  - %s\n", info.Name())
	}
}
