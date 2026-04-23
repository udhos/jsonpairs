// Package bench benchmarks the performance of various implementations.
package bench

import (
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/sugawarayuuta/sonnet"
	"github.com/udhos/jsonpairs/jsonpairs"
	"github.com/valyala/fastjson"
)

/*
go test -bench=. -benchmem ./bench
goos: linux
goarch: amd64
pkg: github.com/udhos/jsonpairs/bench
cpu: 13th Gen Intel(R) Core(TM) i7-1360P
BenchmarkFastJSON-16          	 2279842	       519.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkSonnetJSON-16        	  424488	      3158 ns/op	    1856 B/op	      45 allocs/op
BenchmarkStandardJSON-16      	  213093	      5456 ns/op	    1673 B/op	      47 allocs/op
BenchmarkStandardJSONv2-16    	  330757	      3634 ns/op	    1569 B/op	      30 allocs/op
BenchmarkJsonPairs-16         	 1784907	       675.1 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/udhos/jsonpairs/bench	6.102s
*/

const jsonSample = `
{
	"flatString": "string",
	"flatNumber": 123,
	"flatBoolean": true,
	"flatNull": null,
	"nestedObject": {
		"nestedString": "string",
		"nestedNumber": 123,
		"nestedBoolean": true,
		"nestedNull": null
	},
	"nestedArray": [
		"string",
		123,
		true,
		null,
		{
			"nestedStringInArray": "string",
			"nestedNumberInArray": 123,
			"nestedBooleanInArray": true,
			"nestedNullInArray": null
		}
	]
}
`

// BenchmarkFastJSON measures how fast fastjson parses the top level keys
// go test -bench=. -benchmem ./bench
func BenchmarkFastJSON(b *testing.B) {
	var p fastjson.Parser
	data := []byte(jsonSample)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		v, _ := p.ParseBytes(data)
		o := v.GetObject()

		// Visiting only top-level fields
		o.Visit(func(key []byte, val *fastjson.Value) {
			// This simulates your intended use case:
			// processing only the top-level keys
			_ = key
			_ = val
		})
	}
}

// BenchmarkSonnetJSON measures the overhead of full parsing and allocation
// for the same data structure.
// go test -bench=. -benchmem ./bench
func BenchmarkSonnetJSON(b *testing.B) {
	data := []byte(jsonSample)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		var result map[string]any
		// This is the "expensive" part: Reflection + Heap Allocation
		err := sonnet.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}

		// Accessing keys to force the compiler to keep the work
		_ = result["flatString"]
	}
}

// BenchmarkStandardJSON measures the overhead of full parsing and allocation
// for the same data structure.
// go test -bench=. -benchmem ./bench
func BenchmarkStandardJSON(b *testing.B) {
	data := []byte(jsonSample)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		var result map[string]any
		// This is the "expensive" part: Reflection + Heap Allocation
		err := json.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}

		// Accessing keys to force the compiler to keep the work
		_ = result["flatString"]
	}
}

// BenchmarkStandardJSONv2 measures the overhead of full parsing and allocation
// for the same data structure.
// GOEXPERIMENT=jsonv2 go test -bench=. -benchmem ./bench
func BenchmarkStandardJSONv2(b *testing.B) {
	data := []byte(jsonSample)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		var result map[string]any
		// This is the "expensive" part: Reflection + Heap Allocation
		err := jsonv2.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}

		// Accessing keys to force the compiler to keep the work
		_ = result["flatString"]
	}
}

// BenchmarkJsonPairs measures the speed of your custom zero-alloc iterator
// go test -bench=. -benchmem ./bench
func BenchmarkJsonPairs(b *testing.B) {
	data := []byte(jsonSample)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		it := jsonpairs.NewIterator(data)
		for it.Next() {
			// We access the methods to ensure they are called
			_ = it.Key()
			_ = it.Value()
		}
		if err := it.Err(); err != nil {
			b.Fatal(err)
		}
	}
}
