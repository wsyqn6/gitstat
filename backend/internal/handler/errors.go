package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

const (
	ErrCodeInvalidRequest = "INVALID_REQUEST"
	ErrCodePathRequired   = "PATH_REQUIRED"
	ErrCodeRepoNotFound   = "REPO_NOT_FOUND"
	ErrCodeInternalError  = "INTERNAL_ERROR"
)

func writeError(w http.ResponseWriter, code string, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message, Code: code})
}
