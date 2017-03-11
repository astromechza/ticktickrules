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
	t2 := r.NextAfter(t1)
	t1 = t1.Truncate(time.Minute).Add(time.Minute)
	if t1 != t2 {
		t.Errorf("%s should have matched %s", t1, t2)
	}
}

func TestNaiveNextFarFuture(t *testing.T) {
	r, _ := NewRule("*", "*", "31", "2", "*")
	t1 := time.Now()
	t2 := r.NextAfter(t1)
	if t2.Year() < 3000 {
		t.Errorf("Year should have been max")
	}
}

func TestNaiveMultiple(t *testing.T) {
	start := time.Date(2000, 1, 1, 1, 0, 1, 0, time.UTC)
	r, _ := NewRule("*/25", "*/2", "*", "*", "*")
	n1 := r.NextAfter(start)
	e1 := time.Date(2000, 1, 1, 2, 25, 0, 0, time.UTC)
	if n1 != e1 {
		t.Errorf("n1 %s != %s", n1, e1)
		return
	}
	n2 := r.NextAfter(n1)
	e2 := time.Date(2000, 1, 1, 2, 50, 0, 0, time.UTC)
	if n2 != e2 {
		t.Errorf("n2 %s != %s", n2, e2)
		return
	}
	n3 := r.NextAfter(n2)
	e3 := time.Date(2000, 1, 1, 4, 0, 0, 0, time.UTC)
	if n3 != e3 {
		t.Errorf("n3 %s != %s", n3, e3)
		return
	}
	n4 := r.NextAfter(n3)
	e4 := time.Date(2000, 1, 1, 4, 25, 0, 0, time.UTC)
	if n4 != e4 {
		t.Errorf("n4 %s != %s", n4, e4)
		return
	}
	n5 := r.NextAfter(n4)
	e5 := time.Date(2000, 1, 1, 4, 50, 0, 0, time.UTC)
	if n5 != e5 {
		t.Errorf("n5 %s != %s", n5, e5)
		return
	}
}

func TestNaiveMultipleDays(t *testing.T) {
	start := time.Date(2000, 2, 28, 23, 59, 0, 0, time.UTC)
	r, _ := NewRule("*", "*", "31", "*", "*")
	n1 := r.NextAfter(start)
	e1 := time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC)
	if n1 != e1 {
		t.Errorf("n1 %s != %s", n1, e1)
		return
	}
}
