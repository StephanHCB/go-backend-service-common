package vault

import (
	"context"
	"errors"
	"github.com/StephanHCB/go-autumn-acorn-registry/api"
	auzerolog "github.com/StephanHCB/go-autumn-logging-zerolog"
	"github.com/StephanHCB/go-backend-service-common/acorns/repository"
)

// --- implementing Acorn ---

func New() auacornapi.Acorn {
	return &Impl{
		VaultProtocol: "https",
	}
}

// NewNoAcorn wires up the component, but does not set it up.
//
// You will still need to call Setup(), which injects configuration values before they are validated.
// This means, you must call Setup() for vault before you call Setup() for the configuration.
//
// All in all, the call order with Vault needs to be:
//   - c := config.NewNoAcorn(...)
//   - l := logging.NewNoAcorn(c)
//   - v := vault.NewNoAcorn(c, l)
//   - c.Assemble(l)
//   - l.Setup()
//   - vault.Execute(v)
//   - c.Setup()
func NewNoAcorn(configuration repository.Configuration, logging repository.Logging) repository.Vault {
	return &Impl{
		VaultProtocol: "https",

		Configuration: configuration,
		Logging:       logging,
	}
}

func (v *Impl) IsVault() bool {
	return true
}

func (v *Impl) AcornName() string {
	return repository.VaultAcornName
}

func (v *Impl) AssembleAcorn(registry auacornapi.AcornRegistry) error {
	v.Configuration = registry.GetAcornByName(repository.ConfigurationAcornName).(repository.Configuration)
	v.Logging = registry.GetAcornByName(repository.LoggingAcornName).(repository.Logging)

	return registry.AddSetupOrderRule(v, v.Configuration.(auacornapi.Acorn))
}

func (v *Impl) SetupAcorn(registry auacornapi.AcornRegistry) error {
	if err := registry.SetupAfter(v.Logging.(auacornapi.Acorn)); err != nil {
		return err
	}

	return Execute(v)
}

// setup convenience function for no acorn setup

func Execute(vaultComponent repository.Vault) error {
	ctx := auzerolog.AddLoggerToCtx(context.Background())

	v, ok := vaultComponent.(*Impl)
	if !ok {
		v.Logging.Logger().Ctx(ctx).Error().Print("received invalid component as vault client. You can only call Execute() on a vault instance. BAILING OUT")
		return errors.New("not a vault component instance - cannot run full setup on mocks")
	}

	if err := v.Validate(ctx); err != nil {
		return err
	}
	v.Obtain(ctx)

	if !v.VaultEnabled {
		v.Logging.Logger().Ctx(ctx).Info().Print("vault disabled, local values will be used.")
		return nil
	}

	if err := v.Setup(ctx); err != nil {
		v.Logging.Logger().Ctx(ctx).Error().WithErr(err).Print("failed to set up vault client. BAILING OUT")
		return err
	}
	if err := v.Authenticate(ctx); err != nil {
		v.Logging.Logger().Ctx(ctx).Error().WithErr(err).Print("failed to authenticate to vault. BAILING OUT")
		return err
	}
	if err := v.ObtainSecrets(ctx); err != nil {
		v.Logging.Logger().Ctx(ctx).Error().WithErr(err).Print("failed to get secrets from vault. BAILING OUT")
		return err
	}
	v.Logging.Logger().Ctx(ctx).Info().Print("successfully obtained vault secrets")
	return nil
}

func (v *Impl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	return nil
}
