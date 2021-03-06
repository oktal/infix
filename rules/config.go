package rules

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
	"github.com/oktal/infix/filter"
)

// Config represents a configuration for a rule
type Config interface {
	Sample() string

	Build() (Rule, error)
}

// LoadConfig will load rules from a TOML configuration file
func LoadConfig(path string) ([]Rule, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	table, err := toml.Parse(data)
	if err != nil {
		return nil, err
	}

	var rules []Rule

	for name, val := range table.Fields {
		subTable, ok := val.(*ast.Table)
		if !ok {
			return nil, fmt.Errorf("%s: invalid configuration %s", path, name)
		}

		switch name {
		case "rules":
			for ruleName, ruleVal := range subTable.Fields {
				ruleSubTable, ok := ruleVal.([]*ast.Table)
				if !ok {
					return nil, fmt.Errorf("%s: invalid configuration %s", path, ruleName)
				}

				for _, r := range ruleSubTable {
					rule, err := loadRule(ruleName, r)
					if err != nil {
						return nil, fmt.Errorf("%s: %s: %s", path, ruleName, err)
					}
					rules = append(rules, rule)
				}
			}
		case "filters":
		default:
			return nil, fmt.Errorf("%s: unsupported config file format %s", path, name)
		}
	}

	return rules, nil
}

func loadRule(name string, table *ast.Table) (Rule, error) {
	config, err := NewRule(name)
	if err != nil {
		return nil, err
	}

	if err := unmarshalFilters(table, config); err != nil {
		return nil, err
	}

	if err := toml.UnmarshalTable(table, config); err != nil {
		return nil, err
	}

	return config.Build()
}

func unmarshalFilters(table *ast.Table, config Config) error {
	e := reflect.ValueOf(config).Elem()
	filterType := reflect.TypeOf((*filter.Filter)(nil)).Elem()

	for i := 0; i < e.NumField(); i++ {
		field := e.Type().Field(i)
		varName := field.Name
		varType := field.Type

		if varType.Implements(filterType) {
			f, err := filter.Unmarshal(table, varName)
			if err != nil {
				return err
			}
			e.Field(i).Set(reflect.ValueOf(f))
		}

	}

	return nil
}
