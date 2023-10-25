package vault

import (
	"context"
	"encoding/json"
	"fmt"
	auconfigenv "github.com/StephanHCB/go-autumn-config-env"
	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"
	aurestmock "github.com/StephanHCB/go-autumn-restclient/implementation/mock"
	"github.com/StephanHCB/go-backend-service-common/acorns/repository"
	"github.com/StephanHCB/go-backend-service-common/repository/logging"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const key1 = "key1"
const key2 = "key2"
const key3 = "key3"
const mapKey = "mapKey"

var testValues = map[string]string{
	key1: "value1",
	key2: "value2",
	key3: "value3",
}

func setupTest() *Impl {
	_ = auconfigenv.Setup(nil, nil)
	logger := logging.LoggingImpl{}
	logger.SetupForTesting()

	vaultSecretsConfig := createVaultSecretsConfig()

	cut := &Impl{
		Logging:            &logger,
		VaultClient:        mockVaultClientRequests(),
		VaultSecretsConfig: vaultSecretsConfig,
	}

	return cut
}

func createVaultSecretsConfig() map[string][]repository.VaultSecretConfig {
	simpleKey1 := key1
	mapKey2 := fmt.Sprintf("%s.%s", mapKey, key2)
	mapKey3 := fmt.Sprintf("%s.%s", mapKey, key3)

	vaultSecretsConfig := map[string][]repository.VaultSecretConfig{
		"path/to/secret": {
			{VaultKey: key1, ConfigKey: &simpleKey1},
			{VaultKey: key2, ConfigKey: &mapKey2},
		},
		"path/to/second/secret": {
			{VaultKey: key3, ConfigKey: &mapKey3},
		},
	}
	return vaultSecretsConfig
}

func mockVaultClientRequests() aurestclientapi.Client {
	client := aurestmock.New(map[string]aurestclientapi.ParsedResponse{
		"GET :///v1/system_kv/data/v1/path/to/secret <nil>": {
			Status: http.StatusOK,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: &SecretsResponse{
				Data: &SecretsResponseData{
					Data: map[string]string{
						key1: testValues[key1],
						key2: testValues[key2],
					},
				},
			},
		},
		"GET :///v1/system_kv/data/v1/path/to/second/secret <nil>": {
			Status: http.StatusOK,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: &SecretsResponse{
				Data: &SecretsResponseData{
					Data: map[string]string{
						key3: testValues[key3],
					},
				},
			},
		},
	}, map[string]error{})
	return client
}

func TestImpl_ObtainSecrets(t *testing.T) {
	cut := setupTest()

	assert.NoError(t, cut.ObtainSecrets(context.Background()))

	secretMap := map[string]string{}
	if assert.NoError(t, json.Unmarshal([]byte(auconfigenv.Get(mapKey)), &secretMap)) {
		assert.Equal(t, map[string]string{
			key2: testValues[key2],
			key3: testValues[key3],
		}, secretMap)
	}

	assert.Equal(t, testValues[key1], auconfigenv.Get(key1))
}

func Test_appendSecretToMap(t *testing.T) {
	type args struct {
		secretMapJson string
		secretKey     string
		secretValue   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty key creates new map and appends secrets",
			args: args{
				secretMapJson: "",
				secretKey:     "key1",
				secretValue:   "value1",
			},
			want:    "{\"key1\":\"value1\"}",
			wantErr: false,
		},
		{
			name: "appends key to existing map",
			args: args{
				secretMapJson: "{\"key1\":\"value1\"}",
				secretKey:     "key2",
				secretValue:   "value2",
			},
			want:    "{\"key1\":\"value1\",\"key2\":\"value2\"}",
			wantErr: false,
		},
		{
			name: "throws error on invalid json input and returns empty map",
			args: args{
				secretMapJson: "invalid",
				secretKey:     "key1",
				secretValue:   "value1",
			},
			want:    "{}",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := appendSecretToMap(tt.args.secretMapJson, tt.args.secretKey, tt.args.secretValue)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
