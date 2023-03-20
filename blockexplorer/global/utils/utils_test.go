package utils

import "testing"

func TestLabelMatching(t *testing.T) {
	type labelMatching struct {
		label  string
		rule   string
		result bool
	}
	var matchings = []labelMatching{
		{
			label:  "label",
			rule:   "label",
			result: true,
		}, {
			label:  "aslabel",
			rule:   "*label",
			result: true,
		}, {
			label:  "labelde",
			rule:   "label*",
			result: true,
		}, {
			label:  "asdalabelewe",
			rule:   "*label*",
			result: true,
		},
	}
	for _, m := range matchings {
		result := LabelMatching(m.label, m.rule)
		if result != m.result {
			t.Errorf("test with label: '%s', rule: '%s', got: %v, expected: %v", m.label, m.rule, result, m.result)
		} else {
			t.Log("Test passed")
		}
	}
}
