package response

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)

	return err
}
