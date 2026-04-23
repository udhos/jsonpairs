// Package jsonpairs provides an iterator for extracting key-value pairs
// from JSON objects without fully parsing the JSON structure.
package jsonpairs

// Iterator holds the state of the JSON parsing process.
type Iterator struct {
	data      []byte
	err       error
	pos       int
	depth     int
	currKey   []byte
	currValue []byte
}

// NewIterator creates a new Iterator for the given JSON data.
func NewIterator(data []byte) *Iterator {
	return &Iterator{
		data:  data,
		pos:   0,
		depth: 0,
	}
}

func (it *Iterator) skipWhitespace() {
	for it.pos < len(it.data) {
		switch it.data[it.pos] {
		case ' ', '\t', '\n', '\r':
			it.pos++
		default:
			return
		}
	}
}

func (it *Iterator) skipString() {
	pos := it.pos
	pos++ // Skip the opening quote
	size := len(it.data)
LOOP:
	for pos < size {
		switch it.data[pos] {
		case '\\':
			pos += 2 // Skip escaped character
		case '"':
			pos++ // Skip the closing quote
			break LOOP
		default:
			pos++
		}
	}
	it.pos = pos
}

func (it *Iterator) skipValue() {
	for it.pos < len(it.data) {
		switch it.data[it.pos] {
		case ' ', '\t', '\n', '\r', ',', ':', '}', ']':
			return
		default:
			it.pos++
		}
	}
}

func (it *Iterator) skipCompositeValue() bool {
	pos := it.pos
	size := len(it.data)

	if pos >= size {
		return false
	}

	var open byte
	var closeMarker byte
	switch it.data[pos] {
	case '{':
		open = '{'
		closeMarker = '}'
	case '[':
		open = '['
		closeMarker = ']'
	default:
		return false
	}

	depth := 0
	for pos < size {
		switch it.data[pos] {
		case '"':
			it.pos = pos
			it.skipString()
			pos = it.pos
		case open:
			depth++
			pos++
		case closeMarker:
			depth--
			pos++
			if depth == 0 {
				it.pos = pos
				return true
			}
		case '{', '[':
			depth++
			pos++
		case '}', ']':
			depth--
			pos++
		default:
			pos++
		}
	}

	it.pos = pos

	return false
}

func (it *Iterator) parsePair() bool {
	// 1. Capture the key, but skip the opening quote
	it.pos++ // Skip opening quote
	startKey := it.pos
	it.skipString() // This will stop at the closing quote
	// The key is everything between the quotes
	it.currKey = it.data[startKey : it.pos-1]

	// 2. Expect a colon
	it.skipWhitespace()
	if it.pos < len(it.data) && it.data[it.pos] == ':' {
		it.pos++
	} else {
		return false
	}

	// 3. Capture the value
	it.skipWhitespace()
	if it.pos < len(it.data) {
		switch it.data[it.pos] {
		case '{', '[':
			it.skipCompositeValue()
			return false
		case '"':
			startVal := it.pos
			it.skipString()
			it.currValue = it.data[startVal:it.pos]
			return true
		}
	}
	startVal := it.pos
	it.skipValue()
	it.currValue = it.data[startVal:it.pos]

	return true
}

// Next advances the iterator to the next key-value pair.
func (it *Iterator) Next() bool {
	if it.err != nil {
		return false
	}

	for it.pos < len(it.data) {
		it.skipWhitespace()
		if it.pos >= len(it.data) {
			return false
		}

		b := it.data[it.pos]

		switch b {
		case '{':
			it.depth++
			it.pos++
			continue // Ensure we don't process '{' as a pair
		case '[':
			it.depth++
			it.pos++
			continue
		case '}', ']':
			it.depth--
			it.pos++
			continue
		case '"':
			if it.depth == 1 {
				if it.parsePair() {
					return true
				}
				continue
			}
			it.skipString()
		case ',', ':':
			it.pos++ // Just consume and continue
		default:
			it.skipValue()
		}
	}
	return false
}

// Err returns the first error encountered by the iterator.
func (it *Iterator) Err() error {
	return it.err
}

// Key returns the current key.
func (it *Iterator) Key() []byte { return it.currKey }

// Value returns the current value.
func (it *Iterator) Value() []byte { return it.currValue }

// Type returns the ValueType of the current key-value pair.
func (it *Iterator) Type() ValueType {
	return GetValueType(it.currValue)
}

// ValueType represents the inferred type of a JSON value.
type ValueType int

const (
	TypeString ValueType = iota
	TypeNumber
	TypeBool
	TypeNull
	TypeUnknown
)

// GetValueType inspects the first byte of a value slice to hint at its type.
// It does not parse the value or perform any allocations.
func GetValueType(b []byte) ValueType {
	if len(b) == 0 {
		return TypeUnknown
	}
	switch b[0] {
	case '"':
		return TypeString
	case 't', 'f':
		return TypeBool
	case 'n':
		return TypeNull
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return TypeNumber
	default:
		// Everything else (like '.') is invalid in strict JSON
		return TypeUnknown
	}
}
