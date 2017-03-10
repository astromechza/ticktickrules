package ticktickrules

import (
	"testing"
)

func TestRuleConstruct(t *testing.T) {
	r, err := NewRule("*", "*", "*", "*", "*")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if r.String() != "* * * * *" {
		t.Errorf("'%s' Did not match!", r.String())
	}
}

func TestRuleConstructExtra(t *testing.T) {
	r, err := NewRule("10/20/30", "*/5", "1", "*", "*")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if r.String() != "10/20/30 */5 * * 1" {
		t.Errorf("'%s' Did not match!", r.String())
	}
}
