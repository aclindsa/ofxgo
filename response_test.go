package ofxgo_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/aclindsa/ofxgo"
	"github.com/aclindsa/xml"
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
		t.Fatalf("%s: Expected %s type for %s, found %s\n", t.Name(), expected.Type(), fieldName, actual.Type())
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
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if expected.Uint() != actual.Uint() {
			t.Fatalf("%s: %s expected to be '%s', found '%s'\n", t.Name(), fieldName, valueToString(expected), valueToString(actual))
		}
	default:
		t.Fatalf("%s: %s has unexpected type that didn't provide an Equal() method: %s\n", t.Name(), fieldName, expected.Type().Name())
	}
}

func checkResponsesEqual(t *testing.T, expected, actual *ofxgo.Response) {
	checkEqual(t, "", reflect.ValueOf(expected), reflect.ValueOf(actual))
}

func checkResponseRoundTrip(t *testing.T, response *ofxgo.Response) {
	b, err := response.Marshal()
	if err != nil {
		t.Fatalf("Unexpected error re-marshaling OFX response: %s\n", err)
	}
	roundtripped, err := ofxgo.ParseResponse(b)
	if err != nil {
		t.Fatalf("Unexpected error re-parsing OFX response: %s\n", err)
	}
	checkResponsesEqual(t, response, roundtripped)
}

// Ensure that these samples both parse without errors, and can be converted
// back and forth without changing.
func TestValidSamples(t *testing.T) {
	fn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		} else if ext := filepath.Ext(path); ext != ".ofx" && ext != ".qfx" {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			t.Fatalf("Unexpected error opening %s: %s\n", path, err)
		}
		response, err := ofxgo.ParseResponse(file)
		if err != nil {
			t.Fatalf("Unexpected error parsing OFX response in %s: %s\n", path, err)
		}
		checkResponseRoundTrip(t, response)
		return nil
	}
	filepath.Walk("samples/valid_responses", fn)
	filepath.Walk("samples/busted_responses", fn)
}

func TestInvalidResponse(t *testing.T) {
	// in this example, the severity is invalid due to mixed upper and lower case letters
	resp, err := ofxgo.ParseResponse(bytes.NewReader([]byte(`
OFXHEADER:100
DATA:OFXSGML
VERSION:102
SECURITY:NONE
ENCODING:USASCII
CHARSET:1252
COMPRESSION:NONE
OLDFILEUID:NONE
NEWFILEUID:NONE

<OFX>
	<SIGNONMSGSRSV1>
		<SONRS>
			<STATUS>
				<CODE>0</CODE>
				<SEVERITY>Info</SEVERITY>
			</STATUS>
			<LANGUAGE>ENG</LANGUAGE>
		</SONRS>
	</SIGNONMSGSRSV1>
	<BANKMSGSRSV1>
		<STMTTRNRS>
			<TRNUID>0</TRNUID>
			<STATUS>
				<CODE>0</CODE>
				<SEVERITY>Info</SEVERITY>
			</STATUS>
		</STMTTRNRS>
	</BANKMSGSRSV1>
</OFX>
`)))
	expectedErr := "Validation failed: Invalid STATUS>SEVERITY; Invalid STATUS>SEVERITY"
	if err == nil {
		t.Fatalf("ParseResponse should fail with %q, found nil", expectedErr)
	}
	if _, ok := err.(ofxgo.ErrInvalid); !ok {
		t.Errorf("ParseResponse should return an error with type ErrInvalid, found %T", err)
	}
	if err.Error() != expectedErr {
		t.Errorf("ParseResponse should fail with %q, found %v", expectedErr, err)
	}
	if resp == nil {
		t.Errorf("Response must not be nil if only validation errors are present")
	}
}

func TestErrInvalidError(t *testing.T) {
	expectedErr := `Validation failed: A; B; C`
	actualErr := ofxgo.ErrInvalid{
		errors.New("A"),
		errors.New("B"),
		errors.New("C"),
	}.Error()
	if expectedErr != actualErr {
		t.Errorf("Unexpected invalid error message to be %q, but was: %s", expectedErr, actualErr)
	}
}

func TestErrInvalidAddErr(t *testing.T) {
	t.Run("nil error should be a no-op", func(t *testing.T) {
		var errs ofxgo.ErrInvalid
		errs.AddErr(nil)
		if len(errs) != 0 {
			t.Errorf("Nil err should not be added")
		}
	})

	t.Run("adds an error normally", func(t *testing.T) {
		var errs ofxgo.ErrInvalid
		errs.AddErr(errors.New("some error"))

	})

	t.Run("adding the same type should flatten the errors", func(t *testing.T) {
		var errs ofxgo.ErrInvalid
		errs.AddErr(ofxgo.ErrInvalid{
			errors.New("A"),
			errors.New("B"),
		})
		errs.AddErr(ofxgo.ErrInvalid{
			errors.New("C"),
		})
		if len(errs) != 3 {
			t.Errorf("Errors should be flattened like [A, B, C], but found: %+v", errs)
		}
	})
}

func TestErrInvalidErrOrNil(t *testing.T) {
	var errs ofxgo.ErrInvalid
	if err := errs.ErrOrNil(); err != nil {
		t.Errorf("No added errors should return nil, found: %v", err)
	}
	someError := errors.New("some error")
	errs.AddErr(someError)
	err := errs.ErrOrNil()
	if err == nil {
		t.Fatal("Expected an error, found nil.")
	}
	if _, ok := err.(ofxgo.ErrInvalid); !ok {
		t.Fatalf("Expected err to be of type ErrInvalid, found: %T", err)
	}
	errInvalid := err.(ofxgo.ErrInvalid)
	if len(errInvalid) != 1 || errInvalid[0] != someError {
		t.Errorf("Expected ErrOrNil to return itself, found: %v", err)
	}
}
