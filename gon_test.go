package gon

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	jsoniterator "github.com/json-iterator/go"
)

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func TestDeserialize(t *testing.T) {
	gonFile, err := os.ReadFile("./gon_files/gon.gon")
	assert(err)
	output, err := DeserializeString(string(gonFile))
	assert(err)
	gonJson, err := json.MarshalIndent(output, "+", "\t")
	assert(err)
	jsonFile, err := os.ReadFile("./json_files/json.json")
	hashmap := map[string]any{}
	err = json.Unmarshal(jsonFile, &hashmap)
	assert(err)
	jsonJson, err := json.MarshalIndent(hashmap, "+", "\t")
	if string(jsonJson) != string(gonJson) {
		panic(fmt.Sprintf("jsonJson: %s\n!=\ngsonJson:%s", string(jsonJson), gonJson))
	}
}

// goos: darwin
// goarch: arm64
// pkg: github.com/gabivlj/gon-go
// BenchmarkDeserialize-10    	1000000000	         0.1401 ns/op	       0 B/op	       0 allocs/op
// PASS
// ok  	github.com/gabivlj/gon-go	1.611s
func BenchmarkDeserialize(b *testing.B) {
	gonFile, err := os.ReadFile("./gon_files/gon.gon")

	assert(err)
	for i := 0; i < 100000; i++ {
		_, err := DeserializeString(string(gonFile))
		assert(err)
	}
}

// goos: darwin
// goarch: arm64
// pkg: github.com/gabivlj/gon-go
// BenchmarkDeserializeJSON-10    	1000000000	         0.4111 ns/op	       0 B/op	       0 allocs/op
// PASS
// ok  	github.com/gabivlj/gon-go	9.261s
func BenchmarkDeserializeJSON(b *testing.B) {
	jsonFile, err := os.ReadFile("./json_files/json.json")

	assert(err)
	for i := 0; i < 100000; i++ {
		hashmap := map[string]any{}
		err := json.Unmarshal(jsonFile, &hashmap)
		assert(err)
	}
}

// goos: darwin
// goarch: arm64
// pkg: github.com/gabivlj/gon-go
// BenchmarkDeserializeJSONWithSpecializedLibrary-10    	1000000000	         0.2645 ns/op	       0 B/op	       0 allocs/op
// PASS
// ok  	github.com/gabivlj/gon-go	4.587s
func BenchmarkDeserializeJSONWithSpecializedLibrary(b *testing.B) {
	jsonFile, err := os.ReadFile("./json_files/json.json")

	assert(err)
	for i := 0; i < 100000; i++ {
		hashmap := map[string]any{}
		err := jsoniterator.Unmarshal(jsonFile, &hashmap)
		assert(err)
	}
}
