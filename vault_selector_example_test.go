package goblin_test

import (
	"fmt"
	"os"

	"github.com/aphistic/goblin"
)

func ExampleSelectEnvBool() {
	const envKey = "USE_SELECT_ENV_BOOL"
	os.Setenv(envKey, "true")

	memVault := goblin.NewMemoryVault()
	fsVault := goblin.NewFilesystemVault("/")

	vs := goblin.NewVaultSelector(
		goblin.SelectEnvBool(envKey, memVault),
		goblin.SelectDefault(fsVault),
	)

	fmt.Printf("%s\n", vs)
	// Output: Vault Selector (Memory Vault)
}

func ExampleSelectEnvNotEmpty() {
	const envKey = "USE_SELECT_ENV_NON_EMPTY"
	os.Setenv(envKey, "any value here")

	memVault := goblin.NewMemoryVault()
	fsVault := goblin.NewFilesystemVault("/")

	vs := goblin.NewVaultSelector(
		goblin.SelectEnvNotEmpty(envKey, memVault),
		goblin.SelectDefault(fsVault),
	)

	fmt.Printf("%s\n", vs)
	// Output: Vault Selector (Memory Vault)
}

func ExampleSelectPath() {
	const vaultPath = "/usr"

	rootVault := goblin.NewFilesystemVault("/")
	usrVault := goblin.NewFilesystemVault(vaultPath)

	vs := goblin.NewVaultSelector(
		goblin.SelectPath(vaultPath, usrVault),
		goblin.SelectDefault(rootVault),
	)

	fmt.Printf("%s\n", vs)
	// Output: Vault Selector (Filesystem Vault (/usr))
}

func ExampleSelectDefault() {
	fsVault := goblin.NewFilesystemVault("/")
	memVault := goblin.NewMemoryVault()

	// The first vault will always be selected because it
	// will always return a vault value.
	vs := goblin.NewVaultSelector(
		goblin.SelectDefault(fsVault),
		goblin.SelectDefault(memVault),
	)

	fmt.Printf("%s\n", vs)
	// Output: Vault Selector (Filesystem Vault (/))
}
