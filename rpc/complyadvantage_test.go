// +build all

package rpc

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestJsonParsing(t *testing.T) {
	file, err := os.Open("testca.json")
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	var x CAResponse
	err = x.UnmarshalJSON(data)
	if err != nil {
		t.Fatal(err)
	}
}
