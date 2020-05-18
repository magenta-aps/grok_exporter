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
	textTemplate "text/template"
	"text/template/parse"
    "runtime"
    "io/ioutil"
    "log"
    "path/filepath"
    "plugin"
    "github.com/fstab/grok_exporter/plugins"
)

var funcs functions = make(map[string]plugins.FunctionWithValidator)

func load_plugin(plugin_name string) {
    p, err := plugin.Open("./dynamic/" + plugin_name)
    if err != nil {
        panic(err)
    }
    f, err := p.Lookup("Generate")
    if err != nil {
        panic(err)
    }
    label, func_val := f.(func() (string, plugins.FunctionWithValidator))()
    funcs.add(label, func_val)
}

func read_plugin_folder() {
    files, err := ioutil.ReadDir("./dynamic/")
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files {
        if filepath.Ext(f.Name()) == ".so" {
            load_plugin(f.Name())
        }
    }
}

func init() {
	funcs.add("timestamp", newTimestampFunc())
	funcs.add("gsub", newGsubFunc())
	funcs.add("add", newAddFunc())
	funcs.add("subtract", newSubtractFunc())
	funcs.add("multiply", newMultiplyFunc())
	funcs.add("divide", newDivideFunc())
	funcs.add("base", newBaseFunc())

    if runtime.GOOS != "windows" {
        read_plugin_folder()
    }
}

type functions map[string]plugins.FunctionWithValidator

func (funcs functions) add(name string, f plugins.FunctionWithValidator) {
	funcs[name] = f
}

func (funcs functions) toFuncMap() textTemplate.FuncMap {
	result := make(textTemplate.FuncMap, len(funcs))
	for name, f := range funcs {
		result[name] = f.Function
	}
	return result
}

func (funcs functions) validate(name string, cmd *parse.CommandNode) error {
	f, exists := funcs[name]
	if !exists {
		return nil // not one of our custom functions, skip validation
	}
	return f.StaticValidator(cmd)
}
