package ds

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

type Entry[K string, V any] struct {
	Key   K
	Value V
}

// Map is a custom Map impementation that maintains order when adding new items and when unmarshalling from JSON and YAML.
type Map[K string, V any] struct {
	keys    map[K]int
	entries []Entry[K, V]
}

// NewMap creates a new ordered Map.
func NewMap[K string, V any]() *Map[K, V] {
	return &Map[K, V]{
		keys:    make(map[K]int),
		entries: make([]Entry[K, V], 0),
	}
}

// Set adds a new key/value pair and keeps track of the order in which it was added.
func (om *Map[K, V]) Set(key K, val V) {
	if idx, exists := om.keys[key]; exists {
		om.entries[idx].Value = val
		return
	}
	om.keys[key] = len(om.entries)
	om.entries = append(om.entries, Entry[K, V]{Key: key, Value: val})
}

// Get returns the value associated with a given key in the map.
// It follows the interface of the built in map type returning true/false in the second parameter to indicate if the key was found in the map.
func (om *Map[K, V]) Get(key K) (V, bool) {
	var v V

	idx, ok := om.keys[key]
	if !ok {
		return v, ok
	}

	return om.entries[idx].Value, true
}

func (om *Map[K, V]) Entries() []Entry[K, V] {
	return om.entries
}

// Keys returns the keys int the map as a slice.
func (om *Map[K, V]) Keys() []K {
	keys := []K{}
	for key := range om.keys {
		keys = append(keys, key)
	}
	return keys
}

// UnmarshalJSON implements the json.Unmarshaller interface to enable json.Unarmshal.
func (om *Map[K, V]) UnmarshalJSON(data []byte) error {
	// Initialize or clear internal state
	om.keys = make(map[K]int)
	om.entries = make([]Entry[K, V], 0)

	dec := json.NewDecoder(bytes.NewReader(data))

	// Expect the start of a JSON object '{'
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expected JSON object start '{'")
	}

	// Read tokens until the object ends '}'
	for dec.More() {
		// Read the key token
		tKey, err := dec.Token()
		if err != nil {
			return err
		}

		// Convert token to key type (usually string)
		var key K
		keyStr, ok := tKey.(string)
		if !ok {
			return fmt.Errorf("expected string key, got %v", tKey)
		}

		// This conversion relies on K being a string or compatible type
		key = any(keyStr).(K)

		// Unmarshal the value dynamically
		var val V
		if err := dec.Decode(&val); err != nil {
			return err
		}

		om.Set(key, val)
	}

	// Consume the closing delim '}'
	_, err = dec.Token()
	return err
}

func (om *Map[K, V]) MarshalJSON() ([]byte, error) {
	data := []byte{}
	data = append(data, "{"...)
	numEntries := len(om.Entries())

	for i, entry := range om.Entries() {
		data = append(data, fmt.Sprintf("\"%s\":", entry.Key)...) // TODO: escape key
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(entry.Value)
		if err != nil {
			return nil, err
		}
		data = append(data, buf.Bytes()...)
		// if there are properties after this one, we need a comma
		if i < numEntries-1 {
			data = append(data, ","...)
		}

	}

	data = append(data, "}"...)

	return data, nil
}

// MarshalYAML implements the yaml interface marshaler to preserve insertion order.
//
// It returns a yaml.MapSlice rather than pre-rendered bytes so that the encoder
// walks the entries through its normal reflection-based encode path. That path
// threads the current column/indent level through nested values as it recurses.
// If instead we rendered this map to bytes independently (e.g. via yaml.Marshal)
// and returned those bytes, the resulting fragment would be parsed back into an
// AST at column 0 and then spliced into the parent document.
func (om *Map[K, V]) MarshalYAML() (any, error) {
	ms := make(yaml.MapSlice, 0, len(om.entries))
	for _, entry := range om.Entries() {
		ms = append(ms, yaml.MapItem{Key: entry.Key, Value: entry.Value})
	}

	return ms, nil
}

// UnmarshalYAML implements the yaml.Unmarshaller interface to enable yaml.Unarmshal.
func (om *Map[K, V]) UnmarshalYAML(node ast.Node) error {
	// Initialize or clear internal state
	om.keys = make(map[K]int)
	om.entries = make([]Entry[K, V], 0)

	// If the node points to a full YAML document context, unwrap its body
	if doc, ok := node.(*ast.DocumentNode); ok {
		node = doc.Body
	}

	// Ensure the node is actually a key-value mapping node block
	mapNode, ok := node.(*ast.MappingNode)
	if !ok {
		return fmt.Errorf("expected *ast.MappingNode, got %T", node)
	}

	// Loop through each structured MappingValueNode property pair
	for _, mapValueNode := range mapNode.Values {
		// 1. Decode the Key (e.g. converting ast.StringNode or ast.IntegerNode into K)
		var key K
		keyYAML := mapValueNode.Key.String()
		if err := yaml.Unmarshal([]byte(keyYAML), &key); err != nil {
			return fmt.Errorf("failed to decode map key '%s': %w", keyYAML, err)
		}

		// 2. Decode the Value into type V using its literal underlying YAML representation
		var val V
		valYAML := mapValueNode.Value.String()
		if err := yaml.Unmarshal([]byte(valYAML), &val); err != nil {
			return fmt.Errorf("failed to decode value for key '%s': %w", keyYAML, err)
		}

		om.Set(key, val)
	}

	return nil
}

// ToBuiltInMap converts the custom ordered map into a built in map type, losing order guarantees
func (om *Map[K, V]) ToBuiltInMap() map[K]V {
	m := make(map[K]V)

	for _, entry := range om.Entries() {
		m[entry.Key] = entry.Value
	}

	return m
}
