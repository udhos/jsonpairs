// Package main demonstrates how to use the jsonpairs package to read
// JSON data from stdin and store it in a map. It reads all of stdin
// into a byte slice, initializes an iterator, and iterates through
// the JSON pairs, copying them into a map. Finally, it prints the
// key-value pairs from the map.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/udhos/jsonpairs/jsonpairs"
)

func main() {
	fmt.Fprintf(os.Stderr, "Reading JSON from stdin...\n")

	// 1. Read all of stdin into a byte slice
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize iterator
	it := jsonpairs.NewIterator(data)
	resultMap := make(map[string]string)

	// 3. Iterate and copy data into the map
	for it.Next() {
		// IMPORTANT: Convert to string to force a copy of the underlying bytes
		// because the iterator's slices only point to the original 'data' buffer.
		key := string(it.Key())
		val := string(it.Value())
		it.Type()
		resultMap[key] = val
	}

	if err := it.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// 4. Print results
	for k, v := range resultMap {
		fmt.Printf("Key: %-15s | Value: %-15s | Type: %v\n",
			k, v, jsonpairs.GetValueType([]byte(v)))
	}
}
