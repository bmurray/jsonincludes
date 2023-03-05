package jsonincludes

import (
	"encoding/json"
	"io"
	"os"
)

type JsonData interface {
	Data(string) ([]byte, error)
}
type JsonReader interface {
	Reader(string) (io.Reader, error)
}
type JsonInclude[T any] struct {
	Val T `json:"-"`
}
type jsonInclude struct {
	Include string `json:"include"`
}

type JsonError string

const (
	ErrNoData        JsonError = "no data"
	ErrNoReader      JsonError = "no reader"
	ErrCannotInclude JsonError = "cannot include"
)

func (j *JsonInclude[T]) UnmarshalJSON(data []byte) error {

	var x interface{} = j.Val
	var inc jsonInclude

	err := json.Unmarshal(data, &inc)
	if err == nil && inc.Include != "" {
		if val, ok := x.(JsonData); ok {
			data, err := val.Data(inc.Include)
			if err != nil {
				return err
			}
			return json.Unmarshal(data, &j.Val)
		}

		if val, ok := x.(JsonReader); ok {
			r, err := val.Reader(inc.Include)
			if err != nil {
				return err
			}
			return json.NewDecoder(r).Decode(&j.Val)
		}
		f, err := os.Open(inc.Include)
		if err != nil {
			return err
		}
		defer f.Close()
		return json.NewDecoder(f).Decode(&j.Val)
	}
	return json.Unmarshal(data, &j.Val)
}

func (j JsonInclude[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Val)
}
