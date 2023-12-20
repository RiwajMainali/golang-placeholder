package main
import ("encoding/json")
func StrPtr(s string) *string {
	return &s
}

func toJSON(r Response) string {
	jsonData, _ := json.Marshal(r)
	return string(jsonData)
}
