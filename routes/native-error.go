package routes

import (
	"errors"
	"net/http"

	"github.com/wisdom-oss/common-go/middleware"
)

func BasicWithErrorHandling(w http.ResponseWriter, r *http.Request) {
	// access the error handlers
	errorHandler := r.Context().Value(middleware.ErrorChannelName).(chan<- interface{})
	// now publish an error to each of the wisdom errors
	errorHandler <- errors.New("native test error")
	// now wait for the error to be handled
	return
}
