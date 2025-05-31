package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LeandroDeJesus-S/confectionery/api/models"
	"github.com/LeandroDeJesus-S/confectionery/api/schemas"
	"github.com/LeandroDeJesus-S/confectionery/utils/errorhandling"
	"github.com/LeandroDeJesus-S/confectionery/utils/httphelpers"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type CustomerController struct {
	db        *gorm.DB
	validator *validator.Validate
}

// Initializes the CustomerController structure
func NewCustomerController(db *gorm.DB, validator *validator.Validate) *CustomerController {
	return &CustomerController{db: db, validator: validator}
}

// GetAllCustomers retrieves all customers from the database,
// converts them to the output schema, and encodes the result
// as a JSON response.
func (c *CustomerController) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	var dbCustomers []models.Customer
	var outputCustomers []schemas.CustomerOutputSchema

	c.db.Find(&dbCustomers)

	for _, dbCustomer := range dbCustomers {
		outputCustomer := schemas.CustomerOutputSchema{
			ID:    dbCustomer.ID,
			Fname: dbCustomer.Fname,
			Lname: dbCustomer.Lname,
			Email: dbCustomer.Email,
		}
		outputCustomers = append(outputCustomers, outputCustomer)
	}
	err := json.NewEncoder(w).Encode(outputCustomers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetCustomer retrieves a customer by ID from the database, converts it
// to the output schema, and encodes the result as a JSON response.
func (c *CustomerController) GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var dbCustomer models.Customer
	c.db.First(&dbCustomer, id)

	outputCustomer := schemas.CustomerOutputSchema{
		ID:    dbCustomer.ID,
		Fname: dbCustomer.Fname,
		Lname: dbCustomer.Lname,
		Email: dbCustomer.Email,
	}
	err := json.NewEncoder(w).Encode(outputCustomer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateCustomer creates a new customer in the database and returns it as a JSON response.
//
// If the request body is invalid or the email already exists, the function will return an appropriate
// HTTP status code and a JSON response with an error message.
//
// If the customer is successfully created, the function will return the created customer as a JSON
// response with the HTTP status code 201 Created.
func (c *CustomerController) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer *schemas.CustomerInputSchema
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	ok := errorhandling.CheckOrHttpError(c.validator.Struct(customer), w, http.StatusBadRequest)
	if !ok {
		return
	}

	var emailExists *models.Customer
	if c.db.First(&emailExists, "email = ?", customer.Email).RowsAffected > 0 {
		httphelpers.JsonResponse(
			w,
			http.StatusConflict,
			&schemas.Message{
				Code:    http.StatusConflict,
				Message: "Email already exists",
			},
		)
		return
	}

	dbCustomer := models.Customer{
		Fname:    customer.Fname,
		Lname:    customer.Lname,
		Email:    customer.Email,
		Password: customer.Password,
	}
	res := c.db.Create(&dbCustomer)
	if res.Error != nil {
		httphelpers.JsonResponse(
			w,
			http.StatusInternalServerError,
			&schemas.Message{
				Code:    http.StatusInternalServerError,
				Message: "Error creating customer: " + res.Error.Error(),
			},
		)
		return
	}

	customer = &schemas.CustomerInputSchema{
		Fname:    dbCustomer.Fname,
		Lname:    dbCustomer.Lname,
		Email:    dbCustomer.Email,
		Password: dbCustomer.Password,
	}
	if err := json.NewEncoder(w).Encode(customer); err != nil {
		httphelpers.JsonResponse(
			w,
			http.StatusInternalServerError,
			&schemas.Message{
				Code:    http.StatusInternalServerError,
				Message: "Error to encode customer: " + err.Error(),
			},
		)
		return
	}

	httphelpers.JsonResponse(
		w,
		http.StatusCreated,
		customer,
	)
}


// UpdateCustomer updates a customer by ID in the database and returns it as a JSON response.
//
// If the request body is invalid or the email already exists, the function will return an appropriate
// HTTP status code and a JSON response with an error message.
//
// If the customer is not found, the function will return a 404 Not Found response.
//
// If the customer is successfully updated, the function will return the updated customer as a JSON
// response with the HTTP status code 200 OK.
func (c *CustomerController) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.Atoi(vars["id"])

	isValidNumber := errorhandling.CheckOrHttpError(err, w, http.StatusBadRequest)
	if !isValidNumber {
		return
	}

	if customerID <= 0 {
		httphelpers.JsonResponse(
			w,
			http.StatusBadRequest,
			&schemas.Message{
				Code:    http.StatusBadRequest,
				Message: "Invalid customer ID",
			},
		)
		return
	}

	var dbCustomer models.Customer
	res := c.db.First(&dbCustomer, customerID)
	errMsg := &schemas.Message{
		Code:    http.StatusInternalServerError,
		Message: "Unexpected error",
	}

	switch res.Error {
	default:
		httphelpers.JsonResponse(
			w,
			http.StatusInternalServerError,
			errMsg,
		)
		return

	case nil:
		break

	case gorm.ErrRecordNotFound:
		errMsg.Code, errMsg.Message = http.StatusNotFound, "Customer not found"
		httphelpers.JsonResponse(w, http.StatusNotFound, errMsg)
		return
	}


	var input schemas.CustomerPatchInputSchema
	json.NewDecoder(r.Body).Decode(&input)

	if input.Fname != "" {
		dbCustomer.Fname = input.Fname
	}
	if input.Lname != "" {
		dbCustomer.Lname = input.Lname
	}
	if input.Email != "" {
		dbCustomer.Email = input.Email
	}

	res = c.db.Save(&dbCustomer)
	saveChecked := errorhandling.CheckOrHttpError(
		res.Error, w, http.StatusInternalServerError, 
		"Error updating customer: " + res.Error.Error(),
	)
	if !saveChecked {
		return
	}

	httphelpers.JsonResponse(
		w,
		http.StatusOK,
		&schemas.CustomerOutputSchema{
			ID:    dbCustomer.ID,
			Fname: dbCustomer.Fname,
			Lname: dbCustomer.Lname,
			Email: dbCustomer.Email,
		},
	)
}

// TODO: DeleteCustomer