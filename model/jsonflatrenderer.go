package model

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type JsonFlatRenderer struct{}

func NewJsonFlatRenderer() *JsonFlatRenderer {
	return new(JsonFlatRenderer)
}

func (r *JsonFlatRenderer) flattenHelper(root interface{}, collection *[]string, path []string) {
	switch t := root.(type) {
	case bool, float64:
		v := fmt.Sprint(strings.Join(path, ""), ": ", "[white::b]", t, "[-:-:-]")
		*collection = append(*collection, v)
	case string:
		v := fmt.Sprint(strings.Join(path, ""), ": ", "[white::b]\"", t, "\"[-:-:-]")
		*collection = append(*collection, v)
	case nil:
		v := fmt.Sprint(strings.Join(path, ""), ": ", "[white::b]null[-:-:-]")
		*collection = append(*collection, v)
	case []interface{}:
		for i, v := range t {
			r.flattenHelper(v, collection, append(path, fmt.Sprint("[", i, "[]")))
		}
	case map[string]interface{}:
		for k, v := range t {
			if path == nil {
				r.flattenHelper(v, collection, []string{k})
			} else {
				r.flattenHelper(v, collection, append(path, fmt.Sprint(".", k)))
			}
		}
	default:
		// This shouldn't happen
		*collection = append(*collection, fmt.Sprint("[red]*** Unexpected value[-]: ", strings.Join(path, ""), ".", t))
	}
}

func (r *JsonFlatRenderer) Name() string {
	return "JSON (flat)"
}

func (r *JsonFlatRenderer) Render(data []byte) ([]byte, bool) {
	var root interface{}
	if err := json.Unmarshal(data, &root); err != nil {
		return []byte(fmt.Sprintf("[red]Document is not valid JSON:[-] %s\n\n", err.Error())), true
	}

	collection := []string{}
	r.flattenHelper(root, &collection, nil)
	sort.Slice(collection, func(i, j int) bool {
		return strings.Compare(collection[i], collection[j]) < 0
	})

	return []byte(strings.Join(collection, "\n")), true
}
