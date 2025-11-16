package status_codes

type UpdateUserStatus int

func (u UpdateUserStatus) String() string {
	return UpdateUserStatusToString(u)
}

func (u UpdateUserStatus) Int() int {
	return int(u)
}

type ChangePasswordStatus int

func (u ChangePasswordStatus) String() string {
	return ChangePasswordStatusToString(u)
}

func (u ChangePasswordStatus) Int() int {
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

const (
	ChangePasswordSuccess ChangePasswordStatus = iota
	ChangePasswordInvalidUser
	ChangePasswordInvalidPassword
	ChangePasswordInvalidNewPassword
	ChangePasswordPasswordsDontMatch
	ChangePasswordAlreadyUsedPassword
	ChangePasswordError
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

func ChangePasswordStatusToString(code ChangePasswordStatus) string {
	switch code {
	case ChangePasswordSuccess:
		return "Senha alterada com sucesso!"
	case ChangePasswordInvalidUser:
		return "Usuário inválido"
	case ChangePasswordInvalidPassword:
		return "Senha atual inválida"
	case ChangePasswordInvalidNewPassword:
		return "Nova senha inválida"
	case ChangePasswordPasswordsDontMatch:
		return "As senhas não são iguais"
	case ChangePasswordAlreadyUsedPassword:
		return "A nova senha não pode ser igual à senha anterior"
	case ChangePasswordError:
		return "Erro ao atualizar a senha. Tente novamente mais tarde"
	default:
		return "Erro ao atualizar a senha. Tente novamente mais tarde"
	}
}
