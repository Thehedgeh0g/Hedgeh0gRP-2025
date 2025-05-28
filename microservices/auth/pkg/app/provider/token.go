package provider

type TokenProvider interface {
	GetTokenByLogin(login string) (string, error)
}
