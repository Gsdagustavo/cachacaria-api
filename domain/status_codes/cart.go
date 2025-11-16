package status_codes

type AddProductItemStatus int

const (
	AddProductItemStatusSuccess AddProductItemStatus = iota
	AddProductItemStatusInvalidQuantity
	AddProductItemStatusInvalidProduct
	AddProductItemStatusInvalidUser
	AddProductItemStatusError
)

func (s AddProductItemStatus) String() string {
	switch s {
	case AddProductItemStatusSuccess:
		return "Sucesso"
	case AddProductItemStatusInvalidQuantity:
		return "Quantidade inválida"
	case AddProductItemStatusInvalidProduct:
		return "Produto inválido"
	case AddProductItemStatusInvalidUser:
		return "Usuário inválido"
	case AddProductItemStatusError:
		return "Erro interno no servidor"
	default:
		return "Erro desconhecido"
	}
}

func (s AddProductItemStatus) Int() int {
	return int(s)
}
