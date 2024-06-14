package util

import (
	"bytes"
	"encoding/json"
	"io"
)

func UnpackJSON(obj io.Reader, container any) error {
	bodyBytes, err := io.ReadAll(obj)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, container)
	if err != nil {
		return err
	}
	return nil
}

func PackJSON(obj any) (*bytes.Buffer, error) {
	jsonBody, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}
