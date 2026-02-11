package steamutils

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/iancoleman/orderedmap"
)

// skipWhitespace advances i past all whitespace characters.
func skipWhitespace(s string, i int) int {
	for i < len(s) && unicode.IsSpace(rune(s[i])) {
		i++
	}
	return i
}

// parseString expects a starting double quote and returns the string—note this simple version does not handle escapes.
func parseString(s string, i int) (string, int, error) {
	if s[i] != '"' {
		return "", i, fmt.Errorf("expected '\"' at position %d", i)
	}
	i++ // skip opening quote
	start := i
	for i < len(s) && s[i] != '"' {
		i++
	}
	if i >= len(s) {
		return "", i, errors.New("unterminated string")
	}
	return s[start:i], i + 1, nil // skip closing quote
}

// parseVDFOrdered recursively parses VDF text starting at position i, storing keys in order.
func parseVDFOrdered(s string, i int) (*orderedmap.OrderedMap, int, error) {
	ordMap := orderedmap.New()
	for i < len(s) {
		i = skipWhitespace(s, i)
		if i >= len(s) {
			break
		}
		// A closing brace signals end of the current block.
		if s[i] == '}' {
			return ordMap, i + 1, nil
		}
		// Parse key (should be in quotes)
		key, newIdx, err := parseString(s, i)
		if err != nil {
			return nil, i, err
		}
		i = skipWhitespace(s, newIdx)
		if i >= len(s) {
			return nil, i, fmt.Errorf("unexpected end after key %s", key)
		}

		var value interface{}
		// If the next character is an opening brace, parse a nested object.
		switch s[i] {
		case '{':
			i++ // skip '{'
			nestedMap, newIdx, err := parseVDFOrdered(s, i)
			if err != nil {
				return nil, i, err
			}
			value = nestedMap
			i = newIdx
		case '"':
			// Otherwise expect a string value.
			strVal, newIdx, err := parseString(s, i)
			if err != nil {
				return nil, i, err
			}
			value = strVal
			i = newIdx
		default:
			return nil, i, fmt.Errorf("unexpected character '%c' at position %d", s[i], i)
		}
		ordMap.Set(key, value)
	}
	return ordMap, i, nil
}

// marshalOrderedVDF recursively serializes an ordered map to a VDF-formatted string.
func marshalOrderedVDF(m *orderedmap.OrderedMap, indent int) (string, error) {
	var sb strings.Builder
	spacing := strings.Repeat("\t", indent)
	for _, key := range m.Keys() {
		value, _ := m.Get(key)
		// Write the key.
		sb.WriteString(fmt.Sprintf("%s\"%s\"\n", spacing, key))
		switch v := value.(type) {
		case string:
			sb.WriteString(fmt.Sprintf("%s\"%s\"\n", spacing, v))
		case *orderedmap.OrderedMap:
			sb.WriteString(fmt.Sprintf("%s{\n", spacing))
			inner, err := marshalOrderedVDF(v, indent+1)
			if err != nil {
				return "", err
			}
			sb.WriteString(inner)
			sb.WriteString(fmt.Sprintf("%s}\n", spacing))
		default:
			// Fallback: print the value using fmt.
			s := fmt.Sprintf("%v", v)
			sb.WriteString(fmt.Sprintf("%s\"%s\"\n", spacing, s))
		}
	}
	return sb.String(), nil
}

// Parses and unmarshals VDF file into map
func Unmarshal(data []byte) (*orderedmap.OrderedMap, error) {
	ordMap, _, err := parseVDFOrdered(string(data), 0)
	return ordMap, err
}

func Marshal(m *orderedmap.OrderedMap) ([]byte, error) {
	out, err := marshalOrderedVDF(m, 0)
	return []byte(out), err
}

// The code in this file was made by ChatGPT, use in production is highly discouraged as unexpected results may occur. The code in this file is not vetted for stability or edge cases.
