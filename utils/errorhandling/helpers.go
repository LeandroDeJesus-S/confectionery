package errorhandling

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/LeandroDeJesus-S/confectionery/api/schemas"
)


// CheckOrHttpError checks if the error is not nil and if so, writes a http response
// with the given status code and a JSON response containing the error message.
// If the error is nil, it returns true. Otherwise, it returns false.
func CheckOrHttpError(err error, w http.ResponseWriter, respCode int, msg ...string) bool {
	if err != nil {
		w.WriteHeader(respCode)
		_msg := err.Error()
		if msg != nil {
			_msg = strings.Join(msg, " ")
		}

		out := schemas.Message{
			Code:    respCode,
			Message: _msg,
		}
		err := json.NewEncoder(w).Encode(out)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return false
	}
	return true
}
