package eval

import (
	"reflect"
	"regexp"
	"strings"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type goFunc struct {
	val  reflect.Value
	args []*sysl.Type
	ret  *sysl.Type
}

func kindToPrimitiveType(kind reflect.Kind) (sysl.Type_Primitive, bool) {
	switch kind {
	case reflect.Bool:
		return sysl.Type_BOOL, true
	case reflect.Int:
		return sysl.Type_INT, true
	case reflect.Int16:
		return sysl.Type_INT, true
	case reflect.Int32:
		return sysl.Type_INT, true
	case reflect.Int64:
		return sysl.Type_INT, true
	case reflect.String:
		return sysl.Type_STRING, true
	default:
		return 0, false
	}
}

func valueTypeToPrimitiveType(t valueType) (sysl.Type_Primitive, bool) {
	switch t {
	case ValueBool:
		return sysl.Type_BOOL, true
	case ValueInt:
		return sysl.Type_INT, true
	case ValueFloat:
		return sysl.Type_FLOAT, true
	case ValueString:
		return sysl.Type_STRING, true
	default:
		return 0, false
	}
}

//nolint:gochecknoglobals
var (
	stringType = &sysl.Type{
		Type: &sysl.Type_Primitive_{
			Primitive: sysl.Type_STRING,
		},
	}

	intType = &sysl.Type{
		Type: &sysl.Type_Primitive_{
			Primitive: sysl.Type_INT,
		},
	}

	listStringType = &sysl.Type{
		Type: &sysl.Type_List_{
			List: &sysl.Type_List{
				Type: &sysl.Type{
					Type: &sysl.Type_Primitive_{Primitive: sysl.Type_STRING},
				},
			},
		},
	}

	boolType = &sysl.Type{Type: &sysl.Type_Primitive_{Primitive: sysl.Type_BOOL}}

	// GoFuncMap is Map of Function Names to goFunc{}
	GoFuncMap = map[string]goFunc{
		"Contains":      {reflect.ValueOf(strings.Contains), []*sysl.Type{stringType, stringType}, boolType},
		"Count":         {reflect.ValueOf(strings.Count), []*sysl.Type{stringType, stringType}, intType},
		"Fields":        {reflect.ValueOf(strings.Fields), []*sysl.Type{stringType}, listStringType},
		"FindAllString": {reflect.ValueOf(FindAllString), []*sysl.Type{stringType, stringType, intType}, listStringType},
		"HasPrefix":     {reflect.ValueOf(strings.HasPrefix), []*sysl.Type{stringType, stringType}, boolType},
		"HasSuffix":     {reflect.ValueOf(strings.HasSuffix), []*sysl.Type{stringType, stringType}, boolType},
		"Join":          {reflect.ValueOf(strings.Join), []*sysl.Type{listStringType, stringType}, stringType},
		"LastIndex":     {reflect.ValueOf(strings.LastIndex), []*sysl.Type{stringType, stringType}, intType},
		"MatchString":   {reflect.ValueOf(MatchString), []*sysl.Type{stringType, stringType}, boolType},
		"Replace":       {reflect.ValueOf(strings.Replace), []*sysl.Type{stringType, stringType, stringType, intType}, stringType}, // nolint:lll
		"Split":         {reflect.ValueOf(strings.Split), []*sysl.Type{stringType, stringType}, listStringType},
		"Title":         {reflect.ValueOf(strings.Title), []*sysl.Type{stringType}, stringType},
		"ToLower":       {reflect.ValueOf(strings.ToLower), []*sysl.Type{stringType}, stringType},
		"ToTitle":       {reflect.ValueOf(strings.ToTitle), []*sysl.Type{stringType}, stringType},
		"ToUpper":       {reflect.ValueOf(strings.ToUpper), []*sysl.Type{stringType}, stringType},
		"Trim":          {reflect.ValueOf(strings.Trim), []*sysl.Type{stringType, stringType}, stringType},
		"TrimLeft":      {reflect.ValueOf(strings.TrimLeft), []*sysl.Type{stringType, stringType}, stringType},
		"TrimPrefix":    {reflect.ValueOf(strings.TrimPrefix), []*sysl.Type{stringType, stringType}, stringType},
		"TrimRight":     {reflect.ValueOf(strings.TrimRight), []*sysl.Type{stringType, stringType}, stringType},
		"TrimSpace":     {reflect.ValueOf(strings.TrimSpace), []*sysl.Type{stringType}, stringType},
		"TrimSuffix":    {reflect.ValueOf(strings.TrimSuffix), []*sysl.Type{stringType, stringType}, stringType},
	}
)

func valueToReflectValue(v *sysl.Value, t *sysl.Type) reflect.Value {
	switch x := t.Type.(type) {
	case *sysl.Type_Primitive_:
		switch x.Primitive {
		case sysl.Type_BOOL:
			return reflect.ValueOf(v.GetB())
		case sysl.Type_INT:
			return reflect.ValueOf(int(v.GetI()))
		case sysl.Type_STRING:
			return reflect.ValueOf(v.GetS())
		}
	case *sysl.Type_List_:
		listOf := x.List.Type
		var valueList []*sysl.Value
		if v.GetList() != nil {
			valueList = v.GetList().Value
		} else {
			valueList = v.GetSet().Value
		}

		if p := listOf.GetPrimitive(); p != sysl.Type_NO_Primitive {
			if p == sysl.Type_STRING {
				stringSlice := []string{}
				for _, listItem := range valueList {
					stringSlice = append(stringSlice, listItem.GetS())
				}
				return reflect.ValueOf(stringSlice)
			}
		}
	}
	panic(errors.Errorf("valueToReflectValue: unsupported value type: %v", v))
}

func isValueExpectedType(v *sysl.Value, t *sysl.Type) bool {
	vType := getValueType(v)
	type1, has := valueTypeToPrimitiveType(vType)
	inType := t.GetPrimitive()
	// if both are primitive types
	if has && inType != sysl.Type_NO_Primitive {
		return inType == type1
	}

	if vType == ValueList && len(v.GetList().Value) == 0 {
		return true
	}

	if vType == ValueSet && len(v.GetSet().Value) == 0 {
		return true
	}

	if vType == ValueList && t.GetList() != nil {
		return isValueExpectedType(v.GetList().Value[0], t.GetList().Type)
	}

	if vType == ValueSet && t.GetSet() != nil {
		return isValueExpectedType(v.GetSet().Value[0], t.GetSet())
	}

	// Pass elements of sets as a slice to Go funcs.
	if vType == ValueSet && t.GetList() != nil {
		return isValueExpectedType(v.GetSet().Value[0], t.GetList().Type)
	}
	return false
}

func isReflectValueExpectedType(r reflect.Value, typ *sysl.Type) bool {
	kind := r.Kind()
	_, has := kindToPrimitiveType(kind)
	if has {
		return true
	}
	if kind == reflect.Slice {
		if typ.GetList() == nil || typ.GetSet() != nil {
			return false
		}
		_, has = kindToPrimitiveType(r.Index(0).Kind())
		return has
	}
	return false
}

func reflectToValue(r reflect.Value, typ *sysl.Type) *sysl.Value {
	if !isReflectValueExpectedType(r, typ) {
		logrus.Warnf("Got %s, Expected Value type: %v \n", r.Kind(), typ.Type)
	}

	switch r.Kind() {
	case reflect.Bool:
		return MakeValueBool(r.Bool())
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ.GetPrimitive() != sysl.Type_INT {
			logrus.Warnf("Got int, Expected Value type: %v \n", typ.Type)
		}
		return MakeValueI64(r.Int())
	case reflect.String:
		if typ.GetPrimitive() != sysl.Type_STRING {
			logrus.Warnf("Got string, Expected Value type: %v \n", typ.Type)
		}
		return MakeValueString(r.String())
	case reflect.Slice:
		list := MakeValueList()
		primitiveType := typ.GetList().GetType()

		for i := 0; i < r.Len(); i++ {
			AppendItemToValueList(list.GetList(), reflectToValue(r.Index(i), primitiveType))
		}
		return list
	}
	panic(errors.Errorf("reflectToValue: kind %s not supported\n", r.Kind().String()))
}

func evalGoFunc(name string, list *sysl.Value) *sysl.Value {
	if f, has := GoFuncMap[name]; has {
		var in []reflect.Value
		if len(list.GetList().Value) != len(f.args) {
			logrus.Errorf("Incorrect number of arg function %s\n", name)
			logrus.Errorf("Expected number of arg: (%d)\n", len(f.args))
			logrus.Errorf("Got num args: (%d)\n", len(list.GetList().Value))
			return nil
		}

		for i, l := range list.GetList().Value {
			if isValueExpectedType(l, f.args[i]) {
				in = append(in, valueToReflectValue(l, f.args[i]))
			} else {
				return nil
			}
		}
		result := f.val.Call(in)
		if len(result) != 1 {
			logrus.Errorf("Function %s has %d return values, expecting 1\n", name, len(result))
		}
		return reflectToValue(result[0], f.ret)
	}
	return nil
}

// MatchString exposes regexp.MatchString to sysl transforms
func MatchString(pattern, word string) bool {
	return regexp.MustCompile(pattern).MatchString(word)
}

// FindAllString exposes regexp.FindAllString to sysl transforms
func FindAllString(pattern, word string, n int) []string {
	return regexp.MustCompile(pattern).FindAllString(word, n)
}
