package apicommon_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/go-tk/jroh/go/apicommon/testdata/fooapi"
	"github.com/stretchr/testify/assert"
)

func myStructInt32() fooapi.MyStructInt32 {
	return fooapi.MyStructInt32{
		TheInt32A:                        int32(1),
		TheOptionalInt32A:                func(x int32) *int32 { return &x }(1),
		TheRepeatedInt32A:                []int32{1, 2},
		TheCountLimitedRepeatedInt32A:    []int32{1, 2, 3, 4},
		TheInt32B:                        fooapi.Int32(1),
		TheOptionalInt32B:                func(x fooapi.Int32) *fooapi.Int32 { return &x }(1),
		TheRepeatedInt32B:                []fooapi.Int32{1, 2},
		TheCountLimitedRepeatedInt32B:    []fooapi.Int32{1, 2, 3, 4},
		TheXInt32A:                       int32(101),
		TheOptionalXInt32A:               func(x int32) *int32 { return &x }(101),
		TheRepeatedXInt32A:               []int32{101, 102},
		TheCountLimitedRepeatedXInt32A:   []int32{101, 102, 103, 104},
		TheXInt32B:                       fooapi.XInt32(101),
		TheOptionalXInt32B:               func(x fooapi.XInt32) *fooapi.XInt32 { return &x }(101),
		TheRepeatedXInt32B:               []fooapi.XInt32{101, 102},
		TheCountLimitedRepeatedXInt32B:   []fooapi.XInt32{101, 102, 103, 104},
		TheEnumInt32:                     fooapi.EnumInt32(fooapi.C321),
		TheOptionalEnumInt32:             func(x fooapi.EnumInt32) *fooapi.EnumInt32 { return &x }(fooapi.C321),
		TheRepeatedEnumInt32:             []fooapi.EnumInt32{fooapi.C321, fooapi.C322},
		TheCountLimitedRepeatedEnumInt32: []fooapi.EnumInt32{fooapi.C321, fooapi.C321, fooapi.C322, fooapi.C322},
	}
}

func myStructInt64() fooapi.MyStructInt64 {
	return fooapi.MyStructInt64{TheInt64A: int64(1),
		TheOptionalInt64A:                func(x int64) *int64 { return &x }(1),
		TheRepeatedInt64A:                []int64{1, 2},
		TheCountLimitedRepeatedInt64A:    []int64{1, 2, 3, 4},
		TheInt64B:                        fooapi.Int64(1),
		TheOptionalInt64B:                func(x fooapi.Int64) *fooapi.Int64 { return &x }(1),
		TheRepeatedInt64B:                []fooapi.Int64{1, 2},
		TheCountLimitedRepeatedInt64B:    []fooapi.Int64{1, 2, 3, 4},
		TheXInt64A:                       int64(-101),
		TheOptionalXInt64A:               func(x int64) *int64 { return &x }(-101),
		TheRepeatedXInt64A:               []int64{-101, -102},
		TheCountLimitedRepeatedXInt64A:   []int64{-101, -102, -103, -104},
		TheXInt64B:                       fooapi.XInt64(-101),
		TheOptionalXInt64B:               func(x fooapi.XInt64) *fooapi.XInt64 { return &x }(-101),
		TheRepeatedXInt64B:               []fooapi.XInt64{-101, -102},
		TheCountLimitedRepeatedXInt64B:   []fooapi.XInt64{-101, -102, -103, -104},
		TheEnumInt64:                     fooapi.EnumInt64(fooapi.C641),
		TheOptionalEnumInt64:             func(x fooapi.EnumInt64) *fooapi.EnumInt64 { return &x }(fooapi.C641),
		TheRepeatedEnumInt64:             []fooapi.EnumInt64{fooapi.C641, fooapi.C642},
		TheCountLimitedRepeatedEnumInt64: []fooapi.EnumInt64{fooapi.C641, fooapi.C641, fooapi.C642, fooapi.C642},
	}
}

func myStructFloat32() fooapi.MyStructFloat32 {
	return fooapi.MyStructFloat32{
		TheFloat32A:                            float32(1),
		TheOptionalFloat32A:                    func(x float32) *float32 { return &x }(1),
		TheRepeatedFloat32A:                    []float32{1, 2},
		TheCountLimitedRepeatedFloat32A:        []float32{1, 2, 3, 4},
		TheFloat32B:                            fooapi.Float32(1),
		TheOptionalFloat32B:                    func(x fooapi.Float32) *fooapi.Float32 { return &x }(1),
		TheRepeatedFloat32B:                    []fooapi.Float32{1, 2},
		TheCountLimitedRepeatedFloat32B:        []fooapi.Float32{1, 2, 3, 4},
		TheXClosedFloat32A:                     float32(100),
		TheOptionalXClosedFloat32A:             func(x float32) *float32 { return &x }(100),
		TheRepeatedXClosedFloat32A:             []float32{99, 100},
		TheCountLimitedRepeatedXClosedFloat32A: []float32{97, 98, 99, 100},
		TheXClosedFloat32B:                     fooapi.XClosedFloat32(100),
		TheOptionalXClosedFloat32B:             func(x fooapi.XClosedFloat32) *fooapi.XClosedFloat32 { return &x }(100),
		TheRepeatedXClosedFloat32B:             []fooapi.XClosedFloat32{99, 100},
		TheCountLimitedRepeatedXClosedFloat32B: []fooapi.XClosedFloat32{97, 98, 99, 100},
		TheXOpenFloat32A:                       float32(2),
		TheOptionalXOpenFloat32A:               func(x float32) *float32 { return &x }(2),
		TheRepeatedXOpenFloat32A:               []float32{2, 3},
		TheCountLimitedRepeatedXOpenFloat32A:   []float32{2, 3, 4, 5},
		TheXOpenFloat32B:                       fooapi.XOpenFloat32(2),
		TheOptionalXOpenFloat32B:               func(x fooapi.XOpenFloat32) *fooapi.XOpenFloat32 { return &x }(2),
		TheRepeatedXOpenFloat32B:               []fooapi.XOpenFloat32{2, 3},
		TheCountLimitedRepeatedXOpenFloat32B:   []fooapi.XOpenFloat32{2, 3, 4, 5},
	}
}

func myStructFloat64() fooapi.MyStructFloat64 {
	return fooapi.MyStructFloat64{
		TheFloat64A:                            float64(1),
		TheOptionalFloat64A:                    func(x float64) *float64 { return &x }(1),
		TheRepeatedFloat64A:                    []float64{1, 2},
		TheCountLimitedRepeatedFloat64A:        []float64{1, 2, 3, 4},
		TheFloat64B:                            fooapi.Float64(1),
		TheOptionalFloat64B:                    func(x fooapi.Float64) *fooapi.Float64 { return &x }(1),
		TheRepeatedFloat64B:                    []fooapi.Float64{1, 2},
		TheCountLimitedRepeatedFloat64B:        []fooapi.Float64{1, 2, 3, 4},
		TheXClosedFloat64A:                     float64(-100),
		TheOptionalXClosedFloat64A:             func(x float64) *float64 { return &x }(-100),
		TheRepeatedXClosedFloat64A:             []float64{-99, -100},
		TheCountLimitedRepeatedXClosedFloat64A: []float64{-97, -98, -99, -100},
		TheXClosedFloat64B:                     fooapi.XClosedFloat64(-100),
		TheOptionalXClosedFloat64B:             func(x fooapi.XClosedFloat64) *fooapi.XClosedFloat64 { return &x }(-100),
		TheRepeatedXClosedFloat64B:             []fooapi.XClosedFloat64{-99, -100},
		TheCountLimitedRepeatedXClosedFloat64B: []fooapi.XClosedFloat64{-97, -98, -99, -100},
		TheXOpenFloat64A:                       float64(-2),
		TheOptionalXOpenFloat64A:               func(x float64) *float64 { return &x }(-2),
		TheRepeatedXOpenFloat64A:               []float64{-2, -3},
		TheCountLimitedRepeatedXOpenFloat64A:   []float64{-2, -3, -4, -5},
		TheXOpenFloat64B:                       fooapi.XOpenFloat64(-2),
		TheOptionalXOpenFloat64B:               func(x fooapi.XOpenFloat64) *fooapi.XOpenFloat64 { return &x }(-2),
		TheRepeatedXOpenFloat64B:               []fooapi.XOpenFloat64{-2, -3},
		TheCountLimitedRepeatedXOpenFloat64B:   []fooapi.XOpenFloat64{-2, -3, -4, -5},
	}
}

func myStructString() fooapi.MyStructString {
	return fooapi.MyStructString{
		TheStringA:                        string("abc"),
		TheOptionalStringA:                func(x string) *string { return &x }("abc"),
		TheRepeatedStringA:                []string{"abc", "def"},
		TheCountLimitedRepeatedStringA:    []string{"abc", "def", "ghi", "jkl"},
		TheStringB:                        fooapi.String("abc"),
		TheOptionalStringB:                func(x fooapi.String) *fooapi.String { return &x }("abc"),
		TheRepeatedStringB:                []fooapi.String{"abc", "def"},
		TheCountLimitedRepeatedStringB:    []fooapi.String{"abc", "def", "ghi", "jkl"},
		TheXStringA:                       string("abcde"),
		TheOptionalXStringA:               func(x string) *string { return &x }("abcde"),
		TheRepeatedXStringA:               []string{"abcde", "fghij"},
		TheCountLimitedRepeatedXStringA:   []string{"abcde", "fghij", "klmno", "pqrst"},
		TheXStringB:                       fooapi.XString("abcde"),
		TheOptionalXStringB:               func(x fooapi.XString) *fooapi.XString { return &x }("abcde"),
		TheRepeatedXStringB:               []fooapi.XString{"abcde", "fghij"},
		TheCountLimitedRepeatedXStringB:   []fooapi.XString{"abcde", "fghij", "klmno", "pqrst"},
		TheEnumString:                     fooapi.S1,
		TheOptionalEnumString:             func(x fooapi.EnumString) *fooapi.EnumString { return &x }(fooapi.S1),
		TheRepeatedEnumString:             []fooapi.EnumString{fooapi.S1, fooapi.S2},
		TheCountLimitedRepeatedEnumString: []fooapi.EnumString{fooapi.S1, fooapi.S2, fooapi.S1, fooapi.S2},
	}
}

func TestModelValidation(t *testing.T) {
	{
		vc := apicommon.NewValidationContext(context.Background())
		if m := myStructInt32(); !assert.True(t, m.Validate(vc)) {
			t.Fatal(vc.ErrorDetails())
		}
		if m := myStructInt64(); !assert.True(t, m.Validate(vc)) {
			t.Fatal(vc.ErrorDetails())
		}
		if m := myStructFloat32(); !assert.True(t, m.Validate(vc)) {
			t.Fatal(vc.ErrorDetails())
		}
		if m := myStructFloat64(); !assert.True(t, m.Validate(vc)) {
			t.Fatal(vc.ErrorDetails())
		}
		if m := myStructString(); !assert.True(t, m.Validate(vc)) {
			t.Fatal(vc.ErrorDetails())
		}
	}

	{
		type TC struct {
			M fooapi.MyStructInt32
			E string
		}
		var tcs []TC

		m := myStructInt32()
		m.TheCountLimitedRepeatedInt32A = m.TheCountLimitedRepeatedInt32A[:len(m.TheCountLimitedRepeatedInt32A)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedInt32A: length < 3"})

		m = myStructInt32()
		m.TheCountLimitedRepeatedInt32A = append(m.TheCountLimitedRepeatedInt32A, m.TheCountLimitedRepeatedInt32A...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedInt32A: length > 5"})

		m = myStructInt32()
		m.TheCountLimitedRepeatedInt32B = m.TheCountLimitedRepeatedInt32B[:len(m.TheCountLimitedRepeatedInt32B)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedInt32B: length < 3"})

		m = myStructInt32()
		m.TheCountLimitedRepeatedInt32B = append(m.TheCountLimitedRepeatedInt32B, m.TheCountLimitedRepeatedInt32B...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedInt32B: length > 5"})

		m = myStructInt32()
		m.TheCountLimitedRepeatedXInt32A = m.TheCountLimitedRepeatedXInt32A[:len(m.TheCountLimitedRepeatedXInt32A)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt32A: length < 3"})

		m = myStructInt32()
		m.TheCountLimitedRepeatedXInt32A = append(m.TheCountLimitedRepeatedXInt32A, m.TheCountLimitedRepeatedXInt32A...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt32A: length > 5"})

		m = myStructInt32()
		m.TheCountLimitedRepeatedXInt32B = m.TheCountLimitedRepeatedXInt32B[:len(m.TheCountLimitedRepeatedXInt32B)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt32B: length < 3"})

		m = myStructInt32()
		m.TheCountLimitedRepeatedXInt32B = append(m.TheCountLimitedRepeatedXInt32B, m.TheCountLimitedRepeatedXInt32B...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt32B: length > 5"})

		m = myStructInt32()
		m.TheCountLimitedRepeatedEnumInt32 = m.TheCountLimitedRepeatedEnumInt32[:len(m.TheCountLimitedRepeatedEnumInt32)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedEnumInt32: length < 3"})

		m = myStructInt32()
		m.TheCountLimitedRepeatedEnumInt32 = append(m.TheCountLimitedRepeatedEnumInt32, m.TheCountLimitedRepeatedEnumInt32...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedEnumInt32: length > 5"})

		m = myStructInt32()
		m.TheXInt32A = 10
		tcs = append(tcs, TC{m, "theXInt32A: value < 100"})

		m = myStructInt32()
		m.TheXInt32A = 1001
		tcs = append(tcs, TC{m, "theXInt32A: value > 999"})

		m = myStructInt32()
		*m.TheOptionalXInt32A = 10
		tcs = append(tcs, TC{m, "theOptionalXInt32A: value < 100"})

		m = myStructInt32()
		*m.TheOptionalXInt32A = 1001
		tcs = append(tcs, TC{m, "theOptionalXInt32A: value > 999"})

		m = myStructInt32()
		*m.TheOptionalXInt32B = 10
		tcs = append(tcs, TC{m, "theOptionalXInt32B: value < 100"})

		m = myStructInt32()
		*m.TheOptionalXInt32B = 1001
		tcs = append(tcs, TC{m, "theOptionalXInt32B: value > 999"})

		m = myStructInt32()
		m.TheRepeatedXInt32A[1] = 10
		tcs = append(tcs, TC{m, "theRepeatedXInt32A.1: value < 100"})

		m = myStructInt32()
		m.TheRepeatedXInt32A[1] = 1001
		tcs = append(tcs, TC{m, "theRepeatedXInt32A.1: value > 999"})

		m = myStructInt32()
		m.TheRepeatedXInt32B[1] = 10
		tcs = append(tcs, TC{m, "theRepeatedXInt32B.1: value < 100"})

		m = myStructInt32()
		m.TheRepeatedXInt32B[1] = 1001
		tcs = append(tcs, TC{m, "theRepeatedXInt32B.1: value > 999"})

		m = myStructInt32()
		m.TheCountLimitedRepeatedXInt32A[1] = 10
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt32A.1: value < 100"})

		m = myStructInt32()
		m.TheCountLimitedRepeatedXInt32A[1] = 1001
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt32A.1: value > 999"})

		m = myStructInt32()
		m.TheCountLimitedRepeatedXInt32B[1] = 10
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt32B.1: value < 100"})

		m = myStructInt32()
		m.TheCountLimitedRepeatedXInt32B[1] = 1001
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt32B.1: value > 999"})

		m = myStructInt32()
		*m.TheOptionalEnumInt32 = fooapi.C321 + fooapi.C322
		tcs = append(tcs, TC{m, "theOptionalEnumInt32: value not in (100, 200)"})

		m = myStructInt32()
		m.TheRepeatedEnumInt32[1] = fooapi.C321 + fooapi.C322
		tcs = append(tcs, TC{m, "theRepeatedEnumInt32.1: value not in (100, 200)"})

		m = myStructInt32()
		m2 := myStructInt32()
		m2.TheRepeatedEnumInt32[1] = fooapi.C321 + fooapi.C322
		m.Other = &m2
		tcs = append(tcs, TC{m, "other.theRepeatedEnumInt32.1: value not in (100, 200)"})

		m = myStructInt32()
		m2 = myStructInt32()
		m2.TheRepeatedEnumInt32[1] = fooapi.C321 + fooapi.C322
		m.Others = []fooapi.MyStructInt32{m, m2, m}
		tcs = append(tcs, TC{m, "others.1.theRepeatedEnumInt32.1: value not in (100, 200)"})

		m = myStructInt32()
		m2 = myStructInt32()
		m2.TheRepeatedEnumInt32[1] = fooapi.C321 + fooapi.C322
		m.CountLimitedOthers = []fooapi.MyStructInt32{m, m2, m}
		tcs = append(tcs, TC{m, "countLimitedOthers.1.theRepeatedEnumInt32.1: value not in (100, 200)"})

		m = myStructInt32()
		m.CountLimitedOthers = []fooapi.MyStructInt32{m, m, m, m}
		tcs = append(tcs, TC{m, "countLimitedOthers: length > 3"})

		for _, tc := range tcs {
			vc := apicommon.NewValidationContext(context.Background())
			assert.False(t, tc.M.Validate(vc))
			assert.Equal(t, tc.E, vc.ErrorDetails())
		}
	}

	{
		type TC struct {
			M fooapi.MyStructInt64
			E string
		}
		var tcs []TC

		m := myStructInt64()
		m.TheCountLimitedRepeatedInt64A = m.TheCountLimitedRepeatedInt64A[:len(m.TheCountLimitedRepeatedInt64A)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedInt64A: length < 3"})

		m = myStructInt64()
		m.TheCountLimitedRepeatedInt64A = append(m.TheCountLimitedRepeatedInt64A, m.TheCountLimitedRepeatedInt64A...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedInt64A: length > 5"})

		m = myStructInt64()
		m.TheCountLimitedRepeatedInt64B = m.TheCountLimitedRepeatedInt64B[:len(m.TheCountLimitedRepeatedInt64B)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedInt64B: length < 3"})

		m = myStructInt64()
		m.TheCountLimitedRepeatedInt64B = append(m.TheCountLimitedRepeatedInt64B, m.TheCountLimitedRepeatedInt64B...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedInt64B: length > 5"})

		m = myStructInt64()
		m.TheCountLimitedRepeatedXInt64A = m.TheCountLimitedRepeatedXInt64A[:len(m.TheCountLimitedRepeatedXInt64A)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt64A: length < 3"})

		m = myStructInt64()
		m.TheCountLimitedRepeatedXInt64A = append(m.TheCountLimitedRepeatedXInt64A, m.TheCountLimitedRepeatedXInt64A...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt64A: length > 5"})

		m = myStructInt64()
		m.TheCountLimitedRepeatedXInt64B = m.TheCountLimitedRepeatedXInt64B[:len(m.TheCountLimitedRepeatedXInt64B)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt64B: length < 3"})

		m = myStructInt64()
		m.TheCountLimitedRepeatedXInt64B = append(m.TheCountLimitedRepeatedXInt64B, m.TheCountLimitedRepeatedXInt64B...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt64B: length > 5"})

		m = myStructInt64()
		m.TheCountLimitedRepeatedEnumInt64 = m.TheCountLimitedRepeatedEnumInt64[:len(m.TheCountLimitedRepeatedEnumInt64)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedEnumInt64: length < 3"})

		m = myStructInt64()
		m.TheCountLimitedRepeatedEnumInt64 = append(m.TheCountLimitedRepeatedEnumInt64, m.TheCountLimitedRepeatedEnumInt64...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedEnumInt64: length > 5"})

		m = myStructInt64()
		m.TheXInt64A = -10
		tcs = append(tcs, TC{m, "theXInt64A: value > -100"})

		m = myStructInt64()
		m.TheXInt64A = -1001
		tcs = append(tcs, TC{m, "theXInt64A: value < -999"})

		m = myStructInt64()
		*m.TheOptionalXInt64A = -10
		tcs = append(tcs, TC{m, "theOptionalXInt64A: value > -100"})

		m = myStructInt64()
		*m.TheOptionalXInt64A = -1001
		tcs = append(tcs, TC{m, "theOptionalXInt64A: value < -999"})

		m = myStructInt64()
		*m.TheOptionalXInt64B = -10
		tcs = append(tcs, TC{m, "theOptionalXInt64B: value > -100"})

		m = myStructInt64()
		*m.TheOptionalXInt64B = -1001
		tcs = append(tcs, TC{m, "theOptionalXInt64B: value < -999"})

		m = myStructInt64()
		m.TheRepeatedXInt64A[1] = -10
		tcs = append(tcs, TC{m, "theRepeatedXInt64A.1: value > -100"})

		m = myStructInt64()
		m.TheRepeatedXInt64A[1] = -1001
		tcs = append(tcs, TC{m, "theRepeatedXInt64A.1: value < -999"})

		m = myStructInt64()
		m.TheRepeatedXInt64B[1] = -10
		tcs = append(tcs, TC{m, "theRepeatedXInt64B.1: value > -100"})

		m = myStructInt64()
		m.TheRepeatedXInt64B[1] = -1001
		tcs = append(tcs, TC{m, "theRepeatedXInt64B.1: value < -999"})

		m = myStructInt64()
		m.TheCountLimitedRepeatedXInt64A[1] = -10
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt64A.1: value > -100"})

		m = myStructInt64()
		m.TheCountLimitedRepeatedXInt64A[1] = -1001
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt64A.1: value < -999"})

		m = myStructInt64()
		m.TheCountLimitedRepeatedXInt64B[1] = -10
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt64B.1: value > -100"})

		m = myStructInt64()
		m.TheCountLimitedRepeatedXInt64B[1] = -1001
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXInt64B.1: value < -999"})

		m = myStructInt64()
		*m.TheOptionalEnumInt64 = fooapi.C641 + fooapi.C642
		tcs = append(tcs, TC{m, "theOptionalEnumInt64: value not in (200, 400)"})

		m = myStructInt64()
		m.TheRepeatedEnumInt64[1] = fooapi.C641 + fooapi.C642
		tcs = append(tcs, TC{m, "theRepeatedEnumInt64.1: value not in (200, 400)"})

		m = myStructInt64()
		m2 := myStructInt64()
		m2.TheRepeatedEnumInt64[1] = fooapi.C641 + fooapi.C642
		m.Other = &m2
		tcs = append(tcs, TC{m, "other.theRepeatedEnumInt64.1: value not in (200, 400)"})

		m = myStructInt64()
		m2 = myStructInt64()
		m2.TheRepeatedEnumInt64[1] = fooapi.C641 + fooapi.C642
		m.Others = []fooapi.MyStructInt64{m, m2, m}
		tcs = append(tcs, TC{m, "others.1.theRepeatedEnumInt64.1: value not in (200, 400)"})

		m = myStructInt64()
		m2 = myStructInt64()
		m2.TheRepeatedEnumInt64[1] = fooapi.C641 + fooapi.C642
		m.CountLimitedOthers = []fooapi.MyStructInt64{m, m2, m}
		tcs = append(tcs, TC{m, "countLimitedOthers.1.theRepeatedEnumInt64.1: value not in (200, 400)"})

		m = myStructInt64()
		m.CountLimitedOthers = []fooapi.MyStructInt64{m, m, m, m}
		tcs = append(tcs, TC{m, "countLimitedOthers: length > 3"})

		for _, tc := range tcs {
			vc := apicommon.NewValidationContext(context.Background())
			assert.False(t, tc.M.Validate(vc))
			assert.Equal(t, tc.E, vc.ErrorDetails())
		}
	}

	{
		type TC struct {
			M fooapi.MyStructFloat32
			E string
		}
		var tcs []TC

		m := myStructFloat32()
		m.TheCountLimitedRepeatedFloat32A = m.TheCountLimitedRepeatedFloat32A[:len(m.TheCountLimitedRepeatedFloat32A)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedFloat32A: length < 3"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedFloat32A = append(m.TheCountLimitedRepeatedFloat32A, m.TheCountLimitedRepeatedFloat32A...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedFloat32A: length > 5"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedFloat32B = m.TheCountLimitedRepeatedFloat32B[:len(m.TheCountLimitedRepeatedFloat32B)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedFloat32B: length < 3"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedFloat32B = append(m.TheCountLimitedRepeatedFloat32B, m.TheCountLimitedRepeatedFloat32B...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedFloat32B: length > 5"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXOpenFloat32A = m.TheCountLimitedRepeatedXOpenFloat32A[:len(m.TheCountLimitedRepeatedXOpenFloat32A)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat32A: length < 3"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXOpenFloat32A = append(m.TheCountLimitedRepeatedXOpenFloat32A, m.TheCountLimitedRepeatedXOpenFloat32A...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat32A: length > 5"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXOpenFloat32B = m.TheCountLimitedRepeatedXOpenFloat32B[:len(m.TheCountLimitedRepeatedXOpenFloat32B)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat32B: length < 3"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXOpenFloat32B = append(m.TheCountLimitedRepeatedXOpenFloat32B, m.TheCountLimitedRepeatedXOpenFloat32B...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat32B: length > 5"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXClosedFloat32A = m.TheCountLimitedRepeatedXClosedFloat32A[:len(m.TheCountLimitedRepeatedXClosedFloat32A)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat32A: length < 3"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXClosedFloat32A = append(m.TheCountLimitedRepeatedXClosedFloat32A, m.TheCountLimitedRepeatedXClosedFloat32A...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat32A: length > 5"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXClosedFloat32B = m.TheCountLimitedRepeatedXClosedFloat32B[:len(m.TheCountLimitedRepeatedXClosedFloat32B)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat32B: length < 3"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXClosedFloat32B = append(m.TheCountLimitedRepeatedXClosedFloat32B, m.TheCountLimitedRepeatedXClosedFloat32B...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat32B: length > 5"})

		m = myStructFloat32()
		m.TheXClosedFloat32A = 0.9
		tcs = append(tcs, TC{m, "theXClosedFloat32A: value < 1"})

		m = myStructFloat32()
		m.TheXClosedFloat32A = 100.1
		tcs = append(tcs, TC{m, "theXClosedFloat32A: value > 100"})

		m = myStructFloat32()
		*m.TheOptionalXClosedFloat32A = 0.9
		tcs = append(tcs, TC{m, "theOptionalXClosedFloat32A: value < 1"})

		m = myStructFloat32()
		*m.TheOptionalXClosedFloat32A = 100.1
		tcs = append(tcs, TC{m, "theOptionalXClosedFloat32A: value > 100"})

		m = myStructFloat32()
		m.TheRepeatedXClosedFloat32A[1] = 0.9
		tcs = append(tcs, TC{m, "theRepeatedXClosedFloat32A.1: value < 1"})

		m = myStructFloat32()
		m.TheRepeatedXClosedFloat32A[1] = 100.1
		tcs = append(tcs, TC{m, "theRepeatedXClosedFloat32A.1: value > 100"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXClosedFloat32A[1] = 0.9
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat32A.1: value < 1"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXClosedFloat32A[1] = 100.1
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat32A.1: value > 100"})

		m = myStructFloat32()
		*m.TheOptionalXClosedFloat32B = 0.9
		tcs = append(tcs, TC{m, "theOptionalXClosedFloat32B: value < 1"})

		m = myStructFloat32()
		*m.TheOptionalXClosedFloat32B = 100.1
		tcs = append(tcs, TC{m, "theOptionalXClosedFloat32B: value > 100"})

		m = myStructFloat32()
		m.TheRepeatedXClosedFloat32B[1] = 0.9
		tcs = append(tcs, TC{m, "theRepeatedXClosedFloat32B.1: value < 1"})

		m = myStructFloat32()
		m.TheRepeatedXClosedFloat32B[1] = 100.1
		tcs = append(tcs, TC{m, "theRepeatedXClosedFloat32B.1: value > 100"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXClosedFloat32B[1] = 0.9
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat32B.1: value < 1"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXClosedFloat32B[1] = 100.1
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat32B.1: value > 100"})

		m = myStructFloat32()
		m.TheXOpenFloat32A = 1.0
		tcs = append(tcs, TC{m, "theXOpenFloat32A: value <= 1"})

		m = myStructFloat32()
		m.TheXOpenFloat32A = 100.0
		tcs = append(tcs, TC{m, "theXOpenFloat32A: value >= 100"})

		m = myStructFloat32()
		*m.TheOptionalXOpenFloat32A = 1.0
		tcs = append(tcs, TC{m, "theOptionalXOpenFloat32A: value <= 1"})

		m = myStructFloat32()
		*m.TheOptionalXOpenFloat32A = 100.0
		tcs = append(tcs, TC{m, "theOptionalXOpenFloat32A: value >= 100"})

		m = myStructFloat32()
		m.TheRepeatedXOpenFloat32A[1] = 1.0
		tcs = append(tcs, TC{m, "theRepeatedXOpenFloat32A.1: value <= 1"})

		m = myStructFloat32()
		m.TheRepeatedXOpenFloat32A[1] = 100.0
		tcs = append(tcs, TC{m, "theRepeatedXOpenFloat32A.1: value >= 100"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXOpenFloat32A[1] = 1.0
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat32A.1: value <= 1"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXOpenFloat32A[1] = 100.0
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat32A.1: value >= 100"})

		m = myStructFloat32()
		*m.TheOptionalXOpenFloat32B = 1.0
		tcs = append(tcs, TC{m, "theOptionalXOpenFloat32B: value <= 1"})

		m = myStructFloat32()
		*m.TheOptionalXOpenFloat32B = 100.0
		tcs = append(tcs, TC{m, "theOptionalXOpenFloat32B: value >= 100"})

		m = myStructFloat32()
		m.TheRepeatedXOpenFloat32B[1] = 1.0
		tcs = append(tcs, TC{m, "theRepeatedXOpenFloat32B.1: value <= 1"})

		m = myStructFloat32()
		m.TheRepeatedXOpenFloat32B[1] = 100.0
		tcs = append(tcs, TC{m, "theRepeatedXOpenFloat32B.1: value >= 100"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXOpenFloat32B[1] = 1.0
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat32B.1: value <= 1"})

		m = myStructFloat32()
		m.TheCountLimitedRepeatedXOpenFloat32B[1] = 100.0
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat32B.1: value >= 100"})

		m = myStructFloat32()
		m2 := myStructFloat32()
		m2.TheCountLimitedRepeatedXOpenFloat32B[1] = 100.0
		m.Other = &m2
		tcs = append(tcs, TC{m, "other.theCountLimitedRepeatedXOpenFloat32B.1: value >= 100"})

		m = myStructFloat32()
		m2 = myStructFloat32()
		m2.TheCountLimitedRepeatedXOpenFloat32B[1] = 100.0
		m.Others = []fooapi.MyStructFloat32{m, m2, m}
		tcs = append(tcs, TC{m, "others.1.theCountLimitedRepeatedXOpenFloat32B.1: value >= 100"})

		m = myStructFloat32()
		m2 = myStructFloat32()
		m2.TheCountLimitedRepeatedXOpenFloat32B[1] = 100.0
		m.CountLimitedOthers = []fooapi.MyStructFloat32{m, m2, m}
		tcs = append(tcs, TC{m, "countLimitedOthers.1.theCountLimitedRepeatedXOpenFloat32B.1: value >= 100"})

		m = myStructFloat32()
		m.CountLimitedOthers = []fooapi.MyStructFloat32{m, m, m, m}
		tcs = append(tcs, TC{m, "countLimitedOthers: length > 3"})

		for _, tc := range tcs {
			vc := apicommon.NewValidationContext(context.Background())
			assert.False(t, tc.M.Validate(vc))
			assert.Equal(t, tc.E, vc.ErrorDetails())
		}
	}

	{
		type TC struct {
			M fooapi.MyStructFloat64
			E string
		}
		var tcs []TC

		m := myStructFloat64()
		m.TheCountLimitedRepeatedFloat64A = m.TheCountLimitedRepeatedFloat64A[:len(m.TheCountLimitedRepeatedFloat64A)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedFloat64A: length < 3"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedFloat64A = append(m.TheCountLimitedRepeatedFloat64A, m.TheCountLimitedRepeatedFloat64A...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedFloat64A: length > 5"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedFloat64B = m.TheCountLimitedRepeatedFloat64B[:len(m.TheCountLimitedRepeatedFloat64B)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedFloat64B: length < 3"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedFloat64B = append(m.TheCountLimitedRepeatedFloat64B, m.TheCountLimitedRepeatedFloat64B...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedFloat64B: length > 5"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXOpenFloat64A = m.TheCountLimitedRepeatedXOpenFloat64A[:len(m.TheCountLimitedRepeatedXOpenFloat64A)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat64A: length < 3"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXOpenFloat64A = append(m.TheCountLimitedRepeatedXOpenFloat64A, m.TheCountLimitedRepeatedXOpenFloat64A...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat64A: length > 5"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXOpenFloat64B = m.TheCountLimitedRepeatedXOpenFloat64B[:len(m.TheCountLimitedRepeatedXOpenFloat64B)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat64B: length < 3"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXOpenFloat64B = append(m.TheCountLimitedRepeatedXOpenFloat64B, m.TheCountLimitedRepeatedXOpenFloat64B...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat64B: length > 5"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXClosedFloat64A = m.TheCountLimitedRepeatedXClosedFloat64A[:len(m.TheCountLimitedRepeatedXClosedFloat64A)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat64A: length < 3"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXClosedFloat64A = append(m.TheCountLimitedRepeatedXClosedFloat64A, m.TheCountLimitedRepeatedXClosedFloat64A...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat64A: length > 5"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXClosedFloat64B = m.TheCountLimitedRepeatedXClosedFloat64B[:len(m.TheCountLimitedRepeatedXClosedFloat64B)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat64B: length < 3"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXClosedFloat64B = append(m.TheCountLimitedRepeatedXClosedFloat64B, m.TheCountLimitedRepeatedXClosedFloat64B...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat64B: length > 5"})

		m = myStructFloat64()
		m.TheXClosedFloat64A = -0.9
		tcs = append(tcs, TC{m, "theXClosedFloat64A: value > -1"})

		m = myStructFloat64()
		m.TheXClosedFloat64A = -100.1
		tcs = append(tcs, TC{m, "theXClosedFloat64A: value < -100"})

		m = myStructFloat64()
		*m.TheOptionalXClosedFloat64A = -0.9
		tcs = append(tcs, TC{m, "theOptionalXClosedFloat64A: value > -1"})

		m = myStructFloat64()
		*m.TheOptionalXClosedFloat64A = -100.1
		tcs = append(tcs, TC{m, "theOptionalXClosedFloat64A: value < -100"})

		m = myStructFloat64()
		m.TheRepeatedXClosedFloat64A[1] = -0.9
		tcs = append(tcs, TC{m, "theRepeatedXClosedFloat64A.1: value > -1"})

		m = myStructFloat64()
		m.TheRepeatedXClosedFloat64A[1] = -100.1
		tcs = append(tcs, TC{m, "theRepeatedXClosedFloat64A.1: value < -100"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXClosedFloat64A[1] = -0.9
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat64A.1: value > -1"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXClosedFloat64A[1] = -100.1
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat64A.1: value < -100"})

		m = myStructFloat64()
		*m.TheOptionalXClosedFloat64B = -0.9
		tcs = append(tcs, TC{m, "theOptionalXClosedFloat64B: value > -1"})

		m = myStructFloat64()
		*m.TheOptionalXClosedFloat64B = -100.1
		tcs = append(tcs, TC{m, "theOptionalXClosedFloat64B: value < -100"})

		m = myStructFloat64()
		m.TheRepeatedXClosedFloat64B[1] = -0.9
		tcs = append(tcs, TC{m, "theRepeatedXClosedFloat64B.1: value > -1"})

		m = myStructFloat64()
		m.TheRepeatedXClosedFloat64B[1] = -100.1
		tcs = append(tcs, TC{m, "theRepeatedXClosedFloat64B.1: value < -100"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXClosedFloat64B[1] = -0.9
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat64B.1: value > -1"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXClosedFloat64B[1] = -100.1
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXClosedFloat64B.1: value < -100"})

		m = myStructFloat64()
		m.TheXOpenFloat64A = -1.0
		tcs = append(tcs, TC{m, "theXOpenFloat64A: value >= -1"})

		m = myStructFloat64()
		m.TheXOpenFloat64A = -100.0
		tcs = append(tcs, TC{m, "theXOpenFloat64A: value <= -100"})

		m = myStructFloat64()
		*m.TheOptionalXOpenFloat64A = -1.0
		tcs = append(tcs, TC{m, "theOptionalXOpenFloat64A: value >= -1"})

		m = myStructFloat64()
		*m.TheOptionalXOpenFloat64A = -100.0
		tcs = append(tcs, TC{m, "theOptionalXOpenFloat64A: value <= -100"})

		m = myStructFloat64()
		m.TheRepeatedXOpenFloat64A[1] = -1.0
		tcs = append(tcs, TC{m, "theRepeatedXOpenFloat64A.1: value >= -1"})

		m = myStructFloat64()
		m.TheRepeatedXOpenFloat64A[1] = -100.0
		tcs = append(tcs, TC{m, "theRepeatedXOpenFloat64A.1: value <= -100"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXOpenFloat64A[1] = -1.0
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat64A.1: value >= -1"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXOpenFloat64A[1] = -100.0
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat64A.1: value <= -100"})

		m = myStructFloat64()
		*m.TheOptionalXOpenFloat64B = -1.0
		tcs = append(tcs, TC{m, "theOptionalXOpenFloat64B: value >= -1"})

		m = myStructFloat64()
		*m.TheOptionalXOpenFloat64B = -100.0
		tcs = append(tcs, TC{m, "theOptionalXOpenFloat64B: value <= -100"})

		m = myStructFloat64()
		m.TheRepeatedXOpenFloat64B[1] = -1.0
		tcs = append(tcs, TC{m, "theRepeatedXOpenFloat64B.1: value >= -1"})

		m = myStructFloat64()
		m.TheRepeatedXOpenFloat64B[1] = -100.0
		tcs = append(tcs, TC{m, "theRepeatedXOpenFloat64B.1: value <= -100"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXOpenFloat64B[1] = -1.0
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat64B.1: value >= -1"})

		m = myStructFloat64()
		m.TheCountLimitedRepeatedXOpenFloat64B[1] = -100.0
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXOpenFloat64B.1: value <= -100"})

		m = myStructFloat64()
		m2 := myStructFloat64()
		m2.TheCountLimitedRepeatedXOpenFloat64B[1] = -100.0
		m.Other = &m2
		tcs = append(tcs, TC{m, "other.theCountLimitedRepeatedXOpenFloat64B.1: value <= -100"})

		m = myStructFloat64()
		m2 = myStructFloat64()
		m2.TheCountLimitedRepeatedXOpenFloat64B[1] = -100.0
		m.Others = []fooapi.MyStructFloat64{m, m2, m}
		tcs = append(tcs, TC{m, "others.1.theCountLimitedRepeatedXOpenFloat64B.1: value <= -100"})

		m = myStructFloat64()
		m2 = myStructFloat64()
		m2.TheCountLimitedRepeatedXOpenFloat64B[1] = -100.0
		m.CountLimitedOthers = []fooapi.MyStructFloat64{m, m2, m}
		tcs = append(tcs, TC{m, "countLimitedOthers.1.theCountLimitedRepeatedXOpenFloat64B.1: value <= -100"})

		m = myStructFloat64()
		m.CountLimitedOthers = []fooapi.MyStructFloat64{m, m, m, m}
		tcs = append(tcs, TC{m, "countLimitedOthers: length > 3"})

		for _, tc := range tcs {
			vc := apicommon.NewValidationContext(context.Background())
			assert.False(t, tc.M.Validate(vc))
			assert.Equal(t, tc.E, vc.ErrorDetails())
		}
	}

	{
		type TC struct {
			M fooapi.MyStructString
			E string
		}
		var tcs []TC

		m := myStructString()
		m.TheCountLimitedRepeatedStringA = m.TheCountLimitedRepeatedStringA[:len(m.TheCountLimitedRepeatedStringA)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedStringA: length < 3"})

		m = myStructString()
		m.TheCountLimitedRepeatedStringA = append(m.TheCountLimitedRepeatedStringA, m.TheCountLimitedRepeatedStringA...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedStringA: length > 5"})

		m = myStructString()
		m.TheCountLimitedRepeatedStringB = m.TheCountLimitedRepeatedStringB[:len(m.TheCountLimitedRepeatedStringB)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedStringB: length < 3"})

		m = myStructString()
		m.TheCountLimitedRepeatedStringB = append(m.TheCountLimitedRepeatedStringB, m.TheCountLimitedRepeatedStringB...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedStringB: length > 5"})

		m = myStructString()
		m.TheCountLimitedRepeatedXStringA = m.TheCountLimitedRepeatedXStringA[:len(m.TheCountLimitedRepeatedXStringA)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXStringA: length < 3"})

		m = myStructString()
		m.TheCountLimitedRepeatedXStringA = append(m.TheCountLimitedRepeatedXStringA, m.TheCountLimitedRepeatedXStringA...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXStringA: length > 5"})

		m = myStructString()
		m.TheCountLimitedRepeatedXStringB = m.TheCountLimitedRepeatedXStringB[:len(m.TheCountLimitedRepeatedXStringB)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXStringB: length < 3"})

		m = myStructString()
		m.TheCountLimitedRepeatedXStringB = append(m.TheCountLimitedRepeatedXStringB, m.TheCountLimitedRepeatedXStringB...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXStringB: length > 5"})

		m = myStructString()
		m.TheCountLimitedRepeatedEnumString = m.TheCountLimitedRepeatedEnumString[:len(m.TheCountLimitedRepeatedEnumString)-2]
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedEnumString: length < 3"})

		m = myStructString()
		m.TheCountLimitedRepeatedEnumString = append(m.TheCountLimitedRepeatedEnumString, m.TheCountLimitedRepeatedEnumString...)
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedEnumString: length > 5"})

		m = myStructString()
		m.TheXStringA = "ab"
		tcs = append(tcs, TC{m, "theXStringA: length < 3"})

		m = myStructString()
		m.TheXStringA = "0123456789"
		tcs = append(tcs, TC{m, "theXStringA: length > 9"})

		m = myStructString()
		m.TheXStringA = "a.c"
		tcs = append(tcs, TC{m, "theXStringA: value not matched to \"[a-zA-Z0-9]*\""})

		m = myStructString()
		*m.TheOptionalXStringA = "ab"
		tcs = append(tcs, TC{m, "theOptionalXStringA: length < 3"})

		m = myStructString()
		*m.TheOptionalXStringA = "0123456789"
		tcs = append(tcs, TC{m, "theOptionalXStringA: length > 9"})

		m = myStructString()
		*m.TheOptionalXStringA = "a.c"
		tcs = append(tcs, TC{m, "theOptionalXStringA: value not matched to \"[a-zA-Z0-9]*\""})

		m = myStructString()
		*m.TheOptionalXStringB = "ab"
		tcs = append(tcs, TC{m, "theOptionalXStringB: length < 3"})

		m = myStructString()
		*m.TheOptionalXStringB = "0123456789"
		tcs = append(tcs, TC{m, "theOptionalXStringB: length > 9"})

		m = myStructString()
		*m.TheOptionalXStringB = "a.c"
		tcs = append(tcs, TC{m, "theOptionalXStringB: value not matched to \"[a-zA-Z0-9]*\""})

		m = myStructString()
		m.TheRepeatedXStringA[1] = "ab"
		tcs = append(tcs, TC{m, "theRepeatedXStringA.1: length < 3"})

		m = myStructString()
		m.TheRepeatedXStringA[1] = "0123456789"
		tcs = append(tcs, TC{m, "theRepeatedXStringA.1: length > 9"})

		m = myStructString()
		m.TheRepeatedXStringA[1] = "a.c"
		tcs = append(tcs, TC{m, "theRepeatedXStringA.1: value not matched to \"[a-zA-Z0-9]*\""})

		m = myStructString()
		m.TheRepeatedXStringB[1] = "ab"
		tcs = append(tcs, TC{m, "theRepeatedXStringB.1: length < 3"})

		m = myStructString()
		m.TheRepeatedXStringB[1] = "0123456789"
		tcs = append(tcs, TC{m, "theRepeatedXStringB.1: length > 9"})

		m = myStructString()
		m.TheRepeatedXStringB[1] = "a.c"
		tcs = append(tcs, TC{m, "theRepeatedXStringB.1: value not matched to \"[a-zA-Z0-9]*\""})

		m = myStructString()
		m.TheCountLimitedRepeatedXStringA[1] = "ab"
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXStringA.1: length < 3"})

		m = myStructString()
		m.TheCountLimitedRepeatedXStringA[1] = "0123456789"
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXStringA.1: length > 9"})

		m = myStructString()
		m.TheCountLimitedRepeatedXStringA[1] = "a.c"
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXStringA.1: value not matched to \"[a-zA-Z0-9]*\""})

		m = myStructString()
		m.TheCountLimitedRepeatedXStringB[1] = "ab"
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXStringB.1: length < 3"})

		m = myStructString()
		m.TheCountLimitedRepeatedXStringB[1] = "0123456789"
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXStringB.1: length > 9"})

		m = myStructString()
		m.TheCountLimitedRepeatedXStringB[1] = "a.c"
		tcs = append(tcs, TC{m, "theCountLimitedRepeatedXStringB.1: value not matched to \"[a-zA-Z0-9]*\""})

		m = myStructString()
		*m.TheOptionalEnumString = fooapi.S1 + fooapi.S2
		tcs = append(tcs, TC{m, `theOptionalEnumString: value not in ("abc", "def")`})

		m = myStructString()
		m.TheRepeatedEnumString[1] = fooapi.S1 + fooapi.S2
		tcs = append(tcs, TC{m, `theRepeatedEnumString.1: value not in ("abc", "def")`})

		m = myStructString()
		m2 := myStructString()
		m2.TheRepeatedEnumString[1] = fooapi.S1 + fooapi.S1
		m.Other = &m2
		tcs = append(tcs, TC{m, `other.theRepeatedEnumString.1: value not in ("abc", "def")`})

		m = myStructString()
		m2 = myStructString()
		m2.TheRepeatedEnumString[1] = fooapi.S1 + fooapi.S1
		m.Others = []fooapi.MyStructString{m, m2, m}
		tcs = append(tcs, TC{m, `others.1.theRepeatedEnumString.1: value not in ("abc", "def")`})

		m = myStructString()
		m2 = myStructString()
		m2.TheRepeatedEnumString[1] = fooapi.S1 + fooapi.S1
		m.CountLimitedOthers = []fooapi.MyStructString{m, m2, m}
		tcs = append(tcs, TC{m, `countLimitedOthers.1.theRepeatedEnumString.1: value not in ("abc", "def")`})

		m = myStructString()
		m.CountLimitedOthers = []fooapi.MyStructString{m, m, m, m}
		tcs = append(tcs, TC{m, "countLimitedOthers: length > 3"})

		for _, tc := range tcs {
			vc := apicommon.NewValidationContext(context.Background())
			assert.False(t, tc.M.Validate(vc))
			assert.Equal(t, tc.E, vc.ErrorDetails())
		}
	}
}

func TestModelMarshalingAndUnmarshaling(t *testing.T) {
	tsf := fooapi.TestServerFuncs{
		DoSomething2Func: func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
			results.MyStructInt32 = params.MyStructInt32
			results.MyStructInt64 = params.MyStructInt64
			results.MyStructFloat32 = params.MyStructFloat32
			results.MyStructFloat64 = params.MyStructFloat64
			results.MyStructString = params.MyStructString
			results.MyOnOff = params.MyOnOff
			return nil
		},
	}
	sm := http.NewServeMux()
	fooapi.RegisterTestServer(&tsf, sm, apicommon.ServerOptions{})
	s := http.Server{
		Addr:    "127.0.0.1:7890",
		Handler: sm,
	}
	go s.ListenAndServe()
	defer s.Shutdown(context.Background())
	var err error
	for i := 0; i < 5; i++ {
		_, err := http.Get("http://127.0.0.1:7890")
		if err != nil {
			t.Log(err)
			time.Sleep(time.Second / 10)
			continue
		}
		break
	}
	if err != nil {
		t.Fatal(err)
	}

	c := fooapi.NewTestClient("http://127.0.0.1:7890", apicommon.ClientOptions{})
	t1 := myStructInt32()
	t2 := myStructInt64()
	t3 := myStructFloat32()
	t4 := myStructFloat64()
	t5 := myStructString()
	params := fooapi.DoSomething2Params{
		MyStructInt32:   &t1,
		MyStructInt64:   &t2,
		MyStructFloat32: &t3,
		MyStructFloat64: &t4,
		MyStructString:  &t5,
		MyOnOff:         true,
	}
	results, err := c.DoSomething2(context.Background(), &params)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, params.MyStructInt32, results.MyStructInt32)
	assert.Equal(t, params.MyStructInt64, results.MyStructInt64)
	assert.Equal(t, params.MyStructFloat32, results.MyStructFloat32)
	assert.Equal(t, params.MyStructFloat64, results.MyStructFloat64)
	assert.Equal(t, params.MyStructString, results.MyStructString)
	assert.Equal(t, params.MyOnOff, results.MyOnOff)
}

func TestError(t *testing.T) {
	setup := func(tsf fooapi.TestServerFuncs, port uint16, serverOptions apicommon.ServerOptions, clientOptions apicommon.ClientOptions) (fooapi.TestClient, func()) {
		sm := http.NewServeMux()
		fooapi.RegisterTestServer(&tsf, sm, serverOptions)
		s := http.Server{
			Addr:    fmt.Sprintf("127.0.0.1:%d", port),
			Handler: sm,
		}
		go s.ListenAndServe()
		var err error
		for i := 0; i < 5; i++ {
			_, err = http.Get(fmt.Sprintf("http://127.0.0.1:%d", port))
			if err != nil {
				t.Log(err)
				time.Sleep(time.Second / 10)
				continue
			}
			break
		}
		if err != nil {
			s.Shutdown(context.Background())
			t.Fatal(err)
		}
		c := fooapi.NewTestClient(fmt.Sprintf("http://127.0.0.1:%d", port), clientOptions)
		return c, func() { s.Shutdown(context.Background()) }
	}

	func() {
		c, cleanup := setup(fooapi.TestServerFuncs{}, 7890, apicommon.ServerOptions{}, apicommon.ClientOptions{})
		defer cleanup()
		t1 := myStructString()
		t1.TheXStringA = "a.c"
		t2 := myStructString()
		t2.Others = []fooapi.MyStructString{t1}
		params := fooapi.DoSomething2Params{
			MyStructString: &t2,
		}
		_, err := c.DoSomething2(context.Background(), &params)
		if !assert.Error(t, err) {
			t.FailNow()
		}
		var error *apicommon.Error
		if !assert.ErrorAs(t, err, &error) {
			t.Fatal(err)
		}
		assert.ErrorIs(t, err, apicommon.ErrInvalidParams)
		assert.Equal(t, "myStructString.others.0.theXStringA: value not matched to \"[a-zA-Z0-9]*\"", error.Details)
	}()

	func() {
		c, cleanup := setup(fooapi.TestServerFuncs{}, 7890, apicommon.ServerOptions{
			Middlewares: map[apicommon.MethodIndex][]apicommon.ServerMiddleware{
				fooapi.Test_DoSomething2: {
					func(handler http.Handler) http.Handler {
						return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							incomingRPC := apicommon.MustGetRPCFromContext(r.Context()).IncomingRPC()
							data := incomingRPC.RawParams()
							data[0] = ','
							handler.ServeHTTP(w, r)
						})
					},
				},
			},
		}, apicommon.ClientOptions{})
		defer cleanup()
		t1 := myStructString()
		params := fooapi.DoSomething2Params{
			MyStructString: &t1,
		}
		_, err := c.DoSomething2(context.Background(), &params)
		if !assert.Error(t, err) {
			t.FailNow()
		}
		var error *apicommon.Error
		if !assert.ErrorAs(t, err, &error) {
			t.Fatal(err)
		}
		assert.ErrorIs(t, err, apicommon.ErrParse)
		assert.Equal(t, "invalid character ',' looking for beginning of value", error.Details)
	}()

	func() {
		c, cleanup := setup(fooapi.TestServerFuncs{
			DoSomething3Func: func(context.Context) error {
				err := fooapi.NewSomethingWrongError()
				err.Details = "hello world"
				err.Data.SetValue("foo", "bar")
				return err
			},
		}, 7890, apicommon.ServerOptions{}, apicommon.ClientOptions{})
		defer cleanup()
		err := c.DoSomething3(context.Background())
		if !assert.Error(t, err) {
			t.FailNow()
		}
		var error *apicommon.Error
		if !assert.ErrorAs(t, err, &error) {
			t.Fatal(err)
		}
		assert.ErrorIs(t, err, fooapi.ErrSomethingWrong)
		assert.Equal(t, "hello world", error.Details)
		assert.Equal(t, apicommon.ErrorData{"foo": "bar"}, error.Data)
	}()

	func() {
		apicommon.DebugMode = true
		c, cleanup := setup(fooapi.TestServerFuncs{
			DoSomething3Func: func(context.Context) error {
				err := fooapi.NewSomethingWrongError()
				err.Details = "hello world"
				err.Data.SetValue("foo", "bar")
				return fmt.Errorf("err: %w", err)
			},
		}, 7890, apicommon.ServerOptions{}, apicommon.ClientOptions{})
		defer cleanup()
		err := c.DoSomething3(context.Background())
		if !assert.Error(t, err) {
			t.FailNow()
		}
		var error *apicommon.Error
		if !assert.ErrorAs(t, err, &error) {
			t.Fatal(err)
		}
		assert.ErrorIs(t, err, apicommon.ErrInternal)
		assert.Equal(t, "err: api: something wrong (1): hello world", error.Details)

		apicommon.DebugMode = false
		err = c.DoSomething3(context.Background())
		if !assert.Error(t, err) {
			t.FailNow()
		}
		if !assert.ErrorAs(t, err, &error) {
			t.Fatal(err)
		}
		assert.ErrorIs(t, err, apicommon.ErrInternal)
		assert.Equal(t, "", error.Details)
		assert.Equal(t, apicommon.ErrorData(nil), error.Data)
	}()

	func() {
		apicommon.DebugMode = true
		c, cleanup := setup(fooapi.TestServerFuncs{
			DoSomething3Func: func(context.Context) error {
				panic("NOOOOO!")
			},
		}, 7890, apicommon.ServerOptions{}, apicommon.ClientOptions{})
		defer cleanup()
		err := c.DoSomething3(context.Background())
		if !assert.Error(t, err) {
			t.FailNow()
		}
		var error *apicommon.Error
		if !assert.ErrorAs(t, err, &error) {
			t.Fatal(err)
		}
		assert.ErrorIs(t, err, apicommon.ErrInternal)
		assert.Equal(t, "NOOOOO!", error.Details)
		stackTrace := error.Data["stackTrace"]
		if assert.IsType(t, string(""), stackTrace) {
			stackTrace := stackTrace.(string)
			assert.True(t, strings.HasPrefix(stackTrace, "goroutine "))
		}

		apicommon.DebugMode = false
		err = c.DoSomething3(context.Background())
		if !assert.ErrorAs(t, err, &error) {
			t.Fatal(err)
		}
	}()

	func() {
		c, cleanup := setup(fooapi.TestServerFuncs{
			DoSomething2Func: func(ctx context.Context, _ *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
				t2 := myStructString()
				t2.TheXStringA = "a.c"
				*results = fooapi.DoSomething2Results{
					MyStructString: &t2,
				}
				return nil
			},
		}, 7890, apicommon.ServerOptions{}, apicommon.ClientOptions{})
		defer cleanup()
		t2 := myStructString()
		params := fooapi.DoSomething2Params{
			MyStructString: &t2,
		}
		_, err := c.DoSomething2(context.Background(), &params)
		_ = err
		if !assert.Error(t, err) {
			t.FailNow()
		}
		var error *apicommon.Error
		if errors.As(err, &error) {
			t.Fatal(err)
		}
		assert.True(t, strings.HasSuffix(err.Error(), "myStructString.theXStringA: value not matched to \"[a-zA-Z0-9]*\""))
	}()
}

func TestTraceID(t *testing.T) {
	setup := func(tsf fooapi.TestServerFuncs, port uint16, serverOptions apicommon.ServerOptions, clientOptions apicommon.ClientOptions) (fooapi.TestClient, func()) {
		sm := http.NewServeMux()
		fooapi.RegisterTestServer(&tsf, sm, serverOptions)
		s := http.Server{
			Addr:    fmt.Sprintf("127.0.0.1:%d", port),
			Handler: sm,
		}
		go s.ListenAndServe()
		var err error
		for i := 0; i < 5; i++ {
			_, err = http.Get(fmt.Sprintf("http://127.0.0.1:%d", port))
			if err != nil {
				t.Log(err)
				time.Sleep(time.Second / 10)
				continue
			}
			break
		}
		if err != nil {
			s.Shutdown(context.Background())
			t.Fatal(err)
		}
		c := fooapi.NewTestClient(fmt.Sprintf("http://127.0.0.1:%d", port), clientOptions)
		return c, func() { s.Shutdown(context.Background()) }
	}

	var k int
	var traceIDs []string
	c2, cleanup2 := setup(fooapi.TestServerFuncs{
		DoSomething3Func: func(ctx context.Context) error {
			rpc := apicommon.MustGetRPCFromContext(ctx)
			traceIDs = append(traceIDs, "A-"+rpc.TraceID())
			return nil
		},
	}, 7891, apicommon.ServerOptions{}, apicommon.ClientOptions{})
	defer cleanup2()
	c1, cleanup1 := setup(fooapi.TestServerFuncs{
		DoSomething3Func: func(ctx context.Context) error {
			rpc := apicommon.MustGetRPCFromContext(ctx)
			traceIDs = append(traceIDs, "B1-"+rpc.TraceID())
			err := c2.DoSomething3(ctx)
			traceIDs = append(traceIDs, "B2-"+rpc.TraceID())
			return err
		},
	}, 7890, apicommon.ServerOptions{
		TraceIDGenerator: func() string {
			k++
			return fmt.Sprintf("My-Trace-ID-%d", k)
		},
	}, apicommon.ClientOptions{
		RPCFilters: map[apicommon.MethodIndex][]apicommon.RPCHandler{
			fooapi.Test_DoSomething3: {
				func(ctx context.Context, rpc *apicommon.RPC) error {
					traceIDs = append(traceIDs, "C1-"+rpc.TraceID())
					err := rpc.Do(ctx)
					traceIDs = append(traceIDs, "C2-"+rpc.TraceID())
					return err
				},
			},
		},
	})
	defer cleanup1()
	err := c1.DoSomething3(context.Background())
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
	assert.Equal(t, []string{
		"C1-",
		"B1-My-Trace-ID-1",
		"A-My-Trace-ID-1",
		"B2-My-Trace-ID-1",
		"C2-My-Trace-ID-1",
	}, traceIDs)
}

func TestMiddlewareAndRPCFilter(t *testing.T) {
	setup := func(tsf fooapi.TestServerFuncs, port uint16, serverOptions apicommon.ServerOptions, clientOptions apicommon.ClientOptions) (fooapi.TestClient, func()) {
		sm := http.NewServeMux()
		fooapi.RegisterTestServer(&tsf, sm, serverOptions)
		s := http.Server{
			Addr:    fmt.Sprintf("127.0.0.1:%d", port),
			Handler: sm,
		}
		go s.ListenAndServe()
		var err error
		for i := 0; i < 5; i++ {
			_, err = http.Get(fmt.Sprintf("http://127.0.0.1:%d", port))
			if err != nil {
				t.Log(err)
				time.Sleep(time.Second / 10)
				continue
			}
			break
		}
		if err != nil {
			s.Shutdown(context.Background())
			t.Fatal(err)
		}
		c := fooapi.NewTestClient(fmt.Sprintf("http://127.0.0.1:%d", port), clientOptions)
		return c, func() { s.Shutdown(context.Background()) }
	}
	var k int
	var s string
	c, cleanup := setup(
		fooapi.TestServerFuncs{
			DoSomething2Func: func(ctx context.Context, params *fooapi.DoSomething2Params, results *fooapi.DoSomething2Results) error {
				assert.NotNil(t, params)
				assert.NotNil(t, results)
				s += "XX"
				return nil
			},
		},
		7891,
		apicommon.ServerOptions{
			Middlewares: map[apicommon.MethodIndex][]apicommon.ServerMiddleware{
				apicommon.AnyMethod: {
					func(handler http.Handler) http.Handler {
						return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							incomingRPC := apicommon.MustGetRPCFromContext(r.Context()).IncomingRPC()
							assert.Equal(t, "Foo", incomingRPC.Namespace())
							assert.Equal(t, "Test", incomingRPC.ServiceName())
							assert.Equal(t, "DoSomething2", incomingRPC.MethodName())
							assert.Equal(t, "My-Trace-ID-1", incomingRPC.TraceID())
							assert.NotNil(t, incomingRPC.RawParams())
							assert.NotNil(t, incomingRPC.Params())
							assert.Nil(t, incomingRPC.RawResp())
							assert.NotNil(t, incomingRPC.Results())
							s += "A1"
							handler.ServeHTTP(w, r)
							assert.NotNil(t, incomingRPC.RawResp())
							s += "B1"
						})
					},
					func(handler http.Handler) http.Handler {
						return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							s += "A2"
							handler.ServeHTTP(w, r)
							s += "B2"
						})
					},
				},
				fooapi.Test_DoSomething2: {
					func(handler http.Handler) http.Handler {
						return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							s += "C1"
							handler.ServeHTTP(w, r)
							s += "D1"
						})
					},
					func(handler http.Handler) http.Handler {
						return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							s += "C2"
							handler.ServeHTTP(w, r)
							s += "D2"
						})
					},
				},
			},
			RPCFilters: map[apicommon.MethodIndex][]apicommon.RPCHandler{
				apicommon.AnyMethod: {
					func(ctx context.Context, rpc *apicommon.RPC) error {
						incomingRPC := rpc.IncomingRPC()
						assert.Equal(t, "Foo", incomingRPC.Namespace())
						assert.Equal(t, "Test", incomingRPC.ServiceName())
						assert.Equal(t, "DoSomething2", incomingRPC.MethodName())
						assert.Equal(t, "My-Trace-ID-1", incomingRPC.TraceID())
						assert.NotNil(t, incomingRPC.RawParams())
						assert.NotNil(t, incomingRPC.Params())
						assert.Nil(t, incomingRPC.RawResp())
						assert.NotNil(t, incomingRPC.Results())
						s += "E1"
						err := rpc.Do(ctx)
						assert.Nil(t, incomingRPC.RawResp())
						s += "F1"
						return err
					},
					func(ctx context.Context, rpc *apicommon.RPC) error {
						s += "E2"
						err := rpc.Do(ctx)
						s += "F2"
						return err
					},
				},
				fooapi.Test_DoSomething2: {
					func(ctx context.Context, rpc *apicommon.RPC) error {
						s += "G1"
						err := rpc.Do(ctx)
						s += "H1"
						return err
					},
					func(ctx context.Context, rpc *apicommon.RPC) error {
						s += "G2"
						err := rpc.Do(ctx)
						s += "H2"
						return err
					},
				},
			},
			TraceIDGenerator: func() string {
				k++
				return fmt.Sprintf("My-Trace-ID-%d", k)
			},
		},
		apicommon.ClientOptions{
			RPCFilters: map[apicommon.MethodIndex][]apicommon.RPCHandler{
				apicommon.AnyMethod: {
					func(ctx context.Context, rpc *apicommon.RPC) error {
						outgoingRPC := rpc.OutgoingRPC()
						assert.Equal(t, "Foo", outgoingRPC.Namespace())
						assert.Equal(t, "Test", outgoingRPC.ServiceName())
						assert.Equal(t, "DoSomething2", outgoingRPC.MethodName())
						assert.Equal(t, "", outgoingRPC.TraceID())
						assert.Nil(t, outgoingRPC.RawParams())
						assert.NotNil(t, outgoingRPC.Params())
						assert.Nil(t, outgoingRPC.RawResp())
						assert.NotNil(t, outgoingRPC.Results())
						s += "I1"
						err := rpc.Do(ctx)
						assert.NotNil(t, outgoingRPC.RawParams())
						assert.NotNil(t, outgoingRPC.RawResp())
						s += "J1"
						assert.Equal(t, "My-Trace-ID-1", outgoingRPC.TraceID())
						return err
					},
					func(ctx context.Context, rpc *apicommon.RPC) error {
						s += "I2"
						err := rpc.Do(ctx)
						s += "J2"
						return err
					},
				},
				fooapi.Test_DoSomething2: {
					func(ctx context.Context, rpc *apicommon.RPC) error {
						s += "K1"
						err := rpc.Do(ctx)
						s += "L1"
						return err
					},
					func(ctx context.Context, rpc *apicommon.RPC) error {
						s += "K2"
						err := rpc.Do(ctx)
						s += "L2"
						return err
					},
				},
			},
			Middlewares: map[apicommon.MethodIndex][]apicommon.ClientMiddleware{
				apicommon.AnyMethod: {
					func(transport http.RoundTripper) http.RoundTripper {
						return apicommon.TransportFunc(func(request *http.Request) (*http.Response, error) {
							outgoingRPC := apicommon.MustGetRPCFromContext(request.Context()).OutgoingRPC()
							assert.Equal(t, "Foo", outgoingRPC.Namespace())
							assert.Equal(t, "Test", outgoingRPC.ServiceName())
							assert.Equal(t, "DoSomething2", outgoingRPC.MethodName())
							assert.Equal(t, "", outgoingRPC.TraceID())
							assert.NotNil(t, outgoingRPC.RawParams())
							assert.NotNil(t, outgoingRPC.Params())
							assert.Nil(t, outgoingRPC.RawResp())
							assert.NotNil(t, outgoingRPC.Results())
							s += "M1"
							response, err := transport.RoundTrip(request)
							assert.NotNil(t, outgoingRPC.RawParams())
							assert.NotNil(t, outgoingRPC.RawResp())
							s += "N1"
							assert.Equal(t, "", outgoingRPC.TraceID())
							return response, err
						})
					},
					func(transport http.RoundTripper) http.RoundTripper {
						return apicommon.TransportFunc(func(request *http.Request) (*http.Response, error) {
							s += "M2"
							response, err := transport.RoundTrip(request)
							s += "N2"
							return response, err
						})
					},
				},
				fooapi.Test_DoSomething2: {
					func(transport http.RoundTripper) http.RoundTripper {
						return apicommon.TransportFunc(func(request *http.Request) (*http.Response, error) {
							s += "O1"
							response, err := transport.RoundTrip(request)
							s += "P1"
							return response, err
						})
					},
					func(transport http.RoundTripper) http.RoundTripper {
						return apicommon.TransportFunc(func(request *http.Request) (*http.Response, error) {
							s += "O2"
							response, err := transport.RoundTrip(request)
							s += "P2"
							return response, err
						})
					},
				},
			},
		},
	)
	defer cleanup()
	t1 := myStructString()
	params := fooapi.DoSomething2Params{
		MyStructString: &t1,
	}
	_, err := c.DoSomething2(context.Background(), &params)
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
	assert.Equal(t, "I1I2K1K2M1M2O1O2A1A2C1C2E1E2G1G2XXH2H1F2F1D2D1B2B1P2P1N2N1L2L1J2J1", s)
}
