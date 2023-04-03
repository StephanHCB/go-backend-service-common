package idp

func (r *Impl) GetBasicAuthUsername() string {
	return r.BasicAuthUsername
}

func (r *Impl) GetBasicAuthPassword() string {
	return r.BasicAuthPassword
}

func (r *Impl) GetAuthOidcTokenAudience() string {
	return r.AuthOidcTokenAudience
}

func (r *Impl) GetAuthGroupWrite() string {
	return r.AuthGroupWrite
}

func (r *Impl) GetAuthGroupAdmin() string {
	return r.AuthGroupAdmin
}

func (r *Impl) GetAuthBasicUserGroup() string {
	return r.AuthBasicUserGroup
}

func (r *Impl) GetAuthorName() string {
	return r.AuthorName
}

func (r *Impl) GetAuthorEmail() string {
	return r.AuthorEmail
}
