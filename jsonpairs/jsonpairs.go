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
	it.pos++ // Skip the opening quote
	for it.pos < len(it.data) {
		if it.data[it.pos] == '\\' {
			it.pos += 2 // Skip escaped character
		} else if it.data[it.pos] == '"' {
			it.pos++ // Skip the closing quote
			return
		} else {
			it.pos++
		}
	}
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
			return false
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
