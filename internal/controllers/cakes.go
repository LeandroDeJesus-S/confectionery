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

type CakeController struct {
	db        *gorm.DB
	validator *validator.Validate
}

func NewCakeController(db *gorm.DB, validator *validator.Validate) *CakeController {
	return &CakeController{db: db, validator: validator}
}

// GetCakes retrieves all the cakes from the database,
// converts them to the output schema, and encodes the result
// as a JSON response.
func (c *CakeController) GetCakes(w http.ResponseWriter, r *http.Request) {
	var dbCakes []models.Cake
	c.db.Find(&dbCakes)

	outputCakes := make([]schemas.CakeOutputSchema, 0)
	for _, dbCake := range dbCakes {
		outputCake := schemas.CakeOutputSchema{
			ID:    dbCake.ID,
			Name:  dbCake.Name,
			Price: dbCake.Price,
		}
		outputCakes = append(outputCakes, outputCake)
	}
	httphelpers.JsonResponse(w, http.StatusOK, outputCakes)
}

// GetCake retrieves a cake by ID from the database, converts it
// to the output schema, and encodes the result as a JSON response.
//
// If the ID is invalid, the function will return a 400 Bad Request response.
//
// If the customer is not found, the function will return a 404 Not Found response.
//
// If the customer is successfully retrieved, the function will return the customer details as a
// JSON response with a 200 OK status code. If there are any server errors, it
// returns a 500 Internal Server Error response.
func (c *CakeController) GetCake(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if !errorhandling.CheckOrHttpError(err, w, http.StatusBadRequest, "Invalid user id") {
		return
	}

	var dbCake models.Cake
	result := c.db.First(&dbCake, id)

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
				Detail: []string{"Cake not found"},
			},
		)
		return

	case nil:
		outputCake := schemas.CakeOutputSchema{
			ID:    dbCake.ID,
			Name:  dbCake.Name,
			Price: dbCake.Price,
		}
		httphelpers.JsonResponse(w, http.StatusOK, outputCake)
		return
	}
}

// CreateCake creates a new cake in the database and returns it as a JSON response.
//
// If the request body is invalid or the email already exists, the function will return an appropriate
// HTTP status code and a JSON response with an error message.
//
// If the customer is successfully created, the function will return the created customer as a JSON
// response with the HTTP status code 201 Created.
func (c *CakeController) CreateCake(w http.ResponseWriter, r *http.Request) {
	var inputCake schemas.CakeInputSchema
	hasDecoded := json.NewDecoder(r.Body).Decode(&inputCake)
	if !errorhandling.CheckOrHttpError(hasDecoded, w, http.StatusBadRequest, "Invalid input") {
		return
	}

	isValidStruct := c.validator.Struct(inputCake)
	if !errorhandling.CheckOrHttpError(isValidStruct, w, http.StatusBadRequest, "Invalid input") {
		return
	}

	found := c.db.First(&models.Cake{}, "name = ?", inputCake.Name)
	if found.RowsAffected > 0 {
		m := schemas.Message{Code: http.StatusBadRequest, Detail: []string{"Cake already exists"}}
		httphelpers.JsonResponse(w, http.StatusBadRequest, m)
		return
	}

	dbCake := models.Cake{
		Name:  inputCake.Name,
		Price: inputCake.Price,
	}

	result := c.db.Create(&dbCake)

	switch result.Error {
	case nil:
		outputCake := schemas.CakeOutputSchema{
			ID:    dbCake.ID,
			Name:  dbCake.Name,
			Price: dbCake.Price,
		}
		httphelpers.JsonResponse(w, http.StatusCreated, outputCake)

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

// UpdateCake updates a cake by ID in the database.
//
// It expects the request body to be a JSON object with optional "name" and "price" fields.
// If the request body is invalid, the function will return a 400 Bad Request response.
// If the ID is invalid, the function will return a 400 Bad Request response.
//
// If the customer is not found, the function will return a 404 Not Found response.
//
// If the customer is successfully updated, the function will return the updated customer details as a
// JSON response with a 200 OK status code. If there are any server errors, it
// returns a 500 Internal Server Error response.
func (c *CakeController) UpdateCake(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cakeId, err := strconv.ParseUint(vars["id"], 10, 32)

	castCheck := errorhandling.CheckOrHttpError(err, w, http.StatusBadRequest, "Invalid cake id")
	if !castCheck {
		return
	}

	var inputCake schemas.CakePatchInputSchema
	hasDecoded := json.NewDecoder(r.Body).Decode(&inputCake)
	if !errorhandling.CheckOrHttpError(hasDecoded, w, http.StatusBadRequest, "Invalid input") {
		return
	}

	var dbCake models.Cake

	duplicated := c.db.First(&dbCake, "name = ?", inputCake.Name)
	if duplicated.RowsAffected > 0 {
		m := schemas.Message{Code: http.StatusBadRequest, Detail: []string{"Cake already exists"}}
		httphelpers.JsonResponse(w, http.StatusBadRequest, m)
		return
	}

	result := c.db.First(&dbCake, cakeId)

	switch result.Error {
	case nil:
		break

	case gorm.ErrRecordNotFound:
		httphelpers.JsonResponse(
			w,
			http.StatusNotFound,
			schemas.Message{
				Code:   http.StatusNotFound,
				Detail: []string{"Cake not found"},
			},
		)
		return

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
	}

	c.db.Model(&dbCake).Updates(inputCake)

	out := schemas.CakeOutputSchema{
		ID:    dbCake.ID,
		Name:  dbCake.Name,
		Price: dbCake.Price,
	}
	httphelpers.JsonResponse(w, http.StatusOK, out)
}

// DeleteCake deletes a cake by ID from the database and returns a 204 No Content response.
//
// If the ID is invalid, the function will return a 400 Bad Request response.
//
// If the customer is not found, the function will return a 404 Not Found response.
//
// If the customer is successfully deleted, the function will return a 204 No Content response.
func (c *CakeController) DeleteCake(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cakeId, err := strconv.ParseUint(vars["id"], 10, 32)

	castCheck := errorhandling.CheckOrHttpError(err, w, http.StatusBadRequest, "Invalid cake id")
	if !castCheck {
		return
	}

	var dbCake models.Cake
	result := c.db.First(&dbCake, cakeId)

	switch result.Error {
	case nil:
		c.db.Delete(&dbCake)
		httphelpers.JsonResponse(w, http.StatusNoContent, nil)

	case gorm.ErrRecordNotFound:
		httphelpers.JsonResponse(
			w,
			http.StatusNotFound,
			schemas.Message{
				Code:   http.StatusNotFound,
				Detail: []string{"Cake not found"},
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
