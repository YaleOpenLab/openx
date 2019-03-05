// +build all

package main

import (
	"encoding/json"
	"testing"

	database "github.com/YaleOpenLab/openx/database"
	"github.com/pquerna/ffjson/ffjson"
)

// go test -run=XXX -tags="all" -bench=.

// pit three json encoding libraries against each other to see which one is the fastest.
func BenchmarkFFJsonMarshal(b *testing.B) {
	s := &database.User{
		Index:       1,
		Name:        "testuser",
		PublicKey:   "randompublickey",
		Username:    "myusername",
		Pwhash:      "mypwhash",
		Address:     "myhomeaddress",
		Description: "mydescription",
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ffjson.Marshal(s)
	}
}

func BenchmarkFFJsonUnmarshal(b *testing.B) {
	s := &database.User{
		Index:       1,
		Name:        "testuser",
		PublicKey:   "randompublickey",
		Username:    "myusername",
		Pwhash:      "mypwhash",
		Address:     "myhomeaddress",
		Description: "mydescription",
	}
	data, _ := ffjson.Marshal(s)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ffjson.Unmarshal(data, &s)
	}
}

func BenchmarkEJsonMarshal(b *testing.B) {
	s := &database.User{
		Index:       1,
		Name:        "testuser",
		PublicKey:   "randompublickey",
		Username:    "myusername",
		Pwhash:      "mypwhash",
		Address:     "myhomeaddress",
		Description: "mydescription",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.MarshalJSON()
	}
}

func BenchmarkEJsonUnmarshal(b *testing.B) {
	s := &database.User{
		Index:       1,
		Name:        "testuser",
		PublicKey:   "randompublickey",
		Username:    "myusername",
		Pwhash:      "mypwhash",
		Address:     "myhomeaddress",
		Description: "mydescription",
	}
	data, _ := s.MarshalJSON()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.UnmarshalJSON(data)
	}
}

func BenchmarkJsonMarshal(b *testing.B) {
	s := &database.User{
		Index:       1,
		Name:        "testuser",
		PublicKey:   "randompublickey",
		Username:    "myusername",
		Pwhash:      "mypwhash",
		Address:     "myhomeaddress",
		Description: "mydescription",
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(&s)
	}
}

func BenchmarkJsonUnmarshal(b *testing.B) {
	s := &database.User{
		Index:       1,
		Name:        "testuser",
		PublicKey:   "randompublickey",
		Username:    "myusername",
		Pwhash:      "mypwhash",
		Address:     "myhomeaddress",
		Description: "mydescription",
	}
	data, _ := json.Marshal(s)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		json.Unmarshal(data, &s)
	}
}
