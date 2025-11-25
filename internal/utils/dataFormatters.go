package utils

import "encoding/json"

// dataToJson converts a slice of maps containing data into a JSON string
func DataToJson(data []map[string]interface{}) (string, error) {
	enco, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(enco), nil
}
