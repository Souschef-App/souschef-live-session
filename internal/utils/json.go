package utils

import (
	"bytes"
	"encoding/json"
)

func DecodeJSON(data []byte, cmd interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(cmd); err != nil {
		return err
	}

	return nil
}
