package main

import (
	"text/template/parse"
	"github.com/fstab/grok_exporter/plugins"
)

func newDummyFunc() plugins.FunctionWithValidator {
	return plugins.FunctionWithValidator{
		Function: dummy,
		StaticValidator: validate,
	}
}

func dummy() (string, error) {
    return "dummy value", nil
}

func validate(cmd *parse.CommandNode) error {
    return nil
}

func Generate() (string, plugins.FunctionWithValidator) { return "dummy", newDummyFunc() }
