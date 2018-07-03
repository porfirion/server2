package main

import (
	"testing"
	"reflect"
)

type Typed interface {
	getType() int
}
type TestType struct {
}

func (tt TestType) getType() int {
	return 1
}

type TestType2 struct {
	TestType
}

var (
	result int

	tt  Typed = TestType{}
	tt2 Typed = TestType2{}
	tp        = reflect.TypeOf(tt)
	tp2       = reflect.TypeOf(tt2)
)

func BenchmarkSwitch_type(b *testing.B) {
	var res int
	for i := 0; i < b.N; i++ {
		switch tt.(type) {
		case TestType:
			res = 1
		case TestType2:
			res = 2
		default:
			res = 3
		}
	}
	result = res
}

func BenchmarkSwitch_field(b *testing.B) {
	var res int
	for i := 0; i < b.N; i++ {
		switch tt.getType() {
		case 1:
			res = 1
		case 2:
			res = 2
		default:
			res = 3
		}
	}
	result = res
}

func BenchmarkSwitch_typeof(b *testing.B) {
	var res int
	for i := 0; i < b.N; i++ {
		currentType := reflect.TypeOf(tt)
		if currentType == tp {
			res = 1
		} else if currentType == tp2 {
			res = 2
		} else {
			res = 3
		}
	}
	result = res
}

func BenchmarkSwitch_bool(b *testing.B) {
	var res int
	v := 3
	for i := 0; i < b.N; i++ {
		if v == 1 {
			res = 1
		} else if v == 2 {
			res = 2
		} else if v == 3 {
			res = 3
		}
	}
	result = res
}
