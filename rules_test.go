package ticktickrules

import (
	"testing"
	"time"
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
	r, err := NewRule("10/20/30", "*/5", "1", "2/3", "*")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if r.String() != "10/20/30 */5 1 2/3 *" {
		t.Errorf("'%s' Did not match!", r.String())
	}
}

func TestBadMinute(t *testing.T) {
	_, err := NewRule("-1", "*", "*", "*", "*")
	if err == nil {
		t.Error("should have failed")
		return
	}
	_, err = NewRule("60", "*", "*", "*", "*")
	if err == nil {
		t.Error("should have failed")
		return
	}
}

func TestBadHour(t *testing.T) {
	_, err := NewRule("*", "-1", "*", "*", "*")
	if err == nil {
		t.Error("should have failed")
		return
	}
	_, err = NewRule("*", "60", "*", "*", "*")
	if err == nil {
		t.Error("should have failed")
		return
	}
}

func TestBadDayOfMonth(t *testing.T) {
	_, err := NewRule("*", "*", "0", "*", "*")
	if err == nil {
		t.Error("should have failed")
		return
	}
	_, err = NewRule("*", "*", "32", "*", "*")
	if err == nil {
		t.Error("should have failed")
		return
	}
}

func TestBadMonth(t *testing.T) {
	_, err := NewRule("*", "*", "*", "0", "*")
	if err == nil {
		t.Error("should have failed")
		return
	}
	_, err = NewRule("*", "*", "*", "13", "*")
	if err == nil {
		t.Error("should have failed")
		return
	}
}

func TestBadDayOfWeek(t *testing.T) {
	_, err := NewRule("*", "*", "*", "*", "-1")
	if err == nil {
		t.Error("should have failed")
		return
	}
	_, err = NewRule("*", "*", "*", "*", "8")
	if err == nil {
		t.Error("should have failed")
		return
	}
}

func TestMatchesAll(t *testing.T) {
	r, _ := NewRule("*", "*", "*", "*", "*")
	if !r.Matches(time.Now()) {
		t.Error("should not match")
	}
}

func TestMatchesMinute(t *testing.T) {
	r, _ := NewRule("16", "*", "*", "*", "*")
	if !r.Matches(time.Date(2000, 12, 1, 1, 16, 0, 0, time.UTC)) {
		t.Error("should not match")
	}
	if r.Matches(time.Date(2000, 12, 1, 1, 15, 0, 0, time.UTC)) {
		t.Error("should match")
	}
}

func TestMatchesHour(t *testing.T) {
	r, _ := NewRule("*", "21", "*", "*", "*")
	if !r.Matches(time.Date(2000, 12, 1, 21, 16, 0, 0, time.UTC)) {
		t.Error("should not match")
	}
	if r.Matches(time.Date(2000, 12, 1, 1, 15, 0, 0, time.UTC)) {
		t.Error("should match")
	}
}

func TestNaiveNext(t *testing.T) {
	r, _ := NewRule("*", "*", "*", "*", "*")
	t1 := time.Now()
	t2 := r.NaiveNextFrom(t1)
	t1 = t1.Truncate(time.Minute).Add(time.Minute)
	if t1 != t2 {
		t.Errorf("%s should have matched %s", t1, t2)
	}
}

// This is slow! seems to take about 0.2 seconds!
// TODO: remove when we have a better method
func TestNaiveNextFarFuture(t *testing.T) {
	r, _ := NewRule("*", "*", "31", "2", "*")
	t1 := time.Now()
	t2 := r.NaiveNextFrom(t1)
	if t2.Year() < 3000 {
		t.Errorf("Year should have been max")
	}
}
