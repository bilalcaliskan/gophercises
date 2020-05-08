package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type fruitBasket struct {
	Name string
	Fruit []string
	Id int64
	private string // An unexported field is not encoded
	Created time.Time // time.Time and the numeric types in the math/big package can be automatically encoded as JSON strings.
}

var tmpJsonData = []byte(`
{
    "Name": "Standard",
    "Fruit": [
        "Apple",
        "Banana",
        "Orange"
    ],
    "ref": 999,
    "Created": "2018-04-09T23:00:00Z"
}`)

var tmpJsonData2 = []byte(`
{
    "Name": "Eve",
    "Age": 6,
    "Parents": [
        "Alice",
        "Bob"
    ]
}`)

var tmpJsonData3 = []byte(`
    {"Name": "Alice", "Age": 25}
    {"Name": "Bob", "Age": 22}
`)

func main() {
	/*
	The json.Marshal function in package encoding/json generates JSON data.
	Only data that can be represented as JSON will be encoded; see https://golang.org/pkg/encoding/json/#Marshal for the complete rules.
	Only the exported (public) fields of a struct will be present in the JSON output. Other fields are ignored.
	A field with a json: struct tag is stored with its tag name instead of its variable name.
	Pointers will be encoded as the values they point to, or null if the pointer is nil.
	*/
	basket := fruitBasket{
		Name:    "Standart",
		Fruit:   []string{"Apple", "Banana", "Orange"},
		Id:      999,
		private: "Second-rate",
		Created: time.Now(),
	}
	jsonData, err := json.Marshal(basket)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(jsonData))

	prettyJsonData, err := json.MarshalIndent(basket, "", "    ")
	log.Println(string(prettyJsonData))

	/*
	Unmarshaling
	Note that Unmarshal allocated a new slice all by itself. This is how unmarshaling works for slices, maps and pointers.
	For a given JSON key Foo, Unmarshal will attempt to match the struct fields in this order:
		- an exported (public) field with a struct tag json:"Foo",
		- an exported field named Foo, or
		- an exported field named FOO, FoO, or some other case-insensitive match.
	Only fields thar are found in the destination type will be decoded:
		- This is useful when you wish to pick only a few specific fields.
		- In particular, any unexported fields in the destination struct will be unaffected.
	*/
	var unmarshalledBasket fruitBasket
	err = json.Unmarshal(tmpJsonData, &unmarshalledBasket)
	if err != nil {
		log.Println(err)
	}
	log.Println(unmarshalledBasket.Name, unmarshalledBasket.Fruit, unmarshalledBasket.Id)
	log.Println(unmarshalledBasket.Created)

	/*
	Arbitrary objects and arrays
	The encoding/json package uses
		- map[string]interface{} to store arbitrary JSON objects, and
		- []interface{} to store arbitrary JSON arrays.
	It will unmarshal any valid JSON data into a plain interface{} value.
	Consider this JSON data:
		{
		    "Name": "Eve",
		    "Age": 6,
		    "Parents": [
		        "Alice",
		        "Bob"
		    ]
		}
	The json.Unmarshal function will parse it into a map whose keys are strings, and whose values are themselves stored
	as empty interface values.
	We can iterate through the map with a range statement and use a type switch to access its values.
	*/
	var v interface{}
	json.Unmarshal(tmpJsonData2, &v)
	if data, ok := v.(map[string]interface{}); ok {
		for k, v := range data {
			switch v := v.(type) {
			case string:
				fmt.Println(k, v, "(string)")
			case float64:
				fmt.Println(k, v, "(float64)")
			case []interface{}:
				fmt.Println(k, "(array):")
				for i, u := range v {
					fmt.Println("    ", i, u)
				}
			default:
				fmt.Println(k, v, "(unknown)")
			}
		}
	}

	/*
	Reading JSON files
	The json.Decoder and json.Encoder types in package encoding/json offer support for reading and writing streams, e.g.
	files, of JSON data.
	The code in this example:
		- reads a stream of JSON objects from a Reader (strings.Reader),
		- removes the Age field from each object,
		- and then writes the objects to a Writer (os.Stdout).
	*/
	reader := strings.NewReader(string(tmpJsonData3))
	writer := os.Stdout
	decoder := json.NewDecoder(reader)
	encoder := json.NewEncoder(writer)
	for {
		// Read one JSON object and store it in a map(decode)
		var m map[string]interface{}
		if err := decoder.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln(err)
		}

		// Remove all key-value pairs with key == "Age" from the map
		for k := range m {
			if k == "Age" {
				delete(m, k)
			}
		}

		// Write the map as a JSON object(encode)
		if err := encoder.Encode(&m); err != nil {
			log.Println(err)
		}
	}
}