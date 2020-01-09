package filter

import (
	"github.com/dpakach/gorkin/lexer"
	"github.com/dpakach/gorkin/parser"
	"testing"
)

func areArrayEqual(a, b []string) bool {
	if len(a) == 0 {
		if len(b) == 0 {
			return true
		}
		return false
	}
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func areMapEqual(a, b map[string]string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestLineFilterGetStart(t *testing.T) {
	testData := []struct {
		input         string
		expectedStart int
		expectedEnd   int
	}{
		{"12", 12, 12},
		{"12-34", 12, 34},
		{"100", 100, 100},
		{"1", 1, 1},
	}

	for _, tt := range testData {
		filter := LineFilter{tt.input}
		if filter.getStart() != tt.expectedStart {
			t.Fatalf("Start doesnt Match, expected: %v, got: %v", tt.expectedStart, filter.getStart())
		}
		if filter.getEnd() != tt.expectedEnd {
			t.Fatalf("End doesnt Match, expected: %v, got: %v", tt.expectedEnd, filter.getEnd())
		}
	}
}

func TestLineFilterMatchFeature(t *testing.T) {

	input := `
	# This is feature

	Feature: test

	Scenario: test
		When test
	`

	l := lexer.New(input)
	p := parser.New(l)
	feature := p.ParseFeature()

	testData := []struct {
		input         string
		expectedMatch bool
	}{
		{"100", false},
		{"1-100", true},
		{"2-5", true},
		//{"addd", false},
		//{"1-200-300", false},
		{"4", true},
		{"5", false},
	}

	for _, tt := range testData {
		lf := &LineFilter{tt.input}
		match := lf.MatchFeature(feature)
		if match != tt.expectedMatch {
			t.Fatalf("Match feature for %q incorrect, expected: %v, got: %v", tt.input, tt.expectedMatch, match)
		}
	}
}

func TestGetTags(t *testing.T) {
	testdata := []struct {
		input           string
		expectedTags    []string
		expectedNotTags []string
	}{
		{"@tag1&&@tag2&&~@tag3", []string{"tag1", "tag2"}, []string{"tag3"}},
		{"@tag1&&@tag2&&@tag3", []string{"tag1", "tag2", "tag3"}, []string{}},
		{"~@tag1&&~@tag2&&~@tag3", []string{}, []string{"tag1", "tag2", "tag3"}},
		{"~@tag1&&@tag2&&~@tag3", []string{"tag2"}, []string{"tag1", "tag3"}},
	}

	for _, tt := range testdata {
		filter := &TagFilter{tt.input}
		tags := filter.getTags()
		tagsNot := filter.getTagsNot()

		if !areArrayEqual(tags, tt.expectedTags) {
			t.Fatalf("Tags are not equal, expected: %v, got: %v", tt.expectedTags, tags)
		}
		if !areArrayEqual(tagsNot, tt.expectedNotTags) {
			t.Fatalf("Tags not are not equal, expected: %v, got: %v", tt.expectedNotTags, tagsNot)
		}
	}
}

func TesttagMatchPresent(t *testing.T) {
	testdata := []struct {
		input     string
		inputTags []string
		match     bool
	}{
		{"@tag1&&@tag2&&~@tag3", []string{"tag1", "tag2"}, true},
		{"@tag1&&@tag2&&@tag3", []string{"tag1", "tag2", "tag3"}, true},
		{"~@tag1&&~@tag2&&~@tag3", []string{}, true},
		{"~@tag1&&@tag2&&~@tag3", []string{"tag2"}, true},
		{"@tag1&&@tag2&&~@tag3", []string{"tag3"}, false},
		{"@tag1&&@tag2&&@tag3", []string{}, false},
		{"~@tag1&&~@tag2&&~@tag3", []string{"tag1", "tag2", "tag3"}, false},
		{"~@tag1&&@tag2&&~@tag3", []string{"tag1", "tag3"}, false},

		{"@tag1", []string{"tag1"}, true},
		{"@tag1", []string{}, true},
		{"~@tag1", []string{"tag1"}, false},
		{"@tag1", []string{"tag2"}, false},
	}

	for _, tt := range testdata {
		filter := TagFilter{tt.input}

		if filter.tagsMatchPresent(tt.inputTags) != tt.match {
			t.Fatalf("Invalid match for filter %q, expected: %v, got: %v", tt.input, tt.match, filter.tagsMatchPresent(tt.inputTags))
		}
	}
}
