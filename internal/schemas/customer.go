package schemas

// CustomerInputSchema is the schema for Customers creation
type CustomerInputSchema struct {
	Fname                string `json:"fName" validate:"required"`
	Lname                string `json:"lName" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
}

// CustomerOutputSchema represents the customer schema returned by the API
type CustomerOutputSchema struct {
	ID    uint   `json:"id"`
	Fname string `json:"fName"`
	Lname string `json:"lName"`
	Email string `json:"email"`
}

// CustomerPatchInputSchema is the schema for Customers update
type CustomerPatchInputSchema struct {
	Fname                string `json:"fName,omitempty"`
	Lname                string `json:"lName,omitempty"`
	Email                string `json:"email,omitempty" validate:"email"`
}