package main
import ("encoding/json")
func StrPtr(s string) *string {
	return &s
}
// This is bad, really bad. But anything other than this will be annoying
func toJSON(r Response) string {
	jsonData, err := json.Marshal(r)
	if err!=nil{
		return 	
	}
	return string(jsonData)

