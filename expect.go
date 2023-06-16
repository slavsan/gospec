package gospec

import (
	"reflect"
)

type Chain struct {
	To         *Chain
	Be         *Chain
	Of         *Chain
	Have       *Chain
	Not        *Chain
	Contain    *Chain
	True       func()
	False      func()
	EqualTo    func(expected any)
	LengthOf   func(expected int)
	Property   func(expected any)
	Element    func(expected any)
	Substring  func(expected string)
	MatchError func(message string)
	Type       func(expected any)
	Nil        func()
}

func (suite *Suite) Expect(value any) *Chain {
	suite.t.Helper()
	return &Chain{
		Not: &Chain{
			To: &Chain{
				Be: &Chain{
					Nil: func() {
						// TODO: implement me
					},
				},
			},
		},
		To: &Chain{
			Contain: &Chain{
				Substring: func(sub string) {
					// TODO: implement me
				},
				Element: func(elem any) {
					// TODO: implement me
				},
			},
			MatchError: func(message string) {
				// TODO: check if value is an error
				// TODO: implement me
			},
			Have: &Chain{
				LengthOf: func(length int) {
					suite.t.Helper()

					kind := reflect.TypeOf(value).Kind()

					if kind != reflect.Slice && kind != reflect.Array && kind != reflect.String && kind != reflect.Map {
						suite.t.Errorf("expected target to be slice/array/map/string but it was %s", kind)
						return
					}

					if kind == reflect.String {
						reflectValue := reflect.ValueOf(value)
						if reflectValue.Len() != length {
							suite.t.Errorf("expected %s to have length %d but it has %d", value, length, reflectValue.Len())
						}
						return
					}

					reflectValue := reflect.ValueOf(value)
					if reflectValue.Len() != length {
						suite.t.Errorf("expected %s to have length %d but it has %d", value, length, reflectValue.Len())
					}
				},
				Property: func(prop any) {
					// TODO: implement me
				},
			},
			Be: &Chain{
				Of: &Chain{
					Type: func(expected any) {
						// TODO: implement me
					},
				},
				Nil: func() {
					valueOf := reflect.ValueOf(value)
					if !valueOf.IsValid() {
						return
					}
					switch valueOf.Kind() {
					case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
						return
					}
					suite.t.Errorf("expected %v to be nil but it is not", value)
				},
				True: func() {
					suite.t.Helper()
					v, ok := value.(bool)
					if !ok {
						suite.t.Errorf("expected test target to be bool but it was %s", reflect.TypeOf(value))
						return
					}
					if v == false {
						suite.t.Errorf("expected true but got false")
					}
				},
				False: func() {
					suite.t.Helper()
					v, ok := value.(bool)
					if !ok {
						suite.t.Errorf("expected test target to be bool but it was %s", reflect.TypeOf(value))
						return
					}
					if v != false {
						suite.t.Errorf("expected false but got true")
					}
				},
				EqualTo: func(expected any) {
					suite.t.Helper()
					if !reflect.DeepEqual(expected, value) {
						expectedType := reflect.TypeOf(expected)
						actualType := reflect.TypeOf(value)
						if expectedType != actualType {
							suite.t.Errorf("equality check failed\n\texpected: %v (type: %s)\n\t  actual: %v (type: %s)\n", expected, expectedType, value, actualType)
							return
						}
						suite.t.Errorf("equality check failed\n\texpected: %v\n\t  actual: %v\n", expected, value)
					}
				},
			},
		},
	}
}
