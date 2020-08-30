package goblin_test

import (
	"fmt"
	"os"

	"github.com/aphistic/goblin"
)

func ExampleVaultSelector() {
	// Use a FilesystemVault if a path is provided as an environment variable,
	// but fall back to using a MemoryVault by default. This, for example,
	// could be useful during development when you are iterating on embedded
	// files and don't want to have to rebuild for every change.
	const assetVaultPathKey = "ASSET_VAULT_PATH"

	_ = os.Setenv(assetVaultPathKey, "/")

	// Create a filesystem vault and a memory vault to use.
	fsysAssetVault := goblin.NewFilesystemVault(os.Getenv(assetVaultPathKey))
	memAssetVault := goblin.NewMemoryVault() // Loaded from an embedded vault

	// Create a vault selector that will use the filesystem asset vault
	// if the environment variable ASSET_VAULT_PATH has a non-empty vault,
	// use the embedded asset vault if not.
	assetVaultSelector := goblin.NewVaultSelector(
		goblin.SelectEnvNotEmpty(assetVaultPathKey, fsysAssetVault),
		goblin.SelectDefault(memAssetVault),
	)

	// Since ASSET_VAULT_PATH is not empty, we use the filesystem vault.
	v, _ := assetVaultSelector.GetVault()
	fmt.Printf("Vault: %s\n", v)

	// Output:
	// Vault: Filesystem Vault (/)
}

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
