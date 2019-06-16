package main

import (
	"fmt"
	"time"
)

// clock()
type ClockNativeFunc struct{}

func (f *ClockNativeFunc) Call(i *Interpreter, args []interface{}) (interface{}, *RuntimeError) {
	val := time.Now().UTC().UnixNano() / 1000000
	return float64(val), nil
}

func (f *ClockNativeFunc) Arity() int {
	return 0
}

func (f *ClockNativeFunc) String() string {
	return "<native fn>"
}

// print()
type PrintNativeFunc struct{}

func (f *PrintNativeFunc) Call(i *Interpreter, args []interface{}) (interface{}, *RuntimeError) {
	fmt.Printf("%s\n", i.stringify(args[0]))
	return nil, nil
}

func (f *PrintNativeFunc) Arity() int {
	return 1
}

func (f *PrintNativeFunc) String() string {
	return "<native fn>"
}
