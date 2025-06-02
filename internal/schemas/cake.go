package schemas


type CakeInputSchema struct {
	Name  string `json:"name" validate:"required"`
	Price uint64  `json:"price" validate:"required"`
}

type CakePatchInputSchema struct {
	Name  string `json:"name,omitempty"`
	Price *uint64  `json:"price,omitempty"`
}

type CakeOutputSchema struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price uint64  `json:"price"`
}
