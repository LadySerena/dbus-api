package common

import "fmt"

const (
	JsonContent          = "application/json"
	ContentTypeHeaderKey = "Content-Type"
)

func GenerateErrorResponse(code int32, text string, err error) []byte {
	return []byte(fmt.Sprintf(`{"status": %d, "reason":"%s due to: %s"}`, code, text, err.Error()))
}

func GenerateProblemResponse(code int32, text string) []byte {
	return []byte(fmt.Sprintf(`{"status": %d, "reason":"%s"}`, code, text))
}
