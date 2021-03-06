package rules

import (
	"fmt"
)

func init() {
	RegisterRule("rename-measurement", func() Config {
		return &RenameMeasurementConfig{}
	})
	RegisterRule("old-serie", func() Config {
		return &OldSerieRuleConfig{}
	})
	RegisterRule("drop-serie", func() Config {
		return &DropSerieRuleConfiguration{}
	})
}

// NewRuleFunc represents a callback to register a rule's configuration to be able to load it from toml
type NewRuleFunc func() Config

var newRuleFuncs = make(map[string]NewRuleFunc)

// RegisterRule registers a rule with the given name and config creation callback
func RegisterRule(name string, fn NewRuleFunc) {
	if _, ok := newRuleFuncs[name]; ok {
		panic(fmt.Sprintf("rule %s has already been registered", name))
	}
	newRuleFuncs[name] = fn
}

// NewRule creates a new rule configuration based on its registration name
func NewRule(name string) (Config, error) {
	fn, ok := newRuleFuncs[name]
	if !ok {
		return nil, fmt.Errorf("No registered rule '%s'", name)
	}

	return fn(), nil
}
