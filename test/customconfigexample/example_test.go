package customconfigexample

import (
	"bytes"
	"context"
	auconfigenv "github.com/StephanHCB/go-autumn-config-env"
	goauzerolog "github.com/StephanHCB/go-autumn-logging-zerolog"
	"github.com/StephanHCB/go-backend-service-common/acorns/repository"
	"github.com/StephanHCB/go-backend-service-common/docs"
	"github.com/StephanHCB/go-backend-service-common/repository/config"
	"github.com/StephanHCB/go-backend-service-common/repository/logging"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const basedir = "../resources/"

func TestRead_EnvWins(t *testing.T) {
	docs.Description("when reading config, environment supersedes yaml values")

	os.Setenv("APPLICATION_NAME", "room-service")
	cut := New().(repository.Configuration)
	auconfigenv.LocalConfigFileName = basedir + "valid-config.yaml"
	err := cut.Read()
	cut.(*config.ConfigImpl).ObtainPredefinedValues()

	os.Unsetenv("APPLICATION_NAME")

	require.Nil(t, err)
	require.Equal(t, "room-service", cut.ApplicationName())
}

func tstSetupCutAndLogRecorder(t *testing.T, configfile string) (repository.Configuration, error) {
	cut := New().(repository.Configuration)

	// Phase --- AssembleAcorn ---

	auconfigenv.LocalConfigFileName = basedir + configfile
	err := cut.Read()
	require.Nil(t, err)

	// Phase --- SetupAcorn --- (but without firing up the whole Acorn subsystem)

	// set up log recorder
	logRecorder := logging.New().(repository.Logging)
	goauzerolog.RecordedLogForTesting = new(bytes.Buffer)
	logRecorder.(*logging.LoggingImpl).SetupForTesting()

	cut.(*config.ConfigImpl).Logging = logRecorder

	// validate and obtain values
	ctx := log.Logger.WithContext(context.Background())
	err = cut.Validate(ctx)
	cut.(*config.ConfigImpl).ObtainPredefinedValues()
	cut.(*config.ConfigImpl).CustomConfiguration.Obtain(auconfigenv.Get)

	return cut, err
}

func TestRead_DefaultLoses(t *testing.T) {
	docs.Description("when reading config, defaults do not supersede yaml values, but are set if value empty")

	cut, validationErr := tstSetupCutAndLogRecorder(t, "valid-config.yaml")
	require.Nil(t, validationErr)

	cut.(*config.ConfigImpl).ObtainPredefinedValues()

	require.Equal(t, "/run/secrets/kubernetes.io/serviceaccount/token", cut.VaultKubernetesTokenPath())
	require.Equal(t, "demo-backend", cut.ApplicationName())
}

func TestValidate_LotsOfErrors(t *testing.T) {
	docs.Description("validation of configuration values works")

	_, err := tstSetupCutAndLogRecorder(t, "invalid-config-values.yaml")
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "some configuration values failed to validate or parse. There were 8 error(s). See details above")

	actualLog := goauzerolog.RecordedLogForTesting.String()

	expectedPart := "\"message\":\"failed to validate configuration field ENVIRONMENT: must match ^(feat|"
	require.Contains(t, actualLog, expectedPart)

	expectedPart2 := "\"message\":\"failed to validate configuration field SERVER_PORT: value 122834 is out of range [1024..65535]"
	require.Contains(t, actualLog, expectedPart2)

	expectedPart3 := "METRICS_PORT: value -12387192873invalid is not a valid integer"
	require.Contains(t, actualLog, expectedPart3)

	expectedPart4 := "\"message\":\"failed to validate configuration field CORS_ALLOW_ORIGIN: must match ^(|https?://.*)$"
	require.Contains(t, actualLog, expectedPart4)
}

func TestAccessors(t *testing.T) {
	docs.Description("the config accessors return the correct values")

	cut, err := tstSetupCutAndLogRecorder(t, "valid-config-unique.yaml")
	require.Nil(t, err)

	actualLog := goauzerolog.RecordedLogForTesting.String()
	require.Equal(t, "", actualLog)

	require.Equal(t, "room-service", cut.ApplicationName())
	require.Equal(t, "192.168.150.0", cut.ServerAddress())
	require.Equal(t, uint16(8081), cut.ServerPort())
	require.Equal(t, uint16(9091), cut.MetricsPort())
	require.Equal(t, "dev", cut.Environment())
	require.Equal(t, "platform", cut.Platform())
	require.Equal(t, true, cut.PlainLogging())
	require.Equal(t, "localhost", cut.VaultServer())
	require.Equal(t, "", cut.VaultCertificateFile())
	require.Equal(t, "room-service/secrets", cut.VaultSecretPath())
	require.Equal(t, true, cut.LocalVault())
	require.Equal(t, "not a real token", cut.LocalVaultToken())
	require.Equal(t, "platform_microservice_role_room-service_prod", cut.VaultKubernetesRole())
	require.Equal(t, "/some/thing", cut.VaultKubernetesTokenPath())
	require.Equal(t, "k8s-dev-something", cut.VaultKubernetesBackend())
	require.Equal(t, "http://localhost:12345", cut.CorsAllowOrigin())

	require.Equal(t, "kitty", cut.Custom().(CustomConfigurationWithOneField).MyCustomField())
}
