# jsonincludes
A Go Package to make it easy to add external file includes to a JSON file

This package makes adding external includes easy. You can import parts of a JSON struct from other files based on some simple rules.

###### Install
`go get github.com/bmurray/jsonincludes`



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
    "objectincludestringopener": {
        "include": "string.json"
    },
    "objectincludestringreader": {
        "include": "string.json"
    }
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
	ObjString         JsonInclude[string]     `json:"objstring"`
	ObjInt            JsonInclude[int]        `json:"objint"`
	ObjStruct         JsonInclude[testData]   `json:"objstruct"`
	ObjArray          JsonInclude[[]testData] `json:"objarray"`
	ObjIncludesString JsonInclude[string]     `json:"objincludestring"`
	ObjIncludesInt    JsonInclude[int]        `json:"objincludesint"`
	ObjIncludesStruct JsonInclude[testData]   `json:"objincludestruct"`
	ObjIncludesArray  JsonInclude[[]testData] `json:"objincludearray"`
}
func main() {
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
}
```
