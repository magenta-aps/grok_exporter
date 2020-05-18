// Copyright 2018-2020 The grok_exporter Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package template

import (
	"fmt"
	"reflect"
	"strconv"
	"text/template/parse"
    "github.com/fstab/grok_exporter/plugins"
)

func newAddFunc() plugins.FunctionWithValidator {
	return plugins.FunctionWithValidator{
		Function: add,
		StaticValidator: func(cmd *parse.CommandNode) error {
			return validate("add", cmd)
		},
	}
}

func newSubtractFunc() plugins.FunctionWithValidator {
	return plugins.FunctionWithValidator{
		Function: subtract,
		StaticValidator: func(cmd *parse.CommandNode) error {
			return validate("subtract", cmd)
		},
	}
}

func newMultiplyFunc() plugins.FunctionWithValidator {
	return plugins.FunctionWithValidator{
		Function: multiply,
		StaticValidator: func(cmd *parse.CommandNode) error {
			return validate("multiply", cmd)
		},
	}
}

func newDivideFunc() plugins.FunctionWithValidator {
	return plugins.FunctionWithValidator{
		Function: divide,
		StaticValidator: func(cmd *parse.CommandNode) error {
			return validate("divide", cmd)
		},
	}
}

func add(a, b interface{}) (float64, error) {
	aFloat, bFloat, err := toFloats(a, b)
	if err != nil {
		return 0, fmt.Errorf("error executing add function: %v", err)
	}
	return aFloat + bFloat, nil
}

func subtract(a, b interface{}) (float64, error) {
	aFloat, bFloat, err := toFloats(a, b)
	if err != nil {
		return 0, fmt.Errorf("error executing subtract function: %v", err)
	}
	return aFloat - bFloat, nil
}

func multiply(a, b interface{}) (float64, error) {
	aFloat, bFloat, err := toFloats(a, b)
	if err != nil {
		return 0, fmt.Errorf("error executing multiply function: %v", err)
	}
	return aFloat * bFloat, nil
}

func divide(a, b interface{}) (float64, error) {
	aFloat, bFloat, err := toFloats(a, b)
	if err != nil {
		return 0, fmt.Errorf("error executing divide function: %v", err)
	}
	if bFloat == 0 {
		return 0, fmt.Errorf("error executing divide function: division by zero")
	}
	return aFloat / bFloat, nil
}

func toFloats(a, b interface{}) (float64, float64, error) {
	floatA, err := toFloat(a)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot convert %v to floating point number: %v", a, err)
	}
	floatB, err := toFloat(b)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot convert %v to floating point number: %v", b, err)
	}
	return floatA, floatB, nil
}

func toFloat(f interface{}) (float64, error) {
	val := reflect.ValueOf(f)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(val.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return val.Float(), nil
	case reflect.String:
		return strconv.ParseFloat(val.String(), 64)
	}
	if val, ok := f.(fmt.Stringer); ok {
		return strconv.ParseFloat(val.String(), 64)
	}
	return 0, fmt.Errorf("%T: unknown type", f)
}

func validate(functionName string, cmd *parse.CommandNode) error {
	prefix := fmt.Sprintf("syntax error in %v call", functionName)
	if len(cmd.Args) != 3 {
		return fmt.Errorf("%v: expected two parameters, but found %v parameters", prefix, len(cmd.Args)-1)
	}
	// If a param is a string or number, we check if we can parse it.
	// Otherwise it might be a variable of a function call, we cannot check this statically.
	for _, paramPos := range []int{1, 2} {
		switch param := cmd.Args[paramPos].(type) {
		case *parse.NumberNode:
			if !param.IsFloat {
				return fmt.Errorf("%v: unable to parse %v as a floating point number", prefix, param)
			}
		case *parse.StringNode:
			if _, err := strconv.ParseFloat(param.Text, 64); err != nil {
				return fmt.Errorf("%v: unable to parse %v as a floating point number: %v", prefix, param, err)
			}
		}
	}
	return nil
}
