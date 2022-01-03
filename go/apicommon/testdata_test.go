package apicommon_test

import (
	"context"
	"testing"
	"time"

	. "github.com/go-tk/jroh/go/apicommon"
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
		vc := NewValidationContext(context.Background())
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
			vc := NewValidationContext(context.Background())
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
			vc := NewValidationContext(context.Background())
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
			vc := NewValidationContext(context.Background())
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
			vc := NewValidationContext(context.Background())
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
			vc := NewValidationContext(context.Background())
			assert.False(t, tc.M.Validate(vc))
			assert.Equal(t, tc.E, vc.ErrorDetails())
		}
	}
}

func TestModelFurtherValidation(t *testing.T) {
	{
		xs := fooapi.XString("taboo")
		vc := NewValidationContext(context.Background())
		ok := xs.Validate(vc)
		assert.False(t, ok)
		assert.Equal(t, "this is taboo!", vc.ErrorDetails())
	}
	{
		s := myStructInt32()
		s.TheInt32A = 666666
		vc := NewValidationContext(context.Background())
		ok := s.Validate(vc)
		assert.False(t, ok)
		assert.Equal(t, "theInt32A is evil!", vc.ErrorDetails())
	}
}

func TestClientActorCommunication(t *testing.T) {
	r := NewRouter()
	var tsf fooapi.TestActorFuncs
	ao := ActorOptions{TraceIDGenerator: func() string { return "xyz" }}
	fooapi.RegisterTestActor(&tsf, r, ao)
	co := ClientOptions{
		Transport: MakeInMemoryTransport(r),
	}
	c := fooapi.NewTestClient("https://localhost", co)
	tsf.DoSomething3Func = func(ctx context.Context, params *fooapi.DoSomething3Params, results *fooapi.DoSomething3Results) error {
		results.MyStructInt32 = params.MyStructInt32
		results.MyStructInt64 = params.MyStructInt64
		results.MyStructFloat32 = params.MyStructFloat32
		results.MyStructFloat64 = params.MyStructFloat64
		results.MyStructString = params.MyStructString
		results.MyOnOff = params.MyOnOff
		return nil
	}
	t1 := myStructInt32()
	t2 := myStructInt64()
	t3 := myStructFloat32()
	t4 := myStructFloat64()
	t5 := myStructString()
	params := fooapi.DoSomething3Params{
		MyStructInt32:   &t1,
		MyStructInt64:   &t2,
		MyStructFloat32: &t3,
		MyStructFloat64: &t4,
		MyStructString:  &t5,
		MyOnOff:         true,
	}
	results, err := c.DoSomething3(context.Background(), &params)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, params.MyStructInt32, results.MyStructInt32)
	assert.Equal(t, params.MyStructInt64, results.MyStructInt64)
	assert.Equal(t, params.MyStructFloat32, results.MyStructFloat32)
	assert.Equal(t, params.MyStructFloat64, results.MyStructFloat64)
	assert.Equal(t, params.MyStructString, results.MyStructString)
	assert.Equal(t, params.MyOnOff, results.MyOnOff)

	err = c.DoSomething(context.Background())
	var error *Error
	if !assert.ErrorAs(t, err, &error) {
		t.FailNow()
	}
	assert.Equal(t, ErrorNotImplemented, error.Code)
	assert.Equal(t, `rpc failed; fullMethodName="Foo.Test.DoSomething" traceID="xyz": api: not implemented`, err.Error())
}

func TestClientTimeout(t *testing.T) {
	r := NewRouter()
	var tsf fooapi.TestActorFuncs
	ao := ActorOptions{}
	fooapi.RegisterTestActor(&tsf, r, ao)
	co := ClientOptions{
		Timeout:   200 * time.Millisecond,
		Transport: MakeInMemoryTransport(r),
	}
	c := fooapi.NewTestClient("https://localhost", co)
	tsf.DoSomethingFunc = func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	}
	err := c.DoSomething(context.Background())
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestClientOutgoingRPCInfo(t *testing.T) {
	r := NewRouter()
	var tsf fooapi.TestActorFuncs
	ao := ActorOptions{}
	fooapi.RegisterTestActor(&tsf, r, ao)
	co := ClientOptions{
		Transport: MakeInMemoryTransport(r),
	}
	f := false
	params := fooapi.DoSomething3Params{MyOnOff: true}
	co.AddCommonRPCFilters(
		func(ctx context.Context, outgoingRPC *OutgoingRPC) error {
			f = true
			assert.Equal(t, "Foo", outgoingRPC.Namespace)
			assert.Equal(t, "Test", outgoingRPC.ServiceName)
			assert.Equal(t, "DoSomething3", outgoingRPC.MethodName)
			assert.Equal(t, "Foo.Test.DoSomething3", outgoingRPC.FullMethodName)
			assert.Equal(t, fooapi.Test_DoSomething3, outgoingRPC.MethodIndex)
			assert.Equal(t, &params, outgoingRPC.Params)
			assert.IsType(t, outgoingRPC.Results, (*fooapi.DoSomething3Results)(nil))
			assert.Equal(t, "https://localhost/rpc/Foo.Test.DoSomething3", outgoingRPC.URL)
			return outgoingRPC.Do(ctx)
		},
	)
	c := fooapi.NewTestClient("https://localhost", co)
	tsf.DoSomething3Func = func(ctx context.Context, params *fooapi.DoSomething3Params, results *fooapi.DoSomething3Results) error {
		return nil
	}
	_, err := c.DoSomething3(context.Background(), &params)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.True(t, f)
}

func TestActorIncomingRPCInfo(t *testing.T) {
	r := NewRouter()
	var tsf fooapi.TestActorFuncs
	ao := ActorOptions{}
	f := false
	results := fooapi.DoSomething3Results{MyOnOff: true}
	ao.AddCommonRPCFilters(
		func(ctx context.Context, incomingRPC *IncomingRPC) error {
			f = true
			assert.Equal(t, "Foo", incomingRPC.Namespace)
			assert.Equal(t, "Test", incomingRPC.ServiceName)
			assert.Equal(t, "DoSomething3", incomingRPC.MethodName)
			assert.Equal(t, "Foo.Test.DoSomething3", incomingRPC.FullMethodName)
			assert.Equal(t, fooapi.Test_DoSomething3, incomingRPC.MethodIndex)
			assert.IsType(t, incomingRPC.Params, (*fooapi.DoSomething3Params)(nil))
			assert.IsType(t, incomingRPC.Results, (*fooapi.DoSomething3Results)(nil))
			err := incomingRPC.Do(ctx)
			assert.Equal(t, &results, incomingRPC.Results)
			return err
		},
	)
	fooapi.RegisterTestActor(&tsf, r, ao)
	assert.Equal(t, []RouteInfo{
		{RPCPath: "/rpc/Foo.Test.DoSomething", FullMethodName: "Foo.Test.DoSomething", RPCFilters: []string{"github.com/go-tk/jroh/go/apicommon_test.TestActorIncomingRPCInfo.func1"}},
		{RPCPath: "/rpc/Foo.Test.DoSomething1", FullMethodName: "Foo.Test.DoSomething1", RPCFilters: []string{"github.com/go-tk/jroh/go/apicommon_test.TestActorIncomingRPCInfo.func1"}},
		{RPCPath: "/rpc/Foo.Test.DoSomething2", FullMethodName: "Foo.Test.DoSomething2", RPCFilters: []string{"github.com/go-tk/jroh/go/apicommon_test.TestActorIncomingRPCInfo.func1"}},
		{RPCPath: "/rpc/Foo.Test.DoSomething3", FullMethodName: "Foo.Test.DoSomething3", RPCFilters: []string{"github.com/go-tk/jroh/go/apicommon_test.TestActorIncomingRPCInfo.func1"}},
	}, r.RouteInfos())
	co := ClientOptions{
		Transport: MakeInMemoryTransport(r),
	}
	c := fooapi.NewTestClient("https://localhost", co)
	tsf.DoSomething3Func = func(ctx context.Context, params *fooapi.DoSomething3Params, results2 *fooapi.DoSomething3Results) error {
		*results2 = results
		return nil
	}
	params := fooapi.DoSomething3Params{MyOnOff: true}
	_, err := c.DoSomething3(context.Background(), &params)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.True(t, f)
}

func TestRPCFilters(t *testing.T) {
	r := NewRouter()
	var tsf fooapi.TestActorFuncs
	ao := ActorOptions{}
	var s string
	ao.AddCommonRPCFilters(
		func(ctx context.Context, incomingRPC *IncomingRPC) error {
			s += "a"
			err := incomingRPC.Do(ctx)
			s += "b"
			return err
		},
		func(ctx context.Context, incomingRPC *IncomingRPC) error {
			s += "c"
			err := incomingRPC.Do(ctx)
			s += "d"
			return err
		},
	)
	ao.AddRPCFilters(
		fooapi.Test_DoSomething,
		func(ctx context.Context, incomingRPC *IncomingRPC) error {
			s += "e"
			err := incomingRPC.Do(ctx)
			s += "f"
			return err
		},
		func(ctx context.Context, incomingRPC *IncomingRPC) error {
			s += "g"
			err := incomingRPC.Do(ctx)
			s += "h"
			return err
		},
	)
	fooapi.RegisterTestActor(&tsf, r, ao)
	co := ClientOptions{
		Transport: MakeInMemoryTransport(r),
	}
	co.AddCommonRPCFilters(
		func(ctx context.Context, outgoingRPC *OutgoingRPC) error {
			s += "1"
			err := outgoingRPC.Do(ctx)
			s += "2"
			return err
		},
		func(ctx context.Context, outgoingRPC *OutgoingRPC) error {
			s += "3"
			err := outgoingRPC.Do(ctx)
			s += "4"
			return err
		},
	)
	co.AddRPCFilters(
		fooapi.Test_DoSomething,
		func(ctx context.Context, outgoingRPC *OutgoingRPC) error {
			s += "5"
			err := outgoingRPC.Do(ctx)
			s += "6"
			return err
		},
		func(ctx context.Context, outgoingRPC *OutgoingRPC) error {
			s += "7"
			err := outgoingRPC.Do(ctx)
			s += "8"
			return err
		},
	)
	c := fooapi.NewTestClient("https://localhost", co)
	tsf.DoSomethingFunc = func(ctx context.Context) error {
		return nil
	}
	err := c.DoSomething(context.Background())
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, "1357aceghfdb8642", s)
}
