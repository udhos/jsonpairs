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

func TestIteratorTypes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]string // Key to expected value string
	}{
		{
			name:  "BasicTypes",
			input: `{"str": "hello", "num": 123, "bool": true, "null": null}`,
			want: map[string]string{
				"str":  `"hello"`,
				"num":  "123",
				"bool": "true",
				"null": "null",
			},
		},
		{
			name:  "Floats",
			input: `{"pi": 3.14159, "neg": -0.01}`,
			want: map[string]string{
				"pi":  "3.14159",
				"neg": "-0.01",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := NewIterator([]byte(tt.input))
			got := make(map[string]string)

			for it.Next() {
				key := string(it.Key())
				val := string(it.Value())
				got[key] = val
			}

			if err := it.Err(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for k, expectedVal := range tt.want {
				if got[k] != expectedVal {
					t.Errorf("key %q: got %q, want %q", k, got[k], expectedVal)
				}
			}
		})
	}
}
