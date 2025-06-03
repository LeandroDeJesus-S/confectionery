package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/LeandroDeJesus-S/confectionery/internal/models"
	"github.com/LeandroDeJesus-S/confectionery/internal/schemas"
	"github.com/LeandroDeJesus-S/confectionery/internal/utils/errorhandling"
	"github.com/LeandroDeJesus-S/confectionery/internal/utils/httphelpers"
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

// GetAllCustomers retrieves all the active customers from the database,
// converts them to the output schema, and encodes the result
// as a JSON response.
func (c *CustomerController) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	var dbCustomers []models.Customer
	outputCustomers := make([]schemas.CustomerOutputSchema, 0)

	c.db.Find(&dbCustomers, "active = ?", true)

	for _, dbCustomer := range dbCustomers {
		outputCustomer := schemas.CustomerOutputSchema{
			ID:    dbCustomer.ID,
			Fname: dbCustomer.Fname,
			Lname: dbCustomer.Lname,
			Email: dbCustomer.Email,
		}
		outputCustomers = append(outputCustomers, outputCustomer)
	}

	httphelpers.JsonResponse(
		w,
		http.StatusOK,
		outputCustomers,
	)
}

// GetCustomer retrieves an active customer by ID from the database, converts it
// to the output schema, and encodes the result as a JSON response.
func (c *CustomerController) GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if !errorhandling.CheckOrHttpError(err, w, http.StatusBadRequest, "Invalid user id") {
		return
	}

	var dbCustomer models.Customer
	result := c.db.First(&dbCustomer, "id = ? AND active = ?", id, true)

	switch result.Error {
	default:
		httphelpers.JsonResponse(
			w,
			http.StatusInternalServerError,
			schemas.Message{
				Code:   http.StatusInternalServerError,
				Detail: []string{"Internal server error"},
			},
		)
		return

	case gorm.ErrRecordNotFound:
		httphelpers.JsonResponse(
			w,
			http.StatusNotFound,
			schemas.Message{
				Code:   http.StatusNotFound,
				Detail: []string{"Customer not found"},
			},
		)
		return

	case nil:
		break
	}

	outputCustomer := schemas.CustomerOutputSchema{
		ID:    dbCustomer.ID,
		Fname: dbCustomer.Fname,
		Lname: dbCustomer.Lname,
		Email: dbCustomer.Email,
	}
	httphelpers.JsonResponse(
		w,
		http.StatusOK,
		outputCustomer,
	)
}

// CreateCustomer creates a new customer in the database and returns it as a JSON response.
//
// If the request body is invalid or the email already exists, the function will return an appropriate
// HTTP status code and a JSON response with an error message.
//
// If the customer is successfully created, the function will return the created customer as a JSON
// response with the HTTP status code 201 Created.
func (c *CustomerController) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var inpCustomer schemas.CustomerInputSchema
	err := json.NewDecoder(r.Body).Decode(&inpCustomer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	err = c.validator.Struct(inpCustomer)
	if err != nil {
		httphelpers.JsonResponse(
			w,
			http.StatusBadRequest,
			&schemas.Message{
				Code:   http.StatusBadRequest,
				Detail: strings.Split(err.Error(), "\n"),
			},
		)
		return
	}

	var emailExists *models.Customer
	if c.db.First(&emailExists, "email = ?", inpCustomer.Email).RowsAffected > 0 {
		httphelpers.JsonResponse(
			w,
			http.StatusConflict,
			&schemas.Message{
				Code:   http.StatusConflict,
				Detail: []string{"Email already exists"},
			},
		)
		return
	}

	dbCustomer := models.Customer{
		Fname: inpCustomer.Fname,
		Lname: inpCustomer.Lname,
		Email: inpCustomer.Email,
	}
	res := c.db.Create(&dbCustomer)
	if res.Error != nil {
		httphelpers.JsonResponse(
			w,
			http.StatusInternalServerError,
			&schemas.Message{
				Code:   http.StatusInternalServerError,
				Detail: []string{"Error creating customer: " + res.Error.Error()},
			},
		)
		return
	}

	outCustomer := &schemas.CustomerOutputSchema{
		ID:    dbCustomer.ID,
		Fname: dbCustomer.Fname,
		Lname: dbCustomer.Lname,
		Email: dbCustomer.Email,
	}

	httphelpers.JsonResponse(
		w,
		http.StatusCreated,
		outCustomer,
	)
}

// UpdateCustomer updates an existing customer's details by ID from the database.
// It parses the customer ID from the URL, verifies its validity, and retrieves
// the existing customer record. If the customer is not found, it returns a 404
// Not Found response. The function then decodes the request body for partial
// updates, validates the input, and checks for unique email constraints.
// If any validation fails, it sends a 400 Bad Request response. Upon successful
// update, it saves the changes and returns the updated customer details as a
// JSON response with a 200 OK status code. If there are any server errors, it
// returns a 500 Internal Server Error response.
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
				Code:   http.StatusBadRequest,
				Detail: []string{"Invalid customer ID"},
			},
		)
		return
	}

	var dbCustomer models.Customer
	res := c.db.First(&dbCustomer, customerID)
	errMsg := &schemas.Message{
		Code:   http.StatusInternalServerError,
		Detail: []string{"Unexpected error"},
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
		errMsg.Code, errMsg.Detail = http.StatusNotFound, []string{"Customer not found"}
		httphelpers.JsonResponse(w, http.StatusNotFound, errMsg)
		return
	}

	var input schemas.CustomerPatchInputSchema
	deacoded := json.NewDecoder(r.Body).Decode(&input)

	if !errorhandling.CheckOrHttpError(deacoded, w, http.StatusBadRequest, "Invalid input") {
		return
	}

	updated := c.db.Model(&dbCustomer).Updates(input)
	if !errorhandling.CheckOrHttpError(
		updated.Error, w, http.StatusInternalServerError, "Error updating customer",
	) {
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

// DeleteCustomer deletes a customer by ID changing its active status and returns a 204 No Content response.
//
// If the ID is invalid, the function will return a 400 Bad Request response.
//
// If the customer is not found, the function will return a 404 Not Found response.
//
// If the customer is successfully deleted, the function will return a 204 No Content response.
func (c *CustomerController) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId, err := strconv.Atoi(vars["id"])
	castCheck := errorhandling.CheckOrHttpError(
		err,
		w,
		http.StatusBadRequest,
		"Invalid user id",
	)

	if !castCheck {
		return
	}

	var dbCustomer models.Customer
	result := c.db.First(&dbCustomer, "id = ? AND active = ?", customerId, true)

	switch result.Error {
	case gorm.ErrRecordNotFound:
		httphelpers.JsonResponse(
			w,
			http.StatusNotFound,
			&schemas.Message{
				Code:   http.StatusNotFound,
				Detail: []string{"Customer not found"},
			},
		)
		return

	case nil:
		break

	default:
		httphelpers.JsonResponse(
			w,
			http.StatusInternalServerError,
			&schemas.Message{
				Code:   http.StatusInternalServerError,
				Detail: []string{"Something went wrong"},
			},
		)
		return
	}

	dbCustomer.Active = false
	c.db.Save(&dbCustomer)
	httphelpers.JsonResponse(w, http.StatusNoContent, nil)
}
