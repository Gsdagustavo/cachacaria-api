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
		return "SUCCESS"
	case UserDeleteSuccess:
		return "SUCCESS"
	case UserInvalidCredentials:
		return "INVALID_CREDENTIALS"
	case UserDeleteFailure:
		return "DELETE_FAILURE"
	case UserUpdateFailure:
		return "UPDATE_FAILURE"
	default:
		return "UNKNOWN"
	}
}
