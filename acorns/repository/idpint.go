package repository

import "context"

const IdentityProviderAcornName = "idp"

// IdentityProvider is the central singleton representing an Open ID Connect Identity Provider.
//
// We use this to obtain a JWT keyset and to check its id endpoint to synchronously validate JWT tokens.
type IdentityProvider interface {
	IsIdentityProvider() bool

	// Setup uses the configuration to set up
	Setup(ctx context.Context) error

	// ObtainKeySet calls the key set endpoint and converts the keys to PEM for use with the jwt package
	ObtainKeySet(ctx context.Context) error

	// GetKeySet returns the previously obtained KeySet
	GetKeySet(ctx context.Context) []string

	// VerifyToken ensures synchronously that a token has not been revoked and the account is current.
	//
	// You should do this for critical operations that cannot live with the usual token
	// expiry cycle.
	VerifyToken(ctx context.Context, token string) error

	GetBasicAuthUsername() string
	GetBasicAuthPassword() string
	GetAuthOidcTokenAudience() string
	GetAuthGroupWrite() string
	GetAuthGroupAdmin() string
	GetAuthBasicUserGroup() string
	GetAuthorName() string
	GetAuthorEmail() string
}
