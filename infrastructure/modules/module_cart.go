package modules

type AddToCartRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type CartModule struct {
	name string
	path string
}

func NewCartModule() *CartModule {
	return &CartModule{
		name: "auth",
		path: "/auth",
	}
}

func (a CartModule) Name() string {
	return a.name
}

func (a CartModule) Path() string {
	return a.path
}

func (a CartModule) GetCart() error {
	user :=
}
