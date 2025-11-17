package status_codes

type AddProductItemStatus int
type BuyProductsStatus int

const (
	AddProductItemStatusSuccess AddProductItemStatus = iota
	AddProductItemStatusInvalidQuantity
	AddProductItemStatusInvalidProduct
	AddProductItemStatusInvalidUser
	AddProductItemStatusError
)

const (
	BuyProductsStatusSuccess BuyProductsStatus = iota
	BuyProductsStatusCartEmpty
	BuyProductsStatusInvalidProduct
	BuyProductsStatusOutOfStock
	BuyProductsStatusError
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

func (s BuyProductsStatus) String() string {
	switch s {
	case BuyProductsStatusSuccess:
		return "Produtos comprados com sucesso!"
	case BuyProductsStatusCartEmpty:
		return "Carrinho vazio"
	case BuyProductsStatusInvalidProduct:
		return "Carrinho contém produtos inválidos"
	case BuyProductsStatusOutOfStock:
		return "Produtos sem estoque"
	case BuyProductsStatusError:
		return "Erro interno no servidor"
	default:
		return "Erro desconhecido"
	}
}

func (s BuyProductsStatus) Int() int {
	return int(s)
}
