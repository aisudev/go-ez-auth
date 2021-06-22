package utils

import (
	"fmt"
)

func Response(success bool, message *string, data interface{}, err error) map[string]interface{} {

	resMap := map[string]interface{}{
		"success": success,
	}

	if message != nil {
		resMap["message"] = message
	}

	if data != nil {
		resMap["data"] = data
	}

	if err != nil {
		resMap["error"] = err.Error()
	}

	return resMap
}

func Logger(message string, err error) {
	fmt.Printf("{\nmessage: %s\nerror: %s\n", message, err.Error())
}
