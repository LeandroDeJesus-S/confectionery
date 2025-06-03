package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LeandroDeJesus-S/confectionery/internal/models"
	"github.com/LeandroDeJesus-S/confectionery/internal/schemas"
	"github.com/LeandroDeJesus-S/confectionery/internal/utils/errorhandling"
	"github.com/LeandroDeJesus-S/confectionery/internal/utils/httphelpers"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type OrdersController struct {
	db        *gorm.DB
	validator *validator.Validate
}

func NewOrdersController(db *gorm.DB, validator *validator.Validate) *OrdersController {
	return &OrdersController{db: db, validator: validator}
}

// GetOrders retrieves all the orders from the database and encodes them
// as a JSON response with an HTTP status code 200 OK.
func (c *OrdersController) GetOrders(w http.ResponseWriter, r *http.Request) {
	var dbOrders []models.Order
	c.db.Find(&dbOrders)
	httphelpers.JsonResponse(w, http.StatusOK, dbOrders)
}

// GetOrder retrieves an order by ID from the database and encodes it
// as a JSON response with an HTTP status code 200 OK.
//
// If the ID is invalid, the function will return a 400 Bad Request response.
//
// If the order is not found, the function will return a 404 Not Found response.
//
// If the order is successfully retrieved, the function will return the order details as a
// JSON response with a 200 OK status code. If there are any server errors, it
// returns a 500 Internal Server Error response.
func (c *OrdersController) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)

	if !errorhandling.CheckOrHttpError(err, w, http.StatusBadRequest, "Invalid user id") {
		return
	}

	var dbOrder models.Order
	result := c.db.First(&dbOrder, id)

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

	case gorm.ErrRecordNotFound:
		httphelpers.JsonResponse(
			w,
			http.StatusNotFound,
			schemas.Message{
				Code:   http.StatusNotFound,
				Detail: []string{"Order not found"},
			},
		)

	case nil:
		httphelpers.JsonResponse(w, http.StatusOK, dbOrder)
	}
}

// CreateOrder creates a new order in the database and returns it as a JSON response.
//
// If the request body is invalid or the email already exists, the function will return an appropriate
// HTTP status code and a JSON response with an error message.
//
// If the order is successfully created, the function will return the created order as a JSON
// response with the HTTP status code 201 Created.
func (c *OrdersController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var inputOrder schemas.OrderInputSchema
	hasDecoded := json.NewDecoder(r.Body).Decode(&inputOrder)
	if !errorhandling.CheckOrHttpError(hasDecoded, w, http.StatusBadRequest, "Invalid input") {
		return
	}

	isValidStruct := c.validator.Struct(inputOrder)
	if !errorhandling.CheckOrHttpError(isValidStruct, w, http.StatusBadRequest, "Invalid input") {
		return
	}

	customerExists := c.db.First(&models.Customer{}, "id = ?", inputOrder.CustomerID).RowsAffected > 0
	cakeExists := c.db.First(&models.Cake{}, "id = ?", inputOrder.CakeID).RowsAffected > 0
	if !customerExists || !cakeExists {
		m := schemas.Message{Code: http.StatusBadRequest, Detail: []string{"Customer or Cake not found"}}
		httphelpers.JsonResponse(w, http.StatusBadRequest, m)
		return
	}

	dbOrder := models.Order{
		CustomerID: inputOrder.CustomerID,
		CakeID:     inputOrder.CakeID,
		Qtd:        inputOrder.Qtd,
		Delivered:  inputOrder.Delivered,
	}

	result := c.db.Create(&dbOrder)

	switch result.Error {
	case nil:
		httphelpers.JsonResponse(w, http.StatusCreated, dbOrder)

	case gorm.ErrDuplicatedKey:
		m := schemas.Message{Code: http.StatusBadRequest, Detail: []string{"Order already exists"}}
		httphelpers.JsonResponse(w, http.StatusBadRequest, m)

	default:
		httphelpers.JsonResponse(
			w,
			http.StatusInternalServerError,
			schemas.Message{
				Code:   http.StatusInternalServerError,
				Detail: []string{"Internal server error"},
			},
		)
	}
}

// UpdateOrder updates an order by ID in the database.
// It parses the order ID from the URL, verifies its validity, and retrieves
// the existing order record. If the order is not found, it returns a 404
// Not Found response. The function then decodes the request body for partial
// updates, validates the input, and checks for unique cake and customer
// constraints. If any validation fails, it sends a 400 Bad Request response.
// Upon successful update, it saves the changes and returns the updated order
// details as a JSON response with a 200 OK status code. If there are any
// server errors, it returns a 500 Internal Server Error response.
func (c *OrdersController) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	OrderId, err := strconv.ParseUint(vars["id"], 10, 32)

	castCheck := errorhandling.CheckOrHttpError(err, w, http.StatusBadRequest, "Invalid Order id")
	if !castCheck {
		return
	}

	var inputOrder schemas.OrderPatchInputSchema
	hasDecoded := json.NewDecoder(r.Body).Decode(&inputOrder)
	if !errorhandling.CheckOrHttpError(hasDecoded, w, http.StatusBadRequest, "Invalid input") {
		return
	}

	var dbOrder models.Order
	result := c.db.First(&dbOrder, OrderId)

	isUpdatingCustomerOrCake := inputOrder.CakeID != 0 || inputOrder.CustomerID != 0
	if isUpdatingCustomerOrCake {
		cakeExists := c.db.First(&models.Cake{}, "id = ?", inputOrder.CakeID).RowsAffected > 0
		customerExists := c.db.First(&models.Customer{}, "id = ?", inputOrder.CustomerID).RowsAffected > 0

		if !cakeExists || !customerExists {
			m := schemas.Message{Code: http.StatusBadRequest, Detail: []string{"Customer or Cake not found"}}
			httphelpers.JsonResponse(w, http.StatusBadRequest, m)
			return
		}
	}

	switch result.Error {
	case nil:
		c.db.Model(&dbOrder).Updates(inputOrder)
		httphelpers.JsonResponse(w, http.StatusOK, dbOrder)

	case gorm.ErrRecordNotFound:
		httphelpers.JsonResponse(
			w,
			http.StatusNotFound,
			schemas.Message{
				Code:   http.StatusNotFound,
				Detail: []string{"Order not found"},
			},
		)

	default:
		httphelpers.JsonResponse(
			w,
			http.StatusInternalServerError,
			schemas.Message{
				Code:   http.StatusInternalServerError,
				Detail: []string{"Internal server error"},
			},
		)
	}
}

// DeleteOrder deletes an order by ID and returns a 204 No Content response.
//
// If the ID is invalid, the function will return a 400 Bad Request response.
//
// If the order is not found, the function will return a 404 Not Found response.
//
// If the order is successfully deleted, the function will return a 204 No Content response.
func (c *OrdersController) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	OrderId, err := strconv.ParseUint(vars["id"], 10, 32)

	castCheck := errorhandling.CheckOrHttpError(err, w, http.StatusBadRequest, "Invalid Order id")
	if !castCheck {
		return
	}

	var dbOrder models.Order
	result := c.db.First(&dbOrder, OrderId)

	switch result.Error {
	case nil:
		c.db.Delete(&dbOrder)
		httphelpers.JsonResponse(w, http.StatusNoContent, nil)

	case gorm.ErrRecordNotFound:
		httphelpers.JsonResponse(
			w,
			http.StatusNotFound,
			schemas.Message{
				Code:   http.StatusNotFound,
				Detail: []string{"Order not found"},
			},
		)

	default:
		httphelpers.JsonResponse(
			w,
			http.StatusInternalServerError,
			schemas.Message{
				Code:   http.StatusInternalServerError,
				Detail: []string{"Internal server error"},
			},
		)
	}
}
