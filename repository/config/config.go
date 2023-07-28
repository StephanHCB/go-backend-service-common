package config

import (
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
	auconfigenv "github.com/StephanHCB/go-autumn-config-env"
)

// ApplicationName is only used to set up minimal logging if the configuration cannot be read.
//
// You should set this from your code, so log output can be found in ELK under your service name.
//
// If reading the configuration succeeds, the APPLICATION_NAME configuration is used instead.
var ApplicationName = "set-this-from-your-application"

const (
	KeyApplicationName      = "APPLICATION_NAME"
	KeyServerAddress        = "SERVER_ADDRESS"
	KeyServerPort           = "SERVER_PORT"
	KeyMetricsPort          = "METRICS_PORT"
	KeyEnvironment          = "ENVIRONMENT"
	KeyPlatform             = "PLATFORM"
	KeyLogstyle             = "LOGSTYLE"
	KeyLogLevel             = "LOG_LEVEL"
	KeyVaultServer          = "VAULT_SERVER"
	KeyVaultCertificateFile = "VAULT_CERTIFICATE_FILE"
	// KeyVaultSecretPath deprecated please migrate to KeyVaultSecretsConfig
	//
	// example for a migrated config.yaml
	// VAULT_SECRETS_CONFIG: >-
	//   {
	//     "<path to secrets in vault>": [
	//       {"vaultKey": "<key of secret>"},
	//       {"vaultKey": "<key of secret>", "configKey": "<if code uses a different config key then key of the secret in vault>"}
	//     ]
	//   }
	KeyVaultSecretPath = "VAULT_SECRET_PATH"
	// KeyLocalVaultToken deprecated please use KeyVaultAuthToken
	KeyLocalVaultToken = "LOCAL_VAULT_TOKEN"
	// KeyVaultKubernetesRole deprecated please use KeyVaultAuthKubernetesRole
	KeyVaultKubernetesRole = "VAULT_KUBERNETES_ROLE"
	// KeyVaultKubernetesTokenPath deprecated please use KeyVaultAuthKubernetesTokenPath
	KeyVaultKubernetesTokenPath = "VAULT_KUBERNETES_TOKEN_PATH"
	// KeyVaultKubernetesBackend deprecated please use KeyVaultAuthKubernetesBackend
	KeyVaultKubernetesBackend = "VAULT_KUBERNETES_BACKEND"
	KeyCorsAllowOrigin        = "CORS_ALLOW_ORIGIN"

	KeyVaultEnabled                 = "VAULT_ENABLED"
	KeyVaultAuthToken               = "VAULT_AUTH_TOKEN"
	KeyVaultAuthKubernetesRole      = "VAULT_AUTH_KUBERNETES_ROLE"
	KeyVaultAuthKubernetesTokenPath = "VAULT_AUTH_KUBERNETES_TOKEN_PATH"
	KeyVaultAuthKubernetesBackend   = "VAULT_AUTH_KUBERNETES_BACKEND"
	KeyVaultSecretsConfig           = "VAULT_SECRETS_CONFIG"
)

// PredefinedConfigItems is exposed so you can customize it.
//
// Must be done before calling New().
//
// You must NOT change the keys, but you can change the EnvName, defaults, validators etc.
//
// You must NOT remove any keys.
var PredefinedConfigItems = []auconfigapi.ConfigItem{
	{
		Key:         KeyApplicationName,
		EnvName:     KeyApplicationName,
		Default:     "",
		Description: "the name of the application, lowercase, numbers and - only",
		Validate:    auconfigenv.ObtainPatternValidator("^[a-z][a-z0-9-]*[a-z0-9]$"),
	}, {
		Key:         KeyServerAddress,
		EnvName:     KeyServerAddress,
		Default:     "",
		Description: "address to bind to, one of ip, hostname, [ipv6_ip], [ipv6ip%interface]",
		Validate:    auconfigenv.ObtainPatternValidator("^(|[a-z0-9.-]+|\\[[0-9a-f:]+%?[a-z0-9]*\\])$"),
	}, {
		Key:         KeyServerPort,
		EnvName:     KeyServerPort,
		Default:     "8080",
		Description: "port to listen on, cannot be a privileged port",
		Validate:    auconfigenv.ObtainUintRangeValidator(1024, 65535),
	}, {
		Key:         KeyMetricsPort,
		EnvName:     KeyMetricsPort,
		Default:     "9090",
		Description: "port to provide prometheus metrics on, cannot be a privileged port",
		Validate:    auconfigenv.ObtainUintRangeValidator(1024, 65535),
	}, {
		Key:         KeyEnvironment,
		EnvName:     KeyEnvironment,
		Default:     "dev",
		Description: "environment, used for vault secret lookups etc.",
		Validate:    auconfigenv.ObtainPatternValidator("^(feat|dev|test|acc|livetest|prod)$"),
	}, {
		Key:         KeyPlatform,
		EnvName:     KeyPlatform,
		Default:     "",
		Description: "platform, used for vault secret lookups etc.",
		Validate:    auconfigenv.ObtainPatternValidator("^[a-z]+$"),
	}, {
		Key:         KeyLogstyle,
		EnvName:     KeyLogstyle,
		Default:     "ecs",
		Description: "toggle between json ecs logging and plaintext logging (for local development)",
		Validate:    auconfigenv.ObtainPatternValidator("^(plain|ecs)$"),
	}, {
		Key:         KeyLogLevel,
		EnvName:     KeyLogLevel,
		Default:     "INFO",
		Description: "minimum level of logged messages",
		Validate:    auconfigenv.ObtainPatternValidator("^[a-zA-Z]+$"),
	}, {
		Key:         KeyVaultServer,
		EnvName:     KeyVaultServer,
		Default:     "my-vault-server.packetloss.de",
		Description: "fqdn of the vault server - do not add any other part of the URL",
		Validate:    auconfigenv.ObtainPatternValidator("^[a-z0-9.-]+$"),
	}, {
		Key:         KeyVaultCertificateFile,
		EnvName:     KeyVaultCertificateFile,
		Default:     "",
		Description: "optional: path to a custom ca cert file in PEM format",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	}, {
		Key:     KeyVaultSecretPath,
		EnvName: KeyVaultSecretPath,
		Default: "",
		Description: "the relative path to your secrets in vault (only the part after the environment, " +
			"v1/system_kv/data/v1/<PLATFORM>/microservices/<ENVIRONMENT> is added automatically)",
		Validate: auconfigapi.ConfigNeedsNoValidation,
	}, {
		Key:     KeyLocalVaultToken,
		EnvName: KeyLocalVaultToken,
		Default: "",
		Description: "directly supply a vault access token (for local development). " +
			"Setting this implicitly switches from kubernetes authentication to token mode. " +
			"Note: you can obtain a token in the vault UI by logging in and using the dropdown menu in the top right",
		Validate: auconfigapi.ConfigNeedsNoValidation,
	}, {
		Key:         KeyVaultKubernetesRole,
		EnvName:     KeyVaultKubernetesRole,
		Default:     "",
		Description: "role binding to use for vault kubernetes authentication, usually <PLATFORM>_microservice_role_<APPNAME>_<ENVIRONMENT>",
		Validate:    auconfigenv.ObtainPatternValidator("^(|[a-z]+_microservice_role_.*)$"),
	}, {
		Key:         KeyVaultKubernetesTokenPath,
		EnvName:     KeyVaultKubernetesTokenPath,
		Default:     "/run/secrets/kubernetes.io/serviceaccount/token",
		Description: "the path under which the service account token is injected into your container. The default should work",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	}, {
		Key:         KeyVaultKubernetesBackend,
		EnvName:     KeyVaultKubernetesBackend,
		Default:     "",
		Description: "role binding to use for vault kubernetes authentication, usually <PLATFORM>_microservice_role_<APPNAME>_<ENVIRONMENT>",
		Validate:    auconfigenv.ObtainPatternValidator("^(|k8s-[a-z-]+|aks-[a-z-]+)$"),
	}, {
		Key:         KeyCorsAllowOrigin,
		EnvName:     KeyCorsAllowOrigin,
		Default:     "",
		Description: "setting this enables sending headers to reduce CORS protections. Not usually suitable for production. Leave blank to not disable CORS. Note that this needs to be a single http(s) base URL, or else credentials forwarding will be refused by modern browsers. If you don't need credentials, this can be a comma separated list. Typical example value: 'http://localhost:8000/'",
		Validate:    auconfigenv.ObtainPatternValidator("^(|https?://.*)$"),
	},
}
