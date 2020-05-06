package petshopmodel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONCmpPanics(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() { jsonCmp(struct{}{}, struct{}{}) })
}

func TestJSONCmpNil(t *testing.T) {
	t.Parallel()

	assert.Equal(t, -1, jsonCmp(nil, 1.0))
	assert.Equal(t, 0, jsonCmp(nil, nil))
	assert.Equal(t, 1, jsonCmp(1.0, nil))
}

func TestJSONCmpBool(t *testing.T) {
	t.Parallel()

	assert.Equal(t, -1, jsonCmp(false, true))
	assert.Equal(t, 0, jsonCmp(false, false))
	assert.Equal(t, 0, jsonCmp(true, true))
	assert.Equal(t, 1, jsonCmp(true, false))
}

func TestJSONCmpFloat64(t *testing.T) {
	t.Parallel()

	assert.Equal(t, -1, jsonCmp(-1.0, 2.0))
	assert.Equal(t, -1, jsonCmp(0.0, 2.0))
	assert.Equal(t, -1, jsonCmp(1.0, 2.0))
	assert.Equal(t, 0, jsonCmp(2.0, 2.0))
	assert.Equal(t, 1, jsonCmp(2.0, 1.0))
}

func TestJSONCmpString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 0, jsonCmp("", ""))

	assert.Equal(t, -1, jsonCmp("", "x"))
	assert.Equal(t, 0, jsonCmp("x", "x"))
	assert.Equal(t, 1, jsonCmp("x", ""))

	assert.Equal(t, -1, jsonCmp("foobar", "foobaz"))
	assert.Equal(t, 0, jsonCmp("foobar", "foobar"))
	assert.Equal(t, 1, jsonCmp("foobaz", "foobar"))

	assert.Equal(t, -1, jsonCmp("foo", "foos"))
	assert.Equal(t, 1, jsonCmp("foos", "foo"))

	assert.Equal(t, -1, jsonCmp("Foo", "foo"))
	assert.Equal(t, 1, jsonCmp("foo", "Foo"))
}

func TestJSONCmpArray(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 0, jsonCmp([]interface{}{}, []interface{}{}))

	assert.Equal(t, -1, jsonCmp([]interface{}{}, []interface{}{1.0}))
	assert.Equal(t, 0, jsonCmp([]interface{}{1.0}, []interface{}{1.0}))
	assert.Equal(t, 1, jsonCmp([]interface{}{1.0}, []interface{}{}))

	assert.Equal(t, -1, jsonCmp([]interface{}{"a", "b"}, []interface{}{"a", "c"}))
	assert.Equal(t, 0, jsonCmp([]interface{}{"a", "b"}, []interface{}{"a", "b"}))
	assert.Equal(t, 1, jsonCmp([]interface{}{"a", "b"}, []interface{}{"a", "a"}))
}

func TestJSONCmpObject(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 0, jsonCmp([]interface{}{}, []interface{}{}))

	assert.Equal(t, -1, jsonCmp(map[string]interface{}{}, map[string]interface{}{"a": 1.0}))
	assert.Equal(t, 0, jsonCmp(map[string]interface{}{"a": 1.0}, map[string]interface{}{"a": 1.0}))
	assert.Equal(t, 1, jsonCmp(map[string]interface{}{"a": 1.0}, map[string]interface{}{}))

	assert.Equal(t, -1, jsonCmp(
		map[string]interface{}{"a": 1.0, "b": 1.0},
		map[string]interface{}{"a": 1.0, "b": 2.0},
	))
	assert.Equal(t, 0, jsonCmp(
		map[string]interface{}{"a": 1.0, "b": 1.0},
		map[string]interface{}{"a": 1.0, "b": 1.0},
	))
	assert.Equal(t, 1, jsonCmp(
		map[string]interface{}{"a": 1.0, "b": 2.0},
		map[string]interface{}{"a": 1.0, "b": 1.0},
	))
	assert.Equal(t, -1, jsonCmp(
		map[string]interface{}{"a": 1.0, "b": 1.0},
		map[string]interface{}{"a": 1.0, "c": 1.0},
	))
	assert.Equal(t, 0, jsonCmp(
		map[string]interface{}{"a": 1.0, "b": 1.0},
		map[string]interface{}{"a": 1.0, "b": 1.0},
	))
	assert.Equal(t, 1, jsonCmp(
		map[string]interface{}{"a": 1.0, "c": 2.0},
		map[string]interface{}{"a": 1.0, "b": 1.0},
	))
}

func TestJSONCmp(t *testing.T) {
	t.Parallel()

	assert.Equal(t, -1, jsonCmp(1.0, 2.0))
	assert.Equal(t, 0, jsonCmp(2.0, 2.0))
	assert.Equal(t, 1, jsonCmp(2.0, 1.0))
}
