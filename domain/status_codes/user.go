package status_codes

type UserStatusCode int

func (u UserStatusCode) String() string {
	return UserStatusCodeToString(u)
}

func (u UserStatusCode) Int() int {
	return int(u)
}

const (
	UserUpdateSuccess UserStatusCode = iota
	UserUpdateFailure
	UserInvalidCredentials
	UserDeleteSuccess
	UserDeleteFailure
)

func UserStatusCodeToString(code UserStatusCode) string {
	switch code {
	case UserUpdateSuccess:
		return "Sucesso"
	case UserDeleteSuccess:
		return "Sucesso"
	case UserInvalidCredentials:
		return "Credenciais inválidas"
	case UserDeleteFailure:
		return "Erro interno"
	case UserUpdateFailure:
		return "Erro interno"
	default:
		return "Erro desconhecido"
	}
}
