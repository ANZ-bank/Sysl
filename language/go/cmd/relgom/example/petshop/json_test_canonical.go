package petshopmodel

import (
	"encoding/json"
	"fmt"
	"sort"
)

type jsonSet []interface{}

func (s jsonSet) Len() int {
	return len(s)
}

func (s jsonSet) Less(i, j int) bool {
	return jsonLess(s[i], s[j])
}

func (s jsonSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func jsonLess(a, b interface{}) bool {
	return jsonCmp(a, b) < 0
}

func jsonCmp(a, b interface{}) int {
	if a == nil || b == nil {
		return jsonBoolCmp(a != nil, b != nil)
	}
	switch u := a.(type) {
	case bool:
		if v, ok := b.(bool); ok {
			return jsonBoolCmp(u, v)
		}
	case float64:
		if v, ok := b.(float64); ok {
			return jsonFloat64Cmp(u, v)
		}
	case string:
		if v, ok := b.(string); ok {
			return jsonStringCmp(u, v)
		}
	case []interface{}:
		if v, ok := b.([]interface{}); ok {
			return jsonArrayCmp(u, v)
		}
	case map[string]interface{}:
		if v, ok := b.(map[string]interface{}); ok {
			return jsonObjectCmp(u, v)
		}
	}
	panic(fmt.Sprintf("%T ≮ %T", a, b))
}

func jsonBoolCmp(a, b bool) int {
	if !a && b {
		return -1
	}
	if !b && a {
		return 1
	}
	return 0
}

func jsonFloat64Cmp(a, b float64) int {
	if a < b {
		return -1
	}
	if b < a {
		return 1
	}
	return 0
}

func jsonStringCmp(a, b string) int {
	if a < b {
		return -1
	}
	if b < a {
		return 1
	}
	return 0
}

func jsonArrayCmp(a, b []interface{}) int {
	prefix := len(a)
	if prefix > len(b) {
		prefix = len(b)
	}
	for i := 0; i < prefix; i++ {
		cmp := jsonCmp(a[i], b[i])
		if cmp != 0 {
			return cmp
		}
	}
	if len(a) < len(b) {
		return -1
	}
	if len(b) < len(a) {
		return 1
	}
	return 0
}

func jsonObjectCmp(a, b map[string]interface{}) int {
	return jsonArrayCmp(jsonObjectAsArray(a), jsonObjectAsArray(b))
}

func jsonObjectAsArray(a map[string]interface{}) []interface{} {
	keys := make([]string, 0, len(a))
	for k := range a {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	array := make([]interface{}, 0, 2*len(a))
	for _, k := range keys {
		array = append(array, k, a[k])
	}
	return array
}

func canonicalJSON(json interface{}) interface{} {
	switch v := json.(type) {
	case []interface{}:
		w := make([]interface{}, 0, len(v))
		for _, e := range v {
			w = append(w, canonicalJSON(e))
		}
		sort.Sort(jsonSet(w))
		return w
	case map[string]interface{}:
		w := make(map[string]interface{}, len(v))
		for k, e := range v {
			w[k] = canonicalJSON(e)
		}
		return w
	case nil:
		return v
	default:
		return v
	}
}

func canonicalJSONBytes(b []byte) string {
	var j interface{}
	if err := json.Unmarshal(b, &j); err != nil {
		panic(err)
	}
	c := canonicalJSON(j)
	data, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func canonicalJSONString(s string) string {
	return canonicalJSONBytes([]byte(s))
}
