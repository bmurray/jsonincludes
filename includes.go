package jsonincludes

import (
	"encoding/json"
	"io"
	"os"
)

type JsonError string

const (
	ErrNoData        JsonError = "no data"
	ErrNoReader      JsonError = "no reader"
	ErrCannotInclude JsonError = "cannot include"
)

type JsonData interface {
	Data(string) ([]byte, error)
}
type JsonReader interface {
	Reader(string) (io.Reader, error)
}
type JsonPather interface {
	Path(string) string
}
type JsonInclude[T any] struct {
	Val T `json:"-"`
}
type jsonInclude struct {
	Include string `json:"include"`
}

func (j *JsonInclude[T]) UnmarshalJSON(data []byte) error {

	var x interface{} = j.Val
	var inc jsonInclude

	err := json.Unmarshal(data, &inc)
	if err == nil && inc.Include != "" {
		path := inc.Include
		if val, ok := x.(JsonPather); ok {
			path = val.Path(inc.Include)
		}
		if val, ok := x.(JsonData); ok {
			data, err := val.Data(path)
			if err != nil {
				return err
			}
			return json.Unmarshal(data, &j.Val)
		}

		if val, ok := x.(JsonReader); ok {
			r, err := val.Reader(path)
			if err != nil {
				return err
			}
			return json.NewDecoder(r).Decode(&j.Val)
		}
		f, err := os.Open(path)
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

type JsonIncludeOpener[Y any, T any] struct {
	Val    T `json:"-"`
	Config Y `json:"-"`
}

func (j *JsonIncludeOpener[Y, T]) UnmarshalJSON(data []byte) error {
	var x interface{} = j.Config

	var inc jsonInclude
	err := json.Unmarshal(data, &inc)
	if err == nil && inc.Include != "" {
		path := inc.Include
		if val, ok := x.(JsonPather); ok {
			path = val.Path(inc.Include)
		}

		if val, ok := x.(JsonData); ok {
			data, err := val.Data(path)
			if err != nil {
				return err
			}
			return json.Unmarshal(data, &j.Val)
		}
		if val, ok := x.(JsonReader); ok {
			r, err := val.Reader(path)
			if err != nil {
				return err
			}
			return json.NewDecoder(r).Decode(&j.Val)
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		return json.NewDecoder(f).Decode(&j.Val)
	}
	return json.Unmarshal(data, &j.Val)
}
func (j JsonIncludeOpener[Y, T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Val)
}
