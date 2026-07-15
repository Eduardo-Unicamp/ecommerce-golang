package handler

import (
	"database/sql"
	"errors"
	"first-api/internal/model"
	"log"
	"net/http"
)

func WriteOrderError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrEmptyOrder):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, model.ErrInsufficientStock):
		http.Error(w, err.Error(), http.StatusConflict)

	case errors.Is(err, model.ErrInvalidStockQuantity):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, model.ErrEmailTaken):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, model.ErrInvalidField):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, model.ErrUnableToCancel) || errors.Is(err, model.ErrUnableToPay):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, sql.ErrNoRows) || errors.Is(err, model.CustomerNotFound) || errors.Is(err, model.ErrOrderNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, model.ErrInvalidPassword) || errors.Is(err, model.ErrReadingJSON):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, model.ErrInvalidRefreshToken) || errors.Is(err, model.ErrRefreshTokenRequired) || errors.Is(err, model.ErrAuthorizationFailed):
		http.Error(w, err.Error(), http.StatusUnauthorized)

	default:
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
