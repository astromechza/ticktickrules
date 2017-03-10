package ticktickrules

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Rule is a structure encoding a Cron-like rule
type Rule struct {
	minute         []int
	minuteRule     string
	hour           []int
	hourRule       string
	dayOfWeek      []int
	dayOfWeekRule  string
	dayOfMonth     []int
	dayOfMonthRule string
	month          []int
	monthRule      string
}

// rule to support */10 */0 */1
var ruleType1 = regexp.MustCompile(`^\*/\d+$`)

// rule to support 0/10/20
var ruleType2 = regexp.MustCompile(`^\d+(?:/\d+)+$`)

func parseRuleItem(r string, maxsum int) ([]int, error) {
	var out []int
	if r == "*" {
		// noop
	} else if ruleType1.MatchString(r) {

		i := strings.Split(r, "/")[1]
		v, err := strconv.Atoi(i)
		if err != nil {
			return nil, fmt.Errorf("Rule item '%s' could not be parsed", r)
		}

		if v == 0 {
			return nil, fmt.Errorf("Rule item '%s' cannot be 0", r)
		}

		if v >= maxsum {
			return nil, fmt.Errorf("Rule item '%s' does not divide", r)
		}

		var sum int
		for {
			out = append(out, sum)
			sum += v
			if sum >= maxsum {
				break
			}
		}

	} else if ruleType2.MatchString(r) {

		parts := strings.Split(r, "/")
		for _, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil {
				return nil, fmt.Errorf("Rule item '%s' could not be parsed", r)
			}
			out = append(out, v)
		}

	} else {

		v, err := strconv.Atoi(r)
		if err != nil {
			return nil, fmt.Errorf("Rule item '%s' is not supported", r)
		}
		out = append(out, v)

	}
	return out, nil
}

func validateItemsRange(items []int, min int, max int) error {
	for _, i := range items {
		if i > max {
			return fmt.Errorf("%d is > %d", i, max)
		} else if i < min {
			return fmt.Errorf("%d is < %d", i, min)
		}
	}
	return nil
}

func doesMatch(v int, vs []int) bool {
	for _, i := range vs {
		if v == i {
			return true
		}
	}
	return false
}

// NewRule constructs a new Rule structure from the cron-like arguments provided
func NewRule(minute, hour, dayOfMonth, month, dayOfWeek string) (*Rule, error) {
	output := new(Rule)

	m, err := parseRuleItem(minute, 60)
	if err != nil {
		return nil, err
	}
	output.minute = m
	if err := validateItemsRange(output.minute, 0, 59); err != nil {
		return nil, fmt.Errorf("Minute rule invalid: %s", err.Error())
	}
	output.minuteRule = minute

	h, err := parseRuleItem(hour, 24)
	if err != nil {
		return nil, err
	}
	output.hour = h
	if err := validateItemsRange(output.hour, 0, 23); err != nil {
		return nil, fmt.Errorf("Hour rule invalid: %s", err.Error())
	}
	output.hourRule = hour

	dow, err := parseRuleItem(dayOfWeek, 7)
	if err != nil {
		return nil, err
	}
	output.dayOfWeek = dow
	if err := validateItemsRange(output.dayOfWeek, 0, 7); err != nil {
		return nil, fmt.Errorf("Day of Week rule invalid: %s", err.Error())
	}
	output.dayOfWeekRule = dayOfWeek

	dom, err := parseRuleItem(dayOfMonth, 31)
	if err != nil {
		return nil, err
	}
	output.dayOfMonth = dom
	if err := validateItemsRange(output.dayOfMonth, 1, 31); err != nil {
		return nil, fmt.Errorf("Day of Month rule invalid: %s", err.Error())
	}
	output.dayOfMonthRule = dayOfMonth

	m, err = parseRuleItem(month, 24)
	if err != nil {
		return nil, err
	}
	output.month = m
	if err := validateItemsRange(output.month, 1, 12); err != nil {
		return nil, fmt.Errorf("Month rule invalid: %s", err.Error())
	}
	output.monthRule = month

	return output, nil
}

// String converts the rule back to its native cron expression
func (r *Rule) String() string {
	return fmt.Sprintf("%s %s %s %s %s", r.minuteRule, r.hourRule, r.dayOfMonthRule, r.monthRule, r.dayOfWeekRule)
}

// NextUTC returns the next time this rule is true
func (r *Rule) NextUTC() time.Time {
	return r.NextFrom(time.Now().UTC())
}

// NextFrom returns the next time this rule will run after the given time
func (r *Rule) NextFrom(from time.Time) time.Time {
	return time.Now().UTC()
}

// Matches determines whether the given time is matched by the rule
func (r *Rule) Matches(t time.Time) bool {
	if len(r.month) > 0 {
		if !doesMatch(int(t.Month()), r.month) {
			return false
		}
	}
	if len(r.dayOfWeek) > 0 {
		if !doesMatch(int(t.Weekday()), r.dayOfWeek) {
			return false
		}
	}
	if len(r.dayOfMonth) > 0 {
		if !doesMatch(t.Day(), r.dayOfMonth) {
			return false
		}
	}
	if len(r.hour) > 0 {
		if !doesMatch(t.Hour(), r.hour) {
			return false
		}
	}
	if len(r.minute) > 0 {
		if !doesMatch(t.Minute(), r.minute) {
			return false
		}
	}
	return true
}

const naiveMaxIterations = 31 * 8 * 12 * 60 * 24

// NaiveNextFrom is a very naive method of finding the next time a rule matches
func (r *Rule) NaiveNextFrom(from time.Time) time.Time {
	from = from.Add(time.Minute)
	numIterations := 0
	for {
		if r.Matches(from) {
			return from.Truncate(time.Minute)
		}
		from = from.Add(time.Minute)
		numIterations++
		if numIterations > naiveMaxIterations {
			return time.Unix(1<<62, 0)
		}
	}
}
