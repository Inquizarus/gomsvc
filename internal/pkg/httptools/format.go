package httptools

import (
	"encoding/json"
)

func FormatJSON(container interface{}) ([]byte, error) {
	data, err := json.Marshal(container)
	if err != nil {
		return nil, err
	}
	return FormatJSONData(data)
}

func FormatJSONData(data []byte) ([]byte, error) {
	var err error

	container := map[string]interface{}{}

	if err = json.Unmarshal(data, &container); err != nil {
		return nil, err
	}

	return json.MarshalIndent(container, "", " ")
}
