package status_codes

type AddProductStatus int

const (
	AddProductStatusSuccess AddProductStatus = iota
	AddProductStatusInvalidName
	AddProductStatusInvalidPrice
	AddProductStatusInvalidStock
	AddProductStatusError
)

func (s AddProductStatus) String() string {
	switch s {
	case AddProductStatusSuccess:
		return "Successo"
	case AddProductStatusInvalidName:
		return "Nome inválido"
	case AddProductStatusInvalidPrice:
		return "Preço inválido"
	case AddProductStatusInvalidStock:
		return "Estoque inválido"
	case AddProductStatusError:
		return "Erro interno no servidor"
	default:
		return "Erro desconhecido"
	}
}

func (s AddProductStatus) Int() int {
	return int(s)
}

type DeleteProductStatus int

const (
	DeleteProductStatusSuccess DeleteProductStatus = iota
	DeleteProductStatusNotFound
	DeleteProductStatusError
)

func (s DeleteProductStatus) String() string {
	switch s {
	case DeleteProductStatusSuccess:
		return "Successo"
	case DeleteProductStatusNotFound:
		return "Produto não encontrado"
	case DeleteProductStatusError:
		return "Erro interno no servidor"
	default:
		return "Erro desconhecido"
	}
}

func (s DeleteProductStatus) Int() int {
	return int(s)
}

type UpdateProductStatus int

const (
	UpdateProductStatusSuccess UpdateProductStatus = iota
	UpdateProductStatusNotFound
	UpdateProductStatusInvalidName
	UpdateProductStatusInvalidPrice
	UpdateProductStatusInvalidDescription
	UpdateProductStatusInvalidStock
	UpdateProductStatusError
)

func (s UpdateProductStatus) String() string {
	switch s {
	case UpdateProductStatusSuccess:
		return "Successo"
	case UpdateProductStatusNotFound:
		return "Produto não encontrado"
	case UpdateProductStatusInvalidName:
		return "Nome inválido"
	case UpdateProductStatusInvalidPrice:
		return "Preço inválido"
	case UpdateProductStatusInvalidStock:
		return "Estoque inválido"
	case UpdateProductStatusInvalidDescription:
		return "Descrição inválida"
	case UpdateProductStatusError:
		return "Erro interno no servidor"
	default:
		return "Erro desconhecido"
	}
}

func (s UpdateProductStatus) Int() int {
	return int(s)
}
