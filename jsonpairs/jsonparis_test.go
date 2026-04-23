package jsonpairs

import "testing"

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

// TestIterator tests the basic functionality of the Iterator.
// go test -count 1 ./jsonpairs
func TestIterator(t *testing.T) {
	it := NewIterator([]byte(jsonSample))

	// only the top-level pairs should be returned, nested structures should be skipped
	expectedPairs := []struct {
		key   string
		value string
	}{
		{"flatString", `"string"`},
		{"flatNumber", `123`},
		{"flatBoolean", `true`},
		{"flatNull", `null`},
	}

	for i, expected := range expectedPairs {
		if !it.Next() {
			t.Fatalf("Expected more pairs, got %d", i)
		}
		if string(it.Key()) != expected.key {
			t.Errorf("Expected key %q, got %q", expected.key, it.Key())
		}
		if string(it.Value()) != expected.value {
			t.Errorf("Expected value %q, got %q", expected.value, it.Value())
		}
	}

	if it.Next() {
		t.Fatal("Expected no more pairs, but Next() returned true")
	}
}
