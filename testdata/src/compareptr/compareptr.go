package compareptr

import "fmt"

type MyStruct struct {
	Value *int64
}

func (m *MyStruct) GetValue() int64 {
	return *m.Value
}

type OtherStruct struct {
	Number *int32
}

func (o *OtherStruct) GetNumber() int32 {
	return *o.Number
}

func Demo() {
	val1 := int64(2)
	val2 := int64(1)

	s1 := &MyStruct{Value: &val1}
	s2 := &MyStruct{Value: &val2}

	// This should trigger the linter - comparing struct field pointers with getters
	if s1.Value == s2.Value { // want "comparing pointers with == or != can be error-prone"
		fmt.Println("ok")
	}
}

func DemoNoTrigger() {
	// Simple pointer comparison - should NOT trigger (not struct fields)
	var a int64 = 2
	var b int64 = 1
	var c = &a
	var d = &b
	if c == d {
		fmt.Println("simple pointer comparison")
	}

	// Struct field without getter - should NOT trigger
	type NoGetterStruct struct {
		Field *int64
	}
	ng1 := &NoGetterStruct{Field: &a}
	ng2 := &NoGetterStruct{Field: &b}
	if ng1.Field == ng2.Field {
		fmt.Println("struct field without getter")
	}

	// Non-numeric pointer comparison - should NOT trigger
	type StringStruct struct {
		Text *string
	}
	str1 := "hello"
	str2 := "world"
	ss1 := &StringStruct{Text: &str1}
	ss2 := &StringStruct{Text: &str2}
	if ss1.Text == ss2.Text {
		fmt.Println("string pointer comparison")
	}

	// Getter with wrong return type - should NOT trigger
	type WrongReturnStruct struct {
		Value *int64
	}
	wr1 := &WrongReturnStruct{Value: &a}
	wr2 := &WrongReturnStruct{Value: &b}
	if wr1.Value == wr2.Value {
		fmt.Println("wrong return type")
	}
}
