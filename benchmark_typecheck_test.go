package main

import (
	"reflect"
	"testing"
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

func (tt TestType2) getType() int {
	return 2
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

func BenchmarkSwitch_field2(b *testing.B) {
	var res int
	for i := 0; i < b.N; i++ {
		switch tt2.getType() {
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
		switch currentType {
		case tp:
			res = 1
		case tp2:
			res = 2
		default:
			res = 3
		}
	}
	result = res
}

func BenchmarkSwitch_int(b *testing.B) {
	var res int
	v := 3
	for i := 0; i < b.N; i++ {
		switch v {
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
