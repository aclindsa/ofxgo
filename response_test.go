package ofxgo_test

import (
	"fmt"
	"github.com/aclindsa/go/src/encoding/xml"
	"github.com/aclindsa/ofxgo"
	"reflect"
	"testing"
)

// Attempt to find a method on the provided Value called 'Equal' which is a
// receiver for the Value, takes one argument of the same type, and returns
// one bool. equalMethodOf() returns the nil value if the method couldn't be
// found.
func equalMethodOf(v reflect.Value) reflect.Value {
	if equalMethod, ok := v.Type().MethodByName("Equal"); ok {
		if !equalMethod.Func.IsNil() &&
			equalMethod.Type.NumIn() == 2 &&
			equalMethod.Type.In(0) == v.Type() &&
			equalMethod.Type.In(1) == v.Type() &&
			equalMethod.Type.NumOut() == 1 &&
			equalMethod.Type.Out(0).Kind() == reflect.Bool {
			return v.MethodByName("Equal")
		}
	}
	return reflect.ValueOf(nil)
}

// Attempt to return a string representation of the value appropriate for its
// type by finding a method on the provided Value called 'String' which is a
// receiver for the Value, and returns one string. stringMethodOf() returns
// fmt.Sprintf("%s", v) if it can't find a String method.
func valueToString(v reflect.Value) string {
	if equalMethod, ok := v.Type().MethodByName("String"); ok {
		if !equalMethod.Func.IsNil() &&
			equalMethod.Type.NumIn() == 1 &&
			equalMethod.Type.In(0) == v.Type() &&
			equalMethod.Type.NumOut() == 1 &&
			equalMethod.Type.Out(0).Kind() == reflect.String {
			out := v.MethodByName("String").Call([]reflect.Value{})
			return out[0].String()
		}
	}
	return fmt.Sprintf("%s", v)
}

// Recursively check that the expected and actual Values are equal in value.
// If the two Values are equal in type and contain an appropriate Equal()
// method (see equalMethodOf()), that method is used for comparison. The
// provided testing.T is failed with a message if any inequality is found.
func checkEqual(t *testing.T, fieldName string, expected, actual reflect.Value) {
	if expected.IsValid() && !actual.IsValid() {
		t.Fatalf("%s: %s was unexpectedly nil\n", t.Name(), fieldName)
	} else if !expected.IsValid() && actual.IsValid() {
		t.Fatalf("%s: Expected %s to be nil (it wasn't)\n", t.Name(), fieldName)
	} else if !expected.IsValid() && !actual.IsValid() {
		return
	}

	if expected.Type() != actual.Type() {
		t.Fatalf("%s: Expected %s type for %s, found %s\n", t.Name(), expected.Type().Name(), fieldName, actual.Type().Name())
	}

	equalMethod := equalMethodOf(expected)
	if equalMethod.IsValid() {
		in := []reflect.Value{actual}
		out := equalMethod.Call(in)
		if !out[0].Bool() {
			t.Fatalf("%s: %s !Equal(): expected '%s', got '%s'\n", t.Name(), fieldName, valueToString(expected), valueToString(actual))
		}
		return
	}

	switch expected.Kind() {
	case reflect.Array:
		for i := 0; i < expected.Len(); i++ {
			checkEqual(t, fmt.Sprintf("%s[%d]", fieldName, i), expected.Index(i), actual.Index(i))
		}
	case reflect.Slice:
		if !expected.IsNil() && actual.IsNil() {
			t.Fatalf("%s: %s was unexpectedly nil\n", t.Name(), fieldName)
		} else if expected.IsNil() && !actual.IsNil() {
			t.Fatalf("%s: Expected %s to be nil (it wasn't)\n", t.Name(), fieldName)
		}
		if expected.Len() != actual.Len() {
			t.Fatalf("%s: Expected len(%s) to to be %d, was %d\n", t.Name(), fieldName, expected.Len(), actual.Len())
		}
		for i := 0; i < expected.Len(); i++ {
			checkEqual(t, fmt.Sprintf("%s[%d]", fieldName, i), expected.Index(i), actual.Index(i))
		}
	case reflect.Interface:
		if !expected.IsNil() && actual.IsNil() {
			t.Fatalf("%s: %s was unexpectedly nil\n", t.Name(), fieldName)
		} else if expected.IsNil() && !actual.IsNil() {
			t.Fatalf("%s: Expected %s to be nil (it wasn't)\n", t.Name(), fieldName)
		}
		checkEqual(t, fieldName, expected.Elem(), actual.Elem())
	case reflect.Ptr:
		checkEqual(t, fieldName, expected.Elem(), actual.Elem())
	case reflect.Struct:
		structType := expected.Type()
		for i, n := 0, expected.NumField(); i < n; i++ {
			field := structType.Field(i)
			// skip XMLName fields so we can be lazy and not fill them out in
			// testing code
			var xmlname xml.Name
			if field.Name == "XMLName" && field.Type == reflect.TypeOf(xmlname) {
				continue
			}

			// Construct a new field name for this field, containing the parent
			// fieldName
			newFieldName := fieldName
			if fieldName != "" {
				newFieldName = fieldName + "."
			}
			newFieldName = newFieldName + field.Name
			checkEqual(t, newFieldName, expected.Field(i), actual.Field(i))
		}
	case reflect.String:
		if expected.String() != actual.String() {
			t.Fatalf("%s: %s expected to be '%s', found '%s'\n", t.Name(), fieldName, expected.String(), actual.String())
		}
	default:
		t.Fatalf("%s: %s has unexpected type that didn't provide an Equal() method: %s\n", t.Name(), fieldName, expected.Type().Name())
	}
}

func checkResponsesEqual(t *testing.T, expected, actual *ofxgo.Response) {
	checkEqual(t, "", reflect.ValueOf(expected), reflect.ValueOf(actual))
}
