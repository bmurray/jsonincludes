package jsonincludes

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// JsonData is an interface that allows a type to provide json data based on the name.
//
// This is useful for providing a json file from a zip file or a database.
type JsonData interface {
	Data(string) ([]byte, error)
}

// JsonReader is an interface that allows a type to provide a reader based on the name.
//
// This is useful for providing a json file from a zip file or a database or network connection.
type JsonReader interface {
	Reader(string) (io.Reader, error)
}

// JsonPather is an interface that allows a type to provide a path based on the name.
//
// This is useful for translating the file name to a relative path instead of relative to the current working directory.
type JsonPather interface {
	Path(string) string
}

var rootPath = ""

// SetRootPath sets the root path for the json includes.
//
// This is used by the Resolver type to resolve relative paths, and in the JsonInclude type to resolve relative paths.
func SetRootPath(path string) {
	rootPath = path
}

// Resolver is a type that implements the JsonPather interface.
//
// This is used by the JsonInclude type to resolve relative paths.
type Resolver struct{}

// Path returns the path relative to the root path.
func (Resolver) Path(name string) string {
	return filepath.Join(rootPath, name)
}

// JsonInclude is a type that can be used to include json data from a file.
//
// This is the main interface to the jsonincludes package. Use this if unsure.
type JsonInclude[T any] struct{ JsonBase[Resolver, T] }

// JsonCombo is a type that can be used to include json data from a file.
//
// This is a convience type that allows you to specify the type that implements the JsonPather/JsonData/JsonReader interfaces thats the same as the type that implements the json.Unmarshaler/json.Marshaler interfaces.
type JsonCombo[T any] struct{ JsonBase[T, T] }

type jsonInclude struct {
	Include string `json:"include"`
}

// JsonBase is a type that can be used to include json data from a file.
//
// The first generic type should be the type that implements the JsonPather/JsonData/JsonReader interfaces.
// The second generic type should be the type that implements the json.Unmarshaler/json.Marshaler interfaces.
type JsonBase[Y any, T any] struct {
	Val    T `json:"-"`
	Config Y `json:"-"`
}

// UnmarshalJSON unmarshals the json data.
func (j *JsonBase[Y, T]) UnmarshalJSON(data []byte) error {
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

// MarshalJSON marshals the json data.
//
// Note, this will not marshal out to a file, it will only marshal the data.
func (j JsonBase[Y, T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Val)
}
