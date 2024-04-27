package rules

import (
	"github.com/gainings/tfirg/internal/model/resourceid"
)

const defaultRuleName = ""

var rules map[string]resourceid.Generator

type Rules struct {
	rules map[string]resourceid.Generator
}

func NewRules() resourceid.GeneratorFactory {
	return Rules{
		rules: rules,
	}
}

func (r Rules) Generator(typeName string) resourceid.Generator {
	if _, ok := r.rules[typeName]; ok {
		return r.rules[typeName]
	}
	return r.rules[defaultRuleName]

}

type DefaultRule struct{}

func (r DefaultRule) Generate(attributes map[string]interface{}) *resourceid.ResourceID {
	id := resourceid.ResourceID(attributes["id"].(string))
	return &id
}

func init() {
	rules = make(map[string]resourceid.Generator)
	Register(defaultRuleName, DefaultRule{})
}

func Register(typeName string, g resourceid.Generator) {
	rules[typeName] = g
}
