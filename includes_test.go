package jsonincludes

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
)

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

func TestJsonInclude(t *testing.T) {
	var obj testObject
	f, err := os.Open("testdata/basic.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&obj)
	if err != nil {
		t.Fatal(err)
	}
	if obj.ObjString.Val != "Sphinx" {
		t.Fatal("string not loaded", obj.ObjString.Val)
	}
	if obj.ObjInt.Val != 123 {
		t.Fatal("int not loaded")
	}
	if obj.ObjStruct.Val.Name != "Sphinx" {
		t.Fatal("struct not loaded")
	}
	if obj.ObjArray.Val[0].Name != "Sphinx" {
		t.Fatal("array not loaded")
	}
	if obj.ObjIncludesString.Val != "Sphinx" {
		t.Fatal("include string not loaded", obj.ObjIncludesString.Val)
	}
	if obj.ObjIncludesInt.Val != 123 {
		t.Fatal("include int not loaded")
	}
	if obj.ObjIncludesStruct.Val.Name != "Sphinx" {
		t.Fatal("include struct not loaded")
	}
	if obj.ObjIncludesArray.Val[0].Name != "Sphinx" {
		t.Fatal("include array not loaded")
	}
	output := &bytes.Buffer{}
	err = json.NewEncoder(output).Encode(obj)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output.String())
	var obj2 testObject
	err = json.NewDecoder(output).Decode(&obj2)
	if err != nil {
		t.Fatal(err)
	}
	if obj2.ObjString.Val != "Sphinx" {
		t.Fatal("string not loaded", obj2.ObjString.Val)
	}
	if obj2.ObjInt.Val != 123 {
		t.Fatal("int not loaded")
	}
	if obj2.ObjStruct.Val.Name != "Sphinx" {
		t.Fatal("struct not loaded")
	}
	if obj2.ObjArray.Val[0].Name != "Sphinx" {
		t.Fatal("array not loaded")
	}
	if obj2.ObjIncludesString.Val != "Sphinx" {
		t.Fatal("include string not loaded", obj2.ObjIncludesString.Val)
	}
	if obj2.ObjIncludesInt.Val != 123 {
		t.Fatal("include int not loaded")
	}
	if obj2.ObjIncludesStruct.Val.Name != "Sphinx" {
		t.Fatal("include struct not loaded")
	}
	if obj2.ObjIncludesArray.Val[0].Name != "Sphinx" {
		t.Fatal("include array not loaded")
	}
}

type testOpenerString string

func (t testOpenerString) Data(name string) ([]byte, error) {
	n := filepath.Join(string(t), name)
	f, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, f)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type testOpenerStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (t testOpenerStruct) Data(name string) ([]byte, error) {
	f, err := os.Open("./testdata/" + name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

type testReaderString string

func (t testReaderString) Reader(name string) (io.Reader, error) {
	n := filepath.Join(string(t), name)
	return os.Open(n)
}

type testReaderStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (t testReaderStruct) Reader(name string) (io.Reader, error) {
	return os.Open("./testdata/" + name)
}

type testObjectInterfaces struct {
	testData
	ObjOpenerString JsonInclude[testOpenerString] `json:"objopenerstring"`
	ObjOpenerStruct JsonInclude[testOpenerStruct] `json:"objopenerstruct"`
	ObjReaderString JsonInclude[testReaderString] `json:"objreaderstring"`
	ObjReaderStruct JsonInclude[testReaderStruct] `json:"objreaderstruct"`
}

func TestInterfaces(t *testing.T) {
	var obj testObjectInterfaces
	obj.ObjOpenerString.Val = testOpenerString("./testdata")
	obj.ObjReaderString.Val = testReaderString("./testdata")

	f, err := os.Open("testdata/interfaces.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&obj)
	if err != nil {
		t.Fatal(err)
	}
	if obj.ObjOpenerString.Val != "Sphinx" {
		t.Fatal("opener string not loaded", obj.ObjOpenerString.Val)
	}
	if obj.ObjOpenerStruct.Val.Name != "Sphinx" {
		t.Fatal("opener struct not loaded")
	}
	if obj.ObjReaderString.Val != "Sphinx" {
		t.Fatal("reader string not loaded", obj.ObjReaderString.Val)
	}
	if obj.ObjReaderStruct.Val.Name != "Sphinx" {
		t.Fatal("reader struct not loaded")
	}
	output := &bytes.Buffer{}
	err = json.NewEncoder(output).Encode(obj)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output.String())
	var obj2 testObjectInterfaces
	err = json.NewDecoder(output).Decode(&obj2)
	if err != nil {
		t.Fatal(err)
	}
	if obj2.ObjOpenerString.Val != "Sphinx" {
		t.Fatal("opener string not loaded", obj2.ObjOpenerString.Val)
	}
	if obj2.ObjOpenerStruct.Val.Name != "Sphinx" {
		t.Fatal("opener struct not loaded")
	}
	if obj2.ObjReaderString.Val != "Sphinx" {
		t.Fatal("reader string not loaded")
	}
	if obj2.ObjReaderStruct.Val.Name != "Sphinx" {
		t.Fatal("reader struct not loaded")
	}

}

type testObjectInterfaceBad struct {
	testData
	ObjOpenerString JsonInclude[testOpenerStruct] `json:"objopenerstring"`
}

func TestBadInterfaces(t *testing.T) {
	var obj testObjectInterfaceBad
	f, err := os.Open("testdata/interfaces.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&obj)
	if err == nil {
		t.Fatal("expected error")
	} else {
		t.Log(err)
	}
}

type testReaderHidden struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	rootPath string
}

func (t testReaderHidden) Reader(name string) (io.Reader, error) {
	n := filepath.Join(string(t.rootPath), name)
	return os.Open(n)
}

func TestHiddenFolder(t *testing.T) {
	var obj JsonInclude[testReaderHidden]
	obj.Val.rootPath = "./testdata"
	var x = struct {
		Obj JsonInclude[testReaderHidden] `json:"obj"`
	}{obj}
	err := json.Unmarshal([]byte(`{"obj": {"include": "struct.json"}}`), &x)
	if err != nil {
		t.Fatal(err)
	}
	if x.Obj.Val.Name != "Sphinx" {
		t.Fatal("include struct not loaded")
	}
	output := &bytes.Buffer{}
	err = json.NewEncoder(output).Encode(x)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output.String())
	var y = struct {
		Obj JsonInclude[testReaderHidden] `json:"obj"`
	}{}
	if err := json.Unmarshal(output.Bytes(), &y); err != nil {
		t.Fatal(err)
	}
	if y.Obj.Val.Name != "Sphinx" {
		t.Fatal("include struct not loaded", y.Obj.Val.Name)
	}
}
