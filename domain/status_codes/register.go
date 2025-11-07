package status_codes

type RegisterStatusCode int

func (r RegisterStatusCode) String() string {
	return RegisterStatusCodeToString(r)
}

func (r RegisterStatusCode) Int() int {
	return int(r)
}

const (
	RegisterSuccess RegisterStatusCode = iota
	RegisterFailure
	RegisterUserAlreadyExist
	RegisterInvalidEmail
	RegisterInvalidName
	RegisterInvalidPassword
	RegisterInvalidCredentials
)

func RegisterStatusCodeToString(code RegisterStatusCode) string {
	switch code {
	case RegisterSuccess:
		return "Sucesso!"
	case RegisterFailure:
		return "Erro"
	case RegisterUserAlreadyExist:
		return "Usuário já existente"
	case RegisterInvalidEmail:
		return "Email inválido"
	case RegisterInvalidName:
		return "Nome inválido"
	case RegisterInvalidPassword:
		return "Senha inválida"
	case RegisterInvalidCredentials:
		return "Credenciais inválidas"
	default:
		return "UNKNOWN"
	}
}
