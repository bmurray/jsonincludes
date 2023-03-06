# jsonincludes
A Go Package to make it easy to add external file includes to a JSON file

This package makes adding external includes easy. You can import parts of a JSON struct from other files based on some simple rules.

###### Install
`go get github.com/bmurray/jsonincludes`

###### Usage Notes

The default behavior of the include is to open the file relative to the current directory. This can be overridden using the `JsonData`, `JsonReader`, and `JsonPather` interfaces.

Check the _tests.go file for a deeper usage example.


### Basic Usage

###### rootfile.json
```json
{
    "name": "Sphinx",
    "objstring": "Sphinx",
    "objstruct": {
        "name": "Sphinx",
        "age": 123
    },
    "objint": 123,
    "objarray": [
        {
            "name": "Sphinx",
            "age": 123
        },
        {
            "name": "King",
            "age": 321
        }
    ],
    "objincludestring": {
        "include": "string.json"
    },
    "objincludesint": {
        "include": "int.json"
    },
    "objincludestruct": {
        "include": "struct.json"
    },
    "objincludearray": {
        "include": "array.json"
    },
}

```

Note, any files that are included are using the current working directory. You may change this behavior by implementing the JsonData or JsonReader interfaces. See the _test.go file for details. Check the testdata directory for complete examples.

```go
type testData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type testObject struct {
	testData

	ObjString                 jsonincludes.JsonInclude[string]     `json:"objstring"`
	ObjInt                    jsonincludes.JsonInclude[int]        `json:"objint"`
	ObjStruct                 jsonincludes.JsonInclude[testData]   `json:"objstruct"`
	ObjArray                  jsonincludes.JsonInclude[[]testData] `json:"objarray"`
	ObjIncludesString         jsonincludes.JsonInclude[string]     `json:"objincludesstring"`
	ObjIncludesInt            jsonincludes.JsonInclude[int]        `json:"objincludesint"`
	ObjIncludesStruct         jsonincludes.JsonInclude[testData]   `json:"objincludesstruct"`
	ObjIncludesArray          jsonincludes.JsonInclude[[]testData] `json:"objincludesarray"`

}
func main() {

    // use this to change the relative path of the includes while using JsonInclude
    // jsonincludes.SetRootPath("relative_path")
    var obj testObject
	f, err := os.Open("rootfile.json")
	if err != nil {
		// Handle the error
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&obj)
	if err != nil {
		// Handle the error
	}
    // obj.ObjString.Val and obj.ObjIncludeString.Val contains the data regardless if it was embedded or included
}
```


### Advanced Usage with Global Prefix path

```json
{
    "objstring": {
        "include": "string.json"
    },
}
```

```go
type testData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type testOpener struct {}

// Use a static path, or something in a global or package scope
func (testOpener) Path(name string) string {
    return filepath.Join("./configs", name)
}
type testObject struct {
    ObjString jsonincludes.JsonBase[testOpener, testData] `json:"objstring"`
}
func main() {
    var obj testObject
    f, _ := os.Open("rootfile.json")
    // Handle the error
    json.NewDecoder(f).Decode(&obj)
    // obj.ObjString.Val contains the data regardless of if its imported or embedded
}

```

### Advanced Usage with Specific Prefix path

Use these methods if there needs to be some special rule around each key. This lets you create your own configuration as needed.

```json
{
    "objstring": {
        "include": "string.json"
    },
}
```

```go
type testData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type testOpener string

// Use a static path, or something in a more global scope
func (p testOpener) Path(name string) string {
    return filepath.Join(string(p), name)
}
type testObject struct {
    ObjString JsonBase[testOpener, testData] `json:"objstring"`
}
func main() {
    var obj testObject
    
    obj.ObjString.Config = testOpener("./config")

    f, _ := os.Open("rootfile.json")
    // Handle the error
    json.NewDecoder(f).Decode(&obj)
    // obj.ObjString.Val contains the data regardless of if its imported or embedded
}

```