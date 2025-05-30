package schemas

// CustomerInputSchema is the schema for Customers creation
type CustomerInputSchema struct {
	Fname                string `json:"fName"`
	Lname                string `json:"lName"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

// CustomerOutputSchema represents the customer schema returned by the API
type CustomerOutputSchema struct {
	ID    uint   `json:"id"`
	Fname string `json:"fName"`
	Lname string `json:"lName"`
	Email string `json:"email"`
}
