package responses

import (
	"api/src/apierrors"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Fatal(err)
	}
}

func Error(w http.ResponseWriter, err error) {
	var apiError *apierrors.APIError

	switch {
		case errors.As(err, &apiError):
			JSON(w, apiError.Status, apiError)
		default:
			fmt.Println("unhandled error: ", err)
			JSON(w, http.StatusInternalServerError, &apierrors.APIError{
				Message: "Unexpected error",
			})
	}
}
