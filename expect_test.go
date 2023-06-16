package gospec

import (
	"fmt"
	"testing"
)

type mockTesting struct {
	calls []string
}

var _ testingInterface = (*mockTesting)(nil)

func (m *mockTesting) Helper() {
}

func (m *mockTesting) Errorf(format string, args ...interface{}) {
	m.calls = append(m.calls, fmt.Sprintf(format, args...))
}

func (m *mockTesting) Run(name string, f func(t *testing.T)) bool {
	return false
}

func TestExpect(t *testing.T) {
	mock := new(mockTesting)
	s := NewTestSuite(t)
	s.t = mock
	expect := s.Expect

	t.Run("expect.To.Be.True", func(t *testing.T) {
		testCases := []struct {
			input      interface{}
			callsCount int
			calls      []string
		}{
			{true, 0, []string{}},
			{false, 1, []string{"expected true but got false"}},
			{"foo", 1, []string{"expected test target to be bool but it was string"}},
		}
		for _, tc := range testCases {
			tc := tc
			t.Run("", func(t *testing.T) {
				mock.calls = []string{}
				expect(tc.input).To.Be.True()
				if len(mock.calls) != tc.callsCount {
					t.Errorf("expected Errorf to have been called %d times but it was called %d times", tc.callsCount, len(mock.calls))
				}
				for i, x := range tc.calls {
					if mock.calls[i] != x {
						t.Errorf("expected \"%s\" but got \"%s\"", x, mock.calls[i])
					}
				}
			})
		}
	})
	t.Run("expect.To.Be.False", func(t *testing.T) {
		testCases := []struct {
			input      interface{}
			callsCount int
			calls      []string
		}{
			{false, 0, []string{}},
			{true, 1, []string{"expected false but got true"}},
			{"foo", 1, []string{"expected test target to be bool but it was string"}},
		}
		for _, tc := range testCases {
			tc := tc
			t.Run("", func(t *testing.T) {
				mock.calls = []string{}
				expect(tc.input).To.Be.False()
				if len(mock.calls) != tc.callsCount {
					t.Errorf("expected Errorf to have been called %d times but it was called %d times", tc.callsCount, len(mock.calls))
				}
				for i, x := range tc.calls {
					if mock.calls[i] != x {
						t.Errorf("expected \"%s\" but got \"%s\"", x, mock.calls[i])
					}
				}
			})
		}
	})
	t.Run("expect.To.Be.EqualTo", func(t *testing.T) {
		testCases := []struct {
			left       interface{}
			right      interface{}
			callsCount int
			calls      []string
		}{
			{false, false, 0, []string{}},
			{true, true, 0, []string{}},
			{true, false, 1, []string{"equality check failed\n\texpected: false\n\t  actual: true\n"}},
			{333, 333, 0, []string{}},
			{int32(333), int32(333), 0, []string{}},
			{333, 334, 1, []string{"equality check failed\n\texpected: 334\n\t  actual: 333\n"}},
			{int32(333), int64(333), 1, []string{"equality check failed\n\texpected: 333 (type: int64)\n\t  actual: 333 (type: int32)\n"}},
			{"foo", "foo", 0, []string{}},
			{"foo", "bar", 1, []string{"equality check failed\n\texpected: bar\n\t  actual: foo\n"}},
			{"foo", 333, 1, []string{"equality check failed\n\texpected: 333 (type: int)\n\t  actual: foo (type: string)\n"}},
			{[]string{}, []string{}, 0, []string{}},
			{[]string{}, []string{"foo"}, 1, []string{"equality check failed\n\texpected: [foo]\n\t  actual: []\n"}},
		}
		for _, tc := range testCases {
			tc := tc
			t.Run("", func(t *testing.T) {
				mock.calls = []string{}
				expect(tc.left).To.Be.EqualTo(tc.right)
				if len(mock.calls) != tc.callsCount {
					t.Errorf("expected Errorf to have been called %d times but it was called %d times", tc.callsCount, len(mock.calls))
				}
				for i, x := range tc.calls {
					if mock.calls[i] != x {
						t.Errorf("expected \"%s\" but got \"%s\"", x, mock.calls[i])
					}
				}
			})
		}
	})
	t.Run("expect.To.Be.Nil", func(t *testing.T) {
		type someType struct{ X int }
		var m map[string]int
		var x *someType
		var y []string
		var f func()
		var i interface{}
		testCases := []struct {
			input      interface{}
			callsCount int
			calls      []string
		}{
			{false, 1, []string{"expected false to be nil but it is not"}},
			{nil, 0, []string{}},
			{i, 0, []string{}},
			{m, 0, []string{}},
			{f, 0, []string{}},
			{x, 0, []string{}},
			{y, 0, []string{}},
		}
		for _, tc := range testCases {
			tc := tc
			t.Run("", func(t *testing.T) {
				mock.calls = []string{}
				expect(tc.input).To.Be.Nil()
				if len(mock.calls) != tc.callsCount {
					t.Errorf("expected Errorf to have been called %d times but it was called %d times", tc.callsCount, len(mock.calls))
				}
				for i, x := range tc.calls {
					if mock.calls[i] != x {
						t.Errorf("expected \"%s\" but got \"%s\"", x, mock.calls[i])
					}
				}
			})
		}
	})
	t.Run("expect.To.Have.LengthOf", func(t *testing.T) {
		testCases := []struct {
			input      interface{}
			length     int
			callsCount int
			calls      []string
		}{
			{false, 0, 1, []string{"expected target to be slice/array/map/string but it was bool"}},
			{[]string{}, 0, 0, []string{}},
			{[]string{"foo"}, 1, 0, []string{}},
			{[]string{"foo", "bar"}, 2, 0, []string{}},
			{[]string{"foo", "bar"}, 3, 1, []string{"expected [foo bar] to have length 3 but it has 2"}},
			{[]int{2, 52, 12, 9}, 4, 0, []string{}},
			{"", 0, 0, []string{}},
			{"a", 1, 0, []string{}},
			{"foo bar baz", 11, 0, []string{}},
			{"foo bar baz", 12, 1, []string{"expected foo bar baz to have length 12 but it has 11"}},
			{map[string]int{}, 0, 0, []string{}},
			{map[string]int{"x": 1}, 1, 0, []string{}},
			{map[string]int{"x": 10, "y": 20, "z": 30}, 3, 0, []string{}},
		}
		for _, tc := range testCases {
			tc := tc
			t.Run("", func(t *testing.T) {
				mock.calls = []string{}
				expect(tc.input).To.Have.LengthOf(tc.length)
				if len(mock.calls) != tc.callsCount {
					t.Errorf("expected Errorf to have been called %d times but it was called %d times (with: %#v)", tc.callsCount, len(mock.calls), mock.calls)
				}
				for i, x := range tc.calls {
					if mock.calls[i] != x {
						t.Errorf("expected \"%s\" but got \"%s\"", x, mock.calls[i])
					}
				}
			})
		}
	})
}
