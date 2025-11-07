package status_codes

type LoginStatusCode int

func (l LoginStatusCode) String() string {
	return LoginStatusCodeToString(l)
}

func (l LoginStatusCode) Int() int {
	return int(l)
}

const (
	LoginSuccess LoginStatusCode = iota
	LoginFailure
	LoginInvalidCredentials
	LoginUserNotFound
)

func LoginStatusCodeToString(code LoginStatusCode) string {
	switch code {
	case LoginSuccess:
		return "Sucesso"
	case LoginFailure:
		return "Erro"
	case LoginInvalidCredentials:
		return "Credenciais inválidas"
	case LoginUserNotFound:
		return "Usuário não encontrado"
	default:
		return "Erro desconhecido"
	}
}
