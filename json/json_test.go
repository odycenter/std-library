package json_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"std-library/json"
	"testing"
)

type Outer struct {
	Inner
}

type Inner struct {
	Inner string `json:",omitempty"`
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestMarshal(t *testing.T) {
	p := Person{Name: "Iron", Age: 18}
	assert.Equal(t, `{"name":"Iron","age":18}`, json.String(p))

	outer := Outer{Inner: Inner{Inner: "inner"}}
	assert.Equal(t, `{"Inner":"inner"}`, json.String(outer))
}

func TestUnmarshal(t *testing.T) {
	jsonStr := `{"name":"Iron","age":18}`
	var p Person
	json.Parse(jsonStr, &p)

	assert.Equal(t, "Iron", p.Name)
	assert.Equal(t, 18, p.Age)
}

func TestParse(t *testing.T) {
	var st = struct {
		A int
	}{}
	json.Parse([]byte(`{"A":1}`), &st)
	fmt.Println(st)
	st = struct {
		A int
	}{}
	json.Parse([]byte(`{"A":"2"}`), &st)
	fmt.Println(st)
	st = struct {
		A int
	}{}
	fmt.Println(st)
	json.Parse([]byte(`{"A":{\"A\":3}}`), &st)
	st = struct {
		A int
	}{}
	json.Parse([]byte(`{"A":{\"A\":4}`), &st)
	fmt.Println(st)
}

func TestParseE(t *testing.T) {
	st1 := struct {
		A int
	}{}
	e := json.ParseE(`{"A":1}`, &st1)
	fmt.Println(e, st1)
	st2 := struct {
		A string
	}{}
	e = json.ParseE(`{"A":"2"}`, &st2)
	fmt.Println(e, st2)
	st3 := struct {
		A struct {
			A int
		}
	}{}
	fmt.Println(e, st3)
	e = json.ParseE(`{"A":{\"A\":3}}`, &st3)
	st4 := struct {
		A int
	}{}
	e = json.ParseE(`{"A":{\"A\":4}`, &st4)
	fmt.Println(e, st4)
	st5 := struct {
		A int
	}{}
	e = json.ParseE(`{"A":"5"}`, &st5)
	fmt.Println(e, st5)
}

func TestStringify(t *testing.T) {
	var st1 = struct{ A int }{1}
	fmt.Println(json.Stringify(st1))
	var st2 = struct{ A int }{2}
	fmt.Println(json.Stringify(st2))
	var st3 = struct{ A int }{3}
	fmt.Println(json.Stringify(st3))
	var st4 = struct{ A int }{4}
	fmt.Println(json.Stringify(st4))
	var m1 = map[int]int{1: 2, 2: 3}
	fmt.Println(json.Stringify(m1))
}

func TestStringifyE(t *testing.T) {
	var st1 = struct{ A int }{1}
	fmt.Println(json.StringifyE(st1))
	var st2 = struct{ A int }{2}
	fmt.Println(json.StringifyE(st2))
	var st3 = struct{ A int }{3}
	fmt.Println(json.StringifyE(st3))
	var st4 = struct{ A int }{4}
	fmt.Println(json.StringifyE(st4))
	var m1 = map[int]int{1: 2, 2: 3}
	fmt.Println(json.StringifyE(m1))
}

func TestValid(t *testing.T) {
	fmt.Println(json.Valid(`{"example": 1}`))
	fmt.Println(json.Valid(`{"example":2:]}}`))
}

type ResPlatform struct {
	A int32
	B string
	C json.RawMessage
}

func TestRawMessage(t *testing.T) {
	var v ResPlatform
	s := `{"A":200,"B":"","C":{"CA":[{"T1":1,"T2":"ABC"},{"T1":2,"T2":"DEF"}],"CB":1000}}`
	json.Parse(s, &v)
	fmt.Println(v)
	fmt.Println(string(json.Stringify(v)))
	fmt.Println(json.String(v))

	jsonStr := `{"name":"Iron","details":{"age":37,"city":"Taoyuan City"}}`
	type Person struct {
		Name    string          `json:"name"`
		Details json.RawMessage `json:"details"`
	}

	var p Person
	json.Parse([]byte(jsonStr), &p)

	expected := `{"age":37,"city":"Taoyuan City"}`
	if string(p.Details) != expected {
		t.Errorf("Expected %s, got %s", expected, string(p.Details))
	}
}

func TestNumber(t *testing.T) {
	jsonStr := `{"value":12345678901234567890}`
	var data map[string]big.Int
	json.Parse(jsonStr, &data)
	expected := new(big.Int)
	expected.SetString("12345678901234567890", 10)
	assert.Equal(t, *expected, data["value"])
}

var p = Person{Name: "Iron", Age: 37}

func BenchmarkMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.String(p)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	jsonStr := `{"name":"Iron","age":37}`
	var p Person

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Parse([]byte(jsonStr), &p)
	}
}

//func BenchmarkNumber(b *testing.B) {
//	jsonStr := `{"value":12345678901234567890}`
//	var data map[string]gojson.Number
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		gojson.Unmarshal([]byte(jsonStr), &data)
//	}
//}

func BenchmarkRawMessage(b *testing.B) {
	jsonStr := `{"name":"Iron","details":{"age":37,"city":"Taoyuan City"}}`
	type Person struct {
		Name    string          `json:"name"`
		Details json.RawMessage `json:"details"`
	}
	var p Person

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Parse([]byte(jsonStr), &p)
	}
}
