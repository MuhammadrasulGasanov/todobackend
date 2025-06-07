package handler

import (
	"fmt"
	"net/http"

	"github.com/MuhammadrasulGasanov/go-tasks/internal/middleware"
)

func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey)
	fmt.Fprintf(w, "Hello user %v, you're authenticated!", userID)
}
