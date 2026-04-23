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

	// Define expectations including the new Type field
	expectedPairs := []struct {
		key   string
		value string
		vtype ValueType
	}{
		{"flatString", `"string"`, TypeString},
		{"flatNumber", `123`, TypeNumber},
		{"flatBoolean", `true`, TypeBool},
		{"flatNull", `null`, TypeNull},
	}

	for i, expected := range expectedPairs {
		if !it.Next() {
			t.Fatalf("Expected more pairs, got %d", i)
		}

		// Verify key
		if string(it.Key()) != expected.key {
			t.Errorf("Expected key %q, got %q", expected.key, it.Key())
		}

		// Verify raw value string
		if string(it.Value()) != expected.value {
			t.Errorf("Expected value %q, got %q", expected.value, it.Value())
		}

		// Verify type hint
		if it.Type() != expected.vtype {
			t.Errorf("Expected type %v, got %v", expected.vtype, it.Type())
		}
	}

	if it.Next() {
		t.Fatal("Expected no more pairs, but Next() returned true")
	}
}

func TestIteratorTypes(t *testing.T) {
	type expected struct {
		val   string
		vtype ValueType
	}

	tests := []struct {
		name  string
		input string
		want  map[string]expected
	}{
		{
			name:  "BasicTypes",
			input: `{"str": "hello", "num": 123, "bool": true, "null": null}`,
			want: map[string]expected{
				"str":  {`"hello"`, TypeString},
				"num":  {"123", TypeNumber},
				"bool": {"true", TypeBool},
				"null": {"null", TypeNull},
			},
		},
		{
			name:  "Floats",
			input: `{"pi": 3.14159, "neg": -0.01}`,
			want: map[string]expected{
				"pi":  {"3.14159", TypeNumber},
				"neg": {"-0.01", TypeNumber},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := NewIterator([]byte(tt.input))

			// Track results
			gotVals := make(map[string]string)
			gotTypes := make(map[string]ValueType)

			for it.Next() {
				key := string(it.Key())
				gotVals[key] = string(it.Value())
				gotTypes[key] = it.Type() // Verify the new Type() method
			}

			if err := it.Err(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for k, exp := range tt.want {
				if gotVals[k] != exp.val {
					t.Errorf("key %q: value got %q, want %q", k, gotVals[k], exp.val)
				}
				if gotTypes[k] != exp.vtype {
					t.Errorf("key %q: type got %v, want %v", k, gotTypes[k], exp.vtype)
				}
			}
		})
	}
}

func TestGetValueType(t *testing.T) {
	tests := []struct {
		name string
		val  []byte
		want ValueType
	}{
		{"string", []byte(`"hello"`), TypeString},
		{"number_int", []byte(`123`), TypeNumber},
		{"number_neg", []byte(`-45.6`), TypeNumber},
		{"bool_true", []byte(`true`), TypeBool},
		{"bool_false", []byte(`false`), TypeBool},
		{"null", []byte(`null`), TypeNull},
		{"empty", []byte(``), TypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetValueType(tt.val)
			if got != tt.want {
				t.Errorf("GetValueType() = %v, want %v", got, tt.want)
			}
		})
	}
}
