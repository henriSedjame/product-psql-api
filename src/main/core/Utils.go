package core

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"io"
)

func IsValid(i interface{}) error {
	return validator.New().Struct(i)
}

func ToJson(i interface{}, reader io.Reader) error {
	return json.NewDecoder(reader).Decode(i)
}

func FromJson(i interface{}, writer io.Writer) error {
	return json.NewEncoder(writer).Encode(i)
}
