package status_codes

type UpdateUserStatus int

func (u UpdateUserStatus) String() string {
	return UpdateUserStatusToString(u)
}

func (u UpdateUserStatus) Int() int {
	return int(u)
}

const (
	UpdateUserSuccess UpdateUserStatus = iota
	UpdateUserInvalidUser
	UpdateUserInvalidPassword
	UpdateUserInvalidEmail
	UpdateUserInvalidPhone
	UpdateUserFailure
)

func UpdateUserStatusToString(code UpdateUserStatus) string {
	switch code {
	case UpdateUserSuccess:
		return "Usuário atualizado com sucesso!"
	case UpdateUserInvalidUser:
		return "Usuário inválido"
	case UpdateUserInvalidPassword:
		return "Senha inválida"
	case UpdateUserInvalidEmail:
		return "Email inválido"
	case UpdateUserInvalidPhone:
		return "Telefone inválido"
	case UpdateUserFailure:
		return "Erro ao atualizar usuário"
	default:
		return "UNKNOWN"
	}
}
