package main

import (
	"net/http"
	"fmt"
)

type Response struct {
	Message *string     `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func SendErrorMessage(text string, w http.ResponseWriter) {
	messege := Response{
		Message: StrPtr(text),
		Success: false,
		Data:    nil,
	}
	w.Header().Set("Content-type", "application/json")
	fmt.Fprintf(w, toJSON(messege))
}

