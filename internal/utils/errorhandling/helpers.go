package errorhandling

import (
	"net/http"

	"github.com/LeandroDeJesus-S/confectionery/internal/schemas"
	"github.com/LeandroDeJesus-S/confectionery/internal/utils/httphelpers"
)

// CheckOrHttpError checks if an error is not nil and writes an appropriate error message to the
// http.ResponseWriter with the given response code. If the error is nil, the function returns true.
// If the error is not nil, the function returns false.
//
// If the error is not nil and the msg parameter is not nil, it will write a JSON response with
// the given messages and the given response code.
//
// If the error is not nil and the msg parameter is nil, it will write a JSON response with the
// error message and the given response code.
func CheckOrHttpError(err error, w http.ResponseWriter, respCode int, msg ...string) bool {
	if err != nil {
		if msg != nil {
			messages := schemas.Message{Code: respCode, Detail: msg}
			httphelpers.JsonResponse(w, respCode, messages)
			return false
		}

		out := schemas.Message{
			Code:    respCode,
			Detail: []string{err.Error()},
		}
		httphelpers.JsonResponse(w, respCode, out)
		return false
	}
	return true
}
