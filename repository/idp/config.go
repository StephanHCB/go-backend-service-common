package idp

import (
	"context"
	"fmt"
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
	auconfigenv "github.com/StephanHCB/go-autumn-config-env"
	"github.com/StephanHCB/go-backend-service-common/repository/config"
)

var ConfigItems = []auconfigapi.ConfigItem{
	{
		Key:         config.KeyAuthOidcKeySetUrl,
		EnvName:     config.KeyAuthOidcKeySetUrl,
		Default:     "",
		Description: "keyset url of oidc identity provider",
		Validate:    auconfigenv.ObtainPatternValidator("^https?:.*$"),
	},
	{
		Key:         config.KeyBasicAuthUsername,
		EnvName:     config.KeyBasicAuthUsername,
		Default:     "",
		Description: "username for basic-auth write access to this service",
		Validate:    auconfigenv.ObtainNotEmptyValidator(),
	},
	{
		Key:         config.KeyBasicAuthPassword,
		EnvName:     config.KeyBasicAuthPassword,
		Default:     "",
		Description: "password for basic-auth write access to this service",
		Validate:    auconfigenv.ObtainNotEmptyValidator(),
	},
	{
		Key:         config.KeyAuthOidcTokenAudience,
		EnvName:     config.KeyAuthOidcTokenAudience,
		Default:     "",
		Description: "expected audience of oidc access token",
		Validate:    auconfigenv.ObtainNotEmptyValidator(),
	},
	{
		Key:         config.KeyAuthGroupWrite,
		EnvName:     config.KeyAuthGroupWrite,
		Default:     "",
		Description: "group name or id for write access to this service",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         config.KeyAuthGroupAdmin,
		EnvName:     config.KeyAuthGroupAdmin,
		Default:     "",
		Description: "group name or id for admin access to this service",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         config.KeyAuthBasicUserGroup,
		EnvName:     config.KeyAuthBasicUserGroup,
		Default:     "",
		Description: "group name or id that is used on basic auth user access. defaults to write group.",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         config.KeyAuthorName,
		EnvName:     config.KeyAuthorName,
		Default:     "",
		Description: "name to use if authorised via basic auth",
		Validate:    auconfigenv.ObtainNotEmptyValidator(),
	},
	{
		Key:         config.KeyAuthorEmail,
		EnvName:     config.KeyAuthorEmail,
		Default:     "",
		Description: "email address to use if authorised via basic auth",
		Validate:    auconfigenv.ObtainNotEmptyValidator(),
	},
}

func (v *Impl) Validate(ctx context.Context) error {
	var errorList = make([]error, 0)
	for _, it := range ConfigItems {
		if it.Validate != nil {
			err := it.Validate(it.Key)
			if err != nil {
				v.Logging.Logger().Ctx(ctx).Warn().WithErr(err).Printf("failed to validate configuration field %s", it.EnvName)
				errorList = append(errorList, err)
			}
		}
	}

	if len(errorList) > 0 {
		return fmt.Errorf("some configuration values failed to validate or parse. There were %d error(s). See details above", len(errorList))
	} else {
		return nil
	}
}

func (v *Impl) Obtain(ctx context.Context) {
	v.AuthOidcKeySetUrl = auconfigenv.Get(config.KeyAuthOidcKeySetUrl)
	v.AuthOidcTokenAudience = auconfigenv.Get(config.KeyAuthOidcTokenAudience)
	v.AuthGroupWrite = auconfigenv.Get(config.KeyAuthGroupWrite)
	v.AuthGroupAdmin = auconfigenv.Get(config.KeyAuthGroupAdmin)
	v.AuthBasicUserGroup = auconfigenv.Get(config.KeyAuthBasicUserGroup)
	v.BasicAuthUsername = auconfigenv.Get(config.KeyBasicAuthUsername)
	v.BasicAuthPassword = auconfigenv.Get(config.KeyBasicAuthPassword)
	v.AuthorName = auconfigenv.Get(config.KeyAuthorName)
	v.AuthorEmail = auconfigenv.Get(config.KeyAuthorEmail)
}
