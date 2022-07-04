package config

import "github.com/StephanHCB/go-backend-service-common/acorns/repository"

func (c *ConfigImpl) Custom() repository.CustomConfiguration {
	return c.CustomConfiguration
}

func (c *ConfigImpl) ApplicationName() string {
	return c.VApplicationName
}

func (c *ConfigImpl) ServerAddress() string {
	return c.VServerAddress
}

func (c *ConfigImpl) ServerPort() uint16 {
	return c.VServerPortValue
}

func (c *ConfigImpl) MetricsPort() uint16 {
	return c.VMetricsPortValue
}

func (c *ConfigImpl) Environment() string {
	return c.VEnvironment
}

func (c *ConfigImpl) PlainLogging() bool {
	return c.VLogstyle == "plain"
}

func (c *ConfigImpl) VaultServer() string {
	return c.VVaultServer
}

func (c *ConfigImpl) VaultCertificateFile() string {
	return c.VVaultCertFile
}

func (c *ConfigImpl) VaultSecretPath() string {
	return c.VVaultSecretPath
}

func (c *ConfigImpl) LocalVault() bool {
	return c.VLocalVaultToken != ""
}

func (c *ConfigImpl) LocalVaultToken() string {
	return c.VLocalVaultToken
}

func (c *ConfigImpl) VaultKubernetesRole() string {
	return c.VVaultK8sRole
}

func (c *ConfigImpl) VaultKubernetesTokenPath() string {
	return c.VVaultK8sTokenPath
}

func (c *ConfigImpl) VaultKubernetesBackend() string {
	return c.VVaultK8sBackend
}
