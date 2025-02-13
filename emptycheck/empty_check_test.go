package emptycheck

import (
	"reflect"
	"testing"
	"time"

	"github.com/megur0/testutil"

	"github.com/google/uuid"
)

// go test -v -run ^TestEmptyCheck ./emptycheck

// go test -v -run ^TestEmptyCheckIsStructFieldEmpty$ ./emptycheck
func TestEmptyCheckIsStructFieldEmpty(t *testing.T) {
	for _, v := range []struct {
		v      reflect.Value
		expect bool
	}{
		{
			v:      reflect.ValueOf("test"),
			expect: false,
		},
		{
			v:      reflect.ValueOf(3),
			expect: false,
		},
		{
			v:      reflect.ValueOf(Ptr(3)),
			expect: false,
		},
		{
			v:      reflect.ValueOf(Ptr("test")),
			expect: false,
		},
		{
			v:      reflect.ValueOf(time.Now()),
			expect: false,
		},
		{
			v:      reflect.ValueOf(true),
			expect: false,
		},
		{
			v:      reflect.ValueOf(false),
			expect: true,
		},
		{
			v:      reflect.ValueOf(time.Time{}),
			expect: true,
		},
		{
			v:      reflect.ValueOf(&time.Time{}),
			expect: false,
		},
		{
			v:      reflect.ValueOf((*int)(nil)),
			expect: true,
		},
		{
			v:      reflect.ValueOf((*time.Time)(nil)),
			expect: true,
		},
	} {
		t.Run("", func(t *testing.T) {
			if isStructFieldEmpty(v.v) != v.expect {
				t.Error("unexpected result")
			}
		})
	}
}

// go test -v -run ^TestEmptyCheck$ ./emptycheck
func TestEmptyCheck(t *testing.T) {
	testutil.AssertEqual(t, EmptyCheck(Ptr(struct {
		A *int
	}{})), RequiredFieldError{
		Field: "A",
	})
	testutil.AssertEqual(t, EmptyCheck(Ptr(struct {
		A string
	}{})), RequiredFieldError{
		Field: "A",
	})
	testutil.AssertEqual(t, EmptyCheck(Ptr(struct {
		A int
	}{})), RequiredFieldError{
		Field: "A",
	})
	testutil.AssertEqual(t, EmptyCheck(Ptr(struct {
		A uuid.UUID
	}{})), RequiredFieldError{
		Field: "A",
	})
	testutil.AssertEqual(t, EmptyCheck(Ptr(struct {
		A time.Time
	}{})), RequiredFieldError{
		Field: "A",
	})
	testutil.AssertEqual(t, EmptyCheck(Ptr(struct {
		A testType
	}{})).Error(), RequiredFieldError{
		Field: "A",
	}.Error())
	testutil.AssertEqual(t, EmptyCheck(Ptr(struct {
		A []int
	}{})), RequiredFieldError{
		Field: "A",
	})
	// 配列は長さが0の場合のみemptyとな。要素のゼロバリューはemptyとはならない。
	testutil.AssertUnTypedNil(t, EmptyCheck(Ptr(struct {
		A []int
	}{
		A: []int{0, 0, 0},
	})))
	testutil.AssertEqual(t, EmptyCheck(Ptr(struct {
		A []uuid.UUID
	}{})), RequiredFieldError{
		Field: "A",
	})

	testutil.AssertUnTypedNil(t, EmptyCheck(Ptr(struct {
		A []uuid.UUID
	}{
		A: []uuid.UUID{
			uuid.New(),
		},
	})))
	testutil.AssertEqual(t, EmptyCheck(Ptr(struct {
		A []struct {
			B int
		}
	}{})), RequiredFieldError{
		Field: "A",
	})
	testutil.AssertEqual(t, EmptyCheck(Ptr(struct {
		A []struct {
			B int
		}
	}{
		A: []struct {
			B int
		}{
			{},
		},
	})), RequiredFieldError{
		Field: "B",
	})
	testutil.AssertUnTypedNil(t, EmptyCheck(Ptr(struct {
		A []struct {
			B int
		}
	}{
		A: []struct {
			B int
		}{
			{
				B: 1,
			},
		},
	})))
	testutil.AssertUnTypedNil(t, EmptyCheck(Ptr(struct {
		A testType `require:"noRequired"`
	}{})))
	testutil.AssertEqual(t, EmptyCheck(Ptr(testTypeStruct{
		Child: testTypeStructChild{
			Arr: []string{"", ""},
			Child: testTypeStructChildChild{
				B: false,
			},
			Child2: []testTypeStructChildChild2{
				{
					C: true,
				},
			},
		},
		Num: 1,
		Str: "dummy",
	})), RequiredFieldError{
		Field: "B",
	})

	testutil.AssertEqual(t, EmptyCheck(Ptr(testTypeStruct{
		Child: testTypeStructChild{
			Arr: []string{"", ""},
			Child: testTypeStructChildChild{
				B: true,
			},
			Child2: []testTypeStructChildChild2{
				{
					C: false,
				},
			},
		},
		Num: 1,
		Str: "dummy",
	})), RequiredFieldError{
		Field: "C",
	})
}

type testType string

func (tt testType) IsEmpty() bool {
	return tt == ""
}

type testTypeStruct struct {
	Child testTypeStructChild
	Num   int
	Str   string
}

type testTypeStructChild struct {
	Arr    []string
	Child  testTypeStructChildChild
	Child2 []testTypeStructChildChild2
}

type testTypeStructChildChild struct {
	B bool
}

type testTypeStructChildChild2 struct {
	C bool
}

func Ptr[T any](a T) *T {
	return &a
}
