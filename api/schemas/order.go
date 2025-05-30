package schemas

// Order represents the schema of an order made by a customer
type OrderInputSchema struct {
	CustomerID uint `json:"customerId"`
	CakeID     uint `json:"cakeId"`
	Qtd        uint `json:"qtd"`
	Delivered  bool `json:"delivered"`
}
