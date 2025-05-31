package httphelpers

import (
	"encoding/json"
	"net/http"
)

// JsonResponse writes a JSON response to the HTTP response writer
// with the given response code and data.
func JsonResponse(w http.ResponseWriter, respCode int, data any) {
	w.WriteHeader(respCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
