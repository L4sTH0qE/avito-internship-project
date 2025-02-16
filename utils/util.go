package utils

import (
	"awesomeProject/dto"
	"encoding/json"
	"net/http"
)

// RespondJSON отправляет ответ клиенту с телом в json-формате.
func RespondJSON(w http.ResponseWriter, status int, response interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Кодируем структуру в JSON и отправляем ответ.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RespondWithError отправляет ответ клиенту с описанием возникшей ошибки.
func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, dto.ErrorResponse{Errors: message})
}
