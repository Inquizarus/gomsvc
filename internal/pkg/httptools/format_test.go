package httptools_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/inquizarus/gomsvc/internal/pkg/httptools"
)

func TestFormatJSONData(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected []byte
		err      error
	}{
		{
			name:  "valid JSON",
			input: []byte(`{"name": "Alice","age": 25}`),
			expected: []byte(`{
 "age": 25,
 "name": "Alice"
}`),
			err: nil,
		},
		{
			name:     "invalid JSON",
			input:    []byte(`invalid JSON`),
			expected: nil,
			err:      &json.SyntaxError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := httptools.FormatJSONData(tc.input)
			if err != nil && reflect.TypeOf(err) != reflect.TypeOf(tc.err) {
				t.Fatalf("unexpected error type: expected %v, but got %v", tc.err, err)
			}
			if !bytes.Equal(actual, tc.expected) {
				t.Errorf("unexpected result: expected %s, but got %s", tc.expected, actual)
			}
		})
	}
}
