package schemas

// Order represents the schema of an order made by a customer
type OrderInputSchema struct {
	CustomerID uint `json:"customerId"`
	CakeID     uint `json:"cakeId"`
	Qtd        uint `json:"qtd"`
	Delivered  bool `json:"delivered"`
}

// OrderPatchInputSchema is the schema for Orders update
type OrderPatchInputSchema struct {
	CustomerID uint `json:"customerId,omitempty"`
	CakeID     uint `json:"cakeId,omitempty"`
	Qtd        uint `json:"qtd,omitempty"`
	Delivered  bool `json:"delivered,omitempty"`
}
