package importas

import (
	"fmt"
	"regexp"
)

type Config struct {
	RequiredAlias map[string]string
	Rules         []*Rule
}

func (c *Config) CompileRegexp() error {
	rules := make([]*Rule, 0, len(c.RequiredAlias))
	for path, alias := range c.RequiredAlias {
		reg, err := regexp.Compile(path)
		if err != nil {
			return err
		}

		rules = append(rules, &Rule{
			Regexp: reg,
			Alias:  alias,
		})
	}

	c.Rules = rules
	return nil
}

func (c Config) findRule(path string) *Rule {
	for _, rule := range c.Rules {
		if rule.Regexp.MatchString(path) {
			return rule
		}
	}

	return nil
}

func (c Config) AliasFor(path string) (string, bool) {
	rule := c.findRule(path)
	if rule == nil {
		return "", false
	}

	alias, err := rule.aliasFor(path)
	if err != nil {
		return "", false
	}

	return alias, true
}

type Rule struct {
	Alias  string
	Regexp *regexp.Regexp
}

func (r *Rule) aliasFor(path string) (string, error) {
	str := r.Regexp.FindString(path)
	if len(str) > 0 {
		return r.Regexp.ReplaceAllString(str, r.Alias), nil
	}

	return "", fmt.Errorf("mismatch rule")
}
