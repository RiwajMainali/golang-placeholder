package main

import (
	"encoding/json"
)

type Response struct {
	Message *string     `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func toJSON(r Response) string {
	jsonData, _ := json.Marshal(r)
	return string(jsonData)
}
