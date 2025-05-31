package schemas

// Cake represents the schema of a Cake of the confectionery
type Cake struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}
