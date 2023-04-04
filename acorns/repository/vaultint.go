package repository

import "context"

const VaultAcornName = "vault"

// Vault is the central singleton representing Hashicorp Vault.
//
// We use Vault to obtain sensitive configuration values, called "secrets".
type Vault interface {
	IsVault() bool

	// Setup uses the configuration
	Setup(ctx context.Context) error

	// Authenticate authenticates against vault
	Authenticate(ctx context.Context) error

	// ObtainSecrets fetches the regular secrets from vault
	ObtainSecrets(ctx context.Context) error
}

type VaultConfiguration interface {
	// TODO why is this here? empty interfaces don't do anything useful, they're the same as interface{}
}

type VaultSecretsConfig map[string][]VaultSecretConfig

type VaultSecretConfig struct {
	VaultKey  string  `json:"vaultKey"`
	ConfigKey *string `json:"configKey,omitempty"`
}
